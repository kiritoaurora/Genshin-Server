package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

/*
	唯一任务模块
	1、世界等级突破任务
*/

type TaskInfo struct {
	TaskId int
	State  int
}

type ModUniqueTask struct {
	MyTaskInfo map[int]*TaskInfo

	player *Player
	path   string
}

func (ut *ModUniqueTask) IsTaskFinish(taskId int) bool {
	task, ok := ut.MyTaskInfo[taskId]
	if !ok {
		return false
	}
	return task.State == TASK_STATE_FINISH
}

func (ut *ModUniqueTask) SaveData() {
	content, err := json.Marshal(ut)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(ut.path, content, os.ModePerm)
	if err != nil {
		return
	}
}

func (ut *ModUniqueTask) LoadData(player *Player) {
	ut.player = player
	ut.path = ut.player.localPath + "/uniquetask.json"

	configFile, err := ioutil.ReadFile(ut.path)
	if err != nil {
		fmt.Println("暂无存档")
		// return
	}
	err = json.Unmarshal(configFile, &ut)
	if err != nil {
		ut.InitData()
		return
	}

	// if ut.MyTaskInfo == nil {
	// 	ut.MyTaskInfo = make(map[int]*TaskInfo)
	// }
}

func (ut *ModUniqueTask) InitData() {
	if ut.MyTaskInfo == nil {
		ut.MyTaskInfo = make(map[int]*TaskInfo)
	}
}
