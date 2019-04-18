package a2b

import (
	"errors"
	"log"
	"reflect"
)

var typeError = errors.New("a和b的数据类型不一致，请核对其底层类型是否一致，比如整数类型和整数类型，结构体和结构体,数组和切片")

func A2B(a interface{}, b interface{}) error {
	bv := reflect.ValueOf(b)
	if bv.Kind() != reflect.Ptr || bv.IsNil() {
		return errors.New("b对象必须为指针类型")
	}
	av := reflect.ValueOf(a)

	if !isSimilarType(av.Kind(), bv.Kind()) {
		log.Printf("不能在%s和%s之间赋值\n", av.Kind(), bv.Kind())
		return typeError
	}

	deepCopy(av, bv)
	return nil
}

func deepCopy(av reflect.Value, bv reflect.Value) {
	av = getRealValue(av)
	bv = getRealValue(bv)
	if !isSimilarType(av.Kind(), bv.Kind()) {
		log.Printf("不能在%s和%s之间赋值\n", av.Kind(), bv.Kind())
	}
	if isString(bv.Kind()) {
		if bv.CanSet() {
			bv.Set(av)
		}
		return
	}

	if isFloat(bv.Kind()) {
		if bv.CanSet() {
			bv.SetFloat(av.Float())
		}
		return
	}

	if isBool(bv.Kind()) {
		if bv.CanSet() {
			bv.SetBool(av.Bool())
		}
		return
	}
	if isMap(bv.Kind()) {
		//如果键类型不一样，返回
		if bv.Type().Key() != av.Type().Key() {
			return
		}
		if bv.CanSet() {
			keys := av.MapKeys()
			for _, key := range keys {
				//todo 如果值的类型不一样
				bv.SetMapIndex(key, av.MapIndex(key))
			}
		}
		return
	}

	//2-6是有符号整数，7-11是无符号整数
	if isInteger(bv.Kind()) {
		if bv.CanSet() {
			if isInt(bv.Kind()) {
				if isInt(av.Kind()) {
					bv.SetInt(av.Int())
				} else if isUint(av.Kind()) {
					bv.SetInt(int64(av.Uint()))
				}
			} else if isUint(bv.Kind()) {
				if isInt(av.Kind()) {
					bv.SetUint(uint64(av.Int()))
				} else if isUint(av.Kind()) {
					bv.SetUint(uint64(av.Uint()))
				}
			}
		}
		return
	}

}

//返回底层value，如果不是字符串，数字，bool，结构体，切片数组等，可能出现未知bug
func getRealValue(in reflect.Value) reflect.Value {
	if in.Kind() == reflect.Ptr || in.Kind() == reflect.Interface {
		for {
			if in.Kind() != reflect.Ptr && in.Kind() != reflect.Interface {
				return in
			}
			in = in.Elem()
		}
	}
	return in
}

//获取结构体或者结构体指针的字段名
func getFieldMap(a interface{}) map[string]struct{} {

	at := reflect.TypeOf(a)
	if at.Kind() == reflect.Ptr {
		at = at.Elem()
	}

	fieldMap := make(map[string]struct{})

	c := at.NumField()
	for i := 0; i < c; i++ {
		f := at.Field(i)
		//不记录匿名字段
		if f.Anonymous {
			continue
		}
		fieldMap[f.Name] = struct{}{}
	}
	return fieldMap
}

func isSimilarType(a reflect.Kind, b reflect.Kind) bool {

	return (isInteger(a) && isInteger(b)) || (isBool(a) && isBool(b)) || (isFloat(a) && isFloat(b)) || (isArrayOrSlice(a) && isArrayOrSlice(b)) || (isMap(a) && isMap(b)) || (isString(a) && isString(b)) || (isStruct(a) && isStruct(b))
}

//只能是这几个类型，才能设置其值
func isCanSetValue(kind reflect.Kind) bool {
	return isInteger(kind) || isStruct(kind) || isString(kind) || isMap(kind) || isArrayOrSlice(kind) || isBool(kind) || isFloat(kind)
}

func isInteger(kind reflect.Kind) bool {
	return kind >= 2 && kind <= 11
}

func isInt(kind reflect.Kind) bool {
	return kind >= 2 && kind <= 6
}

func isUint(kind reflect.Kind) bool {
	return kind >= 7 && kind <= 11
}

func isBool(kind reflect.Kind) bool {
	return kind == 1
}

func isFloat(kind reflect.Kind) bool {
	return kind == 13 || kind == 14
}

func isArrayOrSlice(kind reflect.Kind) bool {
	return kind == 17 || kind == 23
}

func isMap(kind reflect.Kind) bool {
	return kind == 21
}

func isString(kind reflect.Kind) bool {
	return kind == 24
}

func isStruct(kind reflect.Kind) bool {
	return kind == 25
}
