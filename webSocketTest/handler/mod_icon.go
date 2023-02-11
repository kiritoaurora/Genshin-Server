package handler

import (
	"fmt"
	"webSocketTest/csvs"
)

type Icon struct {
	IconId int //`json:"iconId"`
}

type ModIcon struct {
	IconInfo map[int]*Icon 
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

func (i *ModIcon) ShowInfo() {
	for _, v := range i.IconInfo {
		fmt.Printf("头像id:%d\t---%s\n",v.IconId, csvs.GetItemConfigName(v.IconId))
	}
}