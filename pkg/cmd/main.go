package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"net/http"
	"template/pkg/api"
	"template/pkg/frame"
	"template/pkg/middleware"
	"template/pkg/services/scan"
)

type Option struct {
}

func main() {
	defer glog.Flush()

	frame.InitFramework()
	defer frame.UnInitFramework()

	service := NewService(Option{})
	service.Init()
	defer service.UnInit()

	debug := viper.GetBool("application.debug")
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())
	//r.Use(static.Serve("/", static.LocalFile("./bin/dist", true)))

	r.Use(middleware.Cors())
	r.Use(middleware.LoggerToFile())

	api.Setup(r)
	addr := viper.GetString("server.addr")

	if len(addr) != 0 {
		err := r.Run(addr)
		if err != nil {
			panic(err.Error())
		}
	} else {
		err := r.Run()
		if err != nil {
			panic(err.Error())
		}
	}

}

func recoverHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, err)
}

type Service struct {
	scan       *scan.Engine
	ctx        context.Context
	shutdownFn context.CancelFunc
}

func NewService(opt Option) *Service {
	rootCtx, shutdownFn := context.WithCancel(context.Background())

	scan, err := scan.NewEngine(rootCtx)
	if err != nil {
		panic(err)
	}

	service := &Service{
		scan:       scan,
		ctx:        rootCtx,
		shutdownFn: shutdownFn,
	}

	return service
}

func (service *Service) Init() {
	go service.scan.Run(service.ctx)
	//go service.scan.Run2(service.ctx)
	//go SubscribeEvent()

}

func (service *Service) UnInit() {
	service.shutdownFn()
}

//func SubscribeEvent() {
//	client, err := ethclient.Dial("wss://snowy-alien-liquid.bsc.discover.quiknode.pro/ad3f1d96711a722a587e4f52b4fdef9eefb2e1c3/")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// 0x Protocol (ZRX) token address
//	contractAddress := common.HexToAddress("0x17d70FF5670742d9bf7609AB3d1C790210ace997")
//	query := ethereum.FilterQuery{
//		Addresses: []common.Address{contractAddress},
//	}
//
//	logs := make(chan types.Log)
//	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for {
//		select {
//		case err := <-sub.Err():
//			log.Fatal(err)
//		case vLog := <-logs:
//			fmt.Println(vLog) // pointer to event log
//		}
//	}
//
//}
