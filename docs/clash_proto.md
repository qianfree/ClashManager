# Clash & Clash.meta (Mihomo) 支持协议列表

本文档详细列出了 Clash 及 Clash.meta (Mihomo) 支持的代理协议及其配置格式。

## 1. 原生 Clash 协议 (Standard Protocols)

这些协议在原版 Clash Premium 和 Clash.meta 中均可使用。

### HTTP / HTTPS
最基础的代理协议。
```yaml
- name: "HTTP"
  type: http
  server: server.com
  port: 8080
  username: "username"
  password: "password"
  # tls: true # 开启 HTTPS
  # skip-cert-verify: true
```

### SOCKS5
```yaml
- name: "Socks5"
  type: socks5
  server: server.com
  port: 1080
  username: "username"
  password: "password"
  # tls: true
  # skip-cert-verify: true
  # udp: true
```

### Shadowsocks (SS)
支持多种加密方式 (aes-128-gcm, aes-256-gcm, chacha20-ietf-poly1305 等)。
```yaml
- name: "SS"
  type: ss
  server: server.com
  port: 443
  cipher: chacha20-ietf-poly1305
  password: "password"
  # plugin: obfs
  # plugin-opts:
  #   mode: tls
  #   host: bing.com
```

### ShadowsocksR (SSR)
部分客户端支持，Clash.meta 完美支持。
```yaml
- name: "SSR"
  type: ssr
  server: server.com
  port: 443
  cipher: chacha20
  password: "password"
  obfs: tls1.2_ticket_auth
  protocol: auth_aes128_md5
  obfs-param: "obfs-param"
  protocol-param: "protocol-param"
```

### VMess (V2Ray)
广泛使用的协议。
```yaml
- name: "Vmess"
  type: vmess
  server: server.com
  port: 443
  uuid: "uuid"
  alterId: 0
  cipher: auto
  # network: ws / h2 / http / grpc / tcp
  network: ws
  ws-opts:
    path: "/path"
    headers:
      Host: v2ray.com
  # tls: true
  # skip-cert-verify: true
  # servername: v2ray.com
```

### Trojan
```yaml
- name: "Trojan"
  type: trojan
  server: server.com
  port: 443
  password: "password"
  # udp: true
  # sni: example.com
  # skip-cert-verify: true
  # network: ws / grpc
```

### Snell
Clash 独有协议。
```yaml
- name: "Snell"
  type: snell
  server: server.com
  port: 44046
  psk: "yourpsk"
  version: 2
  # obfs-opts:
  #   mode: http # or tls
  #   host: bing.com
```

---

## 2. Clash.meta (Mihomo) 独占协议 (Meta-Only)

以下协议仅在 Clash.meta (Mihomo) 内核中支持，原版 Clash 无法识别。

### VLESS (V2Ray)
轻量级协议，性能优于 VMess。
```yaml
- name: "Vless"
  type: vless
  server: server.com
  port: 443
  uuid: "uuid"
  # flow: xtls-rprx-vision # 启用 Reality / Vision 必须
  network: tcp
  tls: true
  udp: true
  servername: example.com
  # reality-opts:
  #   public-key: "your-public-key"
  #   short-id: "short-id"
  client-fingerprint: chrome
```

### Hysteria (v1)
基于 UDP 的高速协议。
```yaml
- name: "Hysteria"
  type: hysteria
  server: server.com
  port: 443
  auth_str: "your-auth-str"
  # up: 100 mbps
  # down: 100 mbps
  # alpn:
  #   - h3
  # protocol: udp
```

### Hysteria2
Hysteria 的升级版，配置更简单。
```yaml
- name: "Hysteria2"
  type: hysteria2
  server: server.com
  port: 443
  password: "password"
  sni: example.com
  skip-cert-verify: true
  # obfs: "salamander"
  # obfs-password: "password"
```

### Tuic (v5)
另一种基于 QUIC 的协议。
```yaml
- name: "Tuic"
  type: tuic
  server: server.com
  port: 443
  uuid: "uuid"
  password: "password"
  congestion-controller: bbr
  # udp-relay-mode: native
  # reduce-rtt: true
```

### WireGuard
著名的 VPN 协议。
```yaml
- name: "WireGuard"
  type: wireguard
  server: server.com
  port: 51820
  ip: "10.0.0.2"
  ipv6: "fd00::2" # 可选
  private-key: "your-private-key"
  public-key: "peer-public-key"
  # pre-shared-key: "psk"
  # mtu: 1280
  # udp: true
```

### SSH
将 SSH 作为代理通过。
```yaml
- name: "SSH"
  type: ssh
  server: server.com
  port: 22
  username: "root"
  password: "password"
  # private-key: "-----BEGIN OPENSSH PRIVATE KEY-----..."
```

> **注意**: 使用 Meta 独占协议时，请确保你的客户端（如 Clash Verge, Clash Meta for Android 等）使用的是 **Clash.meta (Mihomo)** 内核。
