package csvs

const (
	LOGIC_FALSE = 0
	LOGIC_TRUE  = 1

	PERCENT_ALL = 10000 //权重万分制
)

const (
	DROP_ITEM_TYPE_ITEM   = 1 //掉落类型：掉落物
	DROP_ITEM_TYPE_GROUP  = 2 //掉落类型：掉落组
	DROP_ITEM_TYPE_WEIGHT = 3 //权值掉落
)

const (
	EVENT_STATE_START  = 0  //事件开始
	EVENT_STATE_FINISH = 9  //事件完成
	EVENT_STATE_END    = 10 //事件领取

	EVENT_TYPE_NORMAL = 1 //事件类型：其他正常事件
	EVENT_TYPE_REWARD = 2 //事件类型：奖励
)

const (
	MAP_REFRESH_DAY  = 1 //日刷新
	MAP_REFRESH_WEEK = 2 //周刷新
	MAP_REFRESH_SELF = 3 //自刷新
	MAP_REFRESH_CANT = 4 //不可刷新

	MAP_REFRESH_DAY_TIME  = 60 //日刷新时间间隔
	MAP_REFRESH_WEEK_TIME = 40 //周刷新时间间隔
	MAP_REFRESH_SELF_TIME = 20 //自刷新时间间隔

	REFRESH_SYSTEM = 1 //系统刷新
	REFRESH_PLAYER = 2 //玩家刷新
)

const (
	REDUCE_WORLD_LEVEL_START         = 5  //降低世界等级的要求
	REDUCE_WORLD_LEVEL_MAX           = 1  //最多能降低的级数
	REDUCE_WORLD_LEVEL_COOL_TIME     = 10 //冷却时间
	SHOW_CARD_SIZE                   = 9  //展示名片上限
	SHOW_TEAM_SIZE                   = 8  //展示阵容上限
	ADD_ROLE_TIME_NORMAL_MIN         = 2
	ADD_ROLE_TIME_NORMAL_MAX         = 7    //角色1命到到6命区间
	WEAPON_MAX_COUNT                 = 2000 //武器背包最大数量
	RELICS_MAX_COUNT                 = 1500 //圣遗物背包最大数量
	FIVE_STAR_TIMES_LIMIT            = 73   //73抽之后每多一抽，抽到五星概率上涨6.0%
	FIVE_STAR_TIMES_LIMIT_EACH_VALUE = 600  //上涨6.0%  即:万分之600
	FOUR_STAR_TIMES_LIMIT            = 8    //8抽之后每多一抽，抽到五星概率上涨51.0%
	FOUR_STAR_TIMES_LIMIT_EACH_VALUE = 5100 //上涨6.0%  即:万分之5100
	ALL_ENTRY_RATE                   = 2000 //圣遗物出现四个副词条的概率：20%
	WEAPON_MAX_REFINE                = 5    //	武器精炼最大等级
)
