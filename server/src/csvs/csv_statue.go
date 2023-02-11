package csvs

import "server_logic/src/utils"

type ConfigStatue struct {
	StatueId    int `json:"StatueId"`
	Level       int `json:"Level"`
	CostItem    int `json:"CostItem"`
	CostNum     int64 `json:"CostNum"`
	RewardItem1 int `json:"RewardItem1"`
	RewardNum1  int `json:"RewardNum1"`
	RewardItem2 int `json:"RewardItem2"`
	RewardNum2  int `json:"RewardNum2"`
}

var (
	ConfigStatueSlice []*ConfigStatue
)

func init() {
	ConfigStatueSlice = make([]*ConfigStatue, 0)
	utils.GetCsvUtilMgr().LoadCsv("Statue", &ConfigStatueSlice)
}