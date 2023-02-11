package csvs

import "webSocketTest/utils"

type ConfigRelics struct {
	RelicsId      int `json:"RelicsId"`
	Type          int `json:"Type"`
	Pos           int `json:"Pos"`
	Star          int `json:"Star"`
	MainGroup     int `json:"MainGroup"`
	OtherGroup    int `json:"OtherGroup"`
	OtherGroupNum int `json:"OtherGroupNum"`
}

type ConfigRelicsEntry struct {
	Id        int    `json:"Id"`
	Group     int    `json:"Group"`
	AttrType  int    `json:"AttrType"`
	AttrName  string `json:"AttrName"`
	AttrValue int    `json:"AttrValue"`
	Weight    int    `json:"Weight"`
}

type ConfigRelicsLevel struct {
	EntryId   int    `json:"EntryId"`
	Level     int    `json:"Level"`
	AttrType  int    `json:"AttrType"`
	AttrName  string `json:"AttrName"`
	AttrValue int    `json:"AttrValue"`
	NeedExp   int    `json:"NeedExp"`
}

type ConfigRelicsSuit struct {
	Type        int    `json:"Type"`
	Num         int    `json:"Num"`
	SuitSkill   int    `json:"SuitSkill"`
	SkillString string `json:"SkillString"`
}

var (
	ConfigRelicsMap        map[int]*ConfigRelics
	ConfigRelicsEntryMap   map[int]*ConfigRelicsEntry
	ConfigRelicsLevelSlice []*ConfigRelicsLevel
	ConfigRelicsSuitSlice  []*ConfigRelicsSuit
)

func init() {
	ConfigRelicsMap = make(map[int]*ConfigRelics)
	utils.GetCsvUtilMgr().LoadCsv("Relics", &ConfigRelicsMap)

	ConfigRelicsEntryMap = make(map[int]*ConfigRelicsEntry)
	utils.GetCsvUtilMgr().LoadCsv("RelicsEntry", &ConfigRelicsEntryMap)

	ConfigRelicsLevelSlice = make([]*ConfigRelicsLevel, 0)
	utils.GetCsvUtilMgr().LoadCsv("RelicsLevel", &ConfigRelicsLevelSlice)

	ConfigRelicsSuitSlice = make([]*ConfigRelicsSuit, 0)
	utils.GetCsvUtilMgr().LoadCsv("RelicsSuit", &ConfigRelicsSuitSlice)
}

func GetRelicsConfig(relicsId int) *ConfigRelics {
	return ConfigRelicsMap[relicsId]
}
