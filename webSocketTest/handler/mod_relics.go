package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"webSocketTest/csvs"
)

/*
	圣遗物模块
*/

type Relics struct {
	RelicsId   int   //圣遗物Id
	KeyId      int   //圣遗物唯一标识
	MainEntry  int   //主词条
	Level      int   //等级
	Exp        int   //经验
	OtherEntry []int //副词条
	RoleId     int   //穿戴角色
}

type ModRelics struct {
	RelicsInfo map[int]*Relics
}

func (r *ModRelics) AddItem(itemId int, num int64) {
	config := csvs.GetRelicsConfig(itemId)
	if config == nil {
		fmt.Println("配置不存在")
		return
	}
	//圣遗物最大值为1500
	// if len(r.RelicsInfo)+int(num) > csvs.RELICS_MAX_COUNT {
	// 	fmt.Println("超过最大值")
	// 	return
	// }

	// for i := int64(0); i < num; i++ {
	// 	relics := r.NewRelics(itemId)
	// 	r.RelicsInfo[relics.KeyId] = relics
	// 	relics.ShowInfo()
	// }
}

func (r *ModRelics) ShowInfo() {
	for _, v := range r.RelicsInfo {
		fmt.Printf("Keyd:%d ---%s ---等级:%d\n", v.KeyId,
			csvs.GetItemConfigName(v.RelicsId), v.Level)
	}
}

func (r *ModRelics) RelicsUp(keyId int, exp int) {
	params := url.Values{}
	params.Set("userId", strconv.FormatInt(player.UserId, 10))
	params.Set("keyId", strconv.Itoa(keyId))
	params.Set("exp", strconv.Itoa(exp))
	resp, err := http.PostForm("http://127.0.0.1:8888/relicsup", params)
	if err != nil {
		fmt.Println(err)
	}
	msg, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var relics *Relics
	err = json.Unmarshal(msg, &relics)
	if err != nil {
		fmt.Println("解析失败！")
		return
	}
	r.RelicsInfo[keyId] = relics
}
