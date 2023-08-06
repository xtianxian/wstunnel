package wstunnel

type ProxyItem struct {
	Listen string `yaml:"listen"`
	Remote string `yaml:"remote"`
}

type Conf struct {
	ProxyConfig []ProxyItem `yaml:"proxy_config"`
}
