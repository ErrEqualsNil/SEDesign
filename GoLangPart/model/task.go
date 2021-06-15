package model

type TaskStatus int32
const (
	TaskStatusUnknown     = 0
	TaskStatusCreating = 1
	TaskStatusQueueing   = 2
	TaskStatusProcessing = 3
	TaskStatusComplete   = 4
)

// 分析任务
type Task struct {
	Id           uint64 `gorm:"column:id;unique;AUTO_INCREMENT;NOT NULL"`
	ItemName     string `gorm:"column:item_name;NOT NULL"` // 商品名
	Url          string `gorm:"column:url;NOT NULL"`       // 商品url
	Status       TaskStatus    `gorm:"column:status;NOT NULL"`    // 任务状态
	Report       string `gorm:"column:report"`             // 分析结果报告
	WordCloudUrl string `gorm:"column:word_cloud_url"`     // 词云图
	HotWords     string `gorm:"column:hot_words"`          // 高频词列表
}

func (m *Task) TableName() string {
	return "task"
}