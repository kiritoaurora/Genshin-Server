package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"server_logic/src/csvs"
	"time"
)

/*
	地图模块
*/

type Event struct {
	EventId       int
	Stata         int
	NextResetTime int64
}

type Map struct {
	MapId     int
	EventInfo map[int]*Event
}

type StatueInfo struct {
	StatusId int
	Level    int
	ItemInfo map[int]*ItemInfo
}

type ModMap struct {
	MapInfo map[int]*Map
	Statue  map[int]*StatueInfo

	player *Player
	path   string
}

//地图初始化
func (m *ModMap) InitData() {
	m.MapInfo = make(map[int]*Map)
	m.Statue = make(map[int]*StatueInfo)
	//初始化区域地图
	for _, v := range csvs.ConfigMapMap {
		_, ok := m.MapInfo[v.MapId]
		if !ok {
			m.MapInfo[v.MapId] = m.NewMapInfo(v.MapId)
		}
	}
	//加载事件
	for _, v := range csvs.ConfigMapEventMap {
		_, ok := m.MapInfo[v.MapId]
		if !ok {
			continue
		}
		_, ok = m.MapInfo[v.MapId].EventInfo[v.EventId]
		if !ok {
			m.MapInfo[v.MapId].EventInfo[v.EventId] = new(Event)
			m.MapInfo[v.MapId].EventInfo[v.EventId].EventId = v.EventId
			m.MapInfo[v.MapId].EventInfo[v.EventId].Stata = csvs.EVENT_STATE_START
		}
	}
}

func (m *ModMap) NewMapInfo(mapId int) *Map {
	mapInfo := new(Map)
	mapInfo.MapId = mapId
	mapInfo.EventInfo = make(map[int]*Event)
	return mapInfo
}

// 获取地图事件列表
func (m *ModMap) GetEventList(config *csvs.ConfigMap) {
	_, ok := m.MapInfo[config.MapId]
	if !ok {
		return
	}
	for _, v := range m.MapInfo[config.MapId].EventInfo {
		m.CheckEventRefresh(v)                          //检查事件刷新情况
		lastTime := v.NextResetTime - time.Now().Unix() //事件刷新时间 = 刷新时间戳 - 当前时间戳
		noticeTime := ""
		if lastTime <= 0 {
			noticeTime = "已刷新"
		} else {
			noticeTime = fmt.Sprintf("%d秒后刷新", lastTime)
		}
		fmt.Printf("事件Id:%d,名字:%s,状态:%d, %s\n", v.EventId, csvs.GetEventName(v.EventId), v.Stata, noticeTime)

	}
}

func (m *ModMap) SetEventState(mapId int, eventId int, state int) {
	_, ok := m.MapInfo[mapId]
	if !ok {
		fmt.Println("地图不存在！")
	}
	_, ok = m.MapInfo[mapId].EventInfo[eventId]
	if !ok {
		fmt.Println("事件不存在！")
	}
	if m.MapInfo[mapId].EventInfo[eventId].Stata == state {
		fmt.Println("状态异常！")
		return
	}
	eventConfig := csvs.GetEventConfig(m.MapInfo[mapId].EventInfo[eventId].EventId)
	if eventConfig == nil {
		return
	}

	/*
	 验秘境副本是否完成所有事件
	 若未完成事件直接领取奖励，则直接返回，不予以奖励法方
	*/
	configMap := csvs.ConfigMapMap[mapId] //获取地图配置
	if configMap == nil {
		return
	}
	// 验证树脂数量
	if !m.player.GetModBag().HasEnoughItem(eventConfig.CostItem, eventConfig.CostNum) {
		fmt.Printf("%s不足!\n", csvs.GetItemConfig(eventConfig.CostItem).ItemName)
		return
	}
	//地图类型为玩家刷新且事件类型为奖励
	if configMap.MapType == csvs.REFRESH_PLAYER && eventConfig.EventType == csvs.EVENT_TYPE_REWARD {
		//遍历该地图所有事件，判定是否有未完成非奖励的事件
		for _, v := range m.MapInfo[mapId].EventInfo {
			eventConfigNow := csvs.GetEventConfig(v.EventId)
			if eventConfigNow == nil {
				continue
			}
			if eventConfigNow.EventType == csvs.EVENT_TYPE_REWARD {
				continue
			}
			if v.EventId == eventId {
				continue
			}
			if v.Stata != csvs.EVENT_STATE_END {
				fmt.Println("有事件尚未完成:", v.EventId)
				return
			}
		}
	}

	// 设置事件状态
	m.MapInfo[mapId].EventInfo[eventId].Stata = state
	if state == csvs.EVENT_STATE_FINISH {
		fmt.Println("事件完成！")
	}
	if state == csvs.EVENT_STATE_END {
		// 多倍掉落，例：浓缩树脂
		for i := 0; i < eventConfig.EventDropTimes; i++ {
			//获取掉落组，随机掉落物品和掉落数量
			config := csvs.GetRandDropItemNew(eventConfig.EventDrop)
			for _, v := range config {
				randNum := rand.Intn(csvs.PERCENT_ALL)
				if randNum < v.Weight {
					randAll := v.ItemNumMax - v.ItemNumMin + 1
					itemNum := rand.Intn(randAll) + v.ItemNumMin
					//世界等级加成
					worldLevel := m.player.GetModPlayer().GetWorldLevelNow()
					if worldLevel > 0 {
						itemNum = itemNum * (csvs.PERCENT_ALL + worldLevel*v.WorldAdd) / csvs.PERCENT_ALL
					}
					//物品入背包
					m.player.GetModBag().AddItem(v.ItemId, int64(itemNum))
				}
			}
		}
		fmt.Println("事件领取！")
	}
	// 事件状态不再是START
	if state > 0 {
		// 判定事件的刷新类型是否时自刷新，若是，则更新其刷新时间戳
		switch eventConfig.RefreshType {
		case csvs.MAP_REFRESH_SELF:
			m.MapInfo[mapId].EventInfo[eventId].NextResetTime = time.Now().Unix() + csvs.MAP_REFRESH_SELF_TIME
		}
	}
}

//按日刷新
func (m *ModMap) RefreshDay() {
	for _, v := range m.MapInfo {
		for _, v := range m.MapInfo[v.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[v.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_DAY {
				continue
			}
			v.Stata = csvs.EVENT_STATE_START
		}
	}
}

//按周刷新
func (m *ModMap) RefreshWeek() {
	for _, v := range m.MapInfo {
		for _, v := range m.MapInfo[v.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[v.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_WEEK {
				continue
			}
			v.Stata = csvs.EVENT_STATE_START
		}
	}
}

//自刷新
func (m *ModMap) RefreshSelf() {
	for _, v := range m.MapInfo {
		for _, v := range m.MapInfo[v.MapId].EventInfo {
			config := csvs.ConfigMapEventMap[v.EventId]
			if config == nil {
				continue
			}
			if config.RefreshType != csvs.MAP_REFRESH_SELF {
				continue
			}
			if time.Now().Unix() <= v.NextResetTime {
				v.Stata = csvs.EVENT_STATE_START
			}
		}
	}
}

// 刷新检查，若事件下一次刷新时间戳小于当前时间戳，即：需要刷新，则刷新
// 日刷新和周刷新到时间自动刷新，与其状态无关，例：地图小怪，周本
// 自刷新需要状态不为0时才刷新，在这里不做操作，例：地图BOSS
// 一次性事件不刷新，例：神瞳
func (m *ModMap) CheckEventRefresh(event *Event) {
	if event.NextResetTime > time.Now().Unix() {
		return
	}
	eventConfig := csvs.GetEventConfig(event.EventId)
	if eventConfig == nil {
		return
	}
	switch eventConfig.RefreshType {
	case csvs.MAP_REFRESH_DAY:
		count := time.Now().Unix() / csvs.MAP_REFRESH_DAY_TIME
		count++
		event.NextResetTime = count * csvs.MAP_REFRESH_DAY_TIME
	case csvs.MAP_REFRESH_WEEK:
		count := time.Now().Unix() / csvs.MAP_REFRESH_WEEK_TIME
		count++
		event.NextResetTime = count * csvs.MAP_REFRESH_WEEK_TIME
	case csvs.MAP_REFRESH_SELF:
	case csvs.MAP_REFRESH_CANT:
		return
	}
	event.Stata = csvs.EVENT_STATE_START
}

// 玩家触发副本刷新
func (m *ModMap) RefreshByPlayer(mapId int) {
	config := csvs.ConfigMapMap[mapId] //获取地图配置
	if config == nil {
		return
	}
	//判定地图刷新是否为玩家触发类型
	if config.MapType != csvs.REFRESH_PLAYER {
		return
	}
	_, ok := m.MapInfo[mapId]
	if !ok {
		return
	}
	// 设置地图中的所有事件状态为0
	for _, v := range m.MapInfo[config.MapId].EventInfo {
		v.Stata = csvs.EVENT_STATE_START
	}
}

// 初始化神像
func (m *ModMap) NewStatue(statusId int) *StatueInfo {
	data := new(StatueInfo)
	data.StatusId = statusId
	data.Level = 0
	data.ItemInfo = make(map[int]*ItemInfo)
	return data
}

//七天神像升级
func (m *ModMap) UpStatus(statueId int) {
	_, ok := m.Statue[statueId]
	if !ok {
		m.Statue[statueId] = m.NewStatue(statueId)
	}
	info, ok := m.Statue[statueId] //获取神像信息
	if !ok {
		return
	}

	nextLevel := info.Level + 1                                   //神像的下一个等级
	nextStatueConfig := csvs.GetStatueConfig(statueId, nextLevel) //取神像下个等级的配置
	if nextStatueConfig == nil {
		return
	}
	//获取已交数量
	_, ok = info.ItemInfo[nextStatueConfig.CostItem]
	nowNum := int64(0)
	if ok {
		nowNum = info.ItemInfo[nextStatueConfig.CostItem].ItemNum
	}
	//升到下一级所需要的数量
	needNum := nextStatueConfig.CostNum - nowNum

	//神瞳不充足
	if !m.player.GetModBag().HasEnoughItem(nextStatueConfig.CostItem, needNum) {
		num := m.player.GetModBag().GetItemNum(nextStatueConfig.CostItem) //获取背包神瞳数量
		if num <= 0 {
			fmt.Printf("七天神像升级所需物品不足\n")
			return
		}
		_, ok = info.ItemInfo[nextStatueConfig.CostItem]
		if !ok {
			info.ItemInfo[nextStatueConfig.CostItem] = new(ItemInfo)
			info.ItemInfo[nextStatueConfig.CostItem].ItemId = nextStatueConfig.CostItem
			info.ItemInfo[nextStatueConfig.CostItem].ItemNum = 0
		}
		_, ok = info.ItemInfo[nextStatueConfig.CostItem]
		if !ok {
			return
		}
		info.ItemInfo[nextStatueConfig.CostItem].ItemNum = num        //神像持有神瞳数
		m.player.GetModBag().RemoveItemToBag(nextStatueConfig.CostItem, num) //扣除背包中的神瞳
		fmt.Printf("七天神像升级，提交物品:%d,提交数量:%d,当前数量:%d\n", nextStatueConfig.CostItem, num, info.ItemInfo[nextStatueConfig.CostItem].ItemNum)

	} else {
		//神瞳充足，扣除神瞳，等级++，神像持有神瞳信息重置
		m.player.GetModBag().RemoveItemToBag(nextStatueConfig.CostItem, needNum)
		info.Level++
		info.ItemInfo = make(map[int]*ItemInfo)
		fmt.Printf("七天神像升级成功，神像:%d,当前等级:%d\n", info.StatusId, info.Level)
	}
}

func (m *ModMap) SaveData() {
	content, err := json.Marshal(m)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(m.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (m *ModMap) LoadData(player *Player) {
	m.player = player
	m.path = m.player.localPath + "/map.json"
	
	configFile, err := ioutil.ReadFile(m.path)
	if err != nil {
		fmt.Println("暂无存档")
		// return
	}
	err = json.Unmarshal(configFile, &m)
	if err != nil {
		m.InitData()
		return
	}
	// m.InitData()
}
