package es

import (
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

type Conf struct {
	Ip string
	Port int
}

func GetESConn() (*elastic.Client, error) {
	conf := Conf{
		Ip: "121.36.31.113",
		Port: 9200,
	}
	url := fmt.Sprintf("http://%s:%d", conf.Ip, conf.Port)
	client, err := elastic.NewClient(
			elastic.SetURL(url),
			elastic.SetSniff(false),
		)
	if err != nil {
		log.Printf("conn: %v", url)
		log.Fatal("es get conn err: ", err)
		return nil, err
	}
	return client, nil
}