package game

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

var managePlayer *ManagePlayer

type ManagePlayer struct {
	Players map[int64]*Player
	lock    sync.RWMutex
}

func GetManagePlayer() *ManagePlayer {
	if managePlayer == nil {
		managePlayer = new(ManagePlayer)
		managePlayer.Players = make(map[int64]*Player)
		managePlayer.lock = sync.RWMutex{}
	}
	return managePlayer
}

func (mp *ManagePlayer) PlayerLogin(ws *websocket.Conn, userId int64) *Player {
	mp.lock.Lock()
	defer mp.lock.Unlock()

	playerInfo, ok := mp.Players[userId]
	if ok {
		//处理顶号
		if playerInfo.ws != ws {
			oldWs := playerInfo.ws
			playerInfo.ws = ws
			playerInfo.exitTime = 0
			if oldWs != nil {
				oldWs.Write([]byte("账号在别处登陆"))
				oldWs.Close()
			}
		}
	}
	playerInfo = NewTestPlayer(ws, userId)
	mp.Players[userId] = playerInfo
	mp.SendPlayerMsg(ws, userId)
	return playerInfo
}

func (mp *ManagePlayer) PlayerClose(ws *websocket.Conn, userId int64) {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	playerInfo, ok := mp.Players[userId]
	if ok {
		//玩家下线
		if playerInfo.ws == ws {
			playerInfo.ws = nil
			playerInfo.exitTime = time.Now().Unix() + 300
			fmt.Println("websocket连接断开")
		}
	}
}

func (mp *ManagePlayer) Run() {
	ticker := time.NewTicker(time.Second * 30)
	for {
		// select {
		// case <- ticker.C:
		// 	mp.CheckPlayerOff()
		// }
		<-ticker.C
		mp.CheckPlayerOff()
	}
}

func (mp *ManagePlayer) CheckPlayerOff() {
	mp.lock.Lock()
	defer mp.lock.Unlock()
	for k, v := range mp.Players {
		if v.exitTime > time.Now().Unix() {
			fmt.Println("内存中清除玩家：", v.UserId)
			delete(mp.Players, k)
		}
	}
}

func (mp *ManagePlayer) SendPlayerMsg(ws *websocket.Conn, userId int64) {
	var msg MsgPlayer
	msg.MsgId = 1
	msg.UserId = userId
	msg.ModPlayer = *mp.Players[userId].GetModPlayer()
	msg.ModIcon = *mp.Players[userId].GetModIcon()
	msg.ModCard = *mp.Players[userId].GetModCard()
	msg.ModRole = *mp.Players[userId].GetModRole()
	msg.ModBag = *mp.Players[userId].GetModBag()
	msg.ModWeapon = *mp.Players[userId].GetModWeapon()
	msg.ModRelics = *mp.Players[userId].GetModRelics()
	msg.ModCook = *mp.Players[userId].GetModCook()
	msg.ModHome = *mp.Players[userId].GetModHome()

	str, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	ws.Write(str)
}
