package game

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"server_logic/src/csvs"
	"sync"
	"syscall"
	"time"
)

type DBConfig struct {
	DBUser     string `json:"dbuser"`
	DBPassword string `json:"dbpassword"`
}

type ServerConfig struct {
	ServerId      int       `json:"serverid"`
	Host          string    `json:"host"`
	LocalSavePath string    `json:"localsavepath"`
	DBConfig      *DBConfig `json:"database"`
}

type Server struct {
	Wait         sync.WaitGroup
	BanWordBase  []string //配置生成
	BanWordExtra []string //更新
	Lock         sync.RWMutex
	Config       *ServerConfig
}

var server *Server

func GetServer() *Server {
	if server == nil {
		server = new(Server)
	}
	return server
}

func (s *Server) Start() {
	//加载服务器配置
	s.LoadConfig()
	// 加载模块
	rand.Seed(time.Now().Unix())
	csvs.CheckLoadCsv()
	go GetManageBanWord().Run()
	go GetManageHttp().InitData()
	go GetManagePlayer().Run()

	fmt.Printf("数据测试----start\n")
	// playerTest := NewTestPlayer(nil,10000086)
	// go playerTest.Run()
	go s.HandleSignal()

	go http.ListenAndServe(GetServer().Config.Host, nil)

	s.Wait.Wait()
	fmt.Println("服务器关闭成功")
}

func (s *Server) Stop() {
	GetManageBanWord().Close()
}

func (s *Server) AddGo() {
	s.Wait.Add(1)
}

func (s *Server) GoDone() {
	s.Wait.Done()
}

func (s *Server) IsBanWord(txt string) bool {
	s.Lock.RLock()
	defer s.Lock.RUnlock()
	for _, v := range s.BanWordBase {
		match, _ := regexp.MatchString(v, txt)
		if match {
			fmt.Println("发现违禁词:", v)
			return match
		}
	}

	for _, v := range s.BanWordExtra {
		match, _ := regexp.MatchString(v, txt)
		if match {
			fmt.Println("发现违禁词:", v)
			return match
		}
	}

	return false
}

func (s *Server) UpdateBanWord(banBaseWord, banExtraWord []string) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.BanWordBase = banBaseWord
	s.BanWordExtra = banExtraWord
}

func (s *Server) HandleSignal() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)
	<-signalChan
	fmt.Println("get syscall.SIGINT")
	fmt.Println("get syscall.SIGINT")
	fmt.Println("get syscall.SIGINT")
	s.Stop()
}

func (s *Server) LoadConfig() {
	configFile, err := ioutil.ReadFile("./config.json")
	if err != nil {
		fmt.Println("读取错误")
		return
	}
	err = json.Unmarshal(configFile, &s.Config)
	if err != nil {
		fmt.Println("反序列化错误")
		return
	}
	// fmt.Println(s.Config.ServerId,s.Config.Host,s.Config.LocalSavePath,
	// 	s.Config.DBConfig.DBUser,s.Config.DBConfig.DBPassword)
}
