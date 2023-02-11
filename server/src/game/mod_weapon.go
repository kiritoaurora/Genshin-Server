package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server_logic/src/csvs"
)

/*
	武器模块
*/

type Weapon struct {
	WeaponId    int //武器Id
	KeyId       int //武器唯一标识
	Level       int //武器等级
	Exp         int //武器当前拥有的经验
	StarLevel   int //武器突破星级
	RefineLevel int //精炼等级
	RoleId      int //武器装备的角色Id
}

type ModWeapon struct {
	WeaponInfo map[int]*Weapon
	MaxKey     int

	player *Player
	path   string
}

func (w *ModWeapon) AddItem(itemId int, num int64) {

	config := csvs.GetWeaponConfig(itemId)
	if config == nil {
		fmt.Println("配置不存在")
		return
	}

	if len(w.WeaponInfo)+int(num) > csvs.WEAPON_MAX_COUNT {
		fmt.Println("超过最大值")
		return
	}

	for i := int64(0); i < num; i++ {
		weapon := new(Weapon)
		weapon.WeaponId = itemId
		w.MaxKey++
		weapon.KeyId = w.MaxKey
		w.WeaponInfo[weapon.KeyId] = weapon
		fmt.Println("获得武器：", csvs.GetItemConfig(itemId).ItemName, "---武器编号：", weapon.KeyId)
	}
}

// 武器升级
func (w *ModWeapon) WeaponUp(keyId, exp int) {
	weapon := w.WeaponInfo[keyId]
	if weapon == nil {
		return
	}
	weaponConfig := csvs.ConfigWeaponMap[weapon.WeaponId]
	if weaponConfig == nil {
		return
	}
	weapon.Exp += exp
	for {
		//下一级配置
		nextLevelConfig := csvs.GetWeaponLevelConfig(weaponConfig.Star, weapon.Level+1)
		if nextLevelConfig == nil {
			fmt.Println("返还经验:", weapon.Exp)
			weapon.Exp = 0
			break
		}
		if weapon.StarLevel < nextLevelConfig.NeedStarLevel {
			fmt.Println("返还经验:", weapon.Exp)
			weapon.Exp = 0
			break
		}
		if weapon.Exp < nextLevelConfig.NeedExp {
			break
		}
		weapon.Level++
		weapon.Exp -= nextLevelConfig.NeedExp
	}
	weapon.ShowInfo()
}

func (w *Weapon) ShowInfo() {
	fmt.Printf("key:%d,Id:%d\n", w.KeyId, w.WeaponId)
	fmt.Printf("当前等级:%d,当前经验:%d,当前突破等级:%d,当前精炼等级:%d\n", w.Level, w.Exp, w.StarLevel, w.RefineLevel)
}

// 武器突破
func (w *ModWeapon) WeaponUpStar(keyId int) {
	weapon := w.WeaponInfo[keyId]
	if weapon == nil {
		return
	}
	weaponConfig := csvs.ConfigWeaponMap[weapon.WeaponId]
	if weaponConfig == nil {
		return
	}
	nextStarConfig := csvs.GetWeaponStarConfig(weaponConfig.Star, weapon.StarLevel+1)
	if nextStarConfig == nil {
		return
	}

	//TODO:验物品充足并扣除

	if weapon.Level < nextStarConfig.Level {
		fmt.Println("武器等级不够，无法突破")
		return
	}
	weapon.StarLevel++
	weapon.ShowInfo()
}

// 武器精炼
func (w *ModWeapon) WeaponUpRefine(weaponKeyId int, weaponTargetKeyId int) {
	if weaponKeyId == weaponTargetKeyId {
		fmt.Println("错误的材料")
		return
	}
	weapon := w.WeaponInfo[weaponKeyId]
	if weapon == nil {
		return
	}
	weaponTarget := w.WeaponInfo[weaponTargetKeyId]
	if weaponTarget == nil {
		return
	}
	if weapon.WeaponId != weaponTarget.WeaponId {
		fmt.Println("错误的材料")
		return
	}
	if weapon.RefineLevel >= csvs.WEAPON_MAX_REFINE {
		fmt.Println("已达到最大精炼等级")
		return
	}
	// TODO：验摩拉是否充足并扣除
	weapon.RefineLevel++
	delete(w.WeaponInfo, weaponTargetKeyId)
	weapon.ShowInfo()
}

func (w *ModWeapon) SaveData() {
	content, err := json.Marshal(w)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(w.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (w *ModWeapon) LoadData(player *Player) {
	w.player = player
	w.path = w.player.localPath + "/weapon.json"

	configFile, err := ioutil.ReadFile(w.path)
	if err != nil {
		fmt.Println("暂无存档")
		// return
	}
	err = json.Unmarshal(configFile, &w)
	if err != nil {
		w.InitData()
		return
	}

	// if w.WeaponInfo == nil {
	// 	w.WeaponInfo = make(map[int]*Weapon)
	// }
}

func (w *ModWeapon) InitData() {
	if w.WeaponInfo == nil {
		w.WeaponInfo = make(map[int]*Weapon)
	}
}
