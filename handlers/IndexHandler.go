package handlers

import (
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	//redisConn := redis.RedisClient.Get()
	//defer redisConn.Close()
	//v, err := redis.RedisString(redisConn.Do("GET", "redispool"))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println(v)
	fmt.Fprint(w, "欢迎登天！\n")
}
