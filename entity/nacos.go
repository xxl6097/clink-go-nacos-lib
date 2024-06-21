package entity

// clink-common-redis-lettuce.yaml clink-common
// clink-common-rocketmq.yaml
type redis struct {
	DataBase int    `yaml:"database"`
	Host     string `yaml:"host"`
	Port     uint64 `yaml:"port"`
	Password string `yaml:"password"`
}

type rocketmq struct {
	NameServers string `yaml:"name-server"`
	Servers     []string
}

type spring struct {
	Redis    redis    `yaml:"redis"`
	Rocketmq rocketmq `yaml:"rocketmq"`
}

type SpringConfig struct {
	Spring spring `yaml:"spring"`
}
