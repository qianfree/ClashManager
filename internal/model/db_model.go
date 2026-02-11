package model

import (
	"time"

	"gorm.io/gorm"
)

// GlobalSetting stores key-value pairs for global configurations (e.g., DNS, Ports)
type GlobalSetting struct {
	Key   string `gorm:"primaryKey"`
	Value string // JSON string
}

// User represents an administrator
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex" json:"username"`
	Password string `json:"-"`                        // Bcrypt hash, never send to client
	Token    string `gorm:"uniqueIndex" json:"token"` // Subscription token
}

// Node represents a proxy node (SS, Vmess, Trojan, etc.)
type Node struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex"`
	Type        string // ss, vmess, hysteria2, etc.
	Server      string
	Port        int
	UUID        string // VMess, VLESS, Hysteria2 (password sometimes stored here or in Password)
	Password    string // SS, Trojan, Hysteria2
	Cipher      string // SS, VMess
	UDP         bool
	TLS         bool
	SkipCert    bool
	Network     string // ws, grpc, etc.
	Path        string // WebSocket path
	Host        string // HTTP Host / SNI
	ALPN        string // h3
	ExtraConfig string // JSON for other fields like headers, flow, etc.
}

// Rule represents a routing rule
type Rule struct {
	gorm.Model
	Type       string // DOMAIN-SUFFIX, IP-CIDR, etc.
	Payload    string // google.com, US
	Target     string // ProxyGroupName or specifics (Legacy)
	TargetID   uint   `gorm:"default:0"` // ID of the target Node or Group
	TargetType string // "node", "group"
	Priority   int    `gorm:"default:0"` // Lower number = higher priority
	NoResolve  bool
	Remark     string // Remark for display only, not used in config generation
}

// ProxyGroupDB represents a user-defined proxy group (optional if we want dynamic groups)
// For simplicity, we might store groups in GlobalSetting or a separate table
type ProxyGroupModel struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex"`
	Type     string // select, url-test
	ProxyIDs string // JSON array of node IDs: [1, 2, 3]
	Use      string // JSON array of provider names/other groups
	URL      string
	Interval int
}

// ProxyNode represents a node reference in a proxy group
type ProxyNode struct {
	ID   uint   `json:"id"`
	Type string `json:"type"` // "node" or "group"
	Name string `json:"name"` // populated by API for display
}

// SubscriptionLog represents a subscription access log
type SubscriptionLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
	UserID    uint           `gorm:"index" json:"userId"`                     // User who subscribed
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"` // Association to User
	Token     string         `gorm:"index" json:"token"`                      // Token used (without exposing full token)
	IP        string         `gorm:"index" json:"ip"`                         // Client IP
	UserAgent string         `json:"userAgent"`                               // Client user agent
	Success   bool           `json:"success"`                                 // Whether subscription was successful
	Error     string         `json:"error"`                                   // Error message if failed
}

// TableName specifies the table name for SubscriptionLog
func (SubscriptionLog) TableName() string {
	return "subscription_logs"
}
