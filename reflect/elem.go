package reflectutils

import "reflect"

func TypeOf(i any) reflect.Type {
	o := reflect.TypeOf(i)
	for o.Kind() == reflect.Ptr {
		o = o.Elem()
	}
	return o
}

func ValueOf(i any) reflect.Value {
	o := reflect.ValueOf(i)
	for o.Kind() == reflect.Ptr {
		o = o.Elem()
	}
	return o
}
