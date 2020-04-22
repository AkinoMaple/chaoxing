package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/remeh/sizedwaitgroup"
)

var mod Model
var config *Config

func addUser(puid int) bool {
	userInfo := getUserInfo(puid)
	if userInfo == nil {
		log.Printf("[%d]返回为空\n", puid)
		return false
	}

	user, r := parseRes(userInfo, puid)
	if r == 0 {
		log.Printf("[%d]用户信息异常\n", puid)
		return false
	}

	res, err := mod.UpsertOneUser(user)
	if err != nil {
		log.Print(err)
		return false
	}

	log.Printf("[%d]用户数据添加成功 %+v\n", puid, res)
	return true
}

func setConfig() {
	wd, _ := os.Getwd()
	cfp := flag.String("config", wd+"/config.toml", "Config file's path.")
	flag.Parse()

	config = unmarshal(loadConfig(cfp))
}

func main() {
	setConfig()

	mod.Init()

	rand.Seed(time.Now().UnixNano())

	swg := sizedwaitgroup.New(config.Settings.Routines)
	for _, work := range config.Settings.FetchID {
		for i := work[0]; i <= work[1]; i += work[2] {
			swg.Add()
			go func(i int) {
				defer swg.Done()
				addUser(i)
			}(i)
		}
	}

	swg.Wait()
}
