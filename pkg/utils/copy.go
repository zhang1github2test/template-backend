package utils

import "reflect"

// DeepCopyStruct 将 src 的属性深拷贝到 dst
// 规则：名字相同 && 类型相同的字段才会拷贝
func DeepCopyStruct(dst, src interface{}) {
	dstVal := reflect.ValueOf(dst).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Type().Field(i)
		srcValue := srcVal.Field(i)

		dstField := dstVal.FieldByName(srcField.Name)
		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}

		if dstField.Type() != srcValue.Type() {
			continue
		}

		dstField.Set(deepCopyValue(srcValue))
	}
}

func deepCopyValue(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		newPtr := reflect.New(v.Elem().Type())
		newPtr.Elem().Set(deepCopyValue(v.Elem()))
		return newPtr

	case reflect.Slice:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		newSlice := reflect.MakeSlice(v.Type(), v.Len(), v.Cap())
		for i := 0; i < v.Len(); i++ {
			newSlice.Index(i).Set(deepCopyValue(v.Index(i)))
		}
		return newSlice

	case reflect.Map:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		newMap := reflect.MakeMapWithSize(v.Type(), v.Len())
		for _, key := range v.MapKeys() {
			newMap.SetMapIndex(key, deepCopyValue(v.MapIndex(key)))
		}
		return newMap

	case reflect.Struct:
		newStruct := reflect.New(v.Type()).Elem()
		for i := 0; i < v.NumField(); i++ {
			if newStruct.Field(i).CanSet() {
				newStruct.Field(i).Set(deepCopyValue(v.Field(i)))
			}
		}
		return newStruct

	default:
		// 基础类型（int、string、bool 等）直接返回
		return v
	}
}
