package handler

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type WrapFunc struct {
	ctxFunc func(ctx *gin.Context)
}

func Wrap(handlerFunc interface{}) func(*gin.Context) {
	// 提前反射
	paramNum := reflect.TypeOf(handlerFunc).NumIn()
	funcValue := reflect.ValueOf(handlerFunc)
	funcType := reflect.TypeOf(handlerFunc)
	paramType := funcType.In(1).Elem()

	// 判断是否 Func
	if funcType.Kind() != reflect.Func {
		panic("the route handlerFunc must be a function")
	}
	// ... 还可以做一些其他校验确保无误
	return func(context *gin.Context) {
		// 只有一个参数说明是未重构的 HandlerFunc
		if paramNum == 1 {
			funcValue.Call(valOf(context))
			return
		}
		proxyHandlerFunc(funcValue, context, paramType)
	}
}

func proxyHandlerFunc(funcValue reflect.Value, ctx *gin.Context, paramType reflect.Type) {
	// 创建实例
	param := reflect.New(paramType).Interface()
	// ...
	// 调用真实 HandlerFunc
	funcValue.Call(valOf(ctx, param))
}

func valOf(i ...interface{}) []reflect.Value {
	var rt []reflect.Value
	for _, i2 := range i {
		rt = append(rt, reflect.ValueOf(i2))
	}
	return rt
}
