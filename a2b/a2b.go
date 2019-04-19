package a2b

import (
	"errors"
	"log"
	"reflect"
	"time"
)

var typeError = errors.New("a和b的数据类型不一致，请核对其底层类型是否一致，比如整数类型和整数类型，结构体和结构体,数组和切片")

const (
	created = "CreatedAt"
	updated = "UpdatedAt"
	deleted = "DeletedAt"
)

var timeType = reflect.TypeOf(time.Now())

func A2B(a interface{}, b interface{}) error {
	bv := reflect.ValueOf(b)
	if bv.Kind() != reflect.Ptr || bv.IsNil() {
		return errors.New("b对象必须为指针类型")
	}
	bv = getRealValue(bv)
	av := getRealValue(reflect.ValueOf(a))
	if !av.IsValid() || bv.IsNil() {
		return nil
	}
	if !isSimilarType(av, bv) {
		log.Printf("不能在类型或指向不同类型的指针之间赋值,%s和%s\n", av.Type(), bv.Type())
		return typeError
	}

	deepCopy(av, bv)
	return nil
}

func deepCopy(av reflect.Value, bv reflect.Value) {
	//todo
	av = getRealValue(av)
	bv = getRealValue(bv)

	if !av.IsValid() {
		return
	}

	//if !bv.IsValid() {
	//	bv.Set(reflect.New(reflect.TypeOf(bv)).Elem())
	//}

	if !isSimilarType(av, bv) {
		log.Printf("不能在类型或指向不同类型的指针之间赋值,%s和%s\n", getRealValue(av).Kind(), getRealValue(bv).Kind())
		return
	}
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

				len := av.Len()
				for i := 0; i < len; i++ {
					v1 := av.Index(i)
					v2 := reflect.New(bv.Type().Elem())
					deepCopy(v1, v2)
					reflect.Append(bv, v1)
				}
			}
		}
		return
	}

	//map
	if isMap(bv.Kind()) {
		//如果键类型不一样，返回
		if bv.Type().Key() != av.Type().Key() {
			return
		}
		if bv.CanSet() {
			if bv.IsNil() {
				bv.Set(reflect.MakeMap(reflect.MapOf(bv.Type().Key(), bv.Type().Elem())))
			}
			//如果元素类型相同
			keys := av.MapKeys()
			if bv.Type().Elem() == av.Type().Elem() {
				for _, key := range keys {
					//todo 如果值的类型一样
					bv.SetMapIndex(key, av.MapIndex(key))
				}
			} else if bv.Type().Elem().Kind() == reflect.Struct && av.Type().Elem().Kind() == reflect.Struct {
				for _, key := range keys {
					//todo 如果值的类型不一样,且都是结构体
					//数据来源
					v1 := av.MapIndex(key)
					//给目标map创建出一个新的对象,然后设置其中的值
					v2 := reflect.New(bv.Type().Elem())
					deepCopy(v1, v2)
					bv.SetMapIndex(key, v2)
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

	//TODO 这儿有没有必要
	av = getRealValue(av)
	bv = getRealValue(bv)
	if bv.IsNil() {

	}

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
	////todo 特殊处理gorm.model这种匿名字段
	////TODO 因为获取到数据库模型和pb模型的结构,因为这个匿名字段的存在就不一样了,数据库模型就有两层结构,所以去掉一层结构,并加上里面的四个字段
	//_, ok := fieldMap[model]
	//if ok {
	//	delete(fieldMap, model)
	//	fieldMap[id] = struct{}{}
	//	fieldMap[created] = struct{}{}
	//	fieldMap[updated] = struct{}{}
	//	fieldMap[deleted] = struct{}{}
	//}
	return fieldMap
}

func isSimilarType(av reflect.Value, bv reflect.Value) bool {
	a := av.Kind()
	b := bv.Kind()
	//return (isInteger(a) && isInteger(b)) || (isBool(a) && isBool(b)) || (isFloat(a) && isFloat(b)) || (isArrayOrSlice(a) && isArrayOrSlice(b)) || (isMap(a) && isMap(b)) || (isString(a) && isString(b)) || (isStruct(a) && isStruct(b))
	return (isInteger(a) && isInteger(b)) || (isFloat(a) && isFloat(b)) || (isArrayOrSlice(a) && isArrayOrSlice(b)) || (av.Type() != bv.Type())
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
