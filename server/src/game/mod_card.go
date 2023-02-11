package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server_logic/src/csvs"
)

/*
	名片模块
*/

type Card struct {
	CardId int	
}

type ModCard struct {
	CardInfo map[int]*Card	

	player *Player
	path   string
}

func (c *ModCard) IsHasCard(cardId int) bool {
	_, ok := c.CardInfo[cardId]
	return ok
}

func (c *ModCard) AddCard(itemId int, friendliness int) {
	_, ok := c.CardInfo[itemId]
	if ok {
		fmt.Println("已存在该名片")
		return
	}

	config := csvs.GetCardConfig(itemId)
	if config == nil {
		fmt.Println("非法名片")
		return
	}

	if friendliness < config.Friendliness {
		fmt.Println("好感度不足:", itemId)
		return
	}

	c.CardInfo[itemId] = &Card{CardId: itemId}
	fmt.Println("获得名片:", itemId)
}

// 如果已获得该角色，则不再添加名片
func (c *ModCard) CheckGetCard(roleId int, friendliness int) {
	config := csvs.GetCardConfigByRoleId(roleId)
	if config == nil {
		return
	}

	c.AddCard(config.CardId, friendliness)
}

func (c *ModCard) SaveData() {
	content, err := json.Marshal(c)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(c.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (c *ModCard) LoadData(player *Player) {
	c.player = player
	c.path = c.player.localPath + "/card.json"

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

	// if c.CardInfo == nil {
	// 	c.CardInfo = make(map[int]*Card)
	// }
}

func (c *ModCard) InitData() {
	if c.CardInfo == nil {
		c.CardInfo = make(map[int]*Card)
	}
}
