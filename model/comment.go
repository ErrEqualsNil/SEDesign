package model

// 评论
type Comment struct {
	Id       uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ItemName string `gorm:"column:item_name;NOT NULL"` // 商品名
	Comment  string `gorm:"column:comment"`            // 评论全文
	TaskId   uint64 `gorm:"column:task_id;NOT NULL"`   // 对应任务id
}

func (m *Comment) TableName() string {
	return "comment"
}
