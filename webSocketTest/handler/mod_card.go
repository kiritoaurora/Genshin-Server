package handler

import (
	"fmt"
	"webSocketTest/csvs"
)

type Card struct {
	CardId int //`json:"cardId"`
}

type ModCard struct {
	CardInfo map[int]*Card 
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

func (c *ModCard) ShowInfo() {
	for _, v := range c.CardInfo {
		fmt.Printf("名片id:%d\t---%s\n",v.CardId, csvs.GetItemConfigName(v.CardId))
	}
}
