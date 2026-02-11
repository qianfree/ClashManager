package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
			decoded, err := base64.RawURLEncoding.DecodeString(userInfo)
			if err == nil {
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

	decoded, err := base64.RawURLEncoding.DecodeString(raw)
	if err != nil {
		// Try standard encoding
		decoded, err = base64.URLEncoding.DecodeString(raw)
	}
	if err != nil {
		return nil, err
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
