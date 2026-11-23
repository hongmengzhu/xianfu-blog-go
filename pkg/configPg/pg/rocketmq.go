package pg

// 消息队列
type Rocketmq struct {
	NameServer  string   `json:"name-serverPg" value:"${rocketmq:=127.0.0.1:9876}"`
	NameServers []string `json:"nameServers" value:"${rocketmq:=['127.0.0.1:9876']}"`
	//Producer    l.Producer `json:"producer"`
}

// 设置默认配置
func (c *Rocketmq) SetDefault() {
	//	if len(c.NameServer) < 1 {
	//		c.NameServer = "localhost:9876"
	//	}
	//
	//	if len(c.NameServers) < 1 {
	//		c.NameServers = append(c.NameServers, c.NameServer)
	//	}
	//
	//	if len(c.Producer.Group) < 1 {
	//		c.Producer.Group = "default"
	//	}
}
