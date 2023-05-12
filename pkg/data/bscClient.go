package data

import (
	"errors"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"sync"
)

var ErrBscNodeConfig = errors.New("bsc node config error")

var ErrConnect = errors.New("connect to bsc error")

var ErrBscClientUninitialized = errors.New("bsc client uninitialized")

type BscClient struct {
	bscClient     map[string]*ethclient.Client
	mutex         *sync.Mutex
	bscNodeConfig *viper.Viper
}

var bscClient *BscClient

func NewBscClient(conf *viper.Viper) *BscClient {
	bscClient := &BscClient{
		bscClient:     make(map[string]*ethclient.Client),
		mutex:         &sync.Mutex{},
		bscNodeConfig: conf,
	}
	return bscClient
}

func InitBscClient() {
	bscClient = NewBscClient(viper.Sub("data.bsc"))
}

func UnInitBscClient() {
	if bscClient != nil {
		bscClient.close()
		bscClient = nil
	}

}

func GetBscClient(name string) (*ethclient.Client, error) {
	if bscClient == nil {
		panic(ErrBscClientUninitialized)
	}
	return bscClient.getBscClient(name)
}

func MustGetBscClient(name string) *ethclient.Client {
	if bscClient == nil {
		panic(ErrBscClientUninitialized)
	}
	return bscClient.mustGetBscClient(name)
}

func (bsc *BscClient) getBscClient(name string) (*ethclient.Client, error) {
	config := bsc.bscNodeConfig.Sub(name)
	if config == nil {
		return nil, ErrBscNodeConfig
	}
	bsc.mutex.Lock()
	defer bsc.mutex.Unlock()
	bscClient, ok := bsc.bscClient[name]
	if ok {
		return bscClient, nil
	}
	bscClient, err := initBscClient(config, name)
	if err != nil {
		return nil, err
	}
	bsc.bscClient[name] = bscClient
	return bscClient, nil

}

func (bsc *BscClient) mustGetBscClient(name string) *ethclient.Client {
	config := bsc.bscNodeConfig.Sub(name)
	if config == nil {
		panic(ErrBscNodeConfig)
	}
	bsc.mutex.Lock()
	defer bsc.mutex.Unlock()
	bscClient, ok := bsc.bscClient[name]
	if ok {
		return bscClient
	}
	bscClient, err := initBscClient(config, name)
	if err != nil {
		panic(err)
	}
	bsc.bscClient[name] = bscClient
	return bscClient

}

func (bsc *BscClient) close() {
	bsc.mutex.Lock()
	defer bsc.mutex.Unlock()
	for _, client := range bsc.bscClient {
		client.Close()
	}
	bsc.bscClient = make(map[string]*ethclient.Client)
}

func initBscClient(config *viper.Viper, name string) (*ethclient.Client, error) {
	nodes := config.GetStringSlice("node")
	for i := 0; i < len(nodes); i++ {
		client, err := ethclient.Dial(nodes[i])
		if err == nil {
			return client, nil
		}
		glog.Errorf("无法链接节点: %s", nodes[i])
		if i == len(nodes)-1 {
			glog.Errorf("连接bsc节点失败")
		}

	}
	return nil, ErrConnect

}
