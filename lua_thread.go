package lua

import (
	"errors"
)

// #cgo LDFLAGS:  -lluajit-5.1 -ldl -lm
//#include "glua.h"
import "C"

type gLuaThread struct {
	id     int64
	thread *C.struct_lua_State
}

func newGLuaThread(vm *C.struct_lua_State) *gLuaThread {
	gl := &gLuaThread{}
	gl.id, gl.thread = createLuaThread(vm)
	return gl
}

func (t *gLuaThread) destory() {
	cleanDummy(t.thread)
}

func (t *gLuaThread) call(scriptPath string, methodName string, args ...interface{}) (interface{}, error) {
	var ret C.int
	if scriptPath != "" {
		target, err := LoadScript(scriptPath)
		if err != nil {
			return nil, err
		}

		ret = C.gluaL_dostring(t.thread, C.CString(target))
		if ret != C.LUA_OK {
			ExpireScript(scriptPath)
			errStr := C.GoString(C.glua_tostring(t.thread, -1))
			return nil, errors.New(errStr)
		}
	}

	C.glua_getglobal(t.thread, C.CString(methodName))
	pushToLua(t.thread, args...)

	ret = C.lua_resume(t.thread, C.int(len(args)))
	switch ret {
	case C.LUA_OK:
		{
			var (
				res interface{}
				err interface{}
			)
			num := C.lua_gettop(t.thread)
			if num > 1 {
				err = pullFromLua(t.thread, -1)
				C.lua_remove(t.thread, -1)
				res = pullFromLua(t.thread, -1)
			} else {
				res = pullFromLua(t.thread, -1)
			}
			C.glua_pop(t.thread, -1)
			if err != nil {
				return nil, errors.New(err.(string))
			}
			return res, nil
		}
	case C.LUA_YIELD:
		return nil, errors.New("LUA_YIELD")
	default:
		temp := C.GoString(C.glua_tostring(t.thread, -1))
		return nil, errors.New(temp)
	}
}

func (t *gLuaThread) resume(args ...interface{}) (interface{}, error) {
	pushToLua(t.thread, args...)
	num := C.lua_gettop(t.thread)
	ret := C.lua_resume(t.thread, num)
	switch ret {
	case C.LUA_OK:
		err := pullFromLua(t.thread, -1)
		C.lua_remove(t.thread, -1)
		res := pullFromLua(t.thread, -1)
		C.glua_pop(t.thread, -1)
		if err != nil {
			return nil, errors.New(err.(string))
		}
		return res, nil
	case C.LUA_YIELD:
		return nil, errors.New("LUA_YIELD")
	default:
		temp := C.GoString(C.glua_tostring(t.thread, -1))
		return nil, errors.New(temp)
	}
}
