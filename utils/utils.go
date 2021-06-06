package utils

import (
	"fmt"
	"reflect"

	"crypto/md5"
)

// MD5值生成函数
func Md5sum(data []byte) string {
	hash := md5.Sum(data)
	return fmt.Sprintf("%x", hash)
}

// 构造器属性值复制函数
func StructAssign(source interface{}, target interface{}) {
	targetVal := reflect.ValueOf(target).Elem()
	sourceVal := reflect.ValueOf(source).Elem()
	sourceTypeOfT := sourceVal.Type()
	for i := 0; i < sourceVal.NumField(); i++ {
		name := sourceTypeOfT.Field(i).Name
		if ok := targetVal.FieldByName(name).IsValid(); ok {
			targetVal.FieldByName(name).Set(reflect.ValueOf(sourceVal.Field(i).Interface()))
		}
	}
}
