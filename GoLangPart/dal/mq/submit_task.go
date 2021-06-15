package mq

import (
	"SEDesign/model"
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
	conn, err := redis.Dial(conf.Protocol, addr)
	if err != nil {
		log.Printf("Redis get conn err: %v\n", err)
		return nil, err
	}
	return conn, err
}

type MqTask struct {
	Id uint64
	name string
}

func SubmitTask(task *model.Task) error {
	conn, err := GetConn()
	if err != nil {
		log.Printf("redis get conn err: %v\n", err)
		return err
	}
	defer conn.Close()

	_, err = conn.Do("RPUSH", "Tasks", task.Id)
	if err != nil {
		log.Printf("redis rpush task err: %v\n", err)
		return err
	}
	return nil
}