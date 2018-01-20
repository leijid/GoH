##GoH安装步骤
####1、安装GoLang
####2、配置GOROOT（GoLang安装目录）
####3、配置GOPATH（工程开发目录）
####4、安装gvt，后续开发时添加包依赖工具使用（添加新的依赖时：cd $GOPATH/src && gvt fetch ~~~）
	cd $GOPATH
	go get github.com/FiloSottile/gvt
####5、添加依赖
    cd $GOPATH
    gvt fetch github.com/dlintw/goconf
    gvt fetch github.com/gorilla/mux
    gvt fetch github.com/garyburd/redigo/redis
    gvt fetch github.com/op/go-logging
    gvt fetch gopkg.in/mgo.v2
    gvt fetch github.com/go-sql-driver/mysql
    gvt fetch github.com/streadway/amqp
    gvt fetch github.com/robfig/cron
    gvt fetch github.com/mitchellh/mapstructure
####5、下载GoH，放到$GOPATH/src目录中
####6、编译运行GoH
	cd $GOPATH/src/GoH
	go build .
	./GoH
##https开启步骤
####1、修改GoH中goh.ini中https配置
	is_https        = true
####2、生成server.crt和server.key（正式生产环境是需要购买证书的）
	openssl genrsa -out goh_server.key 2048
	openssl req -new -x509 -key goh_server.key -out goh_server.crt -days 365
####3、goh.ini中配置server.crt和server.key路径
	https_server_crt = ./https/goh_server.crt
	https_server_key = ./https/goh_server.key


##已实现内容
####1、http server
####2、https server
####3、url路由
####4、redis链接
####5、log文件记录
####6、mysql操作
####7、session本地存储
####8、session分布式存储
####9、rabbitmq操作
####10、简单的错误处理机制
####11、mongodb操作

##待实现内容
####1、zookeeper操作
####2、RPC通信
####3、mysql/redis/mongodb增删改查功能封装