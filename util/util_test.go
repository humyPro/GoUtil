// CreatedBy Hu Min
// CreatedAt 2019/5/21 16:01
// Description
package util

import (
	"fmt"
	"log"
	"testing"
)

type A struct {
	B *B
	//A string
}

type B struct {
	Name string
	Addr []string
}

type C struct {
	B **D
}

type D struct {
	Name **string
	Addr []*string
}

func TestTransformModel(t *testing.T) {
	//n := "humin"
	//a := A{
	//	B: &B{
	//		Name: &n,
	//		Addr: nil,
	//	},
	//	//A: "start",
	//}

	//b := C{}

	bs := make(map[int8]B)
	//ds := make(map[uint64]D)
	var ds map[uint64]D
	bs[1] = B{
		Name: "A",
		Addr: []string{"成都", "武汉"},
	}
	bs[2] = B{
		Name: "B",
		Addr: []string{"北京", "伤害"},
	}
	log.Println(TransformModel(&bs, &ds))

	//bytes, _ := json.Marshal(a)
	//_ = json.Unmarshal(bytes, &b)

	fmt.Println(ds)
}
