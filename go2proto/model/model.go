package model

import (
	"github.com/jinzhu/gorm"
)

// 社区表
type Community struct {
	gorm.Model
	CommunityId    int32   //社区id
	CommunityName  string  //社区名称
	FkRegionId     int32   //行政地区id
	CoverPeopleNum int32   //覆盖人数
	HomeImageId    string  //主页图片地址
	DescImageId    string  //介绍图片地址
	DescText       string  //介绍文本
	Longitude      float64 //经度
	Latitude       float64 //纬度
}

// 便民信息表
type Convenience struct {
	gorm.Model
	ConvenienceId int32  //便民信息ID
	FkUserId      int32  //发布人ID
	UserName      string //发布人昵称
	FkConfigId    int32  //配置id
	Type          int8   //便民信息类型
	Title         string //标题
	Content       string //内容
	ImageIds      string //相关图片地址,用'|'隔开
}

//网络问政表
type WebQuestion struct {
	gorm.Model
	QuestionId int32 //题问id
	FkUserId int32 //题问人的id
	UserName string //提问人昵称
	FkConfigId int32 //什么配置id
	Type int8 //题问类型
	Title string //标题
	Content string //内容
	ImageIds int32 //图片id
}

// 便民信息评论表
type ConvenienceComment struct {
	gorm.Model
	CommentId       int32  //评论ID
	FkConvenienceId int32  //被评论的便民信息ID
	FkUserId        int32  //评论人ID
	UserName        string //评论人昵称
	Content         string //评论内容
}

// 生活圈互动表
type Interact struct {
	gorm.Model
	InteractId int32  //互动主题ID
	FkUserId   int32  //发布人id
	UserName   string //发布人昵称
	FkPlotId   int32  //发布小区id
	PlotName   string //发布小区名称
	Title      string //标题
	//TopicId    int8   //主题类型ID
	//TopicName  string //主题类型
	Content    string //内容
	ImageIds   string //相关图片地址
	ViewNum    int32  //浏览量
	CommentNum int32  //评论数
}

//生活圈互动评论表
type InteractComment struct {
	gorm.Model
	CommentId  int32  //评论id
	InteractId int32  //被评论的主题id
	FkUserId   int32  //评论id
	UserName   string //评论人昵称
	Content    string //内容
}

//通知表
type Notice struct {
	gorm.Model
	FkPlotId      int32  //通知小区id
	PlotName      string //通知小区名称
	CommunityId   int32  //社区id
	CommunityName string //社区名称
	FkUserId      int32  //发布人id
	UserName      string //发布人姓名
	Inscribe      string //落款单位
	Type          int8   //类型
	Content       string //内容
	Status        int8   //状态
}

// 问政评论表
type QuestionComment struct {
	gorm.Model
	CommentId    int32  //评论id
	FkUserId     int32  //评论人id
	UserName     string //评论人姓名
	FkQuestionId int32  //被评论的问题id
	Content      string //内容
}

//物业表
type Estate struct {
	gorm.Model
	EstateId   int32  //物业id
	EstateName string //物业名称
}

// 小区表
type Plot struct {
	gorm.Model
	PlotId         int32   // 小区id
	PlotName       string  //小区名称
	FkRegionId     int32   //所属行政区域id
	FkEstateId     int32   //物业公司ID
	EstateName     string  //物业公司名称
	FkCommunityId  int32   //所属社区id
	CommunityName  string  //所属社区名称
	CoverPeopleNum int32   //覆盖人群
	HomeImageId    int32   //主页图片id
	DescImageId    int32   //介绍图片id
	DescText       string  //介绍文本
	Longitude      float64 //经度
	Latitude       float64 //纬度
}

//行政区域表
type Region struct {
	gorm.Model
	Name string
	Code string  //通用码
	PostCode string //邮政编码
	ParentId int
	IsTopLevel int //是否为最高级,用于嵌套查询
}

// 用户认证表
type UserIdentify struct {
	gorm.Model
	UserIdentifyId int32  //认证信息ID
	FkUserId       int32  //被认证的用户id
	UserName       string //用户姓名
	Status         string //状态
	Type           string //类型
	TelNumber      string //电话号码
	Province       string //省
	City           string //市
	District       string //区
	Community      string //社区
	Plot           string //小区
}
