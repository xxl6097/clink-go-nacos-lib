package iface

type INacos interface {
	GetConfig(dataid, group string) string
}
