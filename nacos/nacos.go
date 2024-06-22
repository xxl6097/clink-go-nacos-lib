package nacos

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/xxl6097/clink-go-nacos-lib/entity"
	"github.com/xxl6097/clink-go-nacos-lib/iface"
	"github.com/xxl6097/go-glog/glog"
	"gopkg.in/yaml.v3"
	"strings"
)

const (
	DATA_ID       = "clink-go-tcp-server.yaml"
	GROUP         = "clink-go-tcp-server"
	DATA_ID_REDIS = "clink-common-redis-lettuce.yaml"
	GROUP_REDIS   = "clink-common"
	DATA_ID_MQ    = "clink-common-rocketmq.yaml"
	GROUP_MQ      = "clink-common"
)

type Nacos struct {
	namespace    string
	host         string
	port         uint64
	clientConfig constant.ClientConfig
	Spring       *entity.SpringConfig
}

func NewNacos(_namespace, _host string, _port uint64) iface.INacos {
	this := &Nacos{namespace: _namespace, host: _host, port: _port}
	glog.Debug("macos init...")
	if this.namespace == "" {
		glog.Fatal("未配置 Nacos 配置中心命名空间")
	}
	// Nacos 客户端配置
	clientConfig := constant.ClientConfig{
		NamespaceId:         this.namespace, // 如果使用默认命名空间，填空字符串即可
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/tmp/nacos/log",
		CacheDir:            "/tmp/nacos/cache",
		LogLevel:            "debug",
	}
	this.clientConfig = clientConfig
	this.initConfig(clientConfig)
	return this
}

func Get[T any](dataid, group string, inacos iface.INacos) *T {
	if inacos == nil {
		return nil
	}
	context := inacos.GetConfig(dataid, group)
	if context == "" {
		return nil
	}
	var v T
	err := yaml.Unmarshal([]byte(context), &v)
	if err != nil {
		glog.Error(err)
	}
	return &v
}
func (this *Nacos) GetConfig(dataid, group string) string {
	serverConfigs1 := []constant.ServerConfig{
		{
			IpAddr: this.host, // Nacos 服务器 IP
			Port:   this.port, // Nacos 服务器端口
		},
	}
	configClient, err3 := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &this.clientConfig,
			ServerConfigs: serverConfigs1,
		},
	)

	if err3 != nil {
		glog.Fatal("Nacos NewConfigClient:", err3)
	}
	// 从Nacos获取配置
	content, err4 := configClient.GetConfig(vo.ConfigParam{
		DataId: dataid,
		Group:  group,
	})
	if err4 != nil {
		glog.Fatal("获取 Nacos 配置文件失败:", err4)
	}
	return content
}

func (this *Nacos) initConfig(clientConfig constant.ClientConfig) {
	//server := getConfig(DATA_ID, GROUP, clientConfig)
	//glog.Debug("Nacos 配置文件:", server)
	//glog.Debug("Nacos 配置文件:", redis)
	//glog.Debug("Nacos 配置文件:", rokmq)

	var _spring entity.SpringConfig
	err := yaml.Unmarshal([]byte(this.GetConfig(DATA_ID_REDIS, GROUP_REDIS)), &_spring)
	if err != nil {
		glog.Fatal(err)
	}
	err1 := yaml.Unmarshal([]byte(this.GetConfig(DATA_ID_MQ, GROUP_MQ)), &_spring)
	if err1 != nil {
		glog.Fatal(err1)
	} else {
		if _spring.Spring.Rocketmq.NameServers != "" {
			arr := strings.Split(_spring.Spring.Rocketmq.NameServers, ";")
			if len(arr) > 0 {
				servers := make([]string, 0)
				for _, v := range arr {
					if v == "" {
						continue
					}
					if strings.HasPrefix(v, "http") {
						servers = append(servers, v)
					} else {
						servers = append(servers, fmt.Sprintf("http://%s", v))
					}
				}
				_spring.Spring.Rocketmq.Servers = servers
			}
		}
	}
	glog.Info("---->", _spring)
	this.Spring = &_spring
}
