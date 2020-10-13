# go-gin-chat(Gin+websocket 的多人聊天室)

> 练手小项目，为熟悉Gin框架跟websocket使用 💛💛💛💛💛💛

> [在线demo](http://go-gin-chat.hzz.cool) (PS: 请尽量使用Chrome游览器)

## 结构
```
.
|-- conf
|-- controller
|-- models
|-- routes
|-- services
|   |               `-- img_kr
|   |               `-- message_service
|   |               `-- session
|   |               `-- user_service
|                   `-- validator
|-- sql
|-- static
|   |               `-- bootstrap
|   |   |           `-- css
|   |   |           `-- fonts
|   |               `-- js
|   |-- emoji
|   |-- images
|   |   |           `-- face
|   |   |           `-- rooms
|   |   |           `-- theme
|   |               `-- user
|   |-- javascripts
|   |-- rolling
|   |   |           `-- css
|   |               `-- js
|   `-- stylesheets
|-- tmp
|-- tmp_images
|-- views
`-- ws

```

## 界面
![](https://static01.imgkr.com/temp/5c3c9096ef9f4796b404dd2f3e23c36d.png)
![](https://static01.imgkr.com/temp/cd66af62792f4d2e8c2fa974e82d0526.png)
![](https://static01.imgkr.com/temp/099bf697686445d79407962cdfb11e56.png)
![](https://static01.imgkr.com/temp/1e89fdd024de47fa862143fba246d632.png)

## feature
- 登录/注册(防止重复登录)
- 群聊(支持文字、emoji、图片)
- 历史消息查看(暂时仅支持最新100条)
- 心跳检测，来自 https://github.com/zimv/websocket-heartbeat-js
- go mod 包管理
- 静态资源嵌入，运行只依赖编译好的可执行文件与mysql

## database
#### mysql
![](https://static01.imgkr.com/temp/a4b4520607da41f796f13c17a250e70e.png)
 
## Tools
- [模板提供](https://github.com/zfowed/charooms-html) 
- github.com/gin-gonic/gin
- gorm.io/driver/mysql
- gorm.io/gorm
- github.com/gravityblast/fresh
- github.com/valyala/fasthttp
- github.com/spf13/viper
- https://blog.hi917.com/detail/87.html

## 使用

```
# 自行导入数据库文件 sql/go_gin_chat.sql
git clone github.com/hezhizheng/go-gin-chat
cd go-gin-chat
cp conf/config.go.env conf/config.go // 根据实际情况修改配置
go-bindata -o=bindata/bindata.go -pkg=bindata ./static/... ./views/... // 安装请参考 https://blog.hi917.com/detail/87.html
go run main.go 
```

## nginx 部署

```
server {
    listen 80;
    server_name  go-gin-chat.hzz.cool;

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
# 生成文件可直接执行 Linux
gox -osarch="linux/amd64"
......
```

## Tip
- 修改静态文件需要执行 `go-bindata -o=bindata/bindata.go -pkg=bindata ./static/... ./views/...`  重新编译

## todo
- [x] 心跳机制
- [ ] 多频道聊天
- [ ] 私聊
- [ ] 在线用户列表
- [ ] https支持
