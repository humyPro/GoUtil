package util

import (
	"errors"
	"reflect"
	"time"
)

//gorm.model中的字段
const (
	created = "CreatedAt"
	updated = "UpdatedAt"
	deleted = "Undefined"
)

//时间类型
var timeType = reflect.TypeOf(time.Now())

// 该方法会从a中取值,给b中同名的字段并且是可以被赋值的类型的字段赋值,会忽略掉同名但不可赋值的字段
//version-1:修复时间问题,如果数据源的时间是0,那么目标的时间类型为零值
//version-2:所有id-Id-ID都可以匹配
//version-3:支持切片和多级指针,支持map,以及各种嵌套
func TransformModel(a interface{}, b interface{}) error {
	bv := reflect.ValueOf(b)
	if bv.Kind() != reflect.Ptr || bv.IsNil() {
		return errors.New("b对象必须为指针类型")
	}
	av := reflect.ValueOf(a)
	if av.Kind() != reflect.Ptr || bv.IsNil() {
		return errors.New("a对象必须为指针类型")
	}
	bv = bv.Elem()
	av = av.Elem()
	deepCopy(av, bv)
	return nil
}

func deepCopy(av reflect.Value, bv reflect.Value) {
	//如果av是nil,返回
	switch av.Kind() {
	case reflect.Slice, reflect.Map, reflect.Chan, reflect.Interface, reflect.Ptr:
		if av.IsNil() {
			return
		}
	}
	av = getRealValue(av, false)
	bv = getRealValue(bv, false)
	//结构体类型
	if isStruct(bv.Kind()) && isStruct(av.Kind()) {
		bf := getFieldMap(bv)
		for key := range bf {
			//特殊处理id
			fa := av.FieldByName(key)

			if key == "ID" || key == "Id" || key == "id" {
				fa = av.FieldByName("ID")
				if !fa.IsValid() {
					fa = av.FieldByName("Id")
					if !fa.IsValid() {
						fa = av.FieldByName("id")
					}
				}
			}

			if !fa.IsValid() {
				continue
			}
			fb := bv.FieldByName(key)
			if key == created || key == updated || key == deleted {
				//特殊处理 gorm.model和grpc的时间类型
				dealTime(fa, fb)
			} else {
				deepCopy(fa, fb)
			}
		}
		return

	}
	//字符串类型
	if isString(bv.Kind()) && isString(av.Kind()) {
		if bv.CanSet() {
			bv.Set(av)
		}
		return
	}

	//浮点数类型
	if isFloat(bv.Kind()) && isFloat(av.Kind()) {
		if bv.CanSet() {
			bv.SetFloat(av.Float())
		}
		return
	}

	//BOOL
	if isBool(bv.Kind()) && isBool(av.Kind()) {
		if bv.CanSet() {
			//bv.SetBool(av.Bool())
			bv.Set(av)
		}
		return
	}

	//2-6是有符号整数，7-11是无符号整数
	if isInteger(bv.Kind()) && isInteger(av.Kind()) {
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
	//数组或者切片(其实只能是切片)
	if isArrayOrSlice(bv.Kind()) && isArrayOrSlice(av.Kind()) {
		if bv.CanSet() {
			sLen := av.Len()
			if bv.IsNil() {
				bv.Set(reflect.MakeSlice(reflect.SliceOf(bv.Type().Elem()), sLen, sLen))
			}
			//切片,不考虑目标的长度
			if bv.Type().Elem() == av.Type().Elem() {
				reflect.Copy(bv, av)
				/*} else if bv.Type().Elem().Kind() == reflect.Struct && av.Type().Elem().Kind() == reflect.Struct {*/
			} else {
				bv.SetLen(0)
				for i := 0; i < sLen; i++ {
					v1 := av.Index(i)
					v2 := reflect.New(bv.Type().Elem())
					deepCopy(v1, v2)
					bv.Set(reflect.Append(bv, v2.Elem()))
				}
			}
		}
		return
	}
	//处理BV是interface的情况,如果av是map或者struct,用map装,其他情况用av的类型

	//都是map的情况
	if isMap(av.Kind()) && isMap(bv.Kind()) {
		if bv.CanSet() {
			if bv.IsNil() {
				bv.Set(reflect.MakeMap(reflect.MapOf(bv.Type().Key(), bv.Type().Elem())))
			}
			keys := av.MapKeys()
			//ak: 原key值对象
			t := av.Type().Key()
			//只有这些类型 可以做key
			t2 := bv.Type().Key()
			if (isInteger(t.Kind()) && isInteger(t2.Kind())) || (isString(t.Kind()) && isString(t2.Kind())) {
				elem := bv.Type().Elem()
				for _, k2 := range keys {
					k1 := reflect.New(t2)
					deepCopy(k2, k1)
					v2 := av.MapIndex(k2)
					v1 := reflect.New(elem)
					deepCopy(v2, v1)
					bv.SetMapIndex(k1.Elem(), v1.Elem())
				}
			} else {
				panic("暂时只有整数和字符串类型可以做map的key值")
			}
		}
		return
	}

}

func dealTime(av reflect.Value, bv reflect.Value) {
	//如果从数据库模型转换成pb模型,就是time.time->其他类型(可能是整形int64或者uint64,也可能是string)
	//反过来也一样
	//2*2至少有四种类型

	if !bv.CanSet() {
		return
	}
	if av.Type() == timeType {
		//如果都是时间类型,直接设置
		if bv.Type() == timeType {
			bv.Set(av)
			return
		}

		//判断是不是合法时间
		isZero := av.MethodByName("IsZero")
		if !isZero.IsValid() {
			return
		}
		if isZero.Call(nil)[0].Bool() {
			return
		}
		unixNano := av.MethodByName("UnixNano")
		if !unixNano.IsValid() {
			return
		}
		//获得毫秒
		timeInt := unixNano.Call(nil)[0].Int() / 1e6
		//如果目标类型是数值类型

		if isInt(bv.Kind()) {
			bv.SetInt(timeInt)
			return
		}

		if isUint(bv.Kind()) {
			bv.SetUint(uint64(timeInt))
			return
		}
	} else if isInteger(av.Kind()) {
		//如果来源是数字类型
		if isInt(av.Kind()) {
			timeInt := av.Int()
			if timeInt <= 0 {
				return
			}
			if bv.Type() == timeType {
				t := time.Unix(timeInt/1e3, (timeInt%1e3)*1e6)
				bv.Set(reflect.ValueOf(t))
				return
			}
			if isInt(bv.Kind()) {
				bv.SetInt(timeInt)
				return
			}
			if isUint(bv.Kind()) {
				bv.SetUint(uint64(timeInt))
				return
			}
		} else if isUint(av.Kind()) {
			timeUint := av.Uint()
			if timeUint <= 0 {
				return
			}
			if bv.Type() == timeType {
				t := time.Unix(int64(timeUint/1e3), int64((timeUint%1e3)*1e6))
				bv.Set(reflect.ValueOf(t))
				return
			}
			if isInt(bv.Kind()) {
				bv.SetInt(int64(timeUint))
				return
			}
			if isUint(bv.Kind()) {
				bv.SetUint(timeUint)
				return
			}
		}
	}
}

//获取结构体或者结构体指针的字段名
func getFieldMap(at reflect.Value) map[string]struct{} {
	fieldMap := make(map[string]struct{})
	c := at.NumField()
	for i := 0; i < c; i++ {
		f := at.Type().Field(i)
		//不记录非导出字段

		if ([]rune(f.Name)[0] >= 'a' && []rune(f.Name)[0] <= 'z') || []rune(f.Name)[0] == '_' {
			continue
		}
		//匿名字段
		if f.Anonymous {
			fMap := getFieldMap(at.FieldByName(f.Name))
			for key := range fMap {
				fieldMap[key] = struct{}{}
			}
			continue
		}
		fieldMap[f.Name] = struct{}{}
	}
	return fieldMap
}

//返回底层value，参考了Json包下的代码,如果是nil,则会为其创建一个空的value值
func getRealValue(v reflect.Value, decodingNull bool) reflect.Value {
	v0 := v
	haveAddr := false
	// Load value from interface, but only if the result will be
	// usefully addressable.
	for {
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			v = e
			//JSON包下只取到指针一级,这儿要取到指针最终指向的数据
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		//如果decodingNull=true,那么如果一个值是nil指针类型嘛,因为v.Elem().Kind() ==Invalid,这儿就会退出
		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}
		//
		if v.IsNil() {
			//如果是nil。则新建一个Zero value(ptr)类型,取elem赋值给v
			v.Set(reflect.New(v.Type().Elem()))
		}

		if haveAddr {
			v = v0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			v = v.Elem()
		}

	}
	return v
}

//所有整数
func isInteger(kind reflect.Kind) bool {
	return kind >= 2 && kind <= 11
}

//有符号整数
func isInt(kind reflect.Kind) bool {
	return kind >= 2 && kind <= 6
}

//无符号整数
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
