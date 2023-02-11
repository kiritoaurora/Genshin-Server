package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"webSocketTest/csvs"
)

/*
	武器模块
*/

type Weapon struct {
	WeaponId    int //武器Id
	KeyId       int //武器唯一标识
	Level       int //武器等级
	Exp         int //武器当前拥有的经验
	StarLevel   int //武器突破星级
	RefineLevel int //精炼等级
	RoleId      int //武器装备的角色Id
}

type ModWeapon struct {
	WeaponInfo map[int]*Weapon
}

func (w *ModWeapon) AddItem(itemId int, num int64) {

}

// 武器升级
func (w *ModWeapon) WeaponUp(keyId, exp int) {
	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("keyId", strconv.Itoa(keyId))
	params.Set("exp", strconv.Itoa(exp))
	resp, err := http.PostForm("http://127.0.0.1:8888/weaponup", params)
	if err != nil {
		fmt.Println(err)
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	newLevel, _ := strconv.Atoi(string(msg))
	w.WeaponInfo[keyId].Level = newLevel
}

// 武器突破
func (w *ModWeapon) WeaponUpStar(keyId int) {
	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("keyId", strconv.Itoa(keyId))
	resp, err := http.PostForm("http://127.0.0.1:8888/weaponupstar", params)
	if err != nil {
		fmt.Println(err)
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	newStarLevel, _ := strconv.Atoi(string(msg))
	w.WeaponInfo[keyId].StarLevel = newStarLevel
}

// 武器精炼
func (w *ModWeapon) WeaponUpRefine(weaponKeyId int, weaponTargetKeyId int) {
	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("weaponKeyId", strconv.Itoa(weaponKeyId))
	params.Set("weaponTargetKeyId", strconv.Itoa(weaponTargetKeyId))
	resp, err := http.PostForm("http://127.0.0.1:8888/weaponuprefine", params)
	if err != nil {
		fmt.Println(err)
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	newRefineLevel, _ := strconv.Atoi(string(msg))
	w.WeaponInfo[weaponKeyId].RefineLevel = newRefineLevel
}

func (w *Weapon) ShowInfo() {
	fmt.Printf("key:%d,Id:%d\n", w.KeyId, w.WeaponId)
	fmt.Printf("当前等级:%d,当前经验:%d,当前突破等级:%d,当前精炼等级:%d\n", w.Level, 
	w.Exp, w.StarLevel, w.RefineLevel)
}

func (w *ModWeapon) ShowInfo() {
	for _, v := range w.WeaponInfo {
		fmt.Printf("Keyd:%d ---%s ---等级:%d ---突破等级:%d ---精炼等级:%d\n", v.KeyId,
			csvs.GetItemConfigName(v.WeaponId), v.Level, v.StarLevel, v.RefineLevel)
	}
}
