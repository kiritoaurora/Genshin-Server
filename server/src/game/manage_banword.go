package game

import (
	"fmt"
	"regexp"
	"server_logic/src/csvs"
	"sync"
	"time"
)

var manageBanWord *ManageBanWord

type ManageBanWord struct {
	BanWordBase  []string //配置生成
	BanWordExtra []string //更新
	MsgChannel   chan int
	Lock         sync.RWMutex
}

func GetManageBanWord() *ManageBanWord {
	if manageBanWord == nil {
		manageBanWord = new(ManageBanWord)
		manageBanWord.BanWordBase = []string{"外挂", "操"}
		manageBanWord.BanWordExtra = []string{"原神", "淦"}
		manageBanWord.MsgChannel = make(chan int)
	}

	return manageBanWord
}

func (bw *ManageBanWord) IsBanWord(txt string) bool {
	bw.Lock.RLock()
	defer bw.Lock.RUnlock()
	for _, v := range bw.BanWordBase {
		match, _ := regexp.MatchString(v, txt)
		if match {
			fmt.Println("发现违禁词:", v)
			return match
		}
	}

	for _, v := range bw.BanWordExtra {
		match, _ := regexp.MatchString(v, txt)
		if match {
			fmt.Println("发现违禁词:", v)
			return match
		}
	}

	return false
}

func (bw *ManageBanWord) Run() {
	GetServer().AddGo()
	// 基础词库加载
	bw.BanWordBase = csvs.GetBanWordBase()
	// 基础词库的更新

	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			if time.Now().Unix()%10 == 0 {
				fmt.Println("违禁词库更新完成")
				GetServer().UpdateBanWord(bw.BanWordBase, bw.BanWordExtra)
			}
		case _, ok := <-bw.MsgChannel:
			if !ok {
				GetServer().GoDone()
				return
			}
		}
	}
}

func (bw *ManageBanWord) Close() {
	close(bw.MsgChannel)
}
