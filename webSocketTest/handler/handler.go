package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"webSocketTest/csvs"

	"golang.org/x/net/websocket"
)

type Player struct {
	UserId    int64     `json:"userId"`
	ModPlayer ModPlayer `json:"modPlayer"`
	ModIcon   ModIcon   `json:"modIcon"`
	ModCard   ModCard   `json:"modCard"`
	ModRole   ModRole   `json:"modRole"`
	ModBag    ModBag    `json:"modBag"`
	ModWeapon ModWeapon `json:"modWeapon"`
	ModRelics ModRelics `json:"modRelics"`
	ModCook   ModCook   `json:"modCook"`
	ModHome   ModHome   `json:"modHome"`
}

var player *Player

func Login(ws *websocket.Conn) {
	fmt.Println("请输入userId:")
	var userId int64
	fmt.Scan(&userId)

	msg := new(MsgLogin)
	msg.MsgId = 0
	msg.UserId = userId

	str, errStr := json.Marshal(msg)
	if errStr != nil {
		return
	}
	ws.Write([]byte(str))
}

func RecvPlayerMsg(msg []byte) {
	err := json.Unmarshal(msg, &player)
	if err != nil {
		fmt.Println("解析失败！")
		return
	}
}

func ListenMsg(ws *websocket.Conn) {
	var msg []byte
	var msgHead MsgHead
	for {
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			continue
		}
		err = json.Unmarshal(msg, &msgHead)
		if err != nil {
			fmt.Println("解析失败")
			continue
		}
		switch msgHead.MsgId {
		case 1:
			RecvPlayerMsg(msg)
		case 303:
			RecvPoolResponse(msg)

		}
	}
}

func Run(ws *websocket.Conn) {
	fmt.Println("从0开始写原神服务器------测试工具v1.0")
	fmt.Println("↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓")
	for {
		fmt.Println("欢迎来到提瓦特大陆,请选择功能: 1基础信息 2背包 3模拟抽卡(八重神子UP池) 4地图 5圣遗物 6角色 7武器 8下线")
		var modChoose int
		fmt.Scan(&modChoose)
		switch modChoose {
		case 1:
			HandleBase(ws)
		case 2:
			HandleBag(ws)
		case 3:
			HandlePool(ws)
		// case 4:
		// 	HandleMap(ws)
		case 5:
			HandleRelics(ws)
		case 6:
			HandleRole(ws)
		case 7:
			HandleWeapon(ws)
		case 8:
			//下线
			return
		}
	}
}

func HandleBase(ws *websocket.Conn) {
	for {
		fmt.Println("当前处于基础信息界面,请选择操作: 0返回 1查询信息 2设置名字 3设置签名 4头像 5名片 6设置生日")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			HandleBaseGetInfo()
		case 2:
			HandleBagSetName()
		case 3:
			HandleBagSetSign()
		case 4:
			HandleBagSetIcon()
		case 5:
			HandleBagSetCard()
		case 6:
			HandleBagSetBirth()
		}
	}
}

func HandleBaseGetInfo() {
	fmt.Println("Uid:", player.ModPlayer.UserId)
	fmt.Println("名字:", player.ModPlayer.Name)
	fmt.Println("等级:", player.ModPlayer.PlayerLevel)
	fmt.Println("大世界等级:", player.ModPlayer.WorldLevelNow)
	if player.ModPlayer.Sign == "" {
		fmt.Println("签名:", "未设置")
	} else {
		fmt.Println("签名:", player.ModPlayer.Sign)
	}

	if player.ModPlayer.Icon == 0 {
		fmt.Println("头像:", "未设置")
	} else {
		fmt.Println("头像:", player.ModPlayer.Icon)
	}

	if player.ModPlayer.Card == 0 {
		fmt.Println("名片:", "未设置")
	} else {
		fmt.Println("名片:", player.ModPlayer.Card)
	}

	if player.ModPlayer.Birth == 0 {
		fmt.Println("生日:", "未设置")
	} else {
		fmt.Println("生日:", player.ModPlayer.Birth/100, "月", player.ModPlayer.Birth%100, "日")
	}
}

func HandleBagSetName() {
	fmt.Println("请输入名字:")
	var newName string
	fmt.Scan(&newName)
	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("name", newName)
	resp, err := http.PostForm("http://127.0.0.1:8888/correctname", params)
	if err != nil {
		fmt.Println(err)
	}
	name, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	player.ModPlayer.Name = string(name)
}

func HandleBagSetSign() {
	fmt.Println("请输入签名:")
	var newsign string
	fmt.Scan(&newsign)
	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("sign", newsign)
	resp, err := http.PostForm("http://127.0.0.1:8888/correctsign", params)
	if err != nil {
		fmt.Println(err)
	}
	sign, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	player.ModPlayer.Sign = string(sign)
}

func HandleBagSetIcon() {
	for {
		fmt.Println("当前处于基础信息--头像界面,请选择操作: 0返回 1查询头像背包 2设置头像")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			player.ModIcon.ShowInfo()
		case 2:
			HandleBagSetIconSet()
		}
	}
}

func HandleBagSetIconSet() {
	fmt.Println("请输入头像id:")
	var icon int
	fmt.Scan(&icon)
	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("iconId", strconv.Itoa(icon))
	resp, err := http.PostForm("http://127.0.0.1:8888/icon", params)
	if err != nil {
		fmt.Println(err)
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	iconId, _ := strconv.Atoi(string(msg))
	player.ModPlayer.Icon = iconId
	fmt.Println("设置成功!名片Id:", iconId)
}

func HandleBagSetCard() {
	for {
		fmt.Println("当前处于基础信息--名片界面,请选择操作: 0返回 1查询名片背包 2设置名片")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			player.ModCard.ShowInfo()
		case 2:
			HandleBagSetCardSet()
		}
	}
}

func HandleBagSetCardSet() {
	fmt.Println("请输入名片id:")
	var card int
	fmt.Scan(&card)
	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("cardId", strconv.Itoa(card))
	resp, err := http.PostForm("http://127.0.0.1:8888/card", params)
	if err != nil {
		fmt.Println(err)
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	cardId, _ := strconv.Atoi(string(msg))
	player.ModPlayer.Card = cardId
	fmt.Println("设置成功!名片Id:", cardId)
}

func HandleBagSetBirth() {
	if player.ModPlayer.Birth > 0 {
		fmt.Println("已设置过生日!")
		return
	}
	fmt.Println("生日只能设置一次，请慎重填写,输入月:")
	var month, day int
	fmt.Scan(&month)
	fmt.Println("请输入日:")
	fmt.Scan(&day)
	birth := strconv.Itoa(month*100 + day)

	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("birth", birth)
	resp, err := http.PostForm("http://127.0.0.1:8888/birthday", params)
	if err != nil {
		fmt.Println(err)
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	birthday, _ := strconv.Atoi(string(msg))
	if birthday == 0 {
		fmt.Println("设置失败")
		return
	}
	player.ModPlayer.Birth = birthday
	fmt.Printf("设置成功！生日：%d月%d日\n", player.ModPlayer.Birth/100, player.ModPlayer.Birth%100)
}

func HandleBag(ws *websocket.Conn) {
	for {
		fmt.Println("0返回 1查看背包 2使用物品")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			player.ModBag.ShowInfo()
		case 2:
			HandleBagUseItem()
		}
	}
}

func HandleBagUseItem() {
	itemId := 0
	itemNum := 0
	fmt.Println("物品ID")
	fmt.Scan(&itemId)
	fmt.Println("物品数量")
	fmt.Scan(&itemNum)

	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("itemId", strconv.FormatInt(int64(itemId), 10))
	params.Set("num", strconv.FormatInt(int64(itemNum), 10))
	resp, err := http.PostForm("http://127.0.0.1:8888/useitem", params)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	player.ModBag.UseItem(itemId, int64(itemNum))
}

func HandlePool(ws *websocket.Conn) {
	for {
		fmt.Println("0返回 1十连抽 2单抽")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			SendPoolMsg(ws, 301)
		case 2:
			SendPoolMsg(ws, 302)
		}
	}
}

//发送抽卡请求
func SendPoolMsg(ws *websocket.Conn, msgId int) {
	msg := new(MsgPool)
	msg.MsgId = msgId
	str, errStr := json.Marshal(msg)
	if errStr != nil {
		fmt.Println("errStr:", errStr)
		return
	}
	ws.Write([]byte(str))
}

//抽卡结果
func RecvPoolResponse(msg []byte) {
	var msgPool MsgResponsePool
	msgErr := json.Unmarshal(msg, &msgPool)
	if msgErr != nil {
		fmt.Println("解析错误！")
		return
	}
	dropName := csvs.GetItemConfigName(msgPool.DropId)
	player.ModCard.CheckGetCard(msgPool.DropId, 10)
	player.ModIcon.CheckGetIcon(msgPool.DropId)
	player.ModBag.AddItem(msgPool.Stuff, msgPool.StuffNum)
	player.ModBag.AddItem(msgPool.StuffItem, msgPool.StuffItemNum)
	fmt.Printf("获得:%s\n", dropName)
}

// func HandleMap(ws *websocket.Conn) {
// 	fmt.Println("向着星辰与深渊,欢迎来到冒险家协会！")
// 	for {
// 		fmt.Println("请选择互动地图 0返回 1蒙德 2璃月 1001深入风龙废墟 2001无妄引咎密宫")
// 		var action int
// 		fmt.Scan(&action)
// 		switch action {
// 		case 0:
// 			return
// 		default:
// 			HandleMapIn(action)
// 		}
// 	}
// }

// func  HandleMapIn(mapId int) {
// 	config := csvs.ConfigMapMap[mapId]
// 	if config == nil {
// 		fmt.Println("无法识别的地图")
// 		return
// 	}
// 	GetModMap().RefreshByPlayer(mapId)
// 	for {
// 		GetModMap().GetEventList(config)
// 		fmt.Println("请选择触发事件Id(0返回)")
// 		var action int
// 		fmt.Scan(&action)
// 		switch action {
// 		case 0:
// 			return
// 		default:
// 			eventConfig := csvs.ConfigMapEventMap[action]
// 			if eventConfig == nil {
// 				fmt.Println("无法识别的事件")
// 				break
// 			}
// 			GetModMap().SetEventState(mapId, eventConfig.EventId, csvs.EVENT_STATE_END)
// 		}
// 	}
// }

func HandleRelics(ws *websocket.Conn) {
	for {
		fmt.Println("当前处于圣遗物界面,选择功能 0返回 1查看圣遗物 2强化圣遗物")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			player.ModRelics.ShowInfo()
		case 2:
			HandleRelicsUp()
		default:
			fmt.Println("无法识别在操作")
		}
	}
}

func HandleRelicsUp() {
	for {
		player.ModRelics.ShowInfo()
		var action, exp int
		fmt.Println("请输入圣遗物KeyId:")
		fmt.Scan(&action)
		if action == 0 {
			return
		}
		fmt.Println("获得的经验:")
		fmt.Scan(&exp)
		player.ModRelics.RelicsUp(action, exp)
	}
}

func HandleRole(ws *websocket.Conn) {
	for {
		fmt.Println("当前处于角色界面,选择功能 0返回 1查询 2装备圣遗物 3卸下圣遗物 4装备武器 5卸下武器")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			player.ModRole.ShowInfo()
		// case 2:
		// 	HandleWearRelics()
		// case 3:
		// 	HandleTakeOffRelics()
		// case 4:
		// 	HandleWearWeapon()
		// case 5:
		// 	HandleTakeOffWeapon()
		default:
			fmt.Println("无法识别在操作")
		}
	}
}

// func HandleWearRelics() {
// 	for {
// 		fmt.Println("输入操作的目标英雄Id:,0返回")
// 		var roleId int
// 		fmt.Scan(&roleId)

// 		if roleId == 0 {
// 			return
// 		}

// 		RoleInfo := GetModRole().RoleInfo[roleId]
// 		if RoleInfo == nil {
// 			fmt.Println("角色不存在")
// 			continue
// 		}

// 		RoleInfo.ShowInfo(p)
// 		fmt.Println("输入需要穿戴的圣遗物key:,0返回")
// 		var relicsKey int
// 		fmt.Scan(&relicsKey)
// 		if relicsKey == 0 {
// 			return
// 		}
// 		relics := GetModRelics().RelicsInfo[relicsKey]
// 		if relics == nil {
// 			fmt.Println("圣遗物不存在")
// 			continue
// 		}
// 		GetModRole().WearRelics(RoleInfo, relics)
// 	}
// }

// func HandleTakeOffRelics() {
// 	for {
// 		fmt.Println("输入操作的目标英雄Id:,0返回")
// 		var roleId int
// 		fmt.Scan(&roleId)

// 		if roleId == 0 {
// 			return
// 		}

// 		RoleInfo := GetModRole().RoleInfo[roleId]
// 		if RoleInfo == nil {
// 			fmt.Println("英雄不存在")
// 			continue
// 		}

// 		RoleInfo.ShowInfo(p)
// 		fmt.Println("输入需要卸下的圣遗物key:,0返回")
// 		var relicsKey int
// 		fmt.Scan(&relicsKey)
// 		if relicsKey == 0 {
// 			return
// 		}
// 		relics := GetModRelics().RelicsInfo[relicsKey]
// 		if relics == nil {
// 			fmt.Println("圣遗物不存在")
// 			continue
// 		}
// 		GetModRole().TakeOffRelics(RoleInfo, relics)
// 	}
// }

func HandleWeapon(ws *websocket.Conn) {
	for {
		fmt.Println("当前处于武器界面,选择功能 0返回 1强化 2突破 3精炼 4查询")
		var action int
		fmt.Scan(&action)
		switch action {
		case 0:
			return
		case 1:
			HandleWeaponUp()
		case 2:
			HandleWeaponStarUp()
		case 3:
			HandleWeaponRefineUp()
		case 4:
			player.ModWeapon.ShowInfo()
		default:
			fmt.Println("无法识别在操作")
		}
	}
}

func HandleWeaponUp() {
	for {
		player.ModWeapon.ShowInfo()
		fmt.Println("输入操作的目标武器keyId:,0返回")
		var weaponKeyId, exp int
		fmt.Scan(&weaponKeyId)
		if weaponKeyId == 0 {
			return
		}
		fmt.Println("获得的经验:")
		fmt.Scan(&exp)
		player.ModWeapon.WeaponUp(weaponKeyId, exp)
	}
}

func HandleWeaponStarUp() {
	for {
		player.ModWeapon.ShowInfo()
		fmt.Println("输入操作的目标武器keyId:,0返回")
		var weaponKeyId int
		fmt.Scan(&weaponKeyId)
		if weaponKeyId == 0 {
			return
		}
		player.ModWeapon.WeaponUpStar(weaponKeyId)
	}
}

func HandleWeaponRefineUp() {
	for {
		player.ModWeapon.ShowInfo()
		fmt.Println("输入操作的目标武器keyId:,0返回")
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
			player.ModWeapon.WeaponUpRefine(weaponKeyId, weaponTargetKeyId)
		}
	}
}

// func  HandleWearWeapon() {
// 	for {
// 		fmt.Println("输入操作的目标角色Id:,0返回")
// 		var roleId int
// 		fmt.Scan(&roleId)

// 		if roleId == 0 {
// 			return
// 		}

// 		RoleInfo := GetModRole().RoleInfo[roleId]
// 		if RoleInfo == nil {
// 			fmt.Println("角色不存在")
// 			continue
// 		}

// 		RoleInfo.ShowWeaponInfo(p)
// 		fmt.Println("输入需要装备的武器key:,0返回")
// 		var weaponKey int
// 		fmt.Scan(&weaponKey)
// 		if weaponKey == 0 {
// 			return
// 		}
// 		weaponInfo := GetModWeapon().WeaponInfo[weaponKey]
// 		if weaponInfo == nil {
// 			fmt.Println("武器不存在")
// 			continue
// 		}
// 		GetModRole().WearWeapon(RoleInfo, weaponInfo)
// 		RoleInfo.ShowInfo(p)
// 	}
// }

// func  HandleTakeOffWeapon() {
// 	for {
// 		fmt.Println("输入操作的目标英雄Id:,0返回")
// 		var roleId int
// 		fmt.Scan(&roleId)

// 		if roleId == 0 {
// 			return
// 		}

// 		RoleInfo := GetModRole().RoleInfo[roleId]
// 		if RoleInfo == nil {
// 			fmt.Println("英雄不存在")
// 			continue
// 		}

// 		RoleInfo.ShowWeaponInfo(p)
// 		fmt.Println("输入需要卸下的武器key:,0返回")
// 		var weaponKey int
// 		fmt.Scan(&weaponKey)
// 		if weaponKey == 0 {
// 			return
// 		}
// 		weapon := GetModWeapon().WeaponInfo[weaponKey]
// 		if weapon == nil {
// 			fmt.Println("武器不存在")
// 			continue
// 		}
// 		GetModRole().TakeOffWeapon(RoleInfo, weapon)
// 		RoleInfo.ShowInfo(p)
// 	}
// }
