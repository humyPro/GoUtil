package model

import "github.com/jinzhu/gorm"

/**
//设备当前状态
*/
type DeviceCurrentStatus struct {
	gorm.Model                  //id为设备id
	ApkVersion           int    //APK版本号
	PushVersion          int    // 推送批次
	CurrentAdvertiseList string `gorm:"type:text"` //当前播放列表
	CurrentNoticeList    string `gorm:"type:text"` //当前通知公告列表
	Status               int    //状态
}
