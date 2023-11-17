package reflectutils

import "reflect"

// TypeOf 封装TypeOf解除引用
func TypeOf(i any) reflect.Type {
	o := reflect.TypeOf(i)
	for o.Kind() == reflect.Ptr {
		o = o.Elem()
	}
	return o
}

// ValueOf 封装ValueOf解除引用
func ValueOf(i any) reflect.Value {
	o := reflect.ValueOf(i)
	for o.Kind() == reflect.Ptr {
		o = o.Elem()
	}
	return o
}
