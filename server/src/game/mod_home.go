package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server_logic/src/csvs"
)

/*
	家园家具背包
*/

type Home struct {
	HomeItemId  int		`json:""`
	HomeItemNum int64
	KeyId       int
}

type ModHome struct {
	HomeItemInfo map[int]*Home

	player *Player
	path   string
}

func (h *ModHome) AddItem(itemId int, num int64) {
	_, ok := h.HomeItemInfo[itemId]
	if ok {
		h.HomeItemInfo[itemId].HomeItemNum += num
	} else {
		h.HomeItemInfo[itemId] = &Home{HomeItemId: itemId, HomeItemNum: num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("获得家具物品:", config.ItemName, "---数量:", num, "---当前数量:", h.HomeItemInfo[itemId].HomeItemNum)
	}
}

func (h *ModHome) SaveData() {
	content, err := json.Marshal(h)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(h.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (h *ModHome) LoadData(player *Player) {
	h.player = player
	h.path = h.player.localPath + "/home.json"

	configFile, err := ioutil.ReadFile(h.path)
	if err != nil {
		fmt.Println("暂无存档")
		// return
	}
	err = json.Unmarshal(configFile, &h)
	if err != nil {
		h.InitData()
		return
	}

	// if h.HomeItemInfo == nil {
	// 	h.HomeItemInfo = make(map[int]*Home)
	// }
}

func (h *ModHome) InitData() {
	if h.HomeItemInfo == nil {
		h.HomeItemInfo = make(map[int]*Home)
	}
}
