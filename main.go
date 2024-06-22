package main

import (
	"fmt"
	"github.com/xxl6097/clink-go-nacos-lib/nacos"
	"github.com/xxl6097/go-glog/glog"
)

func main() {
	fmt.Println("hello main..")
	group := "clink_itest"
	host := "10.6.14.80"
	port := 8848
	inacos := nacos.NewNacos(group, host, uint64(port))
	glog.Debug(inacos)
}
