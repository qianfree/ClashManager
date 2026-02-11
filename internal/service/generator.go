package service

import (
	"encoding/json"

	"clash-manager/internal/model"
	"clash-manager/internal/repository"

	"gopkg.in/yaml.v3"
)

type ConfigService struct {
	NodeRepo     *repository.NodeRepository
	RuleRepo     *repository.RuleRepository
	GroupRepo    *repository.GroupRepository
	SettingsRepo *repository.SettingsRepository
}

func NewConfigService() *ConfigService {
	return &ConfigService{
		NodeRepo:     &repository.NodeRepository{},
		RuleRepo:     &repository.RuleRepository{},
		GroupRepo:    &repository.GroupRepository{},
		SettingsRepo: &repository.SettingsRepository{},
	}
}

func (s *ConfigService) GenerateConfig() ([]byte, error) {
	// 1. Fetch Data
	nodes, err := s.NodeRepo.FindAll()
	if err != nil {
		return nil, err
	}
	rules, err := s.RuleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	groups, err := s.GroupRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// Build ID Maps
	nodeMap := make(map[uint]string)
	for _, n := range nodes {
		nodeMap[n.ID] = n.Name
	}
	groupMap := make(map[uint]string)
	for _, g := range groups {
		groupMap[g.ID] = g.Name
	}

	// 1. Basic Config
	config := model.ClashConfig{
		Port:               7890,
		SocksPort:          7891,
		MixedPort:          7893,
		AllowLan:           true,
		Mode:               "rule",
		LogLevel:           "info",
		IPv6:               false,
		ExternalController: "0.0.0.0:9090",
	}

	// 2. DNS Config (Fetch from DB)
	dnsVal, err := s.SettingsRepo.Get("dns_config")
	if err == nil && dnsVal != "" {
		var dbDNS model.DNSConfig
		if err := json.Unmarshal([]byte(dnsVal), &dbDNS); err == nil {
			config.DNS = dbDNS
		}
	} else {
		// Default if not set
		config.DNS = model.DNSConfig{
			Enable:            true,
			Listen:            "0.0.0.0:53",
			EnhancedMode:      "fake-ip",
			Nameserver:        []string{"223.5.5.5", "119.29.29.29"},
			Fallback:          []string{"8.8.8.8", "1.1.1.1"},
			DefaultNameserver: []string{"223.5.5.5", "119.29.29.29"},
		}
	}

	// 3. Convert Nodes to Proxies List
	var proxies []map[string]interface{}
	var proxyNames []string
	for _, n := range nodes {
		proxy := make(map[string]interface{})
		proxy["name"] = n.Name
		proxy["type"] = n.Type
		proxy["server"] = n.Server
		proxy["port"] = n.Port

		if n.Password != "" {
			proxy["password"] = n.Password
		}
		if n.UUID != "" {
			proxy["uuid"] = n.UUID
		}
		if n.Cipher != "" {
			proxy["cipher"] = n.Cipher
		}
		if n.UDP {
			proxy["udp"] = true
		}
		if n.TLS {
			proxy["tls"] = true
			if n.SkipCert {
				proxy["skip-cert-verify"] = true
			}
			// Automatic SNI (servername) handling
			if n.Host != "" {
				proxy["servername"] = n.Host
			}
		}

		// VMess Specific: ensure alterId is present (default 0)
		if n.Type == "vmess" {
			proxy["alterId"] = 0
		}

		if n.Network != "" {
			proxy["network"] = n.Network

			// Transport Options
			switch n.Network {
			case "ws":
				wsOpts := make(map[string]interface{})
				if n.Path != "" {
					wsOpts["path"] = n.Path
				}
				if n.Host != "" {
					if wsOpts["headers"] == nil {
						wsOpts["headers"] = make(map[string]interface{})
					}
					wsOpts["headers"].(map[string]interface{})["Host"] = n.Host
				}
				if len(wsOpts) > 0 {
					proxy["ws-opts"] = wsOpts
				}
			case "grpc":
				grpcOpts := make(map[string]interface{})
				if n.Path != "" {
					grpcOpts["serviceName"] = n.Path // gRPC usually uses serviceName
				}
				if len(grpcOpts) > 0 {
					proxy["grpc-opts"] = grpcOpts
				}
			case "h2":
				h2Opts := make(map[string]interface{})
				if n.Path != "" {
					h2Opts["path"] = []string{n.Path}
				}
				if n.Host != "" {
					h2Opts["host"] = []string{n.Host}
				}
				if len(h2Opts) > 0 {
					proxy["h2-opts"] = h2Opts
				}
			}
		}

		if n.ExtraConfig != "" {
			var extra map[string]interface{}
			if err := json.Unmarshal([]byte(n.ExtraConfig), &extra); err == nil {
				for k, v := range extra {
					proxy[k] = v
				}
			}
		}

		proxies = append(proxies, proxy)
		proxyNames = append(proxyNames, n.Name)
	}
	config.Proxies = proxies

	// 4. Build Proxy Groups
	var proxyGroups []model.ProxyGroup

	// Default Auto Group
	autoGroup := model.ProxyGroup{
		Name:     "Auto Select",
		Type:     "url-test",
		URL:      "http://www.gstatic.com/generate_204",
		Interval: 300,
		Proxies:  proxyNames,
	}
	if len(proxyNames) == 0 {
		autoGroup.Proxies = []string{"DIRECT"}
	}
	proxyGroups = append(proxyGroups, autoGroup)

	// Helper for new proxy item format
	type ProxyItem struct {
		ID   uint   `json:"id"`
		Type string `json:"type"` // node, group
	}

	// User Defined Groups
	for _, g := range groups {
		pg := model.ProxyGroup{
			Name:     g.Name,
			Type:     g.Type,
			URL:      g.URL,
			Interval: g.Interval,
		}

		// Parse ProxyIDs JSON array: [1, 2, 3]
		if g.ProxyIDs != "" {
			var nodeIDs []uint
			if err := json.Unmarshal([]byte(g.ProxyIDs), &nodeIDs); err == nil {
				for _, id := range nodeIDs {
					if name, ok := nodeMap[id]; ok {
						pg.Proxies = append(pg.Proxies, name)
					}
				}
			}
		}

		// Parse Use JSON ["Provider1"]
		if g.Use != "" {
			var uList []string
			json.Unmarshal([]byte(g.Use), &uList)
			pg.Use = uList
		}
		proxyGroups = append(proxyGroups, pg)
	}

	// Final Proxy Group (Catch-all)
	finalGroup := model.ProxyGroup{
		Name:    "Proxy",
		Type:    "select",
		Proxies: append([]string{"Auto Select"}, proxyNames...),
	}
	proxyGroups = append(proxyGroups, finalGroup)

	config.ProxyGroups = proxyGroups

	// 5. Build Rules
	var ruleStrings []string
	for _, r := range rules {
		targetName := r.Target

		// Try ID resolution first
		if r.TargetID > 0 {
			var resolved string
			var found bool
			if r.TargetType == "node" {
				resolved, found = nodeMap[r.TargetID]
			} else if r.TargetType == "group" {
				resolved, found = groupMap[r.TargetID]
			}

			if found {
				targetName = resolved
			}
			// If not found, targetName remains r.Target (legacy fallback)
		}

		// Format: TYPE,Payload,Target
		line := r.Type + "," + r.Payload + "," + targetName
		if r.NoResolve {
			line += ",no-resolve"
		}
		ruleStrings = append(ruleStrings, line)
	}
	// Add Catch-All
	ruleStrings = append(ruleStrings, "MATCH,Proxy")

	config.Rules = ruleStrings

	// 6. Marshal to YAML
	return yaml.Marshal(config)
}
