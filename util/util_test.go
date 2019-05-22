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
	A string
}

type B struct {
	Name string
	Addr []string
}

type C struct {
	B **D
	A *****string
}

type D struct {
	Name **string
	Addr *[]*string
}

func TestTransformModel(t *testing.T) {

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

func TestTransformModel2(t *testing.T) {
	a := A{
		B: &B{
			Name: "humin",
			Addr: []string{"成都", "上海"},
		},
		//A: "start",
	}

	b := C{}
	e := TransformModel(&a, &b)
	if e != nil {
		t.Error(e)
	}

	t.Log(b)
}

func TestTransformModel3(t *testing.T) {
	var a interface{}
	a = 1

	var b uint64

	_ = TransformModel(&a, &b)
	t.Log(b)
}

func TestTransformModel4(t *testing.T) {
	var a = &A{
		B: &B{
			Name: "testB",
			Addr: []string{"成都", "上海"},
		},
		A: "testA",
	}

	var c C
	var d = &a

	t.Log(TransformModel(&d, &c))
	t.Log(c)
}

func TestTransformModel5(t *testing.T) {
	s1 := []string{"213", "12312"}

	var s2 []string

	t.Log(TransformModel(&s1, &s2))

	t.Log(s2)

}
