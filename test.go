package main

import (
	"fmt"
	"reflect"
)

func main() {
	as := map[string]A{}

	name := reflect.TypeOf(as).Elem().Name()
	s := reflect.TypeOf(as).Key()
	fmt.Println(name)
	fmt.Println(s)

	elem := reflect.ValueOf(&as).Elem()
	i := &A{}
	v := reflect.ValueOf(i)
	a := A{Aa: "aaa"}
	v.Set(reflect.ValueOf(a))
	elem.SetMapIndex(reflect.ValueOf("aaa"),v)

	fmt.Println(as)

}

type A struct {
	Aa string
	B  int
	C *C
}

type C struct {
	C string
}

//传进f1的参数不能取地址，但是在里面取地址，又是*interface类型
func f1(a interface{}) {
	reflect.ValueOf(a).Elem().Elem().FieldByName("Aa").SetString("asdasdas")

}
