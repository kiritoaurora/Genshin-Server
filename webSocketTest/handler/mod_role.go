package handler

import (
	"fmt"
	"webSocketTest/csvs"
)

/*
	角色模块
*/

type RoleInfo struct {
	RoleId   int
	GetTimes int
	//等级、经验、圣遗物、武器
	RelicsInfo []int
	WeaponInfo int
}

type ModRole struct {
	RoleInfo map[int]*RoleInfo 
	// HpPool    int
	// HpCalTime int64
}

func (r *ModRole) IsHasRole(roleId int) bool {
	return true
}

func (r *ModRole) GetRoleLevel(roleId int) int {
	return 80
}

func (r *ModRole) ShowInfo() {
	for _, v := range r.RoleInfo {
		fmt.Printf("角色:%s \n", csvs.GetItemConfigName(v.RoleId))
	}
}

func (r *ModRole) AddItem(roleId int, num int64) {
	config := csvs.GetRoleConfig(roleId)
	if config == nil {
		fmt.Println("配置不存在roleId:", roleId)
		return
	}

	// for i := 0; i < int(num); i++ {
	// 	// 是否已拥有
	// 	_, ok := r.RoleInfo[roleId]
	// 	if !ok { //未获得
	// 		data := new(RoleInfo)
	// 		data.RoleId = roleId
	// 		data.GetTimes = 1
	// 		r.RoleInfo[roleId] = data
	// 	} else {
	// 		// 判断实际获得的东西，五星：命座*1+星辉*10 | 星辉*25，四星：命座*1+星辉*2 | 星辉*10
	// 		r.RoleInfo[roleId].GetTimes++
	// 		if r.RoleInfo[roleId].GetTimes >= csvs.ADD_ROLE_TIME_NORMAL_MIN &&
	// 			r.RoleInfo[roleId].GetTimes <= csvs.ADD_ROLE_TIME_NORMAL_MAX {
	// 			r.player.GetModBag().AddItemToBag(config.Stuff, config.StuffNum)
	// 			r.player.GetModBag().AddItemToBag(config.StuffItem, config.StuffItemNum)
	// 		} else {
	// 			r.player.GetModBag().AddItemToBag(config.MaxStuffItem, config.MaxStuffItemNum)
	// 		}
	// 	}
	// }
	// itemConfig := csvs.GetItemConfig(roleId)
	// if itemConfig != nil {
	// 	fmt.Println("获得角色：", itemConfig.ItemName, "---", r.RoleInfo[roleId].GetTimes, "次")
	// }
	// // 获取角色头像和名片
	// r.player.GetModIcon().CheckGetIcon(roleId)
	// r.player.GetModCard().CheckGetCard(roleId, 10)
}
