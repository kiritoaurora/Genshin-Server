package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server_logic/src/csvs"
	"time"
)

/*
	角色模块
*/

type RoleInfo struct {
	RoleId   int
	GetTimes int
	//等级、经验、圣遗物、武器
	RelicsInfo []int
	WeaponInfo int
}

type ModRole struct {
	RoleInfo  map[int]*RoleInfo
	HpPool    int
	HpCalTime int64

	player *Player
	path   string
}

func (r *ModRole) IsHasRole(roleId int) bool {
	return true
}

func (r *ModRole) GetRoleLevel(roleId int) int {
	return 80
}

func (r *ModRole) AddItem(roleId int, num int64) {
	config := csvs.GetRoleConfig(roleId)
	if config == nil {
		fmt.Println("配置不存在roleId:", roleId)
		return
	}

	for i := 0; i < int(num); i++ {
		// 是否已拥有
		_, ok := r.RoleInfo[roleId]
		if !ok { //未获得
			data := new(RoleInfo)
			data.RoleId = roleId
			data.GetTimes = 1
			r.RoleInfo[roleId] = data
		} else {
			// 判断实际获得的东西，五星：命座*1+星辉*10 | 星辉*25，四星：命座*1+星辉*2 | 星辉*10
			r.RoleInfo[roleId].GetTimes++
			if r.RoleInfo[roleId].GetTimes >= csvs.ADD_ROLE_TIME_NORMAL_MIN &&
				r.RoleInfo[roleId].GetTimes <= csvs.ADD_ROLE_TIME_NORMAL_MAX {
				r.player.GetModBag().AddItemToBag(config.Stuff, config.StuffNum)
				r.player.GetModBag().AddItemToBag(config.StuffItem, config.StuffItemNum)
			} else {
				r.player.GetModBag().AddItemToBag(config.MaxStuffItem, config.MaxStuffItemNum)
			}
		}
	}
	itemConfig := csvs.GetItemConfig(roleId)
	if itemConfig != nil {
		fmt.Println("获得角色：", itemConfig.ItemName, "---", r.RoleInfo[roleId].GetTimes, "次")
	}
	// 获取角色头像和名片
	r.player.GetModIcon().CheckGetIcon(roleId)
	r.player.GetModCard().CheckGetCard(roleId, 10)
}

func (r *ModRole) HandleSendRoleInfo() {
	fmt.Println("当前拥有角色信息如下：")
	for _, v := range r.RoleInfo {
		v.SendRoleInfo()
	}
}

func (r *RoleInfo) SendRoleInfo() {
	fmt.Printf("%s,角色ID:%d,累计获得次数:%d\n", csvs.GetItemConfig(r.RoleId).ItemName, r.RoleId, r.GetTimes)
}

// 获取角色信息，用于仓检测试
func (r *ModRole) GetRoleInfoForPoolCheck() (map[int]int, map[int]int) {
	fiveStarInfo := make(map[int]int)
	fourStarInfo := make(map[int]int)

	for _, v := range r.RoleInfo {
		roleConfig := csvs.GetRoleConfig(v.RoleId)
		if roleConfig == nil {
			continue
		}
		if roleConfig.Star == 5 {
			fiveStarInfo[roleConfig.RoleId] = v.GetTimes
		} else if roleConfig.Star == 4 {
			fourStarInfo[roleConfig.RoleId] = v.GetTimes
		}
	}

	return fiveStarInfo, fourStarInfo
}

func (r *ModRole) CalHpPool() {
	if r.HpCalTime == 0 {
		r.HpCalTime = time.Now().Unix()
	}

	calTime := time.Now().Unix() - r.HpCalTime
	r.HpPool += int(calTime) * 10
	r.HpCalTime = time.Now().Unix()
	fmt.Println("当前血池回复量:", r.HpPool)
}

//装备圣遗物
func (r *ModRole) WearRelics(roleInfo *RoleInfo, relics *Relics) {
	relicsConfig := csvs.GetRelicsConfig(relics.RelicsId) //获取圣遗物配置表
	if relicsConfig == nil {
		return
	}
	r.CheckRelicsPos(roleInfo, relicsConfig.Pos) //检查圣遗物位置
	if relicsConfig.Pos < 0 || relicsConfig.Pos > len(roleInfo.RelicsInfo) {
		return
	}

	curRelicsKeyId := roleInfo.RelicsInfo[relicsConfig.Pos-1] //当前角色装备的圣遗物ID
	otherRoleId := relics.RoleId                              //待替换的圣遗物所装备的角色ID
	// 角色该位置已有圣遗物，且待替换圣遗物也已装备在其他角色身上
	if curRelicsKeyId != 0 && otherRoleId != 0 {
		otherRoleInfo := r.RoleInfo[otherRoleId]                         //待替换的圣遗物所装备的角色信息
		otherRoleInfo.RelicsInfo[relicsConfig.Pos-1] = curRelicsKeyId    //替换对应角色的圣遗物为当前角色装备的圣遗物
		r.player.GetModRelics().RelicsInfo[curRelicsKeyId].RoleId = otherRoleId //圣遗物装备角色ID改变
	} else if curRelicsKeyId != 0 && otherRoleId == 0 { //角色该位置已装备圣遗物，待替换圣遗物未装备在其他角色身上
		curRelics := r.player.GetModRelics().RelicsInfo[curRelicsKeyId] //已装备的圣遗物
		r.TakeOffRelics(roleInfo, curRelics)             //卸下该圣遗物
	} else if curRelicsKeyId == 0 && otherRoleId != 0 { //角色未装备圣遗物，待装备圣遗物已装备在其他角色身上
		otherRoleInfo := r.RoleInfo[otherRoleId]         //待装备的圣遗物所装备的角色信息
		otherRoleInfo.RelicsInfo[relicsConfig.Pos-1] = 0 //卸下圣遗物
	}

	//装备圣遗物
	roleInfo.RelicsInfo[relicsConfig.Pos-1] = relics.KeyId
	relics.RoleId = roleInfo.RoleId //绑定圣遗物所属角色
	roleInfo.ShowInfo(r.player)
}

//检查圣遗物位置
func (r *ModRole) CheckRelicsPos(roleInfo *RoleInfo, pos int) {
	nowSize := len(roleInfo.RelicsInfo)
	needAdd := pos - nowSize
	//是否需要对roleInfo的RelicsInfo扩容
	for i := 0; i < needAdd; i++ {
		roleInfo.RelicsInfo = append(roleInfo.RelicsInfo, 0)
	}
}

func (r *RoleInfo) ShowInfo(player *Player) {
	fmt.Printf("当前角色:%s,角色ID:%d\n", csvs.GetItemConfig(r.RoleId).ItemName, r.RoleId)
	suitMap := make(map[int]int)
	for _, v := range r.RelicsInfo {
		relicsNow := player.GetModRelics().RelicsInfo[v]
		if relicsNow == nil {
			fmt.Println("未装备")
			continue
		}
		fmt.Printf("%s, key:%d\n", csvs.GetItemConfig(relicsNow.RelicsId).ItemName, v)
		relicsNowConfig := csvs.GetRelicsConfig(relicsNow.RelicsId)
		if relicsNowConfig != nil {
			suitMap[relicsNowConfig.Type]++
		}
	}

	suitSkill := make([]int, 0)
	for suit, num := range suitMap {
		for _, config := range csvs.ConfigRelicsSuitMap[suit] {
			if num >= config.Num {
				suitSkill = append(suitSkill, config.SuitSkill)
			}
		}
	}
	for _, v := range suitSkill {
		fmt.Printf("激活套装效果:%d\n", v)
	}

}

//卸下圣遗物
func (r *ModRole) TakeOffRelics(roleInfo *RoleInfo, relics *Relics) {
	relicsConfig := csvs.GetRelicsConfig(relics.RelicsId) //获取圣遗物配置表
	if relicsConfig == nil {
		return
	}
	r.CheckRelicsPos(roleInfo, relicsConfig.Pos) //检查圣遗物位置
	if relicsConfig.Pos < 0 || relicsConfig.Pos > len(roleInfo.RelicsInfo) {
		return
	}
	if roleInfo.RelicsInfo[relicsConfig.Pos-1] != relics.KeyId {
		fmt.Println("当前角色未装备这个物品")
		return
	}
	//卸下圣遗物
	roleInfo.RelicsInfo[relicsConfig.Pos-1] = 0
	relics.RoleId = 0 //绑定圣遗物所属角色
	roleInfo.ShowInfo(r.player)
}

//装备武器
func (r *ModRole) WearWeapon(roleInfo *RoleInfo, weaponInfo *Weapon) {
	weaponConfig := csvs.GetWeaponConfig(weaponInfo.WeaponId)
	if weaponConfig == nil {
		fmt.Println("数据异常，武器配置不存在")
		return
	}

	// 检查角色与武器匹配
	roleConfig := csvs.GetRoleConfig(roleInfo.RoleId)
	if roleConfig.Type != weaponConfig.Type {
		fmt.Println("武器和角色不匹配")
		return
	}

	// 武器替换
	curWeaponKeyId := roleInfo.WeaponInfo
	otherRoleId := weaponInfo.RoleId
	// 该角色已装备武器，且该武器也已被装备
	if curWeaponKeyId != 0 && otherRoleId != 0 {
		otherRole := r.RoleInfo[otherRoleId]                             //装备该武器的角色
		otherRole.WeaponInfo = curWeaponKeyId                            //更改角色的武器为当前角色的武器
		r.player.GetModWeapon().WeaponInfo[curWeaponKeyId].RoleId = otherRoleId //武器绑定的角色变更
	} else if curWeaponKeyId != 0 && otherRoleId == 0 { //该角色已装备武器，且该武器未被装备
		curWeapon := r.player.GetModWeapon().WeaponInfo[curWeaponKeyId] //装备的武器
		r.TakeOffWeapon(roleInfo, curWeapon)             //卸下武器
	} else if curWeaponKeyId == 0 && otherRoleId != 0 { //角色未装备武器，且该武器被装备
		otherRole := r.RoleInfo[otherRoleId] //装备该武器的角色
		otherRole.WeaponInfo = 0             //卸下武器
	}

	//装备武器，并为武器绑定角色
	roleInfo.WeaponInfo = weaponInfo.KeyId
	weaponInfo.RoleId = roleInfo.RoleId
	roleInfo.ShowWeaponInfo(r.player)
}

//卸下武器
func (r *ModRole) TakeOffWeapon(roleInfo *RoleInfo, weaponInfo *Weapon) {
	weaponConfig := csvs.GetWeaponConfig(weaponInfo.WeaponId)
	if weaponConfig == nil {
		return
	}
	if roleInfo.WeaponInfo != weaponInfo.KeyId {
		fmt.Println("角色未装备该武器")
		return
	}
	//卸下武器，且解除武器的绑定角色
	roleInfo.ShowWeaponInfo(r.player)
	roleInfo.WeaponInfo = 0
	weaponInfo.RoleId = 0
	roleInfo.ShowWeaponInfo(r.player)
}

func (r *RoleInfo) ShowWeaponInfo(player *Player) {
	fmt.Printf("当前角色:%s,角色ID:%d,武器key:%d\n", csvs.GetItemConfig(r.RoleId).ItemName, r.RoleId, r.WeaponInfo)
	if r.WeaponInfo == 0 {
		fmt.Println("未装备武器")
		return
	}
	weaponId := player.GetModWeapon().WeaponInfo[r.WeaponInfo].WeaponId
	fmt.Printf("武器:%s, key:%d\n\n", csvs.GetItemConfig(weaponId).ItemName, r.WeaponInfo)
}

func (r *ModRole) SaveData() {
	content, err := json.Marshal(r)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(r.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (r *ModRole) LoadData(player *Player) {
	r.player = player
	r.path = r.player.localPath + "/role.json"

	configFile, err := ioutil.ReadFile(r.path)
	if err != nil {
		fmt.Println("暂无存档")
		// return
	}
	err = json.Unmarshal(configFile, &r)
	if err != nil {
		r.InitData()
		return
	}
	// if r.RoleInfo == nil {
	// 	r.RoleInfo = make(map[int]*RoleInfo)
	// }	
}

func (r *ModRole) InitData() {
	if r.RoleInfo == nil {
		r.RoleInfo = make(map[int]*RoleInfo)
	}
}
