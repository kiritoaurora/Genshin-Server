package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server_logic/src/csvs"
)

/*
	使用食谱获得烹饪技能
*/

type Cook struct {
	CookId int	//`json:"cookId"`
}

type ModCook struct {
	CookInfo map[int]*Cook	

	player *Player
	path   string
}

func (c *ModCook) AddItem(itemId int) {
	_, ok := c.CookInfo[itemId]
	if ok {
		fmt.Println("已习得：", csvs.GetItemConfig(itemId).ItemName)
		return
	}
	config := csvs.GetCookConfig(itemId)
	if config == nil {
		fmt.Println("没有食谱：", csvs.GetItemConfig(itemId).ItemName)
		return
	}

	c.CookInfo[itemId] = &Cook{CookId: itemId}
	fmt.Println("学会烹饪：", csvs.GetItemConfig(itemId).ItemName)
}

func (c *ModCook) SaveData() {
	content, err := json.Marshal(c)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(c.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (c *ModCook) LoadData(player *Player) {
	c.player = player
	c.path = c.player.localPath + "/cook.json"

	configFile, err := ioutil.ReadFile(c.path)
	if err != nil {
		fmt.Println("暂无存档")
		// return
	}
	err = json.Unmarshal(configFile, &c)
	if err != nil {
		c.InitData()
		return
	}

	// if c.CookInfo == nil {
	// 	c.CookInfo = make(map[int]*Cook)
	// }
}

func (c *ModCook) InitData() {
	if c.CookInfo == nil {
		c.CookInfo = make(map[int]*Cook)
	}
}
