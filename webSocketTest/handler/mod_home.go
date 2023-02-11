package handler

import (
	"fmt"
	"webSocketTest/csvs"
)

/*
	家园家具背包
*/

type Home struct {
	HomeItemId  int
	HomeItemNum int64
	KeyId       int
}

type ModHome struct {
	HomeItemInfo map[int]*Home 
}

func (h *ModHome) AddItem(itemId int, num int64) {
	// _, ok := h.HomeItemInfo[itemId]
	// if ok {
	// 	h.HomeItemInfo[itemId].HomeItemNum += num
	// } else {
	// 	h.HomeItemInfo[itemId] = &Home{HomeItemId: itemId, HomeItemNum: num}
	// }
	// config := csvs.GetItemConfig(itemId)
	// if config != nil {
	// 	fmt.Println("获得家具物品:", config.ItemName, "---数量:", num, "---当前数量:", h.HomeItemInfo[itemId].HomeItemNum)
	// }
}

func (h *ModHome) ShowInfo() {
	for _, v := range h.HomeItemInfo {
		fmt.Printf("KeyId:%d ---家具:%s ---数量:%d\n", v.KeyId, 
		csvs.GetItemConfigName(v.HomeItemId), v.HomeItemNum)
	}
}
