package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server_logic/src/csvs"
)

/*
	头像模块
*/

type Icon struct {
	IconId int	
}

type ModIcon struct {
	IconInfo map[int]*Icon	

	player *Player
	path   string
}

func (i *ModIcon) IsHasIcon(iconId int) bool {
	_, ok := i.IconInfo[iconId]
	return ok
}

func (i *ModIcon) AddItem(itemId int) {
	_, ok := i.IconInfo[itemId]
	if ok {
		fmt.Println("已存在该头像")
		return
	}
	config := csvs.GetIconConfig(itemId)
	if config == nil {
		fmt.Println("非法头像")
		return
	}

	i.IconInfo[itemId] = &Icon{IconId: itemId}
	fmt.Println("获得头像", itemId)
}

// 如果已获得该角色，则不再添加头像
func (i *ModIcon) CheckGetIcon(roleId int) {
	config := csvs.GetIconConfigByRoleId(roleId)
	if config == nil {
		return
	}
	i.AddItem(config.IconId)
}

func (i *ModIcon) SaveData() {
	content, err := json.Marshal(i)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(i.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (i *ModIcon) LoadData(player *Player) {
	i.player = player
	i.path = i.player.localPath + "/icon.json"

	configFile, err := ioutil.ReadFile(i.path)
	if err != nil {
		fmt.Println("暂无存档")
		// return
	}
	err = json.Unmarshal(configFile, &i)
	if err != nil {
		i.InitData()
		return
	}

	// if i.IconInfo == nil {
	// 	i.IconInfo = make(map[int]*Icon)
	// }
}

func (i *ModIcon) InitData() {
	if i.IconInfo == nil {
		i.IconInfo = make(map[int]*Icon)
	}
}
