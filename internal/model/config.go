package model

import "time"

type Config struct {
	ID          int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ConfigKey   string    `json:"configKey" gorm:"size:100;not null;unique"`
	ConfigName  string    `json:"configName" gorm:"size:100;not null"`
	ConfigValue string    `json:"configValue" gorm:"size:500;not null"`
	ConfigType  string    `json:"configType" gorm:"size:1;not null"` // Y=系统内置, N=自定义
	Remark      string    `json:"remark" gorm:"size:500"`
	CreateTime  time.Time `json:"createTime" gorm:"autoCreateTime"`
}
