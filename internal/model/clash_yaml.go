package model

// ClashConfig represents the root of the Clash configuration file
type ClashConfig struct {
	Port               int                      `yaml:"port"`
	SocksPort          int                      `yaml:"socks-port"`
	MixedPort          int                      `yaml:"mixed-port"`
	AllowLan           bool                     `yaml:"allow-lan"`
	Mode               string                   `yaml:"mode"`
	LogLevel           string                   `yaml:"log-level"`
	IPv6               bool                     `yaml:"ipv6"`
	ExternalController string                   `yaml:"external-controller"`
	Secret             string                   `yaml:"secret,omitempty"`
	DNS                DNSConfig                `yaml:"dns"`
	Proxies            []map[string]interface{} `yaml:"proxies"`
	ProxyProviders     map[string]interface{}   `yaml:"proxy-providers,omitempty"`
	ProxyGroups        []ProxyGroup             `yaml:"proxy-groups"`
	RuleProviders      map[string]interface{}   `yaml:"rule-providers,omitempty"`
	Rules              []string                 `yaml:"rules"`
}

// DNSConfig represents the DNS section
type DNSConfig struct {
	Enable            bool       `yaml:"enable" json:"enable"`
	Listen            string     `yaml:"listen" json:"listen"`
	EnhancedMode      string     `yaml:"enhanced-mode" json:"enhancedMode"`
	FakeIPFilter      []string   `yaml:"fake-ip-filter" json:"fakeIPFilter"`
	DefaultNameserver []string   `yaml:"default-nameserver" json:"defaultNameserver"`
	Nameserver        []string   `yaml:"nameserver" json:"nameserver"`
	Fallback          []string   `yaml:"fallback,omitempty" json:"fallback,omitempty"`
	FallbackFilter    ConfigFile `yaml:"fallback-filter,omitempty" json:"fallbackFilter,omitempty"`
}

// ConfigFile generic map for mixed content like fallback-filter
type ConfigFile map[string]interface{}

// ProxyGroup represents a proxy group strategy
type ProxyGroup struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Proxies   []string `yaml:"proxies"`
	URL       string   `yaml:"url,omitempty"`
	Interval  int      `yaml:"interval,omitempty"`
	Tolerance int      `yaml:"tolerance,omitempty"`
	Use       []string `yaml:"use,omitempty"`
}
