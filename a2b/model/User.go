package model

import (
	"github.com/humyPro/GoUtil/a2b/model/gorm"
)

/**
用户表
当前只有小程序用户
*/
type User struct {
	gorm.Model
	Openid      string //openid
	WxId        string `gorm:"index"` //微信id
	HeadImage   string //头像
	NickName    string //昵称
	PhoneNumber string `gorm:"index"`           //电话号码
	Gender      int    `gorm:"type:tinyint(1)"` //1 男 2 女
	Type        int    //类型
	OriginJson  string `gorm:"type:text"` //微信返回的原始数据
}

type Userx struct {
	ID                   uint32   `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	CreatedAt            uint64   `protobuf:"varint,2,opt,name=CreatedAt,proto3" json:"CreatedAt,omitempty"`
	UpdatedAt            uint64   `protobuf:"varint,3,opt,name=UpdatedAt,proto3" json:"UpdatedAt,omitempty"`
	DeletedAt            uint64   `protobuf:"varint,4,opt,name=DeletedAt,proto3" json:"DeletedAt,omitempty"`
	Openid               string   `protobuf:"bytes,5,opt,name=Openid,proto3" json:"Openid,omitempty"`
	WxId                 string   `protobuf:"bytes,12,opt,name=WxId,proto3" json:"WxId,omitempty"`
	HeadImage            string   `protobuf:"bytes,6,opt,name=HeadImage,proto3" json:"HeadImage,omitempty"`
	NickName             string   `protobuf:"bytes,7,opt,name=NickName,proto3" json:"NickName,omitempty"`
	PhoneNumber          string   `protobuf:"bytes,8,opt,name=PhoneNumber,proto3" json:"PhoneNumber,omitempty"`
	Gender               int32    `protobuf:"varint,9,opt,name=Gender,proto3" json:"Gender,omitempty"`
	Type                 int32    `protobuf:"varint,10,opt,name=Type,proto3" json:"Type,omitempty"`
	OriginJson           string   `protobuf:"bytes,11,opt,name=OriginJson,proto3" json:"OriginJson,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type Page struct {
	PageNum  int
	PageSize int
	Total    int
	Data     User
}
type PageX struct {
	PageNum  int
	PageSize int
	Total    int
	Data     Userx
}
