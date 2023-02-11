package game

type MsgHead struct {
	MsgId int `json:"msgId"`
}

//登陆 0
type MsgLogin struct {
	MsgId    int    `json:"msgId"`
	Account  string `json:"account"`
	Password string `json:"password"`
	UserId   int64  `json:"userId"`
}

//玩家基本信息 1
// type MsgPlayer struct {
// 	MsgId          int         `json:"msgId"`
// 	UserId         int64       `json:"userId"`         //唯一id
// 	Icon           int         `json:"icon"`           //头像
// 	Card           int         `json:"card"`           //名片
// 	Name           string      `json:"name"`           //昵称
// 	Sign           string      `json:"sign"`           //签名
// 	PlayerLevel    int         `json:"playerLevel"`    //玩家等级
// 	PlayerExp      int         `json:"PlayerExp"`      //冒险阅历
// 	WorldLevel     int         `json:"worldLevel"`     //世界等级
// 	WorldLevelNow  int         `json:"worldLevelNow"`  //当前世界等级
// 	WorldLevelCool int64       `json:"worldLevelCool"` //操作世界等级的冷却时间
// 	Birth          int         `json:"birth"`          //生日
// 	ShowTeam       []*ShowRole `json:"showTeam"`       //展示阵容
// 	HideShowTeam   int         `json:"hideShowTeam"`   //隐藏阵容开关
// 	ShowCard       []int       `json:"showCard"`       //展示名片
// }

//玩家基本信息 1
type MsgPlayer struct {
	MsgId     int       `json:"msgId"`
	UserId    int64     `json:"userId"`
	ModPlayer ModPlayer `json:"modPlayer"`
	ModIcon   ModIcon   `json:"modIcon"`
	ModCard   ModCard   `json:"modCard"`
	ModRole   ModRole   `json:"modRole"`
	ModBag    ModBag    `json:"modBag"`
	ModWeapon ModWeapon `json:"modWeapon"`
	ModRelics ModRelics `json:"modRelics"`
	ModCook   ModCook   `json:"modCook"`
	ModHome   ModHome   `json:"modHome"`
}

//背包信息 2
type MsgBag struct {
	MsgId   int               `json:"msgId"`
	BagInfo map[int]*ItemInfo `json:"bagInfo"`
}

//名片信息 3
type MsgCard struct {
	MsgId    int           `json:"msgId"`
	CardInfo map[int]*Card `json:"cardInfo"`
}

//头像信息 4
type MsgIcon struct {
	MsgId    int           `json:"msgId"`
	IconInfo map[int]*Icon `json:"iconInfo"`
}

//抽卡消息301(十连) 302(单抽)
type MsgPool struct {
	MsgId    int `json:"msgId"`
	PoolType int `json:"pooltype"`
}

//抽卡回应报文303
type MsgResponsePool struct {
	MsgId        int   `json:"msgId"`
	DropId       int   `json:"dropId"`
	Stuff        int   `json:"stuff"`
	StuffNum     int64 `json:"stuffNum"`
	StuffItem    int   `json:"stuffItem"`
	StuffItemNum int64 `json:"stuffItemNum"`
}
