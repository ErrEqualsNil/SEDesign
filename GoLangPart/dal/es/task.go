package es

import (
	"SEDesign/model"
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"log"
	"strconv"
)

type SearchTaskReqParams struct {
	Name   string
	ItemId int64
	Id     int64
	Offset int
	Limit  int
}

const ES_Index_Name = "spider_tasks"

func AddTask(tasks []*model.Task) error {
	if len(tasks) == 0 {
		log.Printf("No tasks to sync")
		return nil
	}
	conn, err := GetESConn()
	if err != nil {
		log.Fatalf("es get conn err: %v", err)
		return err
	}
	bulkReq := conn.Bulk()
	for _, task := range tasks {
		esTask := model.EsTask{
			Id: int(task.Id),
			Name:   task.ItemName,
			ItemId: int(task.ItemId),
		}
		log.Printf("esTask create: %v", task)
		req := elastic.NewBulkIndexRequest().Index(ES_Index_Name).Id(strconv.Itoa(esTask.Id)).Doc(esTask)
		bulkReq = bulkReq.Add(req)
	}

	_, err = bulkReq.Do(context.Background())
	if err != nil {
		log.Printf("es bulkReq operate err: %v", err)
		return err
	}
	return nil
}

func SearchTaskByName(param *SearchTaskReqParams) (int64, []int64, error) {
	conn, err := GetESConn()
	if err != nil {
		log.Fatalf("es get conn err: %v", err)
		return 0, nil, err
	}
	if param == nil {
		log.Printf("req is nil")
		return 0, nil, errors.New("invalid req")
	}
	query := elastic.NewBoolQuery()
	if len(param.Name) != 0 {
		query.Filter(elastic.NewMatchPhraseQuery("name", param.Name))
	}
	if param.ItemId != 0 {
		query.Filter(elastic.NewTermsQuery("item_id", param.ItemId))
	}
	if param.Id != 0 {
		query.Filter(elastic.NewTermsQuery("id", param.Id))
	}
	resp, err := conn.Search().Index(ES_Index_Name).Query(query).From(param.Offset).Size(param.Limit).
		Sort("id", true).Do(context.Background())
	if err != nil {
		log.Printf("err to search task from es: %v", err)
		return 0, nil, err
	}
	if resp == nil || resp.Hits == nil {
		log.Printf("es resp err, resp= %+v", resp)
	}
	resultIds := make([]int64, 0)
	hist := resp.Hits
	var total int64
	if hist.TotalHits != nil {
		total = hist.TotalHits.Value
	}
	for _, doc := range hist.Hits {
		var esTask model.EsTask
		err := json.Unmarshal(doc.Source, &esTask)
		if err != nil {
			log.Printf("json unmarshall err: %v, doc: %v", err, string(doc.Source))
			continue
		}
		resultIds = append(resultIds, int64(esTask.Id))
	}
	return total, resultIds, nil
}

func DeleteTaskById(id int64) error {
	conn, err := GetESConn()
	if err != nil {
		log.Fatalf("es get conn err: %v", err)
		return err
	}
	query := elastic.NewTermsQuery("id", id)
	_, err = conn.DeleteByQuery().Index(ES_Index_Name).Query(query).Do(context.Background())
	if err != nil{
		log.Printf("es delete by query err: %v", err)
		return err
	}
	return nil
}