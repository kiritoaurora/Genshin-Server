package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server_logic/src/csvs"
)

/*
	背包系统
*/

type ItemInfo struct {
	ItemId  int   //`json:"itemId"`
	ItemNum int64 //`json:"itemNum"`
}

type ModBag struct {
	BagInfo map[int]*ItemInfo	

	player *Player
	path   string
}

func (b *ModBag) AddItem(itemId int, num int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	switch itemConfig.SortType {
	// case csvs.ITEMTYPE_NORMAL:
	// 	b.AddItemToBag(itemId, num)
	case csvs.ITEMTYPE_ROLE:
		b.player.GetModRole().AddItem(itemId, num)
	case csvs.ITEMTYPE_ICON:
		b.player.GetModIcon().AddItem(itemId)
	case csvs.ITEMTYPE_CARD:
		b.player.GetModCard().AddCard(itemId, 10)
	case csvs.ITEMTYPE_WEAPON:
		b.player.GetModWeapon().AddItem(itemId, num)
	case csvs.ITEMTYPE_RELICS:
		b.player.GetModRelics().AddItem(itemId, num)
	case csvs.ITEMTYPE_COOK:
		b.player.GetModCook().AddItem(itemId)
	case csvs.ITEMTYPE_HOME_ITEM:
		b.player.GetModHome().AddItem(itemId, num)
	default: //同普通
		b.AddItemToBag(itemId, num)
	}
}

func (b *ModBag) AddItemToBag(itemId int, num int64) {
	_, ok := b.BagInfo[itemId]
	if ok {
		b.BagInfo[itemId].ItemNum += num
	} else {
		b.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("获得物品：", config.ItemName, "---数量：", num, "---当前数量：", b.BagInfo[itemId].ItemNum)
	}
}

func (b *ModBag) RemoveItem(itemId int, num int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	switch itemConfig.SortType {
	case csvs.ITEMTYPE_NORMAL:
		b.RemoveItemToBagGM(itemId, num)
	default: //同普通
		b.RemoveItemToBag(itemId, num)
	}
}

func (b *ModBag) RemoveItemToBagGM(itemId int, num int64) {
	_, ok := b.BagInfo[itemId]
	if ok {
		b.BagInfo[itemId].ItemNum -= num
	} else {
		b.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: 0 - num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("扣除物品：", config.ItemName, "---数量：", num, "---当前数量：", b.BagInfo[itemId].ItemNum)
	}
}

func (b *ModBag) RemoveItemToBag(itemId int, num int64) {
	if itemId == 0 {
		return
	}
	itemConfig := csvs.GetItemConfig(itemId)
	switch itemConfig.SortType {
	case csvs.ITEMTYPE_ROLE:
		fmt.Println("此物品无法扣除")
		return
	case csvs.ITEMTYPE_ICON:
		fmt.Println("此物品无法扣除")
		return
	case csvs.ITEMTYPE_CARD:
		fmt.Println("此物品无法扣除")
		return
	}

	if !b.HasEnoughItem(itemId, num) {
		config := csvs.GetItemConfig(itemId)
		if config != nil {
			nowNum := int64(0)
			_, ok := b.BagInfo[itemId]
			if ok {
				nowNum = b.BagInfo[itemId].ItemNum
			}
			fmt.Println(config.ItemName, "数量不足", "---当前数量：", nowNum)
		}
		return
	}

	_, ok := b.BagInfo[itemId]
	if ok {
		b.BagInfo[itemId].ItemNum -= num
	} else {
		b.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: 0 - num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("扣除物品：", config.ItemName, "---数量：", num, "---当前数量：", b.BagInfo[itemId].ItemNum)
	}
}

func (b *ModBag) HasEnoughItem(itemId int, num int64) bool {
	if itemId == 0 {
		return true
	}
	_, ok := b.BagInfo[itemId]
	if !ok {
		return false
	} else if b.BagInfo[itemId].ItemNum < num {
		return false
	}

	return true
}

func (b *ModBag) UseItem(itemId int, num int64) {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	if !b.HasEnoughItem(itemId, num) {
		config := csvs.GetItemConfig(itemId)
		if config != nil {
			nowNum := int64(0)
			_, ok := b.BagInfo[itemId]
			if ok {
				nowNum = b.BagInfo[itemId].ItemNum
			}
			fmt.Println(config.ItemName, "数量不足", "---当前数量：", nowNum)
		}
		return
	}

	switch itemConfig.SortType {
	case csvs.ITEMTYPE_COOKBOOK:
		b.UseCookBook(itemId, num)
	case csvs.ITEMTYPE_FOOD:
		// 给角色加BUFF
	default: //同普通
		fmt.Println(itemId, "此物品无法使用")
		return
	}
}

func (b *ModBag) UseCookBook(itemId int, num int64) {
	cookBookConfig := csvs.GetCookBookConfig(itemId)
	if cookBookConfig == nil {
		fmt.Println(itemId, "物品不存在")
		return
	}

	b.RemoveItem(itemId, num)
	b.AddItem(cookBookConfig.Reward, num)
}

func (b *ModBag) GetItemNum(itemId int) int64 {
	itemConfig := csvs.GetItemConfig(itemId)
	if itemConfig == nil {
		return 0
	}
	_, ok := b.BagInfo[itemId]
	if !ok {
		return 0
	}
	return b.BagInfo[itemId].ItemNum
}

func (b *ModBag) SaveData() {
	content, err := json.Marshal(b)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(b.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (b *ModBag) LoadData(player *Player) {
	b.player = player
	b.path = b.player.localPath + "/bag.json"

	configFile, err := ioutil.ReadFile(b.path)
	if err != nil {
		fmt.Println("暂无存档")
		// return
	}
	err = json.Unmarshal(configFile, &b)
	if err != nil {
		b.InitData()
		return
	}

	// if b.BagInfo == nil {
	// 	b.BagInfo = make(map[int]*ItemInfo)
	// }
}

func (b *ModBag) InitData() {
	if b.BagInfo == nil {
		b.BagInfo = make(map[int]*ItemInfo)
	}
}
