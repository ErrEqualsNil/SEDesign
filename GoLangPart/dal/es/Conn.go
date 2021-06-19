package es

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type EsConf struct {
	Ip string
	Port int
}

func GetESConn() (*elastic.Client, error) {
	data, err := ioutil.ReadFile("conf/es_conf.yml")
	if err != nil {
		log.Fatalf("read es_conf err: %v", err)
		return nil, err
	}
	conf := EsConf{}
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		log.Fatalf("yaml unmarshal err: %v", err)
		return nil, err
	}
	url := fmt.Sprintf("http://%s:%d", conf.Ip, conf.Port)
	client, err := elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetSniff(false),
		)
	if err != nil {
		log.Fatal("es get conn err: ", err)
		return nil, err
	}
	return client, nil
}