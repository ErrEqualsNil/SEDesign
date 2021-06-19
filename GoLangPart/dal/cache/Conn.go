package cache

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type RedisConf struct {
	Ip       string
	Port     string
	Protocol string
}

func GetRedisConn() (redis.Conn, error) {
	data, err := ioutil.ReadFile("conf/redis_conf.yml")
	if err != nil {
		log.Fatalf("read redis_conf err: %v\n", err)
		return nil, err
	}
	conf := RedisConf{}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("yaml unmarshal err: %v\n", err)
		return nil, err
	}
	addr := fmt.Sprintf("%s:%s", conf.Ip, conf.Port)
	conn, err := redis.Dial(conf.Protocol, addr)
	if err != nil {
		log.Fatalf("Redis get conn err: %v\n", err)
		return nil, err
	}
	return conn, err
}
