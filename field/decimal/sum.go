package decimal

import (
	"fmt"
	"reflect"
)

func Sum(arg interface{}) Dec {
	slice := getSlice(arg)

	sum := Dec{Number: 0, Prec: 0}
	for i := 0; i < slice.Len(); i++ {
		dec := slice.Index(i).Interface().(Dec)
		fmt.Printf("Adding %s\n", dec)
		sum = sum.Add(dec)
		fmt.Println("New Sum: ", dec.String())
	}

	fmt.Println("Total Sum: ", sum.String())

	return sum
}

func getSlice(arg interface{}) reflect.Value {
	val := reflect.ValueOf(arg)

	if val.Kind() != reflect.Slice && val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Slice {
		panic("not a slice")
	}

	return val
}
