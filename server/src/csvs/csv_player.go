package csvs

import (
	"server_logic/src/utils"
)

type ConfigPlayerLevel struct {
	PlayerLevel int `json:"PlayerLevel"`
	PlayerExp   int `json:"PlayerExp"`
	WorldLevel  int `json:"WorldLevel"`
	ChapterId   int `json:"ChapterId"`
}

var ConfigPlayerLevelSlice []*ConfigPlayerLevel

func init() {

	utils.GetCsvUtilMgr().LoadCsv("PlayerLevel", &ConfigPlayerLevelSlice)
}

func GetNowLevelConfig(level int) *ConfigPlayerLevel {

	if level < 0 || level > len(ConfigPlayerLevelSlice) {
		return nil
	}

	return ConfigPlayerLevelSlice[level-1]
}
