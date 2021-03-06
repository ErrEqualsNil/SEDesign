// Code generated by sql2gorm. DO NOT EDIT.
package model

type TaskStatus int32
const (
	TaskStatusUnknown     = 0
	TaskStatusCreating = 1
	TaskStatusQueueing   = 2
	TaskStatusProcessing = 3
	TaskStatusComplete   = 4
)


// 评论
type Comment struct {
	Id              int64  `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Comment         string `gorm:"column:comment"`                              // 评论全文
	Score           int    `gorm:"column:score;default:0;NOT NULL"`             // 评分
	UsefulVoteCount int    `gorm:"column:useful_vote_count;default:0;NOT NULL"` // 点赞数
	TaskId          int64  `gorm:"column:task_id;NOT NULL"`                     // 对应任务id
}

func (m *Comment) TableName() string {
	return "comment"
}


// 分析任务
type Task struct {
	Id           int64  `gorm:"column:id;unique;AUTO_INCREMENT;NOT NULL"`
	ItemName     string `gorm:"column:item_name;NOT NULL"` // 商品名
	ItemId       int64  `gorm:"column:item_id;NOT NULL"`   // 商品Id
	GoodRate     int    `gorm:"column:good_rate"`          // 好评率
	CommentCount int    `gorm:"column:comment_count"`      // 爬取的评论数
	Status       TaskStatus    `gorm:"column:status;NOT NULL"`    // 任务状态
	Report       string `gorm:"column:report"`             // 分析结果报告
	WordCloudUrl string `gorm:"column:word_cloud_url"`     // 词云图
	HotWords     string `gorm:"column:hot_words"`          // 高频词列表
}

func (m *Task) TableName() string {
	return "task"
}


type EsTask struct {
	Id     int `json:"id"`
	Name   string `json:"name"`
	ItemId int `json:"item_id"`
}