package frame

import (
	"flag"

	"github.com/spf13/viper"
	"template/pkg/data"
	"template/pkg/util"
)

var (
	config string
)

// InitFramework 初始化框架
func InitFramework() {
	flag.StringVar(&config, "config", "./conf/application.yml", "")
	flag.Parse()

	viper.SetConfigFile(config)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	util.CheckError(err)

	data.InitHttpMgr()
	data.InitSQLMgr()
	data.InitRedisMgr()
	data.InitBscClient()
}

// UnInitFramework 反初始化框架
func UnInitFramework() {
	data.UnInitHttpMgr()
	data.UninitSQLMgr()
	data.UninitRedisMgr()
	data.UnInitBscClient()
}
