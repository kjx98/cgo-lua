package lua

import (
	"errors"
	"sync"
)

// #cgo LDFLAGS:  -lluajit-5.1 -ldl -lm
//#include "glua.h"
import "C"

var (
	yieldCache map[int64]*gLuaYieldContext
	yieldRW    sync.RWMutex
)

func init() {
	yieldCache = make(map[int64]*gLuaYieldContext)
}

type gLuaYieldContext struct {
	methodName string
	args       []interface{}
}

//yield method
func storeYieldContext(vm *C.struct_lua_State, methodName string, args ...interface{}) {
	if vm == nil {
		panic(errors.New("Invalid Lua VM"))
	}
	vmKey := generateLuaStateId(vm)

	yieldRW.Lock()
	defer yieldRW.Unlock()
	yieldCache[vmKey] = &gLuaYieldContext{methodName: methodName, args: args}
}

func loadYieldContext(threadId int64) (*gLuaYieldContext, error) {
	yieldRW.RLock()
	defer func() {
		delete(yieldCache, threadId)
		yieldRW.RUnlock()
	}()
	target, ok := yieldCache[threadId]
	if false == ok {
		return nil, errors.New("Invalid Yield Contxt")
	}
	return target, nil
}
