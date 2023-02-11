package main

import (
	"server_logic/src/game"
)

func main() {

	// 基础模块
	// 1 UID
	// 2 头像、名片
	// 3 签名
	// 4 名字
	// 5 冒险等阶、冒险阅历
	// 6 世界等级、冷却时间
	// 7 生日
	// 8 展示阵容、展示名片

	// 背包系统
	// 1、物品识别
	// 2、物品增加
	// 3、物品消耗
	// 4、物品使用
	// 5、角色模块->头像模块

	// 掉落模块
	// 1、保底设计
	// 2、大数据测试 （策划）
	// 3、更新测试工具
	// 4、UP池
	// 5、仓检

	// 地图模块
	// 1 蒙德地图
	// 2 地图上的数据结构：1、采集物（矿物） 2、怪物（低级怪物，北风狼王） 3、传送点 4、七天神像，神瞳 5、宝箱
	// 3 秘境地图：1、圣遗物秘境，为后面的圣遗物模块做支持，2、风魔龙

	// 圣遗物模块
	// 1 生成
	// 2 属性
	// 3 强化
	// 4 装备	(替换)
	// 5 卸下
	// 6 套装
	// 7 洗练(用于测试)

	// 武器模块
	// 1 强化
	// 2 突破
	// 3 精炼
	// 4 装备与卸下
	// 5 替换

	/***************************************************************/
	// 加载配置

	game.GetServer().Start()

	return
}