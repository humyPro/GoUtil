package v2

import (
	"errors"
	"reflect"
	"time"
)

const (
	created = "CreatedAt"
	updated = "UpdatedAt"
	deleted = "Undefined"
)

var timeType = reflect.TypeOf(time.Now())

// 目前只支持最简单的表结构以及其切片类型,
func TrandformModel(a interface{}, b interface{}) error {
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

	if bv.Kind() == reflect.Struct && av.Kind() == reflect.Struct {
		deepCopy(av, bv)
		return nil
	} else if bv.Kind() == reflect.Slice && av.Kind() == reflect.Slice {
		deepCopy(av, bv)
		return nil
	}

	return errors.New("a和b的类型必须都是结构体或者切片")

}
func deepCopy(av reflect.Value, bv reflect.Value) {
	av = getRealValue(av)
	bv = getRealValue(bv)
	//结构体类型
	if isStruct(bv.Kind()) {
		bf := getFieldMap(bv)
		for key := range bf {
			fa := av.FieldByName(key)
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
	if isString(bv.Kind()) {
		if bv.CanSet() {
			bv.Set(av)
		}
		return
	}

	//浮点数类型
	if isFloat(bv.Kind()) {
		if bv.CanSet() {
			bv.SetFloat(av.Float())
		}
		return
	}

	//BOOL
	if isBool(bv.Kind()) {
		if bv.CanSet() {
			bv.SetBool(av.Bool())
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
	//数组或者切片(其实只能是切片)
	if isArrayOrSlice(bv.Kind()) {
		if bv.CanSet() {
			if bv.IsNil() {
				bv.Set(reflect.MakeSlice(reflect.SliceOf(bv.Type().Elem()), 0, 0))
			}
			//切片,不考虑目标的长度
			if bv.Type().Elem() == av.Type().Elem() {
				reflect.Copy(bv, av)
			} else if bv.Type().Elem().Kind() == reflect.Struct && av.Type().Elem().Kind() == reflect.Struct {

				l := av.Len()
				for i := 0; i < l; i++ {
					v1 := av.Index(i)
					v2 := reflect.New(bv.Type().Elem())
					deepCopy(v1, v2)
					bv.Set(reflect.Append(bv, v2.Elem()))
				}
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

//返回底层value，如果不是字符串，数字，bool，结构体，切片数组等，可能出现未知bug
func getRealValue(in reflect.Value) reflect.Value {
	//if in.Kind() == reflect.Ptr || in.Kind() == reflect.Interface {
	for {
		if in.Kind() != reflect.Ptr && in.Kind() != reflect.Interface {
			return in
		}
		in = in.Elem()
	}
	//}
	return in
}

////只能是这几个类型，才能设置其值
//func isCanSetValue(kind reflect.Kind) bool {
//	return isInteger(kind) || isStruct(kind) || isString(kind) || isMap(kind) || isArrayOrSlice(kind) || isBool(kind) || isFloat(kind)
//}

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

//func isMap(kind reflect.Kind) bool {
//	return kind == 21
//}

func isString(kind reflect.Kind) bool {
	return kind == 24
}

func isStruct(kind reflect.Kind) bool {
	return kind == 25
}
