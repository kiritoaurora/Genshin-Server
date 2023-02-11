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
	基础信息模块
*/

type ShowRole struct {
	RoleId    int
	RoleLevel int
}

type ModPlayer struct {
	UserId         int64       //唯一id
	Icon           int         //头像
	Card           int         //名片
	Name           string      //昵称
	Sign           string      //签名
	PlayerLevel    int         //玩家等级
	PlayerExp      int         //冒险阅历
	WorldLevel     int         //世界等级
	WorldLevelNow  int         //当前世界等级
	WorldLevelCool int64       //操作世界等级的冷却时间
	Birth          int         //生日
	ShowTeam       []*ShowRole //展示阵容
	HideShowTeam   int         //隐藏阵容开关
	ShowCard       []int       //展示名片
	//看不见的字段
	Prohibit int //封禁状态
	IsGM     int //GM账号标志

	player *Player
	path   string
}

//对外接口
// 设置头像
func (p *ModPlayer) SetIcon(iconId int) {
	if !p.player.GetModIcon().IsHasIcon(iconId) {
		// 通知客户端，操作非法
		fmt.Println("未获得头像：", iconId)
		return
	}
	p.Icon = iconId
	fmt.Println("当前图标", p.Icon)
}

// 设置名片
func (p *ModPlayer) SetCard(cardId int) {
	if !p.player.GetModCard().IsHasCard(cardId) {
		//通知客户端
		return
	}
	p.Card = cardId
	fmt.Println("当前名片", p.Card)
}

// 设置昵称
func (p *ModPlayer) SetName(name string) {
	if GetServer().IsBanWord(name) {
		return
	}
	p.Name = name
	fmt.Println("当前名字", p.Name)
}

//设置签名
func (p *ModPlayer) SetSign(sign string) {
	if GetServer().IsBanWord(sign) {
		return
	}
	p.Sign = sign
	fmt.Println("当前签名", p.Sign)
}

// 内置接口
// 加冒险阅历
func (p *ModPlayer) AddExp(exp int) {
	p.PlayerExp += exp

	for {
		//获取当前等级信息
		config := csvs.GetNowLevelConfig(p.PlayerLevel)
		//等级不正常 <0  >60
		if config == nil {
			break
		}
		//已达到60级
		if config.PlayerExp == 0 {
			break
		}
		//是否完成任务
		if config.ChapterId > 0 && !p.player.GetModUniqueTask().IsTaskFinish(config.ChapterId) {
			break
		}
		if p.PlayerExp >= config.PlayerExp {
			p.PlayerLevel += 1
			p.PlayerExp -= config.PlayerExp
		} else {
			break
		}
	}
	fmt.Println("当前等级：", p.PlayerLevel, "---当前经验：", p.PlayerExp)

}

// 降低世界等级
func (p *ModPlayer) ReduceWorldLevel() {
	if p.WorldLevel < csvs.REDUCE_WORLD_LEVEL_START {
		fmt.Println("操作失败，---当前世界等级：", p.WorldLevel)
		return
	}

	if p.WorldLevel-p.WorldLevelNow >= csvs.REDUCE_WORLD_LEVEL_MAX {
		fmt.Println("降级操作失败，---当前世界等级：", p.WorldLevelNow, "---真实世界等级", p.WorldLevel)
		return
	}

	if time.Now().Unix() < p.WorldLevelCool {
		fmt.Println("操作失败，冷却中")
		return
	}

	p.WorldLevelNow -= 1
	p.WorldLevelCool = time.Now().Unix() + csvs.REDUCE_WORLD_LEVEL_COOL_TIME
	fmt.Println("操作成功，---当前世界等级：", p.WorldLevelNow, "---真实世界等级", p.WorldLevel)

}

// 恢复世界等级
func (p *ModPlayer) ReturnWorldLevel() {
	if p.WorldLevel == p.WorldLevelNow {
		fmt.Println("恢复操作失败，---当前世界等级：", p.WorldLevelNow, "---真实世界等级", p.WorldLevel)
		return
	}

	if time.Now().Unix() < p.WorldLevelCool {
		fmt.Println("操作失败，冷却中")
		return
	}

	p.WorldLevelNow += 1
	p.WorldLevelCool = time.Now().Unix() + csvs.REDUCE_WORLD_LEVEL_COOL_TIME
	fmt.Println("操作成功，---当前世界等级：", p.WorldLevelNow, "---真实世界等级", p.WorldLevel)

}

// 设置生日
func (p *ModPlayer) SetBirth(birth int) {
	if p.Birth > 0 {
		fmt.Println("已设置过生日！")
		return
	}
	// birth: 1010->10月10日、520->5月20日
	month := birth / 100
	day := birth % 100
	switch month {
	case 1, 3, 5, 7, 8, 10, 12:
		if day <= 0 || day > 31 {
			fmt.Println(month, "月没有", day, "日！")
			return
		}
	case 4, 6, 9, 11:
		if day <= 0 || day > 30 {
			fmt.Println(month, "月没有", day, "日！")
			return
		}
	case 2:
		if day <= 0 || day > 29 {
			fmt.Println(month, "月没有", day, "日！")
			return
		}
	default:
		fmt.Println("没有", month, "月！")
		return
	}

	p.Birth = birth
	fmt.Println("设置成功，生日为：", month, "月", day, "日")

	if p.IsBirthDay() {
		fmt.Println("生日快乐！")
	} else {
		fmt.Println("期待你生日的到来！")
	}
}

// 今天是否时生日
func (p *ModPlayer) IsBirthDay() bool {
	month := time.Now().Month()
	day := time.Now().Day()
	if int(month) == p.Birth/100 && day == p.Birth%100 {
		return true
	}
	return false
}

// 设置展示名片
func (p *ModPlayer) SetShowCard(showCard []int) {
	// 验展示名片最大长度 9
	if len(showCard) > csvs.SHOW_CARD_SIZE {
		fmt.Println("消息结构错误")
		return
	}

	// 验重
	cardExist := make(map[int]int)
	newList := make([]int, 0)
	for _, cardId := range showCard {
		_, ok := cardExist[cardId]
		if ok {
			continue
		}
		if !p.player.GetModCard().IsHasCard(cardId) {
			continue
		}

		newList = append(newList, cardId)
		cardExist[cardId] = 1
	}

	p.ShowCard = newList
	fmt.Println(p.ShowCard)
}

// 设置展示阵容
func (p *ModPlayer) SetShowTeam(showRole []int) {
	// 验展示阵容最大长度 8
	if len(showRole) > csvs.SHOW_TEAM_SIZE {
		fmt.Println("消息结构错误")
		return
	}

	roleExist := make(map[int]int)
	newList := make([]*ShowRole, 0)
	for _, roleId := range showRole {
		_, ok := roleExist[roleId]
		if ok {
			continue
		}
		if !p.player.GetModRole().IsHasRole(roleId) {
			continue
		}
		showRole := new(ShowRole)
		showRole.RoleId = roleId
		showRole.RoleLevel = p.player.GetModRole().GetRoleLevel(roleId)
		newList = append(newList, showRole)
		roleExist[roleId] = 1
	}

	p.ShowTeam = newList
	for _, role := range p.ShowTeam {
		fmt.Printf("role: %v, roleLevel: %v\n", role.RoleId, role.RoleLevel)
	}
}

// 设置隐藏展示阵容
func (p *ModPlayer) SetHideShowTeam(isHide int) {
	if isHide != csvs.LOGIC_FALSE && isHide != csvs.LOGIC_TRUE {
		return
	}
	p.HideShowTeam = isHide
}

// 内置接口
// 设置封禁
func (p *ModPlayer) SetProhibit(prohibit int) {
	p.Prohibit = prohibit
}

// GM账号
func (p *ModPlayer) SetIsGM(isGM int) {
	p.IsGM = isGM
}

// 玩家是否能登陆
func (p *ModPlayer) IsCanEnter() bool {
	return int64(p.Prohibit) < time.Now().Unix()
}

func (p *ModPlayer) GetWorldLevelNow() int {
	return p.WorldLevelNow
}

func (p *ModPlayer) SaveData() {
	content, err := json.Marshal(p)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(p.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (p *ModPlayer) LoadData(player *Player) {
	p.player = player
	p.path = p.player.localPath + "/player.json"

	configFile, err := ioutil.ReadFile(p.path)
	if err != nil {
		fmt.Println("暂无存档")
		// return
	}
	err = json.Unmarshal(configFile, &p)
	if err != nil {
		p.InitData()
		return
	}
}

func (p *ModPlayer) InitData() {
	p.UserId = p.player.UserId
	p.PlayerLevel = 1
	p.Name = "旅行者"
	p.WorldLevel = 1
	p.WorldLevelNow = 1
}
