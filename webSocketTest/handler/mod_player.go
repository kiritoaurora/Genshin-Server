package handler

/*
	基础信息模块
*/

type ShowRole struct {
	RoleId    int
	RoleLevel int
}

//玩家基本信息 1
type ModPlayer struct {
	UserId         int64       //`json:"userId"`         //唯一id
	Icon           int         //`json:"icon"`           //头像
	Card           int         //`json:"card"`           //名片
	Name           string      //`json:"name"`           //昵称
	Sign           string      //`json:"sign"`           //签名
	PlayerLevel    int         //`json:"playerLevel"`    //玩家等级
	PlayerExp      int         //`json:"PlayerExp"`      //冒险阅历
	WorldLevel     int         //`json:"worldLevel"`     //世界等级
	WorldLevelNow  int         //`json:"worldLevelNow"`  //当前世界等级
	WorldLevelCool int64       //`json:"worldLevelCool"` //操作世界等级的冷却时间
	Birth          int         //`json:"birth"`          //生日
	ShowTeam       []*ShowRole //`json:"showTeam"`       //展示阵容
	HideShowTeam   int         //`json:"hideShowTeam"`   //隐藏阵容开关
	ShowCard       []int       //`json:"showCard"`       //展示名片
}
