package gobilibili

import (
	"fmt"
)

//Pkg 包名称，通过覆写这个值可以改变Me等附带完整错误信息的FullErr里的Pkg信息
var Pkg = "packageName"

//FullErr 顾名思义，完整的error包装，包含一些详细信息
type FullErr struct {
	Pkg  string
	Info string
	Prev error
}

func (e *FullErr) Error() string {
	if e.Prev == nil {
		return fmt.Sprintf("%s: %s", e.Pkg, e.Info)
	}
	return fmt.Sprintf("%s: %s\n%v", e.Pkg, e.Info, e.Prev)
}

//Err 简单的error包装，不包含其他信息
type Err string

func (e Err) Error() string { return string(e) }

//Me 包装一个可追溯错误链的错误
func Me(err error, format string, args ...interface{}) *FullErr {
	if len(args) > 0 {
		return &FullErr{
			Pkg:  Pkg,
			Info: fmt.Sprintf(format, args...),
			Prev: err,
		}
	}
	return &FullErr{
		Pkg:  Pkg,
		Info: format,
		Prev: err,
	}
}

//CatchErr Catch error
func CatchErr(err error, format string, args ...interface{}) {
	if err != nil {
		panic(Me(err, format, args...))
	}
}

//CatchAny Catch any error
func CatchAny(anyRes ...interface{}) {
	for _, obj := range anyRes {
		if e, ok := obj.(error); ok && e != nil {
			panic(e)
		}
	}
}

//CatchThrow Catch throw
func CatchThrow(err *error) {
	if p := recover(); p != nil {
		if e, ok := p.(error); ok {
			*err = e
		} else {
			panic(p)
		}
	}
}

//CatchThrowHandle Catch throw and call handle
func CatchThrowHandle(handle func(err error)) {
	if p := recover(); p != nil {
		if e, ok := p.(error); ok {
			handle(e)
		} else {
			panic(p)
		}
	}
}

//OrginErr Return the orgin error
func OrginErr(e error) error {
	ret := e
	for err, ok := ret.(*FullErr); ok && err.Prev != nil; err, ok = ret.(*FullErr) {
		ret = err.Prev
	}
	return ret
}

//MustTrue True or panic a error
func MustTrue(b bool, panicErr string) {
	if !b {
		panic(Err(panicErr))
	}
}
