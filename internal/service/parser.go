package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"clash-manager/internal/model"
)

// ParseLink parses proxy links like ss://, vmess://, trojan:// into a Node model
func ParseLink(link string) (*model.Node, error) {
	link = strings.TrimSpace(link)
	if strings.HasPrefix(link, "ss://") {
		return parseSS(link)
	}
	if strings.HasPrefix(link, "vmess://") {
		return parseVmess(link)
	}
	if strings.HasPrefix(link, "trojan://") {
		return parseTrojan(link)
	}
	return nil, errors.New("unsupported protocol")
}

func parseSS(link string) (*model.Node, error) {
	// ss://<base64(method:password@server:port)>#<tag>
	// OR ss://<base64(method:password)>@<server>:<port>#<tag>

	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	node := &model.Node{Type: "ss"}
	node.Name = u.Fragment
	if node.Name == "" {
		node.Name = "Imported SS"
	}

	// Case 1: user info is base64 encoded string "method:password"
	// and host is plain in URL
	if u.User.String() != "" {
		// This is actually "method:password" encoded or plain?
		// Standard SIP002: ss://base64(method:password)@hostname:port
		// But often full block is base64 encoded.

		// Check if host is present. If yes, it's likely SIP002
		node.Server = u.Hostname()
		port, _ := strconv.Atoi(u.Port())
		node.Port = port

		// Decode user info
		userInfo := u.User.String()
		// Try decoding if it looks like base64 (no colon)
		if !strings.Contains(userInfo, ":") {
			// Try multiple base64 decodings
			// 1. RawURLEncoding (no padding) - SIP002 standard
			if decoded, err := base64.RawURLEncoding.DecodeString(userInfo); err == nil {
				userInfo = string(decoded)
			} else if decoded, err := base64.URLEncoding.DecodeString(userInfo); err == nil {
				// 2. URLEncoding (with padding)
				userInfo = string(decoded)
			} else if decoded, err := base64.RawStdEncoding.DecodeString(userInfo); err == nil {
				// 3. RawStdEncoding (no padding, standard chars)
				userInfo = string(decoded)
			} else if decoded, err := base64.StdEncoding.DecodeString(userInfo); err == nil {
				// 4. StdEncoding (with padding, standard chars)
				userInfo = string(decoded)
			}
		}

		parts := strings.SplitN(userInfo, ":", 2)
		if len(parts) == 2 {
			node.Cipher = parts[0]
			node.Password = parts[1]
		}
		return node, nil
	}

	// Case 2: Everything in host is base64 encoded "method:password@hostname:port" (Legacy)
	// u.Host might contain the base64 string
	raw := u.Host
	if u.Path != "" { // sometimes / is appended
		raw += u.Path
	}

	var decoded []byte
	var decodeErr error
	// Try multiple base64 decodings
	if decoded, decodeErr = base64.RawURLEncoding.DecodeString(raw); decodeErr != nil {
		if decoded, decodeErr = base64.URLEncoding.DecodeString(raw); decodeErr != nil {
			if decoded, decodeErr = base64.RawStdEncoding.DecodeString(raw); decodeErr != nil {
				if decoded, decodeErr = base64.StdEncoding.DecodeString(raw); decodeErr != nil {
					return nil, decodeErr
				}
			}
		}
	}

	// "method:password@hostname:port"
	fullStr := string(decoded)
	parts := strings.Split(fullStr, "@")
	if len(parts) != 2 {
		return nil, errors.New("invalid ss format")
	}

	auth := strings.SplitN(parts[0], ":", 2)
	if len(auth) != 2 {
		return nil, errors.New("invalid ss auth")
	}
	node.Cipher = auth[0]
	node.Password = auth[1]

	serverParts := strings.Split(parts[1], ":")
	if len(serverParts) != 2 {
		return nil, errors.New("invalid ss server")
	}
	node.Server = serverParts[0]
	node.Port, _ = strconv.Atoi(serverParts[1])

	return node, nil
}

func parseVmess(link string) (*model.Node, error) {
	// vmess://<base64(json)>
	b64 := strings.TrimPrefix(link, "vmess://")
	decoded, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		decoded, err = base64.RawStdEncoding.DecodeString(b64)
	}
	if err != nil {
		return nil, err
	}

	var vMap map[string]interface{}
	if err := json.Unmarshal(decoded, &vMap); err != nil {
		return nil, err
	}

	node := &model.Node{
		Type:    "vmess",
		Name:    getString(vMap, "ps"),
		Server:  getString(vMap, "add"),
		UUID:    getString(vMap, "id"),
		Cipher:  "auto",
		Network: getString(vMap, "net"),
		Path:    getString(vMap, "path"),
		Host:    getString(vMap, "host"),
	}

	// Port can be string or int in JSON
	if p, ok := vMap["port"].(float64); ok {
		node.Port = int(p)
	} else if pStr, ok := vMap["port"].(string); ok {
		node.Port, _ = strconv.Atoi(pStr)
	}

	// Construct extra config for TLS/Network
	if getString(vMap, "tls") != "" {
		node.TLS = true
	}

	// Capture extra fields like alterId into ExtraConfig
	extra := make(map[string]interface{})

	// alterId (aid) - try different types
	if v, ok := vMap["aid"]; ok {
		extra["alterId"] = v // Keep original type (int/string) or convert? Clash usually accepts int.
		// If it's a string number "0", json.Unmarshal later in generator handles it fine if it fits interface{}
	}

	if len(extra) > 0 {
		if b, err := json.Marshal(extra); err == nil {
			node.ExtraConfig = string(b)
		}
	}

	return node, nil
}

func parseTrojan(link string) (*model.Node, error) {
	// trojan://password@host:port#name
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	node := &model.Node{
		Type:     "trojan",
		Name:     u.Fragment,
		Server:   u.Hostname(),
		Password: u.User.Username(),
	}
	if node.Name == "" {
		node.Name = "Imported Trojan"
	}
	node.Port, _ = strconv.Atoi(u.Port())

	// Check params like snell/tls
	q := u.Query()
	if q.Get("sni") != "" {
		// store in extra?
	}

	return node, nil
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

// ExportLink converts a Node model to a shareable link
func ExportLink(node *model.Node) (string, error) {
	switch node.Type {
	case "ss", "shadowsocks":
		return exportSS(node)
	case "vmess":
		return exportVmess(node)
	case "trojan":
		return exportTrojan(node)
	case "vless":
		return exportVless(node)
	case "hysteria2", "hysteria":
		return exportHysteria2(node)
	default:
		return "", errors.New("unsupported node type: " + node.Type)
	}
}

func exportSS(node *model.Node) (string, error) {
	// Format: ss://base64(method:password@server:port)#name
	if node.Cipher == "" || node.Password == "" {
		return "", errors.New("missing required SS fields")
	}

	userInfo := fmt.Sprintf("%s:%s", node.Cipher, node.Password)
	serverPart := fmt.Sprintf("%s:%d", node.Server, node.Port)
	full := fmt.Sprintf("%s@%s", userInfo, serverPart)

	encoded := base64.RawURLEncoding.EncodeToString([]byte(full))
	link := fmt.Sprintf("ss://%s", encoded)

	if node.Name != "" {
		link += "#" + url.QueryEscape(node.Name)
	}

	return link, nil
}

func exportVmess(node *model.Node) (string, error) {
	// vmess://base64(json)
	if node.UUID == "" {
		return "", errors.New("missing UUID for VMess")
	}

	// Parse extra config for alterId
	alterId := 0
	if node.ExtraConfig != "" {
		var extra map[string]interface{}
		if err := json.Unmarshal([]byte(node.ExtraConfig), &extra); err == nil {
			if v, ok := extra["alterId"]; ok {
				switch val := v.(type) {
				case float64:
					alterId = int(val)
				case int:
					alterId = val
				case string:
					alterId, _ = strconv.Atoi(val)
				}
			}
		}
	}

	vMap := map[string]interface{}{
		"v":    "2",
		"ps":   node.Name,
		"add":  node.Server,
		"port": node.Port,
		"id":   node.UUID,
		"aid":  alterId,
		"net":  node.Network,
		"type": "none",
		"host": node.Host,
		"path": node.Path,
		"tls":  "",
	}

	if node.TLS {
		vMap["tls"] = "tls"
	}

	jsonBytes, _ := json.Marshal(vMap)
	encoded := base64.StdEncoding.EncodeToString(jsonBytes)
	return fmt.Sprintf("vmess://%s", encoded), nil
}

func exportTrojan(node *model.Node) (string, error) {
	// trojan://password@host:port#name
	if node.Password == "" {
		return "", errors.New("missing password for Trojan")
	}

	link := fmt.Sprintf("trojan://%s@%s:%d", node.Password, node.Server, node.Port)

	// Add query parameters for SNI
	if node.Host != "" {
		link += fmt.Sprintf("?sni=%s", url.QueryEscape(node.Host))
	}

	if node.Name != "" {
		link += "#" + url.QueryEscape(node.Name)
	}

	return link, nil
}

func exportVless(node *model.Node) (string, error) {
	// vless://uuid@server:port?params#name
	if node.UUID == "" {
		return "", errors.New("missing UUID for VLESS")
	}

	params := url.Values{}
	params.Set("type", node.Network)
	if node.Network == "ws" || node.Network == "grpc" {
		if node.Path != "" {
			params.Set("path", node.Path)
		}
		if node.Host != "" {
			params.Set("host", node.Host)
		}
	}
	if node.TLS {
		params.Set("security", "tls")
		if node.Host != "" {
			params.Set("sni", node.Host)
		}
	}

	link := fmt.Sprintf("vless://%s@%s:%d?%s", node.UUID, node.Server, node.Port, params.Encode())

	if node.Name != "" {
		link += "#" + url.QueryEscape(node.Name)
	}

	return link, nil
}

func exportHysteria2(node *model.Node) (string, error) {
	// hysteria2://password@server:port?params#name
	password := node.Password
	if password == "" && node.UUID != "" {
		password = node.UUID
	}
	if password == "" {
		return "", errors.New("missing password for Hysteria2")
	}

	params := url.Values{}
	if node.Host != "" {
		params.Set("sni", node.Host)
	}

	link := fmt.Sprintf("hysteria2://%s@%s:%d", url.QueryEscape(password), node.Server, node.Port)
	if len(params) > 0 {
		link += "?" + params.Encode()
	}

	if node.Name != "" {
		link += "#" + url.QueryEscape(node.Name)
	}

	return link, nil
}
