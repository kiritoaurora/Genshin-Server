package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"server_logic/src/csvs"
)

/*
	卡池模块
*/

type Pool struct {
	PoolId        int
	FiveStarTimes int
	FourStarTimes int
	IsMustUp      int
}

type ModPool struct {
	UpPoolInfo *Pool

	player *Player
	path   string
}

func (p *ModPool) AddTimes() {
	p.UpPoolInfo.FiveStarTimes++
	p.UpPoolInfo.FourStarTimes++
}

func (p *ModPool) DoUpPool() {
	fiveStarNum := 0
	fourStarNum := 0

	result := make(map[int]int)         //统计抽中次数
	resultEach := make(map[int]int)     //统计第几抽抽到五星
	resultEachTest := make(map[int]int) //统计10连多黄的次数
	fiveTest := 0                       //十连中有几黄
	for i := 0; i < 10000000; i++ {
		p.AddTimes() //次数++
		if i%10 == 0 {
			fiveTest = 0
		}
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}
		//抽取五星的概率在73抽之后每一抽增加6.0%，抽取四星的概率在8抽之后每一抽增加51%
		if p.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || p.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGrop := new(csvs.DropGroup)
			newDropGrop.DropId = dropGroup.DropId
			newDropGrop.WeightAll = dropGroup.WeightAll

			//addFiveWeight = (当前抽卡次数 - 73) * 600
			addFiveWeight := (p.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			//addFiveWeight = (当前抽卡次数 - 8) * 5100
			addFourWeight := (p.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}
			//构建新的dropGroup
			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.DropId = config.DropId
				newConfig.Result = config.Result
				newConfig.IsEnd = config.IsEnd
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGrop.DropConfigs = append(newDropGrop.DropConfigs, newConfig)
			}

			dropGroup = newDropGrop
		}
		//抽卡
		roleIdConfig := csvs.GetRandDropNew(dropGroup)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result) //获取角色
			if roleConfig != nil {
				if roleConfig.Star == 5 { //五星角色
					fiveTest++
					resultEach[p.UpPoolInfo.FiveStarTimes]++
					p.UpPoolInfo.FiveStarTimes = 0 //重置次数
					fiveStarNum++
					if p.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE { //是否是大保底
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = csvs.GetRandDropNew(dropGroup)
							if roleIdConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					//抽中Up五星角色
					if roleIdConfig.DropId == 100012 {
						p.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						p.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}

				} else if roleConfig.Star == 4 { //四星角色
					p.UpPoolInfo.FourStarTimes = 0
					fourStarNum++
				}
			}
			result[roleIdConfig.Result]++
		}
		if i%10 == 9 {
			resultEachTest[fiveTest]++
		}
	}

	for k, v := range result {
		fmt.Printf("抽中%s次数:%d\n", csvs.GetItemConfig(k).ItemName, v)
	}
	fmt.Printf("抽中四星次数:%d\n", fourStarNum)
	fmt.Printf("抽中五星次数:%d\n", fiveStarNum)

	for k, v := range resultEach {
		fmt.Printf("第%d抽抽出5星的次数:%d\n", k, v)
	}

	for k, v := range resultEachTest {
		fmt.Printf("10连%d黄次数:%d\n", k, v)
	}
}

//十连
func (p *ModPool) HandleUpPoolTen() {
	for i := 0; i < 10; i++ {
		p.AddTimes() //次数++
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}
		//抽取五星的概率在73抽之后每一抽增加6.0%，抽取四星的概率在8抽之后每一抽增加51%
		if p.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || p.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGroup)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll
			//addFiveWeight = (当前抽卡次数 - 73) * 600
			addFiveWeight := (p.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			//addFourWeight = (当前抽卡次数 - 8) * 5100
			addFourWeight := (p.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}
			//构建新的dropGroup
			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.Result = config.Result
				newConfig.DropId = config.DropId
				newConfig.IsEnd = config.IsEnd
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
			}
			dropGroup = newDropGroup
		}
		// 抽卡
		roleIdConfig := csvs.GetRandDropNew(dropGroup)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
			if roleConfig != nil {
				if roleConfig.Star == 5 { //抽中五星
					p.UpPoolInfo.FiveStarTimes = 0
					// 是否是大保底
					if p.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = csvs.GetRandDropNew(dropGroup)
							if roleIdConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					//抽中Up五星
					if roleIdConfig.DropId == 100012 {
						p.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						p.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}
				} else if roleConfig.Star == 4 { //抽中四星
					p.UpPoolInfo.FourStarTimes = 0
				}
			}
			// fmt.Printf("第%d抽抽中:%s\n", i+1, csvs.GetItemConfig(roleIdConfig.Result).ItemName)
			p.player.GetModBag().AddItem(roleIdConfig.Result, 1)
		}
	}
	if p.UpPoolInfo.IsMustUp == csvs.LOGIC_FALSE {
		fmt.Println("当前处于小保底区间！")
	} else {
		fmt.Println("当前处于大保底区间！")
	}
	fmt.Printf("当前累计未出5星次数:%d\n", p.UpPoolInfo.FiveStarTimes)
	fmt.Printf("当前累计未出4星次数:%d\n", p.UpPoolInfo.FourStarTimes)
}

// 单抽
func (p *ModPool) HandleUpPoolSingle(times int) {
	if times <= 0 || times > 100000000 {
		fmt.Println("请输入正确的数值(1~100000000)")
		return
	} else {
		fmt.Printf("累计抽取%d次,结果如下:\n", times)
	}
	// result := make(map[int]int)
	// fourStarNum := 0
	// fiveStarNum := 0
	for i := 0; i < times; i++ {
		p.AddTimes() //次数++
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}
		//抽取五星的概率在73抽之后每一抽增加6.0%，抽取四星的概率在8抽之后每一抽增加51%
		if p.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || p.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGroup)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll
			//addFiveWeight = (当前抽卡次数 - 73) * 600
			addFiveWeight := (p.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			//addFiveWeight = (当前抽卡次数 - 8) * 5100
			addFourWeight := (p.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}
			//构建新的dropGroup
			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.Result = config.Result
				newConfig.DropId = config.DropId
				newConfig.IsEnd = config.IsEnd
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
			}
			dropGroup = newDropGroup
		}
		// 抽卡
		roleIdConfig := csvs.GetRandDropNew(dropGroup)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
			if roleConfig != nil {
				if roleConfig.Star == 5 { //抽中五星
					p.UpPoolInfo.FiveStarTimes = 0
					// fiveStarNum++
					// 是否是大保底
					if p.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = csvs.GetRandDropNew(dropGroup)
							if roleIdConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					//抽中Up五星
					if roleIdConfig.DropId == 100012 {
						p.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						p.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}
				} else if roleConfig.Star == 4 { //抽中四星
					p.UpPoolInfo.FourStarTimes = 0
					// fourStarNum++
				}
			}
			// result[roleIdConfig.Result]++
			// fmt.Printf("第%d抽抽中:%s\n", i+1, csvs.GetItemConfig(roleIdConfig.Result).ItemName)
			p.player.GetModBag().AddItem(roleIdConfig.Result, 1)
		}

	}
	// for k, v := range result {
	// 	fmt.Printf("抽中%s次数:%d\n", csvs.GetItemConfig(k).ItemName, v)
	// }
	// fmt.Printf("抽中4星:%d\n", fourStarNum)
	// fmt.Printf("抽中5星:%d\n", fiveStarNum)

	if p.UpPoolInfo.IsMustUp == csvs.LOGIC_FALSE {
		fmt.Println("当前处于小保底区间！")
	} else {
		fmt.Println("当前处于大保底区间！")
	}
	fmt.Printf("当前累计未出5星次数:%d\n", p.UpPoolInfo.FiveStarTimes)
	fmt.Printf("当前累计未出4星次数:%d\n", p.UpPoolInfo.FourStarTimes)
}

func (p *ModPool) HandleUpPoolSingleCheck1(times int) {
	if times <= 0 || times > 100000000 {
		fmt.Println("请输入正确的数值(1~100000000)")
		return
	} else {
		fmt.Printf("累计抽取%d次,结果如下:\n", times)
	}
	// result := make(map[int]int)
	// fourStarNum := 0
	// fiveStarNum := 0
	for i := 0; i < times; i++ {
		p.AddTimes() //次数++
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}
		//抽取五星的概率在73抽之后每一抽增加6.0%，抽取四星的概率在8抽之后每一抽增加51%
		if p.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || p.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGroup)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll
			//addFiveWeight = (当前抽卡次数 - 73) * 600
			addFiveWeight := (p.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			//addFiveWeight = (当前抽卡次数 - 8) * 5100
			addFourWeight := (p.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}
			//构建新的dropGroup
			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.Result = config.Result
				newConfig.DropId = config.DropId
				newConfig.IsEnd = config.IsEnd
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
			}
			dropGroup = newDropGroup
		}

		fiveStartInfo, fourStarInfo := p.player.GetModRole().GetRoleInfoForPoolCheck()

		// 抽卡
		roleIdConfig := csvs.GetRandDropNew2(dropGroup, fiveStartInfo, fourStarInfo)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
			if roleConfig != nil {
				if roleConfig.Star == 5 { //抽中五星
					p.UpPoolInfo.FiveStarTimes = 0
					// fiveStarNum++
					// 是否是大保底
					if p.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = csvs.GetRandDropNew(dropGroup)
							if roleIdConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					//抽中Up五星
					if roleIdConfig.DropId == 100012 {
						p.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						p.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}
				} else if roleConfig.Star == 4 { //抽中四星
					p.UpPoolInfo.FourStarTimes = 0
					// fourStarNum++
				}
			}
			// result[roleIdConfig.Result]++
			// fmt.Printf("第%d抽抽中:%s\n", i+1, csvs.GetItemConfig(roleIdConfig.Result).ItemName)
			p.player.GetModBag().AddItem(roleIdConfig.Result, 1)
		}

	}
	// for k, v := range result {
	// 	fmt.Printf("抽中%s次数:%d\n", csvs.GetItemConfig(k).ItemName, v)
	// }
	// fmt.Printf("抽中4星:%d\n", fourStarNum)
	// fmt.Printf("抽中5星:%d\n", fiveStarNum)

	if p.UpPoolInfo.IsMustUp == csvs.LOGIC_FALSE {
		fmt.Println("当前处于小保底区间！")
	} else {
		fmt.Println("当前处于大保底区间！")
	}
	fmt.Printf("当前累计未出5星次数:%d\n", p.UpPoolInfo.FiveStarTimes)
	fmt.Printf("当前累计未出4星次数:%d\n", p.UpPoolInfo.FourStarTimes)
}

func (p *ModPool) SaveData() {
	content, err := json.Marshal(p)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(p.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (p *ModPool) LoadData(player *Player) {
	p.player = player
	p.path = p.player.localPath + "/pool.json"

	configFile, err := ioutil.ReadFile(p.path)
	if err != nil {
		fmt.Println("暂无存档")
		// return
	}
	err = json.Unmarshal(configFile, &p)
	if err != nil {
		p.InitData()
		return
	}

	// if p.UpPoolInfo == nil {
	// 	p.UpPoolInfo = new(Pool)
	// }
}

func (p *ModPool) InitData() {
	if p.UpPoolInfo == nil {
		p.UpPoolInfo = new(Pool)
	}
}

//十连抽
func (p *ModPool) HandleUpPoolTenByMsg(msg []byte) {

	var msgInfo MsgPool
	msgErr := json.Unmarshal(msg, &msgInfo)
	if msgErr != nil {
		fmt.Println("消息解析失败！")
		return
	}

	// msgInfo.PoolType

	for i := 0; i < 10; i++ {
		p.AddTimes() //次数++
		dropGroup := csvs.ConfigDropGroupMap[1000]
		if dropGroup == nil {
			return
		}
		//抽取五星的概率在73抽之后每一抽增加6.0%，抽取四星的概率在8抽之后每一抽增加51%
		if p.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || p.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
			newDropGroup := new(csvs.DropGroup)
			newDropGroup.DropId = dropGroup.DropId
			newDropGroup.WeightAll = dropGroup.WeightAll
			//addFiveWeight = (当前抽卡次数 - 73) * 600
			addFiveWeight := (p.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
			if addFiveWeight < 0 {
				addFiveWeight = 0
			}
			//addFourWeight = (当前抽卡次数 - 8) * 5100
			addFourWeight := (p.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
			if addFourWeight < 0 {
				addFourWeight = 0
			}
			//构建新的dropGroup
			for _, config := range dropGroup.DropConfigs {
				newConfig := new(csvs.ConfigDrop)
				newConfig.Result = config.Result
				newConfig.DropId = config.DropId
				newConfig.IsEnd = config.IsEnd
				if config.Result == 10001 {
					newConfig.Weight = config.Weight + addFiveWeight
				} else if config.Result == 10002 {
					newConfig.Weight = config.Weight + addFourWeight
				} else if config.Result == 10003 {
					newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
				}
				newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
			}
			dropGroup = newDropGroup
		}
		// 抽卡
		roleIdConfig := csvs.GetRandDropNew(dropGroup)
		if roleIdConfig != nil {
			roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
			if roleConfig != nil {
				if roleConfig.Star == 5 { //抽中五星
					p.UpPoolInfo.FiveStarTimes = 0
					// 是否是大保底
					if p.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
						dropGroup := csvs.ConfigDropGroupMap[100012]
						if dropGroup != nil {
							roleIdConfig = csvs.GetRandDropNew(dropGroup)
							if roleIdConfig == nil {
								fmt.Println("数据异常")
								return
							}
						}
					}
					//抽中Up五星
					if roleIdConfig.DropId == 100012 {
						p.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
					} else {
						p.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
					}
				} else if roleConfig.Star == 4 { //抽中四星
					p.UpPoolInfo.FourStarTimes = 0
				}
			}  else {
				weaponConfig := csvs.GetWeaponConfig(roleIdConfig.Result)
				if weaponConfig.Star == 4 {
					p.UpPoolInfo.FourStarTimes = 0
				}
			}
			p.player.GetModBag().AddItem(roleIdConfig.Result, 1)	
		}
		//封装回应报文
		msgResponsePool := BuildMsg(roleIdConfig.Result)
		//发送抽卡结果
		str, errStr := json.Marshal(msgResponsePool)
		if errStr != nil {
			fmt.Println("errStr:", errStr)
			return
		}
		p.player.ws.Write(str)
	}
	if p.UpPoolInfo.IsMustUp == csvs.LOGIC_FALSE {
		fmt.Println("当前处于小保底区间！")
	} else {
		fmt.Println("当前处于大保底区间！")
	}
	fmt.Printf("当前累计未出5星次数:%d\n", p.UpPoolInfo.FiveStarTimes)
	fmt.Printf("当前累计未出4星次数:%d\n", p.UpPoolInfo.FourStarTimes)
}

// 单抽
func (p *ModPool) HandleUpPoolSingleByMsg(msg []byte) {
	
	var msgInfo MsgPool
	msgErr := json.Unmarshal(msg, &msgInfo)
	if msgErr != nil {
		fmt.Println("消息解析失败！")
		return
	}

	// msgInfo.PoolType

	p.AddTimes() //次数++
	dropGroup := csvs.ConfigDropGroupMap[1000]
	if dropGroup == nil {
		return
	}
	//抽取五星的概率在73抽之后每一抽增加6.0%，抽取四星的概率在8抽之后每一抽增加51%
	if p.UpPoolInfo.FiveStarTimes > csvs.FIVE_STAR_TIMES_LIMIT || p.UpPoolInfo.FourStarTimes > csvs.FOUR_STAR_TIMES_LIMIT {
		newDropGroup := new(csvs.DropGroup)
		newDropGroup.DropId = dropGroup.DropId
		newDropGroup.WeightAll = dropGroup.WeightAll
		//addFiveWeight = (当前抽卡次数 - 73) * 600
		addFiveWeight := (p.UpPoolInfo.FiveStarTimes - csvs.FIVE_STAR_TIMES_LIMIT) * csvs.FIVE_STAR_TIMES_LIMIT_EACH_VALUE
		if addFiveWeight < 0 {
			addFiveWeight = 0
		}
		//addFourWeight = (当前抽卡次数 - 8) * 5100
		addFourWeight := (p.UpPoolInfo.FourStarTimes - csvs.FOUR_STAR_TIMES_LIMIT) * csvs.FOUR_STAR_TIMES_LIMIT_EACH_VALUE
		if addFourWeight < 0 {
			addFourWeight = 0
		}
		//构建新的dropGroup
		for _, config := range dropGroup.DropConfigs {
			newConfig := new(csvs.ConfigDrop)
			newConfig.Result = config.Result
			newConfig.DropId = config.DropId
			newConfig.IsEnd = config.IsEnd
			if config.Result == 10001 {
				newConfig.Weight = config.Weight + addFiveWeight
			} else if config.Result == 10002 {
				newConfig.Weight = config.Weight + addFourWeight
			} else if config.Result == 10003 {
				newConfig.Weight = config.Weight - addFiveWeight - addFourWeight
			}
			newDropGroup.DropConfigs = append(newDropGroup.DropConfigs, newConfig)
		}
		dropGroup = newDropGroup
	}
	// 抽卡
	roleIdConfig := csvs.GetRandDropNew(dropGroup)
	if roleIdConfig != nil {
		roleConfig := csvs.GetRoleConfig(roleIdConfig.Result)
		if roleConfig != nil {
			if roleConfig.Star == 5 { //抽中五星
				p.UpPoolInfo.FiveStarTimes = 0
				// 是否是大保底
				if p.UpPoolInfo.IsMustUp == csvs.LOGIC_TRUE {
					dropGroup := csvs.ConfigDropGroupMap[100012]
					if dropGroup != nil {
						roleIdConfig = csvs.GetRandDropNew(dropGroup)
						if roleIdConfig == nil {
							fmt.Println("数据异常")
							return
						}
					}
				}
				//抽中Up五星
				if roleIdConfig.DropId == 100012 {
					p.UpPoolInfo.IsMustUp = csvs.LOGIC_FALSE
				} else {
					p.UpPoolInfo.IsMustUp = csvs.LOGIC_TRUE
				}
			} else if roleConfig.Star == 4 { //抽中四星
				p.UpPoolInfo.FourStarTimes = 0
			}
		} else {
			weaponConfig := csvs.GetWeaponConfig(roleIdConfig.Result)
			if weaponConfig.Star == 4 {
				p.UpPoolInfo.FourStarTimes = 0
			}
		}
		// fmt.Printf("第%d抽抽中:%s\n", i+1, csvs.GetItemConfig(roleIdConfig.Result).ItemName)
		p.player.GetModBag().AddItem(roleIdConfig.Result, 1)
	}
	//封装回应报文
	msgResponsePool := BuildMsg(roleIdConfig.Result)
	//发送抽卡结果
	str, errStr := json.Marshal(msgResponsePool)
	if errStr != nil {
		fmt.Println("errStr:", errStr)
		return
	}
	p.player.ws.Write(str)

	if p.UpPoolInfo.IsMustUp == csvs.LOGIC_FALSE {
		fmt.Println("当前处于小保底区间！")
	} else {
		fmt.Println("当前处于大保底区间！")
	}
	fmt.Printf("当前累计未出5星次数:%d\n", p.UpPoolInfo.FiveStarTimes)
	fmt.Printf("当前累计未出4星次数:%d\n", p.UpPoolInfo.FourStarTimes)
}

//封装回应报文
func BuildMsg(dropId int) MsgResponsePool {
	var msgResponsePool MsgResponsePool
	msgResponsePool.MsgId = 303
	msgResponsePool.DropId = dropId
	config := csvs.GetRoleConfig(dropId)
	if config != nil {
		if player.GetModRole().RoleInfo[dropId].GetTimes >= csvs.ADD_ROLE_TIME_NORMAL_MIN &&
		player.GetModRole().RoleInfo[dropId].GetTimes <= csvs.ADD_ROLE_TIME_NORMAL_MAX {
			msgResponsePool.Stuff = config.Stuff
			msgResponsePool.StuffNum = config.StuffNum
			msgResponsePool.StuffItem = config.StuffItem
			msgResponsePool.StuffItemNum = config.StuffItemNum
		} else if player.GetModRole().RoleInfo[dropId].GetTimes > csvs.ADD_ROLE_TIME_NORMAL_MAX {
			msgResponsePool.StuffItem = config.MaxStuffItem
			msgResponsePool.StuffItemNum = config.MaxStuffItemNum
		} else {
			msgResponsePool.Stuff = 0
			msgResponsePool.StuffNum = 0
			msgResponsePool.StuffItem = 0
			msgResponsePool.StuffItemNum = 0
		}
	} else {
		msgResponsePool.Stuff = 0
		msgResponsePool.StuffNum = 0
		msgResponsePool.StuffItem = 0
		msgResponsePool.StuffItemNum = 0
	}
	return msgResponsePool
}