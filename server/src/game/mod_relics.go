package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"server_logic/src/csvs"
)

/*
	圣遗物模块
*/

type Relics struct {
	RelicsId   int   //圣遗物Id
	KeyId      int   //圣遗物唯一标识
	MainEntry  int   //主词条
	Level      int   //等级
	Exp        int   //经验
	OtherEntry []int //副词条
	RoleId     int   //穿戴角色
}

type ModRelics struct {
	RelicsInfo map[int]*Relics
	MaxKey     int

	player *Player
	path   string
}

func (r *ModRelics) AddItem(itemId int, num int64) {
	//验圣遗物是否存在
	config := csvs.GetRelicsConfig(itemId)
	if config == nil {
		fmt.Println("配置不存在")
		return
	}
	//圣遗物最大值为1500
	if len(r.RelicsInfo)+int(num) > csvs.RELICS_MAX_COUNT {
		fmt.Println("超过最大值")
		return
	}

	for i := int64(0); i < num; i++ {
		relics := r.NewRelics(itemId)
		r.RelicsInfo[relics.KeyId] = relics
		relics.ShowInfo()
	}
}

//初始化圣遗物
func (r *ModRelics) NewRelics(itemId int) *Relics {
	relics := new(Relics)
	relics.RelicsId = itemId //圣遗物ID
	r.MaxKey++
	relics.KeyId = r.MaxKey //圣遗物唯一标识
	// 生成主词条
	config := csvs.ConfigRelicsMap[itemId]
	if config == nil {
		return nil
	}
	relics.MainEntry = r.MakeMainEntry(config.MainGroup)
	// 生成副词条
	for i := 0; i < config.OtherGroupNum; i++ {
		// 第四个词条有20%的几率生成
		if i == config.OtherGroupNum-1 {
			randNum := rand.Intn(csvs.PERCENT_ALL)
			if randNum < csvs.ALL_ENTRY_RATE {
				otherEntryId := r.MakeOtherEntry(relics, config.OtherGroup)
				if otherEntryId != 0 {
					relics.OtherEntry = append(relics.OtherEntry, otherEntryId)
				}
			}
		} else {
			otherEntryId := r.MakeOtherEntry(relics, config.OtherGroup)
			if otherEntryId != 0 {
				relics.OtherEntry = append(relics.OtherEntry, otherEntryId)
			}
		}
	}
	return relics
}

// 生成主词条
func (r *ModRelics) MakeMainEntry(mainGroup int) int {
	configs, ok := csvs.ConfigRelicsEntryGroupMap[mainGroup]
	if !ok {
		return 0
	}
	// 总权值
	rateAll := 0
	for _, v := range configs {
		rateAll += v.Weight
	}
	rateNow := 0
	randNum := rand.Intn(rateAll)
	for _, v := range configs {
		rateNow += v.Weight
		if randNum < rateNow {
			return v.Id
		}
	}
	return 0
}

// 生成副词条
func (r *ModRelics) MakeOtherEntry(relics *Relics, otherGroup int) int {
	configs, ok := csvs.ConfigRelicsEntryGroupMap[otherGroup]
	if !ok {
		return 0
	}
	configNow := csvs.GetRelicsConfig(relics.RelicsId)
	if configNow == nil {
		return 0
	}
	if len(relics.OtherEntry) >= configNow.OtherGroupNum {
		//去重样本
		allEntry := make(map[int]int) //已出现过的词条Map
		for _, v := range relics.OtherEntry {
			otherConfig := csvs.ConfigRelicsEntryMap[v]
			if otherConfig != nil {
				allEntry[otherConfig.AttrType] = csvs.LOGIC_TRUE
			}
		}
		// 总权值
		rateAll := 0
		for _, v := range configs {
			_, ok := allEntry[v.AttrType] //已获取副词条的总权重
			if !ok {
				continue
			}
			rateAll += v.Weight
		}
		rateNow := 0
		randNum := rand.Intn(rateAll)
		for _, v := range configs {
			_, ok := allEntry[v.AttrType] //在已获取的副词条类型中选择
			if !ok {
				continue
			}
			rateNow += v.Weight
			if randNum < rateNow {
				return v.Id
			}
		}
	} else {
		//去重样本
		allEntry := make(map[int]int) //已出现过的词条Map
		mainConfig := csvs.ConfigRelicsEntryMap[relics.MainEntry]
		if mainConfig != nil {
			allEntry[mainConfig.AttrType] = csvs.LOGIC_TRUE
		}
		for _, v := range relics.OtherEntry {
			otherConfig := csvs.ConfigRelicsEntryMap[v]
			if otherConfig != nil {
				allEntry[otherConfig.AttrType] = csvs.LOGIC_TRUE
			}
		}
		// 总权值
		rateAll := 0
		for _, v := range configs {
			_, ok := allEntry[v.AttrType] //去重
			if ok {
				continue
			}
			rateAll += v.Weight
		}
		rateNow := 0
		randNum := rand.Intn(rateAll)
		for _, v := range configs {
			_, ok := allEntry[v.AttrType] //去重
			if ok {
				continue
			}
			rateNow += v.Weight
			if randNum < rateNow {
				return v.Id
			}
		}
	}
	return 0
}

func (r *Relics) ShowInfo() {
	fmt.Printf("获得圣遗物:key:%d,Id:%d\n", r.KeyId, r.RelicsId)
	fmt.Printf("当前等级:%d,当前经验:%d\n", r.Level, r.Exp)
	mainEntryConfig := csvs.GetRelicsLevelConfig(r.MainEntry, r.Level)
	if mainEntryConfig != nil {
		fmt.Printf("主词条属性:%s,值:%d\n", mainEntryConfig.AttrName, mainEntryConfig.AttrValue)
	}

	for _, v := range r.OtherEntry {
		otherEntryConfig := csvs.ConfigRelicsEntryMap[v]
		fmt.Printf("副词条属性:%s,值:%d\n", otherEntryConfig.AttrName, otherEntryConfig.AttrValue)
	}
	fmt.Println()
}

// 圣遗物升级
func (r *ModRelics) RelicsUp(keyId int, exp int) {
	relics := r.RelicsInfo[keyId]
	if relics == nil {
		fmt.Println("找不到对应圣遗物")
		return
	}
	relics.Exp += exp
	for {
		//下一级配置
		nextLevelConfig := csvs.GetRelicsLevelConfig(relics.MainEntry, relics.Level+1)
		if nextLevelConfig == nil {
			break
		}
		if relics.Exp < nextLevelConfig.NeedExp {
			break
		}
		relics.Level++
		relics.Exp -= nextLevelConfig.NeedExp
		// 副词条升级
		if relics.Level%4 == 0 {
			relicsConfig := csvs.ConfigRelicsMap[relics.RelicsId]
			if relicsConfig != nil {
				relics.OtherEntry = append(relics.OtherEntry, r.MakeOtherEntry(relics, relicsConfig.OtherGroup))
			}
		}
	}
	relics.ShowInfo()
}

// 模拟一个满级圣遗物
func (r *ModRelics) RelicsTop() {
	relics := r.NewRelics(7000005)
	relics.Level = 20
	config := csvs.GetRelicsConfig(relics.RelicsId)
	if config == nil {
		return
	}
	for i := 0; i < 5; i++ {
		relics.OtherEntry = append(relics.OtherEntry, r.MakeOtherEntry(relics, config.OtherGroup))
	}
	relics.ShowInfo()
}

// 极品双爆头测试
func (r *ModRelics) RelicsTestBest() {
	config := csvs.GetRelicsConfig(7000005)
	if config == nil {
		return
	}
	allTime := 500000
	relicsBestInfo := make([]*Relics, 0)
	for i := 0; i < allTime; i++ {
		relics := r.NewRelics(7000005)
		relics.Level = 20
		config := csvs.GetRelicsConfig(relics.RelicsId)
		if config == nil {
			continue
		}
		for i := 0; i < 5; i++ {
			relics.OtherEntry = append(relics.OtherEntry, r.MakeOtherEntry(relics, config.OtherGroup))
		}

		configMain := csvs.ConfigRelicsEntryMap[relics.MainEntry]
		if configMain == nil {
			continue
		}
		if configMain.AttrType != 4 && configMain.AttrType != 5 {
			continue
		}
		bestEntryConut := 0
		for _, v := range relics.OtherEntry {
			configOther := csvs.ConfigRelicsEntryMap[v]
			if configOther == nil {
				continue
			}
			if configOther.AttrType != 4 && configOther.AttrType != 5 {
				continue
			}
			bestEntryConut++
		}

		if bestEntryConut < 6 {
			continue
		}
		relicsBestInfo = append(relicsBestInfo, relics)
	}
	fmt.Printf("生成了圣遗物头不为:%d个,极品数量%d\n", allTime, len(relicsBestInfo))

	for _, v := range relicsBestInfo {
		v.ShowInfo()
	}
}

func (r *ModRelics) SaveData() {
	content, err := json.Marshal(r)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(r.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (r *ModRelics) LoadData(player *Player) {
	r.player = player
	r.path = r.player.localPath + "/relics.json"

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

	// if r.RelicsInfo == nil {
	// 	r.RelicsInfo = make(map[int]*Relics)
	// }
}

func (r *ModRelics) InitData() {
	if r.RelicsInfo == nil {
		r.RelicsInfo = make(map[int]*Relics)
	}
}
