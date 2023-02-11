package game

import (
	"encoding/json"
	"fmt"
	"os"
	"server_logic/src/csvs"

	"golang.org/x/net/websocket"
)

const (
	TASK_STATE_INIT   = 0
	TASK_STATE_DOING  = 1
	TASK_STATE_FINISH = 2
)

const (
	MOD_PLAYER     = "player"
	MOD_ICON       = "icon"
	MOD_CARD       = "card"
	MOD_UNIQUETASK = "uniquetask"
	MOD_ROLE       = "role"
	MOD_BAG        = "bag"
	MOD_WEAPON     = "weapon"
	MOD_RELICS     = "relics"
	MOD_COOK       = "cook"
	MOD_HOME       = "home"
	MOD_POOL       = "pool"
	MOD_MAP        = "map"
)

type ModBase interface {
	LoadData(player *Player)
	SaveData()
	InitData()
}

type Player struct {
	// ModPlayer     *ModPlayer
	// ModIcon       *ModIcon
	// ModCard       *ModCard
	// ModUniqueTask *ModUniqueTask
	// ModRole       *ModRole
	// ModBag        *ModBag
	// ModWeapon     *ModWeapon
	// ModRelics     *ModRelics
	// ModCook       *ModCook
	// ModHome       *ModHome
	// ModPool       *ModPool
	// ModMap        *ModMap

	UserId    int64
	modManage map[string]ModBase
	localPath string
	ws        *websocket.Conn
	exitTime  int64
	chanLogic chan []byte
}

var player *Player

func NewTestPlayer(ws *websocket.Conn, userId int64) *Player {
	player = new(Player)
	// player.ModPlayer = new(ModPlayer)
	// player.ModIcon = new(ModIcon)
	// player.ModIcon.IconInfo = make(map[int]*Icon)
	// player.ModCard = new(ModCard)
	// player.ModCard.CardInfo = make(map[int]*Card)
	// player.ModUniqueTask = new(ModUniqueTask)
	// player.ModUniqueTask.MyTaskInfo = make(map[int]*TaskInfo)
	// player.ModUniqueTask.Locker = new(sync.RWMutex)
	// player.ModRole = new(ModRole)
	// player.ModRole.RoleInfo = make(map[int]*RoleInfo)
	// player.ModBag = new(ModBag)
	// player.ModBag.BagInfo = make(map[int]*ItemInfo)
	// player.ModWeapon = new(ModWeapon)
	// player.ModWeapon.WeaponInfo = make(map[int]*Weapon)
	// player.ModRelics = new(ModRelics)
	// player.ModRelics.RelicsInfo = make(map[int]*Relics)
	// player.ModCook = new(ModCook)
	// player.ModCook.CookInfo = make(map[int]*Cook)
	// player.ModHome = new(ModHome)
	// player.ModHome.HomeItemInfo = make(map[int]*Home)
	// player.ModPool = new(ModPool)
	// player.ModPool.UpPoolInfo = new(Pool)
	// player.ModMap = new(ModMap)
	// player.ModMap.InitData()

	/*************************泛型架构***********************/
	player.UserId = userId
	player.chanLogic = make(chan []byte)
	player.modManage = map[string]ModBase{
		MOD_PLAYER:     new(ModPlayer),
		MOD_ICON:       new(ModIcon),
		MOD_CARD:       new(ModCard),
		MOD_UNIQUETASK: new(ModUniqueTask),
		MOD_ROLE:       new(ModRole),
		MOD_BAG:        new(ModBag),
		MOD_WEAPON:     new(ModWeapon),
		MOD_RELICS:     new(ModRelics),
		MOD_COOK:       new(ModCook),
		MOD_HOME:       new(ModHome),
		MOD_POOL:       new(ModPool),
		MOD_MAP:        new(ModMap),
	}

	player.InitData()
	player.InitMod()
	player.ws = ws
	return player
}

func (p *Player) InitData() {
	path := GetServer().Config.LocalSavePath
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			return
		}
	}
	p.localPath = path + fmt.Sprintf("/%d/", p.UserId)
	_, err = os.Stat(p.localPath)
	if err != nil {
		err = os.Mkdir(p.localPath, os.ModePerm)
		if err != nil {
			return
		}
	}
}

func (p *Player) InitMod() {
	for _, v := range p.modManage {
		v.LoadData(p)
	}
}

func (p *Player) GetMod(modName string) ModBase {
	return p.modManage[modName]
}

func (p *Player) GetModPlayer() *ModPlayer {
	return p.modManage[MOD_PLAYER].(*ModPlayer)
}

func (p *Player) GetModIcon() *ModIcon {
	return p.modManage[MOD_ICON].(*ModIcon)
}

func (p *Player) GetModCard() *ModCard {
	return p.modManage[MOD_CARD].(*ModCard)
}

func (p *Player) GetModUniqueTask() *ModUniqueTask {
	return p.modManage[MOD_UNIQUETASK].(*ModUniqueTask)
}

func (p *Player) GetModRole() *ModRole {
	return p.modManage[MOD_ROLE].(*ModRole)
}

func (p *Player) GetModBag() *ModBag {
	return p.modManage[MOD_BAG].(*ModBag)
}

func (p *Player) GetModWeapon() *ModWeapon {
	return p.modManage[MOD_WEAPON].(*ModWeapon)
}

func (p *Player) GetModRelics() *ModRelics {
	return p.modManage[MOD_RELICS].(*ModRelics)
}

func (p *Player) GetModCook() *ModCook {
	return p.modManage[MOD_COOK].(*ModCook)
}

func (p *Player) GetModHome() *ModHome {
	return p.modManage[MOD_HOME].(*ModHome)
}

func (p *Player) GetModPool() *ModPool {
	return p.modManage[MOD_POOL].(*ModPool)
}

func (p *Player) GetModMap() *ModMap {
	return p.modManage[MOD_MAP].(*ModMap)
}

func (p *Player) SendLogic(msg []byte) {
	p.chanLogic <- msg
}

func (p *Player) LogicRun() {
	for {
		// select {
		// case msg := <- p.chanLogic:
		// 	p.ProcessMsg(msg)
		// }
		msg := <-p.chanLogic
		p.ProcessMsg(msg)
	}
}

func (p *Player) ProcessMsg(msg []byte) {
	var msgHead MsgHead
	msgErr := json.Unmarshal(msg, &msgHead)
	if msgErr != nil {
		fmt.Println("消息解析失败！")
		return
	}
	switch msgHead.MsgId {
	case 301:		//十连抽
		p.GetModPool().HandleUpPoolTenByMsg(msg)
	case 302:		//单抽
		p.GetModPool().HandleUpPoolSingleByMsg(msg)
	default:
		fmt.Println("无法识别的消息！")
	}
}

//对外接口
func (p *Player) RecvSetIcon(iconId int) {
	p.GetMod(MOD_PLAYER).(*ModPlayer).SetIcon(iconId)
}

func (p *Player) RecvSetCard(cardId int) {
	p.GetMod(MOD_PLAYER).(*ModPlayer).SetCard(cardId)
}

func (p *Player) RecvSetName(name string) {
	p.GetMod(MOD_PLAYER).(*ModPlayer).SetName(name)
}

func (p *Player) RecvSetSign(sign string) {
	p.GetMod(MOD_PLAYER).(*ModPlayer).SetSign(sign)
}

func (p *Player) ReduceWorldLevel() {
	p.GetMod(MOD_PLAYER).(*ModPlayer).ReduceWorldLevel()
}

func (p *Player) ReturnWorldLevel() {
	p.GetMod(MOD_PLAYER).(*ModPlayer).ReturnWorldLevel()
}

func (p *Player) SetBirth(birth int) {
	p.GetMod(MOD_PLAYER).(*ModPlayer).SetBirth(birth)
}

func (p *Player) SetShowCard(showCard []int) {
	p.GetMod(MOD_PLAYER).(*ModPlayer).SetShowCard(showCard)
}

func (p *Player) SetShowTeam(showRole []int) {
	p.GetMod(MOD_PLAYER).(*ModPlayer).SetShowTeam(showRole)
}

func (p *Player) SetHideShowTeam(isHide int) {
	p.GetMod(MOD_PLAYER).(*ModPlayer).SetHideShowTeam(isHide)
}

func (p *Player) Run() {
	fmt.Println("从0开始写原神服务器------测试工具v0.9")
	fmt.Println("作者:B站------golang大海葵")
	fmt.Println("模拟用户创建成功OK------开始测试")
	fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	for {
		fmt.Println("欢迎来到提瓦特大陆,请选择功能: 1基础信息 2背包 3模拟抽卡(神里绫华UP池) 4地图 5圣遗物 6角色 7武器 8存档")
		var modChoose int
		fmt.Scan(&modChoose)
		switch modChoose {
		case 1:
			p.HandleBase()
		case 2:
			p.HandleBag()
		case 3:
			p.HandlePool()
		case 4:
			p.HandleMap()
		case 5:
			p.HandleRelics()
		case 6:
			p.HandleRole()
		case 7:
			p.HandleWeapon()
		case 8:
			for _, v := range p.modManage {
				v.SaveData()
			}
		}
	}
}

//基础信息
func (p *Player) HandleBase() {
	for {
		fmt.Println("当前处于基础信息界面,请选择操作: 0返回 1查询信息 2设置名字 3设置签名 4头像 5名片 6设置生日")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			p.HandleBaseGetInfo()
		case 2:
			p.HandleBagSetName()
		case 3:
			p.HandleBagSetSign()
		case 4:
			p.HandleBagSetIcon()
		case 5:
			p.HandleBagSetCard()
		case 6:
			p.HandleBagSetBirth()
		}
	}
}

func (p *Player) HandleBaseGetInfo() {
	fmt.Println("名字:", p.GetMod(MOD_PLAYER).(*ModPlayer).Name)
	fmt.Println("等级:", p.GetMod(MOD_PLAYER).(*ModPlayer).PlayerLevel)
	fmt.Println("大世界等级:", p.GetMod(MOD_PLAYER).(*ModPlayer).WorldLevelNow)
	if p.GetMod(MOD_PLAYER).(*ModPlayer).Sign == "" {
		fmt.Println("签名:", "未设置")
	} else {
		fmt.Println("签名:", p.GetMod(MOD_PLAYER).(*ModPlayer).Sign)
	}

	if p.GetMod(MOD_PLAYER).(*ModPlayer).Icon == 0 {
		fmt.Println("头像:", "未设置")
	} else {
		fmt.Println("头像:", csvs.GetItemConfig(p.GetMod(MOD_PLAYER).(*ModPlayer).Icon), p.GetMod(MOD_PLAYER).(*ModPlayer).Icon)
	}

	if p.GetMod(MOD_PLAYER).(*ModPlayer).Card == 0 {
		fmt.Println("名片:", "未设置")
	} else {
		fmt.Println("名片:", csvs.GetItemConfig(p.GetMod(MOD_PLAYER).(*ModPlayer).Card), p.GetMod(MOD_PLAYER).(*ModPlayer).Card)
	}

	if p.GetMod(MOD_PLAYER).(*ModPlayer).Birth == 0 {
		fmt.Println("生日:", "未设置")
	} else {
		fmt.Println("生日:", p.GetMod(MOD_PLAYER).(*ModPlayer).Birth/100, "月", p.GetMod(MOD_PLAYER).(*ModPlayer).Birth%100, "日")
	}
}

func (p *Player) HandleBagSetName() {
	fmt.Println("请输入名字:")
	var name string
	fmt.Scan(&name)
	p.RecvSetName(name)
}

func (p *Player) HandleBagSetSign() {
	fmt.Println("请输入签名:")
	var sign string
	fmt.Scan(&sign)
	p.RecvSetSign(sign)
}

func (p *Player) HandleBagSetIcon() {
	for {
		fmt.Println("当前处于基础信息--头像界面,请选择操作: 0返回 1查询头像背包 2设置头像")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			p.HandleBagSetIconGetInfo()
		case 2:
			p.HandleBagSetIconSet()
		}
	}
}

func (p *Player) HandleBagSetIconGetInfo() {
	fmt.Println("当前拥有头像如下:")
	for _, v := range p.GetModIcon().IconInfo {
		config := csvs.GetItemConfig(v.IconId)
		if config != nil {
			fmt.Println(config.ItemName, ":", config.ItemId)
		}
	}
}

func (p *Player) HandleBagSetIconSet() {
	fmt.Println("请输入头像id:")
	var icon int
	fmt.Scan(&icon)
	p.RecvSetIcon(icon)
}

func (p *Player) HandleBagSetCard() {
	for {
		fmt.Println("当前处于基础信息--名片界面,请选择操作: 0返回 1查询名片背包 2设置名片")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			p.HandleBagSetCardGetInfo()
		case 2:
			p.HandleBagSetCardSet()
		}
	}
}

func (p *Player) HandleBagSetCardGetInfo() {
	fmt.Println("当前拥有名片如下:")
	for _, v := range p.GetModCard().CardInfo {
		config := csvs.GetItemConfig(v.CardId)
		if config != nil {
			fmt.Println(config.ItemName, ":", config.ItemId)
		}
	}
}

func (p *Player) HandleBagSetCardSet() {
	fmt.Println("请输入名片id:")
	var card int
	fmt.Scan(&card)
	p.RecvSetCard(card)
}

func (p *Player) HandleBagSetBirth() {
	if p.GetMod(MOD_PLAYER).(*ModPlayer).Birth > 0 {
		fmt.Println("已设置过生日!")
		return
	}
	fmt.Println("生日只能设置一次，请慎重填写,输入月:")
	var month, day int
	fmt.Scan(&month)
	fmt.Println("请输入日:")
	fmt.Scan(&day)
	p.GetMod(MOD_PLAYER).(*ModPlayer).SetBirth(month*100 + day)
}

//背包
func (p *Player) HandleBag() {
	for {
		fmt.Println("当前处于基础信息界面,请选择操作: 0返回 1增加物品 2扣除物品 3使用物品 4升级七天神像(风)")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			p.HandleBagAddItem()
		case 2:
			p.HandleBagRemoveItem()
		case 3:
			p.HandleBagUseItem()
		case 4:
			p.HandleBagWindStatue()
		}
	}
}

//抽卡
func (p *Player) HandlePool() {
	for {
		fmt.Println("当前处于模拟抽卡界面,请选择操作: 0返回 1角色信息 2十连抽 3单抽(可选次数) 4抽卡1000万次 5单抽(带仓检)")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			p.GetModRole().HandleSendRoleInfo()
		case 2:
			p.GetModPool().HandleUpPoolTen()
		case 3:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			p.GetModPool().HandleUpPoolSingle(times)
		case 4:
			p.GetModPool().DoUpPool()
		case 5:
			fmt.Println("请输入抽卡次数,最大值1亿(最大耗时约30秒):")
			var times int
			fmt.Scan(&times)
			p.GetModPool().HandleUpPoolSingleCheck1(times)
		}
	}
}

func (p *Player) HandleBagAddItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	p.GetModBag().AddItem(itemId, int64(itemNum))
}

func (p *Player) HandleBagRemoveItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	p.GetModBag().RemoveItemToBag(itemId, int64(itemNum))
}

func (p *Player) HandleBagUseItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)
	p.GetModBag().UseItem(itemId, int64(itemNum))
}

func (p *Player) HandleBagWindStatue() {
	fmt.Println("开始升级七天神像")
	p.GetModMap().UpStatus(1)
	p.GetModRole().CalHpPool()
}

//地图
func (p *Player) HandleMap() {
	fmt.Println("向着星辰与深渊,欢迎来到冒险家协会！")
	for {
		fmt.Println("请选择互动地图 0返回 1蒙德 2璃月 1001深入风龙废墟 2001无妄引咎密宫")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		default:
			p.HandleMapIn(action)
		}
	}
}

func (p *Player) HandleMapIn(mapId int) {

	config := csvs.ConfigMapMap[mapId]
	if config == nil {
		fmt.Println("无法识别的地图")
		return
	}
	p.GetModMap().RefreshByPlayer(mapId)
	for {
		p.GetModMap().GetEventList(config)
		fmt.Println("请选择触发事件Id(0返回)")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		default:
			eventConfig := csvs.ConfigMapEventMap[action]
			if eventConfig == nil {
				fmt.Println("无法识别的事件")
				break
			}
			p.GetModMap().SetEventState(mapId, eventConfig.EventId, csvs.EVENT_STATE_END)
		}
	}
}

// 圣遗物
func (p *Player) HandleRelics() {
	for {
		fmt.Println("当前处于圣遗物界面,选择功能 0返回 1强化圣遗物 2满级圣遗物 3极品头测试")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			p.HandleRelicsUp()
		case 2:
			p.GetModRelics().RelicsTop()
		case 3:
			p.GetModRelics().RelicsTestBest()
		default:
			fmt.Println("无法识别在操作")
		}
	}
}

func (p *Player) HandleRelicsUp() {
	for {
		for _, v := range p.GetModRelics().RelicsInfo {
			fmt.Printf("圣遗物KeyId:%d, 圣遗物等级:%d\n\n", v.KeyId, v.Level)
		}
		var action,  exp int
		fmt.Println("请输入圣遗物KeyId:")
		fmt.Scan(&action)
		if action == 0 {
			return
		}
		fmt.Println("获得的经验:")
		fmt.Scan(&exp)
		p.GetModRelics().RelicsUp(action, exp)
	}
}

// 角色
func (p *Player) HandleRole() {
	for {
		fmt.Println("当前处于角色界面,选择功能 0返回 1查询 2装备圣遗物 3卸下圣遗物 4装备武器 5卸下武器")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			p.GetModRole().HandleSendRoleInfo()
		case 2:
			p.HandleWearRelics()
		case 3:
			p.HandleTakeOffRelics()
		case 4:
			p.HandleWearWeapon()
		case 5:
			p.HandleTakeOffWeapon()
		default:
			fmt.Println("无法识别在操作")
		}
	}
}

func (p *Player) HandleWearRelics() {
	for {
		fmt.Println("输入操作的目标英雄Id:,0返回")
		var roleId int
		fmt.Scan(&roleId)

		if roleId == 0 {
			return
		}

		RoleInfo := p.GetModRole().RoleInfo[roleId]
		if RoleInfo == nil {
			fmt.Println("角色不存在")
			continue
		}

		RoleInfo.ShowInfo(p)
		fmt.Println("输入需要穿戴的圣遗物key:,0返回")
		var relicsKey int
		fmt.Scan(&relicsKey)
		if relicsKey == 0 {
			return
		}
		relics := p.GetModRelics().RelicsInfo[relicsKey]
		if relics == nil {
			fmt.Println("圣遗物不存在")
			continue
		}
		p.GetModRole().WearRelics(RoleInfo, relics)
	}
}

func (p *Player) HandleTakeOffRelics() {
	for {
		fmt.Println("输入操作的目标英雄Id:,0返回")
		var roleId int
		fmt.Scan(&roleId)

		if roleId == 0 {
			return
		}

		RoleInfo := p.GetModRole().RoleInfo[roleId]
		if RoleInfo == nil {
			fmt.Println("英雄不存在")
			continue
		}

		RoleInfo.ShowInfo(p)
		fmt.Println("输入需要卸下的圣遗物key:,0返回")
		var relicsKey int
		fmt.Scan(&relicsKey)
		if relicsKey == 0 {
			return
		}
		relics := p.GetModRelics().RelicsInfo[relicsKey]
		if relics == nil {
			fmt.Println("圣遗物不存在")
			continue
		}
		p.GetModRole().TakeOffRelics(RoleInfo, relics)
	}
}

func (p *Player) HandleWeapon() {
	for {
		fmt.Println("当前处于武器界面,选择功能 0返回 1强化测试 2突破测试 3精炼测试")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			p.HandleWeaponUp()
		case 2:
			p.HandleWeaponStarUp()
		case 3:
			p.HandleWeaponRefineUp()
		default:
			fmt.Println("无法识别在操作")
		}
	}
}

func (p *Player) HandleWeaponUp() {
	for {
		fmt.Println("输入操作的目标武器keyId:,0返回")
		for _, v := range p.GetModWeapon().WeaponInfo {
			fmt.Printf("武器keyId:%d,等级:%d,突破等级:%d,精炼:%d\n\n", v.KeyId, v.Level, v.StarLevel, v.RefineLevel)
		}
		var weaponKeyId, exp int
		fmt.Scan(&weaponKeyId)
		if weaponKeyId == 0 {
			return
		}
		fmt.Println("获得的经验:")
		fmt.Scan(&exp)
		p.GetModWeapon().WeaponUp(weaponKeyId, exp)
	}
}

func (p *Player) HandleWeaponStarUp() {
	for {
		fmt.Println("输入操作的目标武器keyId:,0返回")
		for _, v := range p.GetModWeapon().WeaponInfo {
			fmt.Printf("武器keyId:%d,等级:%d,突破等级:%d,精炼:%d\n\n", v.KeyId, v.Level, v.StarLevel, v.RefineLevel)
		}
		var weaponKeyId int
		fmt.Scan(&weaponKeyId)
		if weaponKeyId == 0 {
			return
		}
		p.GetModWeapon().WeaponUpStar(weaponKeyId)
	}
}

func (p *Player) HandleWeaponRefineUp() {
	for {
		fmt.Println("输入操作的目标武器keyId:,0返回")
		for _, v := range p.GetModWeapon().WeaponInfo {
			fmt.Printf("武器keyId:%d,等级:%d,突破等级:%d,精炼:%d\n\n", v.KeyId, v.Level, v.StarLevel, v.RefineLevel)
		}
		var weaponKeyId int
		fmt.Scan(&weaponKeyId)
		if weaponKeyId == 0 {
			return
		}
		for {
			fmt.Println("输入作为材料的武器keyId:,0返回")
			var weaponTargetKeyId int
			fmt.Scan(&weaponTargetKeyId)
			if weaponTargetKeyId == 0 {
				return
			}
			p.GetModWeapon().WeaponUpRefine(weaponKeyId, weaponTargetKeyId)
		}
	}
}

func (p *Player) HandleWearWeapon() {
	for {
		fmt.Println("输入操作的目标角色Id:,0返回")
		var roleId int
		fmt.Scan(&roleId)

		if roleId == 0 {
			return
		}

		RoleInfo := p.GetModRole().RoleInfo[roleId]
		if RoleInfo == nil {
			fmt.Println("角色不存在")
			continue
		}

		RoleInfo.ShowWeaponInfo(p)
		fmt.Println("输入需要装备的武器key:,0返回")
		var weaponKey int
		fmt.Scan(&weaponKey)
		if weaponKey == 0 {
			return
		}
		weaponInfo := p.GetModWeapon().WeaponInfo[weaponKey]
		if weaponInfo == nil {
			fmt.Println("武器不存在")
			continue
		}
		p.GetModRole().WearWeapon(RoleInfo, weaponInfo)
		RoleInfo.ShowInfo(p)
	}
}

func (p *Player) HandleTakeOffWeapon() {
	for {
		fmt.Println("输入操作的目标英雄Id:,0返回")
		var roleId int
		fmt.Scan(&roleId)

		if roleId == 0 {
			return
		}

		RoleInfo := p.GetModRole().RoleInfo[roleId]
		if RoleInfo == nil {
			fmt.Println("英雄不存在")
			continue
		}

		RoleInfo.ShowWeaponInfo(p)
		fmt.Println("输入需要卸下的武器key:,0返回")
		var weaponKey int
		fmt.Scan(&weaponKey)
		if weaponKey == 0 {
			return
		}
		weapon := p.GetModWeapon().WeaponInfo[weaponKey]
		if weapon == nil {
			fmt.Println("武器不存在")
			continue
		}
		p.GetModRole().TakeOffWeapon(RoleInfo, weapon)
		RoleInfo.ShowInfo(p)
	}
}
