package main

import (
	"fmt"
	"reflect"
)

func main() {
	x := X{}

	//m := reflect.ValueOf(&x).Elem().FieldByName("Map")
	//
	//fmt.Println(m.IsNil())
	//fmt.Println(m.IsValid())
	//
	//v := reflect.MakeMap(reflect.MapOf(m.Type().Key(),m.Type().Elem()))
	//m.Set(v)
	//m.SetMapIndex(reflect.ValueOf("asd"),reflect.ValueOf(1))
	//fmt.Println(x)
	//fmt.Println(m.IsNil())
	name := reflect.ValueOf(&x).Elem().FieldByName("XX")
	fmt.Println(name.IsNil())
	fmt.Println(name.IsValid())
	fmt.Println(name.Type())
	v := reflect.New(reflect.TypeOf(name.Elem()))
	fmt.Println(v.Type())
	reflect.ValueOf(name).Elem().Set(v)
}

func f1(in interface{}) {
	fmt.Println(reflect.ValueOf(in).Elem().IsNil())

	bv := reflect.ValueOf(in).Elem().Elem()
	if !bv.IsValid() {
		bv = (reflect.New(reflect.TypeOf(bv)))
		fmt.Println("in")
		fmt.Println(bv)
	}
	fmt.Println("end")

}

type X struct {
	Map map[string]int
	XX  *XX
}
type XX struct {
	XX string
}
