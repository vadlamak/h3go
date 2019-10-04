package utils

import "github.com/gistao/RedisGo-Async/redis"

func GetConn() (redis.AsynConn,error) {
//create conn to standalone redis
return redis.AsyncDial("tcp", ":6379")

}

func CloseConn(c redis.AsynConn) {
defer c.Close()
}