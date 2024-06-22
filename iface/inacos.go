package iface

import "github.com/xxl6097/clink-go-nacos-lib/entity"

type INacos interface {
	GetConfig(dataid, group string) string
	GetSpring() *entity.SpringConfig
}
