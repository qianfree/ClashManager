package service

import (
	"encoding/json"
	"fmt"
	"strings"

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
				missingNodes := []uint{}
				for _, id := range nodeIDs {
					if name, ok := nodeMap[id]; ok {
						pg.Proxies = append(pg.Proxies, name)
					} else {
						missingNodes = append(missingNodes, id)
					}
				}
				if len(missingNodes) > 0 {
					fmt.Printf("警告: 策略组 \"%s\" 引用了 %d 个不存在的节点: %v\n", g.Name, len(missingNodes), missingNodes)
				}
			}
		}

		// Parse Use JSON ["Provider1"]
		if g.Use != "" {
			var uList []string
			json.Unmarshal([]byte(g.Use), &uList)
			pg.Use = uList
		}

		// Add all groups, even empty ones (with DIRECT fallback for empty groups)
		if len(pg.Proxies) == 0 && len(pg.Use) == 0 {
			// Empty group - add DIRECT as fallback to avoid Clash error
			pg.Proxies = []string{"DIRECT"}
			fmt.Printf("警告: 策略组 \"%s\" 没有有效节点，已添加 DIRECT 作为后备\n", g.Name)
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
		var targetName string

		// Resolve target based on TargetType
		if r.TargetType == "node" {
			// Target may store node ID or node name
			var nodeID uint
			if _, err := fmt.Sscanf(r.Target, "%d", &nodeID); err == nil {
				// Target is a numeric ID, lookup by ID
				if name, ok := nodeMap[nodeID]; ok {
					targetName = name
				} else {
					// Node ID not found, skip this rule
					continue
				}
			} else {
				// Target is not a number, treat as node name directly
				// Check if the node exists
				found := false
				for _, name := range nodeMap {
					if name == r.Target {
						targetName = r.Target
						found = true
						break
					}
				}
				if !found {
					// Node not found by name, skip this rule
					continue
				}
			}
		} else if r.TargetType == "group" {
			// Target may store group ID or group name
			var groupID uint
			if _, err := fmt.Sscanf(r.Target, "%d", &groupID); err == nil {
				// Target is a numeric ID, lookup by ID
				if name, ok := groupMap[groupID]; ok {
					targetName = name
				} else {
					// Group ID not found, skip this rule
					continue
				}
			} else {
				// Target is not a number, treat as group name directly
				// Check if the group exists
				found := false
				for _, name := range groupMap {
					if name == r.Target {
						targetName = r.Target
						found = true
						break
					}
				}
				if !found {
					// Group not found by name, skip this rule
					continue
				}
			}
		} else {
			// builtin type means built-in target (DIRECT, PROXY, REJECT, etc.)
			targetName = r.Target
		}

		// Normalize payload for IP-CIDR type (ensure subnet mask is present)
		payload := r.Payload
		if r.Type == "IP-CIDR" && !strings.Contains(payload, "/") {
			// Single IP without subnet mask, add /32
			payload = payload + "/32"
		}

		// Format: TYPE,Payload,Target
		line := r.Type + "," + payload + "," + targetName
		if r.NoResolve {
			line += ",no-resolve"
		}
		ruleStrings = append(ruleStrings, line)
	}
	// Add Catch-All (默认直连，未匹配规则的流量不经过代理)
	ruleStrings = append(ruleStrings, "MATCH,Proxy")

	config.Rules = ruleStrings

	// 6. Marshal to YAML
	return yaml.Marshal(config)
}

// ValidationError represents a single validation error
type ValidationError struct {
	Type    string `json:"type"`    // error, warning
	Message string `json:"message"` // error message
	Field   string `json:"field"`   // field name (optional)
	Index   int    `json:"index"`   // index for arrays (optional)
}

// ValidationResult represents the result of config validation
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors"`
}

// ValidateConfig validates the generated Clash configuration
func (s *ConfigService) ValidateConfig() (*ValidationResult, error) {
	result := &ValidationResult{
		Valid:  true,
		Errors: []ValidationError{},
	}

	// 1. Fetch all data
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

	// 2. Validate nodes
	for i, node := range nodes {
		if node.Name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Type:    "error",
				Message: fmt.Sprintf("节点 #%d: 名称不能为空", i+1),
				Field:   "name",
				Index:   i,
			})
		}
		if node.Server == "" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Type:    "error",
				Message: fmt.Sprintf("节点 \"%s\": 服务器地址不能为空", node.Name),
				Field:   "server",
				Index:   i,
			})
		}
		if node.Port <= 0 || node.Port > 65535 {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Type:    "error",
				Message: fmt.Sprintf("节点 \"%s\": 端口无效 (%d)", node.Name, node.Port),
				Field:   "port",
				Index:   i,
			})
		}

		// Type-specific validation
		switch node.Type {
		case "ss", "shadowsocks":
			if node.Cipher == "" {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Type:    "error",
					Message: fmt.Sprintf("节点 \"%s\": Shadowsocks 需要加密方式", node.Name),
					Field:   "cipher",
					Index:   i,
				})
			}
			if node.Password == "" {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Type:    "error",
					Message: fmt.Sprintf("节点 \"%s\": Shadowsocks 需要密码", node.Name),
					Field:   "password",
					Index:   i,
				})
			}
		case "vmess", "vless":
			if node.UUID == "" {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Type:    "error",
					Message: fmt.Sprintf("节点 \"%s\": %s 需要 UUID", node.Name, strings.ToUpper(node.Type)),
					Field:   "uuid",
					Index:   i,
				})
			}
		case "trojan":
			if node.Password == "" {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Type:    "error",
					Message: fmt.Sprintf("节点 \"%s\": Trojan 需要密码", node.Name),
					Field:   "password",
					Index:   i,
				})
			}
		case "hysteria2":
			if node.Password == "" && node.UUID == "" {
				result.Valid = false
				result.Errors = append(result.Errors, ValidationError{
					Type:    "error",
					Message: fmt.Sprintf("节点 \"%s\": Hysteria2 需要密码", node.Name),
					Field:   "password",
					Index:   i,
				})
			}
		}
	}

	// 3. Validate rules
	supportedRuleTypes := map[string]bool{
		"DOMAIN": true, "DOMAIN-SUFFIX": true, "DOMAIN-KEYWORD": true,
		"IP-CIDR": true, "GEOIP": true, "SRC-IP-CIDR": true,
		"SRC-PORT": true, "DST-PORT": true, "PROCESS-NAME": true,
		"RULE-SET": true, "MATCH": true,
	}

	for i, rule := range rules {
		if rule.Type == "" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Type:    "error",
				Message: fmt.Sprintf("规则 #%d: 类型不能为空", i+1),
				Field:   "type",
				Index:   i,
			})
			continue
		}

		if !supportedRuleTypes[rule.Type] {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Type:    "error",
				Message: fmt.Sprintf("规则 #%d: 不支持的规则类型 \"%s\"", i+1, rule.Type),
				Field:   "type",
				Index:   i,
			})
		}

		if rule.Payload == "" && rule.Type != "MATCH" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Type:    "error",
				Message: fmt.Sprintf("规则 #%d: Payload 不能为空", i+1),
				Field:   "payload",
				Index:   i,
			})
		}

		// IP-CIDR validation
		if rule.Type == "IP-CIDR" && rule.Payload != "" {
			if !strings.Contains(rule.Payload, "/") {
				// Warning: will be auto-fixed to /32
				result.Errors = append(result.Errors, ValidationError{
					Type:    "warning",
					Message: fmt.Sprintf("规则 #%d: IP-CIDR \"%s\" 缺少子网掩码，将自动添加 /32", i+1, rule.Payload),
					Field:   "payload",
					Index:   i,
				})
			}
		}

		// Validate target exists
		targetExists := false
		if rule.TargetID > 0 {
			if rule.TargetType == "node" {
				for _, n := range nodes {
					if n.ID == rule.TargetID {
						targetExists = true
						break
					}
				}
			} else if rule.TargetType == "group" {
				for _, g := range groups {
					if g.ID == rule.TargetID {
						targetExists = true
						break
					}
				}
			}
		} else if rule.Target != "" {
			// Check if target is a built-in or valid name
			if rule.Target == "DIRECT" || rule.Target == "REJECT" {
				targetExists = true
			} else if rule.TargetType == "node" || rule.TargetType == "group" {
				// For node/group types, try to parse target as ID first (legacy data support)
				var id uint
				if _, err := fmt.Sscanf(rule.Target, "%d", &id); err == nil {
					if rule.TargetType == "node" {
						for _, n := range nodes {
							if n.ID == id {
								targetExists = true
								break
							}
						}
					} else if rule.TargetType == "group" {
						for _, g := range groups {
							if g.ID == id {
								targetExists = true
								break
							}
						}
					}
				} else {
					// Not an ID, check if it's a name
					if rule.TargetType == "group" {
						for _, g := range groups {
							if g.Name == rule.Target {
								targetExists = true
								break
							}
						}
					} else {
						// For node type, check by name
						for _, n := range nodes {
							if n.Name == rule.Target {
								targetExists = true
								break
							}
						}
					}
				}
			} else {
				// For builtin type or when targetType is not set, check groups by name
				for _, g := range groups {
					if g.Name == rule.Target {
						targetExists = true
						break
					}
				}
			}
		}

		if !targetExists && rule.Type != "MATCH" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Type:    "error",
				Message: fmt.Sprintf("规则 #%d: 目标 \"%s\" 不存在", i+1, getTargetName(rule)),
				Field:   "target",
				Index:   i,
			})
		}
	}

	// 4. Validate proxy groups
	for i, group := range groups {
		if group.Name == "" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Type:    "error",
				Message: fmt.Sprintf("策略组 #%d: 名称不能为空", i+1),
				Field:   "name",
				Index:   i,
			})
		}

		if group.Type != "select" && group.Type != "url-test" && group.Type != "fallback" && group.Type != "load-balance" {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Type:    "error",
				Message: fmt.Sprintf("策略组 \"%s\": 不支持的类型 \"%s\"", group.Name, group.Type),
				Field:   "type",
				Index:   i,
			})
		}

		if group.Type == "url-test" || group.Type == "fallback" {
			if group.URL == "" {
				result.Errors = append(result.Errors, ValidationError{
					Type:    "warning",
					Message: fmt.Sprintf("策略组 \"%s\": 未设置测速 URL，将使用默认值", group.Name),
					Field:   "url",
					Index:   i,
				})
			}
		}

		// Validate proxy IDs
		if group.ProxyIDs != "" {
			var nodeIDs []uint
			if err := json.Unmarshal([]byte(group.ProxyIDs), &nodeIDs); err == nil {
				nodeMap := make(map[uint]bool)
				for _, n := range nodes {
					nodeMap[n.ID] = true
				}
				for _, id := range nodeIDs {
					if !nodeMap[id] {
						result.Valid = false
						result.Errors = append(result.Errors, ValidationError{
							Type:    "error",
							Message: fmt.Sprintf("策略组 \"%s\": 引用的节点 ID %d 不存在", group.Name, id),
							Field:   "proxyIDs",
							Index:   i,
						})
					}
				}
			}
		}
	}

	// 5. General warnings
	if len(nodes) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Type:    "warning",
			Message: "没有配置任何节点",
			Field:   "nodes",
		})
	}
	if len(rules) == 0 {
		result.Errors = append(result.Errors, ValidationError{
			Type:    "warning",
			Message: "没有配置任何规则",
			Field:   "rules",
		})
	}

	return result, nil
}

// getTargetName returns the display name of a rule's target
func getTargetName(rule model.Rule) string {
	if rule.Target != "" {
		return rule.Target
	}
	if rule.TargetType == "node" {
		return fmt.Sprintf("节点 #%d", rule.TargetID)
	}
	if rule.TargetType == "group" {
		return fmt.Sprintf("策略组 #%d", rule.TargetID)
	}
	return "未知"
}
