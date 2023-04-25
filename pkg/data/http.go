package data

import (
	"errors"
	"github.com/golang/glog"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"sync"
	"template/pkg/util"
	"time"
)

var httpMgr *HttpMgr

// ErrHttpConfig 配置错误
var ErrHttpConfig = errors.New("http config error")

// ErrHttpUninitialized Http 还未初始化
var ErrHttpUninitialized = errors.New("http uninitialized")

// InitHttpMgr 初始化 Http
func InitHttpMgr() {
	httpMgr = newHttpMgr(viper.Sub("data.http"))
}

// UnInitHttpMgr 反初始化 Http 相关
func UnInitHttpMgr() {
	if httpMgr != nil {
		httpMgr.Close()
		httpMgr = nil
	}
}

// GetHttp 获取 Http 实例
func GetHttp(name string) (*http.Client, error) {
	if httpMgr == nil {
		panic(ErrHttpUninitialized)
	}

	return httpMgr.getHttp(name)
}

// MustGetHttp 获取 Http 实例，如果获取失败，直接报错
func MustGetHttp(name string) *http.Client {
	if httpMgr == nil {
		panic(ErrHttpUninitialized)
	}

	return httpMgr.mustGetHttp(name)
}

// newHttpMgr 根据配置创建新的 Http 连接管理
func newHttpMgr(conf *viper.Viper) *HttpMgr {
	httpMgr := &HttpMgr{
		httpMap:    make(map[string]*http.Client),
		mutex:      &sync.Mutex{},
		httpConfig: conf,
	}
	return httpMgr
}

// HttpMgr Http 连接管理
type HttpMgr struct {
	httpMap    map[string]*http.Client
	mutex      *sync.Mutex
	httpConfig *viper.Viper
}

// getHttp 根据名称获取 Http 连接
func (mgr *HttpMgr) getHttp(name string) (*http.Client, error) {
	config := mgr.httpConfig.Sub(name)
	if config == nil {
		return nil, ErrHttpConfig
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	http, ok := mgr.httpMap[name]
	if ok {
		return http, nil
	}

	http, err := initHttp(config, name)
	if err != nil {
		return nil, err
	}
	mgr.httpMap[name] = http
	return http, nil
}

// mustGetHttp 根据名称获取数据库连接
func (mgr *HttpMgr) mustGetHttp(name string) *http.Client {
	config := mgr.httpConfig.Sub(name)
	if config == nil {
		panic(ErrHttpConfig)
	}

	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	http, ok := mgr.httpMap[name]
	if ok {
		return http
	}

	http, err := initHttp(config, name)
	util.CheckError(err)

	mgr.httpMap[name] = http
	return http
}

// Close 关闭管理器，释放 Http 连接
func (mgr *HttpMgr) Close() {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	for _, http := range mgr.httpMap {
		http.CloseIdleConnections()
	}
	mgr.httpMap = make(map[string]*http.Client)
}

func initHttp(config *viper.Viper, name string) (*http.Client, error) {
	http := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns: config.GetInt("max-idle-conn"),
			//MaxIdleConnsPerHost: config.GetInt("max-idle-conn-per-host"),
			IdleConnTimeout:    90 * time.Second,
			DisableCompression: config.GetBool("disable-compression"),
		},
		Timeout: 30 * time.Second,
	}

	glog.Infof("%s http: maxIdleConn:%d, disableCompression: %t",
		name, config.GetInt("max-idle-conn"), config.GetBool("disable-compressed"))

	return http, nil
}
