package util_test

import (
	"fmt"
	"reflect"
	"testing"
	"umx/tools/pressure/server/util"
)

func TestIp(t *testing.T) {
	ip1 := "10.116.192.31"
	ip2 := "0.0.0.0"
	ip3 := "255.255.255.255"
	ip4 := "eeeeeeeeee"
	fmt.Println(util.IsLegalIpv4(ip1))
	fmt.Println(util.IsLegalIpv4(ip2))
	fmt.Println(util.IsLegalIpv4(ip3))
	fmt.Println(util.IsLegalIpv4(ip4))
}

func TestUuid(t *testing.T) {
	fmt.Println(util.UUID())
}

type AA struct {
	Name string
	Age  int
}

func TestReflection(t *testing.T) {
	aa := AA{
		Name: "eee",
		Age:  11,
	}
	curT := reflect.TypeOf(aa)
	p := &aa
	curTP := reflect.TypeOf(p)
	fmt.Println(curT, curTP)
	fmt.Println(curT.Kind() == reflect.Struct, curTP.Kind() == reflect.Ptr)
	xx := ee(p)
	fmt.Println(xx.Kind() == reflect.Ptr)
}

func ee(i interface{}) reflect.Type {
	v := &i
	t := reflect.TypeOf(v)
	return t
}
