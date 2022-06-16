# go-gin-chat(Gin+websocket 的多人聊天室)

> 练手小项目，为熟悉Gin框架跟websocket使用 💛💛💛💛💛💛

> [在线demo](http://go-gin-chat.hzz.cool) (PS: 请尽量使用Chrome游览器，开启多个不同用户游览器即可体验效果)

> [github地址](https://github.com/hezhizheng/go-gin-chat)

## feature
- 登录/注册(防止重复登录)
- 群聊(多房间、支持文字、emoji、文件(图片)上传，使用 [freeimage.host](https://freeimage.host/) 做图床 )
- 私聊(消息提醒)
- 历史消息查看(点击加载更多)
- 心跳检测，来自 https://github.com/zimv/websocket-heartbeat-js
- go mod 包管理
- 使用 Golang 1.16 embed 内嵌静态资源(html、js、css等)，运行只依赖编译好的可执行文件与mysql
- 支持 http/ws 、 https/wss

## 结构
```
.
|-- LICENSE.txt
|-- conf #配置文件
|   |-- config.go
|   `-- config.go.env
|-- controller
|   |-- ImageController.go
|   `-- IndexController.go
|-- main.go
|-- models
|   |-- message.go
|   |-- mysql.go
|   `-- user.go
|-- routes
|   `-- route.go
|-- services # 简单逻辑处理服务层
|   |-- helper
|   |   `-- helper.go
|   |-- img_kr
|   |   `-- imgKr.go
|   |-- message_service
|   |   `-- message.go
|   |-- session
|   |   `-- session.go
|   |-- user_service
|   |   `-- user.go
|   `-- validator
|       `-- validator.go
|-- sql
|   `-- go_gin_chat.sql
|-- static #静态文件 js 、css 、image 目录
|-- views
|   |-- index.html
|   |-- login.html
|   |-- private_chat.html
|   `-- room.html
`-- ws websocket 服务端主要逻辑
    |-- ServeInterface.go 
    |-- go_ws
    |   `-- serve.go # websocket服务端处理代码
    |-- primary
    |   `-- start.go # 为了兼容新旧版 websocket服务端 的调用策略
    |-- serve.go # 初版websocket服务端逻辑代码，可以忽略
    `-- ws_test #本地测试代码
        |-- exec.go
        `-- mock_ws_client_coon.go
```

## 伪代码，详情可参考 [serve.go](./ws/go_ws/serve.go)
- 定义客户端信息的结构体
```go
type wsClients struct {
Conn *websocket.Conn `json:"conn"`

RemoteAddr string `json:"remote_addr"`

Uid float64 `json:"uid"`

Username string `json:"username"`

RoomId string `json:"room_id"`

AvatarId string `json:"avatar_id"`
}

// 
```
- 定义全局变量
```go

// client & serve 的消息体
type msg struct {
Status int             `json:"status"`
Data   interface{}     `json:"data"`
Conn   *websocket.Conn `json:"conn"`
}

// 上线、离线、消息发送事件 的 无缓冲区的 channel
var (
clientMsg = msg{}

enterRooms = make(chan wsClients)

sMsg = make(chan msg)

offline = make(chan *websocket.Conn)

chNotify = make(chan int ,1)
)
```  
- 使用 make 创建一个全局的 `map slice` 用于存放房间与用户的信息，用户上线、离线实际上是对map的 append 跟 remove
```go
var (
rooms = make(map[int][]wsClients)
)
```
- 开启`goroutine`处理用户的连接、离线、消息发送等各个事件
```go
go read(c)
go write()
select {}
```



## 界面
![](https://static01.imgkr.com/temp/5c3c9096ef9f4796b404dd2f3e23c36d.png)
![](https://static01.imgkr.com/temp/cd66af62792f4d2e8c2fa974e82d0526.png)
![](https://static01.imgkr.com/temp/099bf697686445d79407962cdfb11e56.png)
![](https://static01.imgkr.com/temp/1e89fdd024de47fa862143fba246d632.png)

## database
#### mysql
```
CREATE TABLE `messages`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL COMMENT '用户ID',
  `room_id` int(11) NOT NULL COMMENT '房间ID',
  `to_user_id` int(11) NULL DEFAULT 0 COMMENT '私聊用户ID',
  `content` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL COMMENT '聊天内容',
  `image_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '图片URL',
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_user_id`(`user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

CREATE TABLE `users`  (
  `id` int(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '昵称',
  `password` varchar(125) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '' COMMENT '密码',
  `avatar_id` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT '1' COMMENT '头像ID',
  `created_at` datetime(0) NULL DEFAULT NULL,
  `updated_at` datetime(0) NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP(0),
  `deleted_at` datetime(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `username`(`username`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = DYNAMIC;

```
 
## Tools
- [模板提供](https://github.com/zfowed/charooms-html) 
- github.com/gin-gonic/gin
- gorm.io/driver/mysql
- gorm.io/gorm
- github.com/gravityblast/fresh
- github.com/valyala/fasthttp
- github.com/spf13/viper

## 使用 (go version >= 1.16)

```
# 自行导入数据库文件 sql/go_gin_chat.sql
git clone github.com/hezhizheng/go-gin-chat
cd go-gin-chat
cp conf/config.go.env conf/config.go // 根据实际情况修改配置
go run main.go 
```

## nginx 部署

```
server {
    listen 80;
    listen 443 ssl http2;
    server_name  go-gin-chat.hzz.cool;

    #ssl on;  
    ssl_certificate xxxpath\cert.pem;   
    ssl_certificate_key xxxpath\key.pem;   
    ssl_session_timeout  5m;  
    ssl_protocols TLSv1 TLSv1.1 TLSv1.2;  
    ssl_ciphers  ECDHE-RSA-AES128-GCM-SHA256:HIGH:!aNULL:!MD5:!RC4:!DHE;  
    ssl_prefer_server_ciphers  on;

    location ~ .*\.(gif|jpg|png|css|js)(.*) {
                proxy_pass http://127.0.0.1:8322;
                proxy_redirect off;
                proxy_set_header Host $host;
                proxy_cache cache_one;
                proxy_cache_valid 200 302 24h;
                proxy_cache_valid 301 30d;
                proxy_cache_valid any 5m;
                expires 90d;
                add_header wall  "Big brother is watching you";
    }
  

   location / {
       try_files /_not_exists_ @backend;
   }
  
   location @backend {
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header Host            $http_host;

        proxy_pass http://127.0.0.1:8322;
    }
  
   location /ws {
        proxy_pass http://127.0.0.1:8322;
        proxy_redirect off;
        proxy_http_version 1.1;

        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";

        proxy_set_header Host $host:$server_port;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;    
        proxy_read_timeout 6000s;
   }
```
## 编译可执行文件(跨平台)

```
# 用法参考 https://github.com/mitchellh/gox
# go install github.com/mitchellh/gox@latest (go 1.18)
# 生成文件可直接执行 Linux
gox -osarch="linux/amd64" -ldflags "-s -w" -gcflags="all=-trimpath=${PWD}" -asmflags="all=-trimpath=${PWD}"
......
```


## todo
- [x] 心跳机制
- [x] 多频道聊天
- [x] 私聊
- [x] 在线用户列表
- [x] https支持

## License
[MIT](./LICENSE.txt)
