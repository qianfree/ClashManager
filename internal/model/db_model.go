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
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Username  string         `gorm:"uniqueIndex" json:"username"`
	Password  string         `json:"-"`                        // Bcrypt hash, never send to client
	Token     string         `gorm:"uniqueIndex" json:"token"` // Subscription token
}

// Node represents a proxy node (SS, Vmess, Trojan, etc.)
type Node struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Server      string         `json:"server"`
	Port        int            `json:"port"`
	UUID        string         `json:"uuid,omitempty"`
	Password    string         `json:"password,omitempty"`
	Cipher      string         `json:"cipher,omitempty"`
	UDP         bool           `json:"udp"`
	TLS         bool           `json:"tls"`
	SkipCert    bool           `json:"skip_cert"`
	Network     string         `json:"network,omitempty"`
	Path        string         `json:"path,omitempty"`
	Host        string         `json:"host,omitempty"`
	ALPN        string         `json:"alpn,omitempty"`
	ExtraConfig string         `json:"extra_config,omitempty"`
	Source      string         `json:"source"`
	SourceID    uint           `json:"source_id"`
	SourceName  string         `json:"source_name"`
}

// Rule represents a routing rule
type Rule struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Type       string         `json:"type"`
	Payload    string         `json:"payload"`
	Target     string         `json:"target"`
	TargetID   uint           `json:"target_id"`
	TargetType string         `json:"target_type"`
	Priority   int            `json:"priority"`
	NoResolve  bool           `json:"no_resolve"`
	Tag        string         `json:"tag"`
	Remark     string         `json:"remark"`
}

// ProxyGroupDB represents a user-defined proxy group (optional if we want dynamic groups)
// For simplicity, we might store groups in GlobalSetting or a separate table
type ProxyGroupModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	ProxyIDs  string         `json:"proxy_ids"`
	Use       string         `json:"use"`
	URL       string         `json:"url"`
	Interval  int            `json:"interval"`
	Source    string         `gorm:"default:user" json:"source"`
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
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	UserID    uint           `gorm:"index" json:"user_id"`                    // User who subscribed
	User      *User          `gorm:"foreignKey:UserID" json:"user,omitempty"` // Association to User
	Token     string         `gorm:"index" json:"token"`                      // Token used (without exposing full token)
	IP        string         `gorm:"index" json:"ip"`                         // Client IP
	UserAgent string         `json:"user_agent"`                              // Client user agent
	Success   bool           `json:"success"`                                 // Whether subscription was successful
	Error     string         `json:"error"`                                   // Error message if failed
}

// TableName specifies the table name for SubscriptionLog
func (SubscriptionLog) TableName() string {
	return "subscription_logs"
}

// SubscriptionSource represents a third-party subscription source
type SubscriptionSource struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name           string         `json:"name"`
	URL            string         `json:"url"`
	Enabled        bool           `json:"enabled"`
	UpdateInterval int            `json:"update_interval"`
	LastSync       *time.Time     `json:"last_sync"`
	NodeTag        string         `json:"node_tag"`
	SyncMode       string         `json:"sync_mode"`
	Error          string         `json:"error"`
}

// TableName specifies the table name for SubscriptionSource
func (SubscriptionSource) TableName() string {
	return "subscription_sources"
}
