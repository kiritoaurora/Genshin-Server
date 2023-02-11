package handler

import (
	"fmt"
	"webSocketTest/csvs"
)

/*
	使用食谱获得烹饪技能
*/

type Cook struct {
	CookId int //`json:"cookId"`
}

type ModCook struct {
	CookInfo map[int]*Cook 
}

func (c *ModCook) AddItem(itemId int) {

}

func (c *ModCook) ShowInfo() {
	for _, v := range c.CookInfo {
		fmt.Printf("物品:%s ---id:%d\n", csvs.GetItemConfigName(v.CookId), v)
	}
}
