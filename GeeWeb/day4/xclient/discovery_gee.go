package xclient

import "time"

type GeeRegistryDiscovery struct {
	*MultiServersDiscovery
	registry   string // 注册中心地
	timeout    time.Duration
	lastUpdate time.Time // 更新服务列表的时间
}

const defaultUpdateTimeout = time.Second * 10

func NewGeeRegistryDiscovery(registryAddr string, timeout time.Duration) *GeeRegistryDiscovery {
	if timeout == 0 {
		timeout = defaultUpdateTimeout
	}
	d := &GeeRegistryDiscovery{
		MultiServersDiscovery: NewMultiServerDiscovery(make([]string, 0)),
		registry:              registryAddr,
		timeout:               timeout,
	}
	return d
}
