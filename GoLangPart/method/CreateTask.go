package method

import (
	"SEDesign/dal/cache"
	"SEDesign/dal/db"
	"SEDesign/dal/es"
	"SEDesign/dal/mq"
	"SEDesign/model"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type CreateTaskReqParam struct {
	Name string `form:"name" json:"name"`
	ItemId int64 `form:"itemId" json:"itemId"`
}

type CreateTaskHandler struct {
	Ctx *gin.Context
	req CreateTaskReqParam
}

func (handler CreateTaskHandler) checkValid () bool {
	if len(handler.req.Name) == 0 && handler.req.ItemId == 0{
		log.Printf("req Illegal, req: %v", handler.req)
		return false
	}
	return true
}

func (handler CreateTaskHandler) Run () error {
	err := handler.Ctx.ShouldBind(&handler.req)
	if err != nil {
		log.Printf("Invalid err: %v", err)
		return errors.New("invalid req")
	}
	//参数检查
	if !handler.checkValid() {
		log.Printf("Invalid params")
		return errors.New("invalid params")
	}

	//写入mysql
	tasks, err := handler.fillTask()
	if err != nil {
		log.Printf("fill task err: %v", err)
		return err
	}
	for _, task := range tasks{
		err := db.CreateTask(task)
		if err != nil {
			log.Printf("db create comment err: %v", err)
			return err
		}

		//写入mq
		err = mq.SubmitTask(task)
		if err != nil {
			log.Printf("mq submit task err: %v", err)
			return err
		}
	}
	err = es.AddTask(tasks)
	if err != nil {
		log.Printf("es add task err: %v", err)
		return err
	}

	handler.Ctx.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"resp": "create task success!",
	})
	return nil
}

func (handler CreateTaskHandler) fillTask() ([]*model.Task, error) {
	if len(handler.req.Name) != 0 && handler.req.ItemId != 0 {
		exist, err := cache.CheckTaskExists(handler.req.ItemId)
		if err != nil{
			log.Printf("cache CheckTaskExists err: %v", err)
		}
		if exist{
			log.Printf("task exist, itemId: %v", handler.req.ItemId)
			return nil, nil
		}
		err = cache.AddTaskToCache(handler.req.ItemId)
		if err != nil {
			log.Printf("cache AddTaskToCache err: %v", err)
		}
		task := &model.Task{
			ItemName: handler.req.Name,
			ItemId: handler.req.ItemId,
			Status: model.TaskStatusQueueing,
			CommentCount: 0,
		}
		return []*model.Task{task}, nil
	}
	if len(handler.req.Name) != 0 {
		encode_name := url.QueryEscape(handler.req.Name)
		searchUrl := fmt.Sprintf("https://search.jd.com/Search?keyword=%v&enc=utf-8", encode_name)
		client := &http.Client{}
		req, err := http.NewRequest("GET", searchUrl, nil)
		if err != nil{
			log.Printf("http err: %v", err)
			return nil, err
		}
		req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
		resp, err := client.Do(req)
		if err != nil{
			log.Printf("http err: %v", err)
			return nil, err
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil{
			log.Printf("goquery new doc err: %v", err)
			return nil, err
		}
		itemIds := make([]int64, 0)
		doc.Find(".p-img a").Each(func(i int, selection *goquery.Selection) {
			tmp, ok := selection.Attr("href")
			if !ok {
				log.Printf("Cannot find href at selection: %s", selection.Text())
				return
			}
			reg, err := regexp.Compile("(\\d\\d*)")
			if err != nil{
				log.Printf("regexp compile err: %v", err)
				return
			}
			matchgroup := reg.FindAllStringIndex(tmp, 1)[0]
			id_string := tmp[matchgroup[0] : matchgroup[1]]
			id, err := strconv.Atoi(id_string)
			if err != nil {
				log.Printf("string to int err, string: %v", id_string)
				return
			}
			itemIds = append(itemIds, int64(id))
		})
		itemIds = itemIds[0:10]
		result := make([]*model.Task, 0)
		for _, itemId := range itemIds {
			exist, err := cache.CheckTaskExists(itemId)
			if err != nil{
				log.Printf("cache CheckTaskExists err: %v", err)
			}
			if exist{
				log.Printf("task exist, itemId: %v", itemId)
				continue
			}
			err = cache.AddTaskToCache(itemId)
			if err != nil {
				log.Printf("cache AddTaskToCache err: %v", err)
			}
			name, err := GetItemNameByItemId(itemId)
			if err != nil {
				log.Printf("GetItemNameByItemId err: %v", err)
				continue
			}
			task := &model.Task{
				ItemName: name,
				ItemId: itemId,
				Status: model.TaskStatusQueueing,
				CommentCount: 0,
			}
			result = append(result, task)
		}
		return result, nil
	}

	if handler.req.ItemId != 0 {
		exist, err := cache.CheckTaskExists(handler.req.ItemId)
		if err != nil{
			log.Printf("cache CheckTaskExists err: %v", err)
		}
		if exist{
			log.Printf("task exist, itemId: %v", handler.req.ItemId)
			return nil, nil
		}
		err = cache.AddTaskToCache(handler.req.ItemId)
		if err != nil {
			log.Printf("cache AddTaskToCache err: %v", err)
		}
		name, err := GetItemNameByItemId(handler.req.ItemId)
		if err != nil {
			log.Printf("GetItemNameByItemId err: %v", err)
			return nil, err
		}
		task := &model.Task{
			ItemName: name,
			ItemId: handler.req.ItemId,
			Status: model.TaskStatusQueueing,
			CommentCount: 0,
		}
		return []*model.Task{task}, nil
	}
	return nil, errors.New("Invalid Param")
}


func GetItemNameByItemId(id int64) (string, error) {
	searchUrl := fmt.Sprintf("https://item.jd.com/%v.html#comment", id)
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil{
		log.Printf("http err: %v", err)
		return "", err
	}
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil{
		log.Printf("http err: %v", err)
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil{
		log.Printf("goquery new doc err: %v", err)
		return "", err
	}
	name := strings.TrimSpace(doc.Find(".sku-name").Text())
	return name, nil
}