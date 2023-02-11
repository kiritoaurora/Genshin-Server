package csvs

import (
	"fmt"
	"math/rand"
)

var (
	ConfigDropGroupMap        map[int]*DropGroup
	ConfigDropItemGroupMap    map[int]*DropItemGroup
	ConfigStatueMap           map[int]map[int]*ConfigStatue
	ConfigRelicsEntryGroupMap map[int]map[int]*ConfigRelicsEntry
	ConfigRelicsLevelGroupMap map[int]map[int]*ConfigRelicsLevel
	ConfigRelicsSuitMap       map[int][]*ConfigRelicsSuit
	ConfigWeaponLevelMap      map[int]map[int]*ConfigWeaponLevel
	ConfigWeaponStarMap       map[int]map[int]*ConfigWeaponStar
)

type DropGroup struct {
	DropId      int
	WeightAll   int
	DropConfigs []*ConfigDrop
}

type DropItemGroup struct {
	DropId          int
	DropItemConfigs []*ConfigDropItem
}

func CheckLoadCsv() {
	// 二次处理
	MakeDropGroupMap()
	MakeDropItemGroupMap()
	MakeConfigStatusMap()
	MakeConfigRelicsEntryGroupMap()
	MakeConfigRelicsLevelGroupMap()
	MakeConfigRelicsSuitMap()
	MakeConfigWeaponLevelMap()
	MakeConfigWeaponStarMap()
	fmt.Println("csv配置读取完成")
}

// 抽卡掉落组
func MakeDropGroupMap() {
	ConfigDropGroupMap = make(map[int]*DropGroup)
	for _, v := range ConfigDropSlice {
		dropGroup, ok := ConfigDropGroupMap[v.DropId]
		if !ok {
			dropGroup = new(DropGroup)
			dropGroup.DropId = v.DropId
			ConfigDropGroupMap[v.DropId] = dropGroup
		}

		dropGroup.WeightAll += v.Weight
		dropGroup.DropConfigs = append(dropGroup.DropConfigs, v)
	}
	//RandDropTest()
}

// 物品掉落组
func MakeDropItemGroupMap() {
	ConfigDropItemGroupMap = make(map[int]*DropItemGroup)
	for _, v := range ConfigDropItemSlice {
		dropItemGroup, ok := ConfigDropItemGroupMap[v.DropId]
		if !ok {
			dropItemGroup = new(DropItemGroup)
			dropItemGroup.DropId = v.DropId
			ConfigDropItemGroupMap[v.DropId] = dropItemGroup
		}

		dropItemGroup.DropItemConfigs = append(dropItemGroup.DropItemConfigs, v)
	}
	//RandDropItemTest()
}

// 七天神像
func MakeConfigStatusMap() {
	ConfigStatueMap = make(map[int]map[int]*ConfigStatue)
	for _, v := range ConfigStatueSlice {
		statusMap, ok := ConfigStatueMap[v.StatueId]
		if !ok {
			statusMap = make(map[int]*ConfigStatue)
			ConfigStatueMap[v.StatueId] = statusMap
		}
		statusMap[v.Level] = v
	}
}

// 圣遗物词条组
func MakeConfigRelicsEntryGroupMap() {
	ConfigRelicsEntryGroupMap = make(map[int]map[int]*ConfigRelicsEntry)
	for _, v := range ConfigRelicsEntryMap {
		groupMap, ok := ConfigRelicsEntryGroupMap[v.Group]
		if !ok {
			groupMap = make(map[int]*ConfigRelicsEntry)
			ConfigRelicsEntryGroupMap[v.Group] = groupMap
		}
		groupMap[v.Id] = v
	}
}

// 圣遗物等级组
func MakeConfigRelicsLevelGroupMap() {
	ConfigRelicsLevelGroupMap = make(map[int]map[int]*ConfigRelicsLevel)
	for _, v := range ConfigRelicsLevelSlice {
		levelMap, ok := ConfigRelicsLevelGroupMap[v.EntryId]
		if !ok {
			levelMap = make(map[int]*ConfigRelicsLevel)
			ConfigRelicsLevelGroupMap[v.EntryId] = levelMap
		}
		levelMap[v.Level] = v
	}
}

// 圣遗物套装Map
func MakeConfigRelicsSuitMap() {
	ConfigRelicsSuitMap = make(map[int][]*ConfigRelicsSuit)
	for _, v := range ConfigRelicsSuitSlice {
		ConfigRelicsSuitMap[v.Type] = append(ConfigRelicsSuitMap[v.Type], v)
	}
}

//武器等级
func MakeConfigWeaponLevelMap() {
	ConfigWeaponLevelMap = make(map[int]map[int]*ConfigWeaponLevel)
	for _, v := range ConfigWeaponLevelSlice {
		levelMap, ok := ConfigWeaponLevelMap[v.WeaponStar]
		if !ok {
			levelMap = make(map[int]*ConfigWeaponLevel)
			ConfigWeaponLevelMap[v.WeaponStar] = levelMap
		}
		levelMap[v.Level] = v
	}
}

//武器突破星级	
func MakeConfigWeaponStarMap() {
	ConfigWeaponStarMap = make(map[int]map[int]*ConfigWeaponStar)
	for _, v := range ConfigWeaponStarSlice {
		starMap, ok := ConfigWeaponStarMap[v.WeaponStar]
		if !ok {
			starMap = make(map[int]*ConfigWeaponStar)
			ConfigWeaponStarMap[v.WeaponStar] = starMap
		}
		starMap[v.StarLevel] = v
	}
}

// func RandDropItemTest() {
// 	dropItemGroup := ConfigDropItemGroupMap[1]
// 	if dropItemGroup == nil {
// 		return
// 	}
// 	for _, v := range dropItemGroup.DropItemConfigs {
// 		randNum := rand.Intn(PERCENT_ALL)
// 		if randNum < v.Weight {
// 			fmt.Println(v.ItemId)
// 		}
// 	}
// }

// func RandDropTest() {
// 	dropGroup := ConfigDropGroupMap[1000]
// 	if dropGroup == nil {
// 		return
// 	}
// 	num := 0
// 	for {
// 		config := GetRandDropNew(dropGroup)
// 		if config.IsEnd == LOGIC_TRUE {
// 			fmt.Println(GetItemConfig(config.Result).ItemName)
// 			num++
// 			dropGroup = ConfigDropGroupMap[1000]
// 			if num >= 100 {
// 				break
// 			} else {
// 				continue
// 			}
// 		}
// 		dropGroup = ConfigDropGroupMap[config.Result]
// 		if dropGroup == nil {
// 			break
// 		}
// 	}
// }

// func GetRandDrop(dropGroup *DropGroup) *ConfigDrop {
// 	randNum := rand.Intn(dropGroup.WeightAll)
// 	randNow := 0
// 	for _, v := range dropGroup.DropConfigs {
// 		randNow += v.Weight
// 		if randNum < randNow {
// 			return v
// 		}
// 	}
// 	return nil
// }

//抽卡
func GetRandDropNew(dropGroup *DropGroup) *ConfigDrop {
	randNum := rand.Intn(dropGroup.WeightAll)
	randNow := 0
	for _, v := range dropGroup.DropConfigs {
		randNow += v.Weight
		if randNum < randNow {
			if v.IsEnd == LOGIC_TRUE {
				return v
			}
			dropGroup = ConfigDropGroupMap[v.Result]
			if dropGroup == nil {
				return nil
			}
			return GetRandDropNew(dropGroup)
		}
	}
	return nil
}

//仓检抽卡
func GetRandDropNew1(dropGroup *DropGroup, fiveStarInfo, fourStarInfo map[int]int) *ConfigDrop {
	for _, v := range dropGroup.DropConfigs {
		_, ok := fiveStarInfo[v.Result]
		if ok {
			index := 0
			maxGetTimes := 0
			for k, config := range dropGroup.DropConfigs {
				_, nowOK := fiveStarInfo[config.Result]
				if !nowOK {
					continue
				}
				if maxGetTimes < fiveStarInfo[config.Result] {
					maxGetTimes = fiveStarInfo[config.Result]
					index = k
				}
			}
			return dropGroup.DropConfigs[index]
		}
		_, ok = fourStarInfo[v.Result]
		if ok {
			index := 0
			maxGetTimes := 0
			for k, config := range dropGroup.DropConfigs {
				_, nowOK := fourStarInfo[config.Result]
				if !nowOK {
					continue
				}
				if maxGetTimes < fourStarInfo[config.Result] {
					maxGetTimes = fourStarInfo[config.Result]
					index = k
				}
			}
			return dropGroup.DropConfigs[index]
		}

	}

	randNum := rand.Intn(dropGroup.WeightAll)
	randNow := 0
	for _, v := range dropGroup.DropConfigs {
		randNow += v.Weight
		if randNum < randNow {
			if v.IsEnd == LOGIC_TRUE {
				return v
			}
			dropGroup = ConfigDropGroupMap[v.Result]
			if dropGroup == nil {
				return nil
			}
			return GetRandDropNew1(dropGroup, fiveStarInfo, fourStarInfo)
		}
	}
	return nil
}

//仓检抽卡
func GetRandDropNew2(dropGroup *DropGroup, fiveStarInfo, fourStarInfo map[int]int) *ConfigDrop {
	for _, v := range dropGroup.DropConfigs {
		_, ok := fiveStarInfo[v.Result]
		if ok {
			index := 0
			minGetTimes := 0
			for k, config := range dropGroup.DropConfigs {
				_, nowOK := fiveStarInfo[config.Result]
				if !nowOK {
					index = k
					break
				}
				if minGetTimes == 0 || minGetTimes > fiveStarInfo[config.Result] {
					minGetTimes = fiveStarInfo[config.Result]
					index = k
				}
			}
			return dropGroup.DropConfigs[index]
		}
		_, ok = fourStarInfo[v.Result]
		if ok {
			index := 0
			minGetTimes := 0
			for k, config := range dropGroup.DropConfigs {
				_, nowOK := fourStarInfo[config.Result]
				if !nowOK {
					index = k
					break
				}
				if minGetTimes == 0 || minGetTimes > fourStarInfo[config.Result] {
					minGetTimes = fourStarInfo[config.Result]
					index = k
				}
			}
			return dropGroup.DropConfigs[index]
		}

	}

	randNum := rand.Intn(dropGroup.WeightAll)
	randNow := 0
	for _, v := range dropGroup.DropConfigs {
		randNow += v.Weight
		if randNum < randNow {
			if v.IsEnd == LOGIC_TRUE {
				return v
			}
			dropGroup = ConfigDropGroupMap[v.Result]
			if dropGroup == nil {
				return nil
			}
			return GetRandDropNew2(dropGroup, fiveStarInfo, fourStarInfo)
		}
	}
	return nil
}

// 返回掉落组
func GetDropItemGroup(dropId int) *DropItemGroup {
	return ConfigDropItemGroupMap[dropId]
}

// 获取掉落组，掉落物品
func GetRandDropItemNew(dropId int) []*ConfigDropItem {
	rel := make([]*ConfigDropItem, 0) //初始化掉落物品切片
	//无掉落物z返回空切片
	if dropId == 0 {
		return rel
	}
	config := GetDropItemGroup(dropId) //掉落组
	configAll := make([]*ConfigDropItem, 0)
	for _, v := range config.DropItemConfigs {
		if v.DropType == DROP_ITEM_TYPE_ITEM { //掉落类型为掉落物
			rel = append(rel, v)
		} else if v.DropType == DROP_ITEM_TYPE_GROUP { //掉落类型为掉落组
			randNum := rand.Intn(PERCENT_ALL)
			if randNum < v.Weight { //判定几率
				config := GetRandDropItemNew(v.ItemId) //递归选取掉落物
				rel = append(rel, config...)
			}
		} else if v.DropType == DROP_ITEM_TYPE_WEIGHT { //掉落类型为权值掉落
			configAll = append(configAll, v)
		}
	}
	// 权值掉落组中有掉落物,选取一个作为掉落物
	if len(configAll) > 0 {
		allRate := 0 //掉落物总权重
		for _, v := range configAll {
			allRate += v.Weight
		}
		randNum := rand.Intn(allRate)
		nowRate := 0
		for _, v := range configAll {
			nowRate += v.Weight
			if randNum < nowRate { //抽中该掉落物，则构建一个新的掉落物配置，权重为PERCENT_ALL
				newConfig := new(ConfigDropItem)
				newConfig.DropId = v.DropId
				newConfig.DropType = v.DropType
				newConfig.ItemId = v.ItemId
				newConfig.ItemNumMax = v.ItemNumMax
				newConfig.ItemNumMin = v.ItemNumMin
				newConfig.Weight = PERCENT_ALL
				newConfig.WorldAdd = v.WorldAdd
				rel = append(rel, newConfig)
				break
			}
		}
	}
	return rel
}

// 获取七天神像配置
func GetStatueConfig(statueId int, level int) *ConfigStatue {
	_, ok := ConfigStatueMap[statueId]
	if !ok {
		return nil
	}

	_, ok = ConfigStatueMap[statueId][level]
	if !ok {
		return nil
	}
	return ConfigStatueMap[statueId][level]
}

// 获取词条配置
func GetRelicsLevelConfig(entryId int, level int) *ConfigRelicsLevel {
	_, ok := ConfigRelicsLevelGroupMap[entryId]
	if !ok {
		return nil
	}

	_, ok = ConfigRelicsLevelGroupMap[entryId][level]
	if !ok {
		return nil
	}
	return ConfigRelicsLevelGroupMap[entryId][level]
}

// 获取武器等级配置
func GetWeaponLevelConfig(star int, level int) *ConfigWeaponLevel {
	_, ok := ConfigWeaponLevelMap[star]
	if !ok {
		return nil
	}

	_, ok = ConfigWeaponLevelMap[star][level]
	if !ok {
		return nil
	}
	return ConfigWeaponLevelMap[star][level]
}

// 获取武器突破星级配置
func GetWeaponStarConfig(weaponStar int, starLevel int) *ConfigWeaponStar {
	_, ok := ConfigWeaponStarMap[weaponStar]
	if !ok {
		return nil
	}

	_, ok = ConfigWeaponStarMap[weaponStar][starLevel]
	if !ok {
		return nil
	}
	return ConfigWeaponStarMap[weaponStar][starLevel]
}