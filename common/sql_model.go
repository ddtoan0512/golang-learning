package common

import "time"

type SQLModel struct {
	Id        int        `json:"id" gorm:"column:id"`
	Status    int        `json:"status" gorm:"column:status;default:1"`
	CreatedAt *time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:" ,omitempty" gorm:"column: "`
}

//func (m *SQLModel) GenUID(dbType int) {
//	uid := NewUID(uint32(m.Id), dbType, 1)
//	m.FakeId = &uid
//}
