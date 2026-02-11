# Clash & Clash.meta (Mihomo) 配置指南

这份指南将帮助你理解 Clash 及 Clash.meta (Mihomo) 的核心配置结构。

## 1. 基本配置结构 (Basic Config)

这是配置文件的基础部分，定义了运行模式、端口等。

```yaml
# HTTP/HTTPS 代理端口
port: 7890
# SOCKS5 代理端口
socks-port: 7891
# 混合端口 (HTTP + SOCKS5)
mixed-port: 7892

# 允许局域网连接
allow-lan: true
# 局域网绑定地址 (ipv4)
bind-address: "*"

# 运行模式
# rule: 按照规则路由 (最常用)
# global: 全局代理 (所有流量走代理)
# direct: 直连 (所有流量不走代理)
mode: rule

# 日志级别 (info / warning / error / debug / silent)
log-level: info

# 启用 IPv6 支持
ipv6: false

# RESTful API 控制器 (用于面板连接)
external-controller: 0.0.0.0:9090
# API 访问密钥 (可选)
secret: "your-secret-key"
```

## 2. DNS 配置 (DNS)

DNS 是防止 DNS 污染的关键。Clash 使用 Fake-IP 或 Redir-Host 模式。推荐使用 **Fake-IP**。

```yaml
dns:
  enable: true
  listen: 0.0.0.0:53
  # 增强模式: fake-ip 或 redir-host
  enhanced-mode: fake-ip
  # Fake-IP 过滤范围 (这些域名将返回真实 IP)
  fake-ip-filter:
    - "*.lan"
    - "localhost.ptlogin2.qq.com"

  # 默认 DNS (用于解析下方的 nameserver 域名)
  default-nameserver:
    - 223.5.5.5
    - 119.29.29.29

  # 主要 DNS 服务器
  nameserver:
    - https://dns.alidns.com/dns-query
    - https://doh.pub/dns-query

  # 备用 DNS (Fallback)
  # 当 nameserver 返回的 IP 是非 CN IP 时，会使用 fallback 的结果
  # 只有当 fallback-filter 启用时才生效
  fallback:
    - https://1.1.1.1/dns-query
    - https://8.8.8.8/dns-query

  # Fallback 过滤器
  fallback-filter:
    geoip: true
    ipcidr:
      - 240.0.0.0/4
```

## 3. 代理集合与节点 (Proxies & Providers)

### 3.1 手动定义节点 (Proxies)

```yaml
proxies:
  - name: "Shadowsocks节点"
    type: ss
    server: server.com
    port: 443
    cipher: chacha20-ietf-poly1305
    password: "password"

  - name: "Vmess节点"
    type: vmess
    server: server.com
    port: 443
    uuid: "uuid"
    alterId: 0
    cipher: auto
    network: ws
    ws-opts:
      path: "/path"

  # Clash.meta 特有协议 (Hysteria2)
  - name: "Hysteria2节点"
    type: hysteria2
    server: server.com
    port: 443
    password: "password"
    sni: example.com
    skip-cert-verify: true
```

### 3.2 引用外部节点订阅 (Proxy Providers)

推荐使用 Provider 管理机场订阅。

```yaml
proxy-providers:
  MyAirport:
    type: http
    url: "https://机场订阅链接"
    interval: 3600
    path: ./proxies/airport.yaml
    health-check:
      enable: true
      interval: 600
      url: http://www.gstatic.com/generate_204
```

## 4. 策略组 (Proxy Groups)

策略组用于对节点进行分类和自动选择。

```yaml
proxy-groups:
  # 手动选择
  - name: "🚀 节点选择"
    type: select
    proxies:
      - "♻️ 自动选择"
      - "🇭🇰 香港节点"
      - "🇺🇸 美国节点"
      - DIRECT

  # 自动测速选择最快节点
  - name: "♻️ 自动选择"
    type: url-test
    url: http://www.gstatic.com/generate_204
    interval: 300
    tolerance: 50
    use:
      - MyAirport # 引用上面的 Provider

  # 故障转移 (如果第一个节点挂了，自动切第二个)
  - name: "🛡️ 故障转移"
    type: fallback
    url: http://www.gstatic.com/generate_204
    interval: 300
    use:
      - MyAirport
```

## 5. 规则配置 (Rules)

规则由上至下匹配，一旦匹配成功即停止。

### 规则类型
- `DOMAIN`: 精确匹配域名
- `DOMAIN-SUFFIX`: 匹配域名后缀 (推荐)
- `DOMAIN-KEYWORD`: 匹配域名关键字
- `IP-CIDR`: 匹配 IP 段
- `GEOIP`: 匹配 IP 地理位置
- `MATCH`: 兜底规则 (全不匹配时)

### 配置示例

```yaml
rules:
  # 5.1 制定域名走固定节点
  # 让 openai.com 及其子域名走 "🇺🇸 美国节点"
  - DOMAIN-SUFFIX,openai.com,🇺🇸 美国节点
  # 让包含 google 的域名走 "🚀 节点选择"
  - DOMAIN-KEYWORD,google,🚀 节点选择
  
  # 5.2 常见应用规则
  - DOMAIN-SUFFIX,github.com,🚀 节点选择
  - DOMAIN-SUFFIX,microsoft.com,DIRECT
  
  # 5.3 局域网直连
  - IP-CIDR,192.168.0.0/16,DIRECT
  - IP-CIDR,10.0.0.0/8,DIRECT
  
  # 5.4 国内流量直连 (基于 GeoIP)
  - GEOIP,CN,DIRECT
  
  # 5.5 兜底规则 (必须存在)
  # 其他未匹配流量走 "🐟 漏网之鱼" 或 "🚀 节点选择"
  - MATCH,🚀 节点选择
```

## 6. 进阶：如何让指定域名走固定节点

这是一个非常常见的需求。假设你想让 `netflix.com` 强制走 `🇸🇬 新加坡节点`。

1.  **定义节点或策略组**: 确保你有一个包含新加坡节点的策略组，或者直接引用节点名称。
    ```yaml
    proxy-groups:
      - name: "🇸🇬 新加坡组"
        type: select
        proxies:
          - "Singleton SG Node" # 具体节点名
          - "SG Auto"           # 或者另一个自动测速组
    ```

2.  **添加规则**: 在 `rules` 的 **上方** 添加规则（越靠前优先级越高）。
    ```yaml
    rules:
      - DOMAIN-SUFFIX,netflix.com,🇸🇬 新加坡组
      - DOMAIN-SUFFIX,netflix.net,🇸🇬 新加坡组
      - DOMAIN-SUFFIX,nflximg.net,🇸🇬 新加坡组
      # ... 其他规则
    ```

## 7. Clash.meta (Mihomo) 独有特性

Clash.meta 支持更多新协议和规则类型：

- **GEOSITE**: 支持引用 `geosite.dat` 数据库，比简单的 DOMAIN 列表更强大。
  ```yaml
  rules:
    - GEOSITE,category-ads-all,REJECT
    - GEOSITE,bilibili,DIRECT
    - GEOSITE,youtube,🚀 节点选择
  ```
- **Rule Providers**: 引用外部规则集，方便管理大量规则。
  ```yaml
  rule-providers:
    reject:
      type: http
      behavior: domain
      url: "https://rules.com/reject.yaml"
      path: ./rules/reject.yaml
      interval: 86400
  
  rules:
    - RULE-SET,reject,REJECT
  ```
