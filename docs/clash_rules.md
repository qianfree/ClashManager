下面是 **Clash Premium / Meta（当前主流）** 完整穷举可用的规则类型，并带有**清晰示例**与用途说明。

---

## 一、域名类规则（最常用）

### 1️⃣ `DOMAIN`

**完全匹配域名**

```yaml
- DOMAIN,example.com,Proxy
```

✔ 仅匹配 `example.com`
✖ 不匹配 `www.example.com`

---

### 2️⃣ `DOMAIN-SUFFIX`

**匹配域名后缀（子域名）**

```yaml
- DOMAIN-SUFFIX,google.com,Proxy
```

✔ `www.google.com`
✔ `mail.google.com`

---

### 3️⃣ `DOMAIN-KEYWORD`

**域名包含关键词**

```yaml
- DOMAIN-KEYWORD,google,Proxy
```

✔ `googleusercontent.com`
✔ `mygoogletest.net`

---

### 4️⃣ `DOMAIN-REGEX`

**正则匹配域名**

```yaml
- DOMAIN-REGEX,^.*\.google\.com$,Proxy
```

✔ 适合复杂匹配
⚠ 性能略差，不建议大量使用

---

## 二、IP 类规则

### 5️⃣ `IP-CIDR`

**IPv4 地址或网段**

```yaml
- IP-CIDR,192.168.1.100/32,DIRECT
- IP-CIDR,10.0.0.0/8,DIRECT
```

✔ 常用于 **局域网设备分流**

---

### 6️⃣ `IP-CIDR6`

**IPv6 地址或网段**

```yaml
- IP-CIDR6,2001:db8::/32,Proxy
```

---

### 7️⃣ `GEOIP`

**按国家/地区 IP 库**

```yaml
- GEOIP,CN,DIRECT
- GEOIP,US,Proxy
```

✔ 非常常见
⚠ 基于 MaxMind IP 库，非 100% 准确

---

## 三、端口类规则

### 8️⃣ `DST-PORT`

**目标端口**

```yaml
- DST-PORT,443,Proxy
- DST-PORT,22,DIRECT
```

---

### 9️⃣ `SRC-PORT`

**源端口（较少用）**

```yaml
- SRC-PORT,12345,Proxy
```

---

## 四、进程 / 应用规则（仅 Premium / Meta）

### 🔟 `PROCESS-NAME`

**进程名（强烈推荐）**

```yaml
- PROCESS-NAME,Telegram.exe,Proxy
- PROCESS-NAME,WeChat.exe,DIRECT
```

✔ Windows / macOS / Linux
✔ 分应用代理神器

---

### 1️⃣1️⃣ `PROCESS-PATH`

**进程路径**

```yaml
- PROCESS-PATH,C:\Program Files\Google\Chrome\Application\chrome.exe,Proxy
```

✔ 更精确
⚠ 路径变更会失效

---

### 1️⃣2️⃣ `PROCESS-NAME-REGEX`

**进程名正则**

```yaml
- PROCESS-NAME-REGEX,.*chrome.*,Proxy
```

---

## 五、网络层 / 协议规则

### 1️⃣3️⃣ `NETWORK`

**按网络类型**

```yaml
- NETWORK,TCP,Proxy
- NETWORK,UDP,DIRECT
```

---

### 1️⃣4️⃣ `IN-TYPE`

**入站类型**

```yaml
- IN-TYPE,SOCKS,Proxy
- IN-TYPE,HTTP,DIRECT
```

---

### 1️⃣5️⃣ `IN-PORT`

**入站端口**

```yaml
- IN-PORT,7890,Proxy
```

---

### 1️⃣6️⃣ `IN-NAME`

**入站名称**

```yaml
- IN-NAME,lan,DIRECT
```

---

## 六、DNS 专用规则（Meta 支持）

### 1️⃣7️⃣ `DNS`

**DNS 查询规则**

```yaml
- DNS,local,DIRECT
```

---

### 1️⃣8️⃣ `DNS-SUFFIX`

```yaml
- DNS-SUFFIX,lan,DIRECT
```

---

## 七、规则集（Rule Provider）

### 1️⃣9️⃣ `RULE-SET`

**引用规则集文件**

```yaml
- RULE-SET,telegram,Proxy
- RULE-SET,apple,DIRECT
```

✔ 适合维护大型规则
✔ 配合 `rule-providers`

---

## 八、兜底规则（必须有）

### 2️⃣0️⃣ `MATCH`

**匹配所有未命中的流量（最后一条）**

```yaml
- MATCH,Proxy
```

⚠ **只能放在最后**

---

## 九、完整规则示例（推荐结构）

```yaml
rules:
  - PROCESS-NAME,Telegram.exe,Proxy
  - IP-CIDR,192.168.1.100/32,DIRECT
  - DOMAIN-SUFFIX,google.com,Proxy
  - GEOIP,CN,DIRECT
  - MATCH,Proxy
```

---

## 十、一句话总结

> **Clash 规则 = 从“最精确”到“最宽泛”依次匹配，命中即停**

**推荐优先级：**

```
进程 → IP → 域名 → GEOIP → MATCH
```
