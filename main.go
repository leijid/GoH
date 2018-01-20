package main

import (
	"GoH/core/error"
	"GoH/core/log"
	"GoH/core/mongodb"
	"GoH/core/mysql"
	"GoH/core/rabbitmq"
	"GoH/core/redis"
	"GoH/core/session"
	"flag"
	"github.com/dlintw/goconf"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"github.com/robfig/cron"
)

func main() {
	//加载配置文件
	conf := loadConfig()
	//初始化日志
	initLogger(conf)
	
	if err := error.OutError("100", "测试错误"); err != nil {
		log.Error(err)
	}
	//日志使用样例
	log.Notice("Test Notice")
	log.Debug("Test Debug")
	log.Warning("Test Warning")
	log.Error("Test Error")
	
	//初始化Redis
	redis.InitRedisPool(conf)
	
	//mysql使用样例
	mysql.InitConnectionInfo(conf)
	
	//初始化mongodb
	mongodb.InitMongoConnection(conf)
	
	//利用cpu多核来处理http/https请求
	multCpuHandler(conf)
	
	//初始化路由器
	router := NewRouter()
	
	//初始化session
	session.InitSession(conf)
	
	//初始化生产者rabbitmq
	rabbitmq.InitProducerMqConnection(conf)
	//rabbitmq样例
	//go producer.SendMQMessage("测试。。。。。")
	//初始化消费者rabbitmq
	rabbitmq.InitConsumerMqConnection(conf)
	//go consumer.ReceiveMQMessage()
	
	//初始化定时任务
	initQuartz()
	
	//启动http/https服务
	httpServerStart(conf, router)
	
}

func httpServerStart(conf *goconf.ConfigFile, router *mux.Router) {
	//判断是否为https协议
	isHttps, err := conf.GetBool("https", "is_https")
	if err != nil {
		log.Errorf("读取https配置失败，%s\n", err)
		os.Exit(1)
	} else {
		if isHttps { //如果为https协议需要配置server.crt和server.key
			serverCrt, _ := conf.GetString("https", "https_server_crt")
			serverKey, _ := conf.GetString("https", "https_server_key")
			httpsPort, _ := conf.GetInt("https", "https_port")
			log.Debug(http.ListenAndServeTLS(":"+strconv.Itoa(httpsPort), serverCrt, serverKey, router))
		} else {
			httpPort, _ := conf.GetInt("http", "http_port")
			log.Debug(http.ListenAndServe(":"+strconv.Itoa(httpPort), router))
		}
	}
}

func multCpuHandler(conf *goconf.ConfigFile) {
	isNumCPU, err := conf.GetBool("server", "use_max_cpu")
	if err != nil {
		log.Errorf("读取多核配置失败，%s\n", err)
		os.Exit(1)
	} else {
		if isNumCPU {
			runtime.GOMAXPROCS(runtime.NumCPU())
		}
	}
}

func initLogger(conf *goconf.ConfigFile) {
	logPath, err := conf.GetString("server", "log_path")
	logLevel, err := conf.GetString("server", "log_level")
	if err != nil {
		log.Errorf("日志文件配置有误, %s\n", err)
		os.Exit(1)
	}
	log.NewLogger(logPath, logLevel)
}

func loadConfig() *goconf.ConfigFile {
	conf_file := flag.String("config", "./goh.ini", "设置配置文件.")
	flag.Parse()
	conf, err := goconf.ReadConfigFile(*conf_file)
	if err != nil {
		log.Errorf("加载配置文件失败，无法打开%q，%s\n", conf_file, err)
		os.Exit(1)
	}
	return conf
}

func initQuartz() {
	i := 0
	c := cron.New()
	spec := "*/5 * * * * ?"
	c.AddFunc(spec, func() {
		i++
		log.Debug("cron running:", i)
	})
	c.Start()
}
