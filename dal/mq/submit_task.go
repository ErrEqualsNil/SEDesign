package mq

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type RedisConf struct {
	Ip       string
	Port     string
	Protocol string
}

func GetConn() (redis.Conn, error) {
	data, err := ioutil.ReadFile("conf/redis_conf.yml")
	if err != nil {
		log.Printf("read redis_conf err: %v\n", err)
		return nil, err
	}
	conf := RedisConf{}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Printf("yaml unmarshal err: %v\n", err)
		return nil, err
	}
	addr := fmt.Sprintf("%s:%s", conf.Ip, conf.Port)
	log.Printf("redis conn, protocol: %v, addr: %v", conf.Protocol, addr)
	conn, err := redis.Dial(conf.Protocol, addr)
	if err != nil {
		log.Printf("Redis get conn err: %v\n", err)
		return nil, err
	}
	return conn, err
}

type MqTaskParam struct {
	Id uint64
	Name string
}

func SubmitTask(task *MqTaskParam) error {
	conn, err := GetConn()
	if err != nil {
		log.Printf("redis get conn err: %v\n", err)
		return err
	}
	defer conn.Close()

	data, err := json.Marshal(task)
	if err != nil {
		log.Printf("json unmarshall task err: %v\n", err)
		return err
	}
	_, err = conn.Do("RPUSH", "TaskMQ", data)
	if err != nil {
		log.Printf("redis rpush task err: %v\n", err)
		return err
	}
	return nil
}