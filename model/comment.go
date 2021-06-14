package model

type TaskStatus int32
const (
	TaskstatusUnknown     = 0
	TaskStatusQueueing   = 1
	TaskStatusProcessing = 2
	TaskStatusComplete   = 3
)
// 评论表
type Comment struct {
	Id          uint64 `gorm:"column:id;primary_key;unique;AUTO_INCREMENT"`
	ItemName    string `gorm:"column:item_name;NOT NULL"`        // 商品名
	Comment     string `gorm:"column:comment"`                   // 评论全文
	HotWordList string `gorm:"column:hot_word_list"`             // 高频词列表
	Status      TaskStatus    `gorm:"column:status;default:0;NOT NULL"` // 状态; 0.Unknown 1.爬取中 2.统计中 3.完成
}

func (m *Comment) TableName() string {
	return "comment"
}

