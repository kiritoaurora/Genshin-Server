package handler

type MsgHead struct {
	MsgId int `json:"msgId"`
}

//登录信息 0
type MsgLogin struct {
	MsgId    int    `json:"msgId"`
	Account  string `json:"account"`
	Password string `json:"password"`
	UserId   int64  `json:"userId"`
}

//抽卡信息301十连， 302单抽
type MsgPool struct {
	MsgId    int `json:"msgId"`
	PoolType int `json:"pooltype"`
}

//抽卡回应报文303
type MsgResponsePool struct {
	// MsgId        int   `json:"msgId"`
	DropId       int   `json:"dropId"`
	Stuff        int   `json:"stuff"`
	StuffNum     int64 `json:"stuffNum"`
	StuffItem    int   `json:"stuffItem"`
	StuffItemNum int64 `json:"stuffItemNum"`
}
