package model

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type BaseModel struct {
	ID        uint64                `gorm:"primarykey;" json:"id"`     // 主键
	CreatedAt time.Time             `gorm:"autoCreateTime" json:"-"`   // 创建时间
	CreatedBy string                `json:"-"`                         // 创建人
	UpdatedAt time.Time             `gorm:"autoUpdateTime" json:"-"`   // 修改时间
	UpdatedBy string                `json:"-"`                         // 修改人
	Version   int16                 `json:"column:version"`            // 乐观锁版本号
	IsDeleted soft_delete.DeletedAt `gorm:"softDelete:flag" json:"-" ` // 删除时间 库表中对应字段 is_deleted bool
}
