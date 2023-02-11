package game

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/websocket"
)

var manageHttp *ManageHttp

type ManageHttp struct {}

func GetManageHttp() *ManageHttp {
	if manageHttp == nil {
		manageHttp = new(ManageHttp)
	}
	return manageHttp
}

func (bw *ManageHttp) InitData() {
	http.Handle("/", websocket.Handler(bw.WebsocketHandler))
	http.HandleFunc("/correctname", bw.CorrectName)
	http.HandleFunc("/correctsign", bw.CorrectSign)
	http.HandleFunc("/birthday", bw.Birthday)
	http.HandleFunc("/card", bw.SetCard)
	http.HandleFunc("/icon", bw.SetIcon)
	http.HandleFunc("/useitem", bw.UseItem)
	http.HandleFunc("/relicsup", bw.RelicsUp)
	http.HandleFunc("/weaponup", bw.WeaponUp)
	http.HandleFunc("/weaponupstar", bw.WeaponUpStar)
	http.HandleFunc("/weaponuprefine", bw.WeaponUpRefine)
}

func (bw *ManageHttp) CorrectName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.ParseInt(r.FormValue("userId"), 10, 64)
	name := r.FormValue("name")
	fmt.Println(userId, name)
	if player, ok := managePlayer.Players[userId]; ok {
		player.RecvSetName(name)
		newName := managePlayer.Players[userId].GetModPlayer().Name
		w.Write([]byte(newName))
	}
}

func (bw *ManageHttp) CorrectSign(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.ParseInt(r.FormValue("userId"), 10, 64)
	sign := r.FormValue("sign")
	fmt.Println(userId, sign)
	if player, ok := managePlayer.Players[userId]; ok {
		player.RecvSetSign(sign)
		newSign := managePlayer.Players[userId].GetModPlayer().Sign
		w.Write([]byte(newSign))
	}
}

func (bw *ManageHttp) Birthday(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.ParseInt(r.FormValue("userId"), 10, 64)
	birth, _ := strconv.Atoi(r.FormValue("birth"))
	fmt.Println(userId, birth)
	if player, ok := managePlayer.Players[userId]; ok {
		player.SetBirth(birth)
		birth = managePlayer.Players[userId].GetModPlayer().Birth
		birthday := strconv.Itoa(birth)
		w.Write([]byte(birthday))
	}
}

func (bw *ManageHttp) SetCard(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.ParseInt(r.FormValue("userId"), 10, 64)
	cardId, _ := strconv.Atoi(r.FormValue("cardId"))
	fmt.Println(userId, cardId)
	if player, ok := managePlayer.Players[userId]; ok {
		player.RecvSetCard(cardId)
		cardId = managePlayer.Players[userId].GetModPlayer().Card
		card := strconv.Itoa(cardId)
		w.Write([]byte(card))
	}
}

func (bw *ManageHttp) SetIcon(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.ParseInt(r.FormValue("userId"), 10, 64)
	iconId, _ := strconv.Atoi(r.FormValue("iconId"))
	fmt.Println(userId, iconId)
	if player, ok := managePlayer.Players[userId]; ok {
		player.RecvSetIcon(iconId)
		iconId = managePlayer.Players[userId].GetModPlayer().Icon
		icon := strconv.Itoa(iconId)
		w.Write([]byte(icon))
	}
}

func (bw *ManageHttp) UseItem(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.ParseInt(r.FormValue("userId"), 10, 64)
	itemId, _ := strconv.Atoi(r.FormValue("itemId"))
	num, _ := strconv.ParseInt(r.FormValue("num"), 10, 64)
	fmt.Println(userId, itemId, num)
	if player, ok := managePlayer.Players[userId]; ok {
		player.GetModBag().UseItem(itemId, num)
		nowNum := player.GetModBag().GetItemNum(itemId)
		str := strconv.FormatInt(nowNum, 10)
		w.Write([]byte(str))
	}
}

func (bw *ManageHttp) RelicsUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.ParseInt(r.FormValue("userId"), 10, 64)
	keyId, _ := strconv.Atoi(r.FormValue("keyId"))
	exp, _ := strconv.Atoi(r.FormValue("exp"))
	fmt.Println(userId, keyId, exp)
	if player, ok := managePlayer.Players[userId]; ok {
		player.GetModRelics().RelicsUp(keyId, exp)
		relics := player.GetModRelics().RelicsInfo[keyId]
		str, err := json.Marshal(relics)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		w.Write(str)
	}
}

func (bw *ManageHttp) WeaponUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.ParseInt(r.FormValue("userId"), 10, 64)
	keyId, _ := strconv.Atoi(r.FormValue("keyId"))
	exp, _ := strconv.Atoi(r.FormValue("exp"))
	fmt.Println(userId, keyId, exp)
	if player, ok := managePlayer.Players[userId]; ok {
		player.GetModWeapon().WeaponUp(keyId, exp)
		level := player.GetModWeapon().WeaponInfo[keyId].Level
		str := strconv.Itoa(level)
		w.Write([]byte(str))
	}
}

func (bw *ManageHttp) WeaponUpStar(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.ParseInt(r.FormValue("userId"), 10, 64)
	keyId, _ := strconv.Atoi(r.FormValue("keyId"))
	fmt.Println(userId, keyId)
	if player, ok := managePlayer.Players[userId]; ok {
		player.GetModWeapon().WeaponUpStar(keyId)
		level := player.GetModWeapon().WeaponInfo[keyId].StarLevel
		str := strconv.Itoa(level)
		w.Write([]byte(str))
	}
}

func (bw *ManageHttp) WeaponUpRefine(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, _ := strconv.ParseInt(r.FormValue("userId"), 10, 64)
	weaponKeyId, _ := strconv.Atoi(r.FormValue("weaponKeyId"))
	weaponTargetKeyId, _ := strconv.Atoi(r.FormValue("weaponTargetKeyId"))
	fmt.Println(userId, weaponKeyId, weaponTargetKeyId)
	if player, ok := managePlayer.Players[userId]; ok {
		player.GetModWeapon().WeaponUpRefine(weaponKeyId, weaponTargetKeyId)
		level := player.GetModWeapon().WeaponInfo[weaponKeyId].RefineLevel
		str := strconv.Itoa(level)
		w.Write([]byte(str))
	}
}

func (bw *ManageHttp) WebsocketHandler(ws *websocket.Conn) {
	defer ws.Close()

	var player *Player

	for {
		var msg []byte
		ws.SetReadDeadline(time.Now().Add(3 * time.Second))
		err := websocket.Message.Receive(ws, &msg)
		fmt.Println(err)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				continue
			}
			if player != nil {
				//存档
				for _, v := range player.modManage {
					v.SaveData()
				}
				GetManagePlayer().PlayerClose(ws, player.UserId)
			}
			break
		}
		fmt.Println(string(msg))

		if player == nil {
			var loginMsg MsgLogin
			msgErr := json.Unmarshal(msg, &loginMsg)
			if msgErr == nil {
				player = GetManagePlayer().PlayerLogin(ws, loginMsg.UserId)
				go player.LogicRun()
			}
		} else {
			player.SendLogic(msg)
		}
	}
}
