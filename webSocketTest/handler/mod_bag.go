package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"webSocketTest/csvs"
)

type ItemInfo struct {
	ItemId  int   //`json:"itemId"`
	ItemNum int64 //`json:"itemNum"`
}

//背包信息 2
type ModBag struct {
	BagInfo map[int]*ItemInfo
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
		player.ModRole.AddItem(itemId, num)
	case csvs.ITEMTYPE_ICON:
		player.ModIcon.AddItem(itemId)
	case csvs.ITEMTYPE_CARD:
		player.ModCard.AddCard(itemId, 10)
	case csvs.ITEMTYPE_WEAPON:
		player.ModWeapon.AddItem(itemId, num)
	case csvs.ITEMTYPE_RELICS:
		player.ModRelics.AddItem(itemId, num)
	case csvs.ITEMTYPE_COOK:
		player.ModCook.AddItem(itemId)
	case csvs.ITEMTYPE_HOME_ITEM:
		player.ModHome.AddItem(itemId, num)
	default: //同普通
		b.AddItemToBag(itemId, num)
	}
}

func (b *ModBag) AddItemToBag(itemId int, num int64) {
	if num == 0 {
		return
	}
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

func (b *ModBag) ShowInfo() {
	for _, v := range b.BagInfo {
		fmt.Printf("物品:%s ---数量:%d\n", csvs.GetItemConfigName(v.ItemId), v.ItemNum)
	}
}

func (b *ModBag) HasEnoughItem(itemId int, num int64) bool {
	if itemId == 0 {
		return false
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

func (b *ModBag) RemoveItem(itemId int, num int64) {
	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("itemId", strconv.Itoa(itemId))
	params.Set("num", strconv.FormatInt(num, 10))
	resp, err := http.PostForm("http://127.0.0.1:8888/useitem", params)
	if err != nil {
		fmt.Println(err)
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	nowNum, _ := strconv.ParseInt(string(msg), 10, 64)

	_, ok := b.BagInfo[itemId]
	if ok {
		b.BagInfo[itemId].ItemNum = nowNum
	} else {
		b.BagInfo[itemId] = &ItemInfo{ItemId: itemId, ItemNum: 0 - num}
	}
	config := csvs.GetItemConfig(itemId)
	if config != nil {
		fmt.Println("扣除物品：", config.ItemName, "---数量：", num, "---当前数量：", b.BagInfo[itemId].ItemNum)
	}
}
