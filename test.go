package main

import (
	"encoding/json"
	"fmt"
	"github.com/humyPro/GoUtil/a2b/model"
	"github.com/humyPro/GoUtil/v2"
	"time"
)

func main() {
	user := model.User{
		Openid:      "213123",
		WxId:        "123123",
		HeadImage:   "12321312",
		NickName:    "asdasd",
		PhoneNumber: "asd123213",
		Gender:      1,
		Type:        0,
		OriginJson:  "asdasd",
	}
	user.ID = 12312312
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	now := time.Now()
	user.DeletedAt = &now

	var users []model.User
	users = append(users, user)
	users = append(users, user)
	users = append(users, user)
	users = append(users, user)
	users = append(users, user)

	//page := model.Page{
	//	PageNum:  12,
	//	PageSize: 13,
	//	Total:    144,
	//	Data:     users,
	//}

	bytes, _ := json.Marshal(users)
	fmt.Println(string(bytes))

	//page2 := model.Page{
	//	Data: []model.Userx{},
	//}
	x := []model.Userx{}
	_ = v2.TrandformModel(&users, &x)

	bytes, _ = json.Marshal(x)
	fmt.Println(string(bytes))

	//page3 := model.Page{
	//	Data: []model.User{},
	//}
	p := []model.User{}
	b := v2.TrandformModel(&x, &p)
	fmt.Println(b)
	bytes, _ = json.Marshal(p)
	fmt.Println(string(bytes))
	//
	//i := reflect.TypeOf(user).NumField()
	//for a:=0;a<i;a++{
	//	fmt.Println(reflect.TypeOf(user).Field(a).Name)
	//}
	//
	//fmt.Println("---------------------")
	//field, b := reflect.TypeOf(user).FieldByName("CreatedAt")
	//fmt.Println(b)
	//fmt.Println(field.Name)
	//fmt.Println(field.Type.Name())
	//
	//fmt.Println("-----------")
	//
	//structField, i2 := reflect.TypeOf(user).FieldByName("Model")
	//fmt.Println(i2)
	//fmt.Println(structField)

	//a := A{}
	//o := reflect.TypeOf(a).NumField()
	//for i := 0; i < o; i++ {
	//	field := reflect.TypeOf(a).Field(i)
	//	fmt.Println("-------------------")
	//	fmt.Println(field.Name + ":" + strconv.FormatBool(field.Anonymous))
	//
	//	if field.Anonymous {
	//		x := reflect.ValueOf(a).FieldByName(field.Name)
	//		numField := x.NumField()
	//		for i := 0; i < numField; i++ {
	//			field := x.Type().Field(i)
	//			fmt.Println("-------------------")
	//			fmt.Println(field.Name + ":" + strconv.FormatBool(field.Anonymous))
	//		}
	//	}
	//}
	//
	//s := "哈dasd"
	//fmt.Println(len([]rune(s)))
}
