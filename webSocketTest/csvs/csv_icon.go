package csvs

import "webSocketTest/utils"

type ConfigIcon struct {
	IconId int `json:"IconId"`
	Check  int `json:"Check"`
}

var (
	ConfigIconMap         map[int]*ConfigIcon
	ConfigIconMapByRoleId map[int]*ConfigIcon //空间换时间，
)

func init() {
	ConfigIconMap = make(map[int]*ConfigIcon)
	utils.GetCsvUtilMgr().LoadCsv("Icon", &ConfigIconMap)
	ConfigIconMapByRoleId = make(map[int]*ConfigIcon)
	for _, v := range ConfigIconMap {
		ConfigIconMapByRoleId[v.Check] = v
	}
}

func GetIconConfig(iconId int) *ConfigIcon {
	return ConfigIconMap[iconId]
}

func GetIconConfigByRoleId(roleId int) *ConfigIcon {
	return ConfigIconMapByRoleId[roleId]
}
