package lua

import (
	"math"
	"unsafe"
)

// #cgo LDFLAGS:  -lluajit-5.1 -ldl -lm
//#include "glua.h"
import "C"

func pushToLua(L *C.struct_lua_State, args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case string:
			cStr := C.CString(arg.(string))
			defer C.free(unsafe.Pointer(cStr))
			C.lua_pushstring(L, cStr)
		case float64:
			C.lua_pushnumber(L, C.lua_Number(arg.(float64)))
		case float32:
			C.lua_pushnumber(L, C.lua_Number(arg.(float32)))
		case uint64:
			C.lua_pushnumber(L, C.lua_Number(arg.(uint64)))
		case int64:
			C.lua_pushnumber(L, C.lua_Number(arg.(int64)))
		case uint32:
			C.lua_pushnumber(L, C.lua_Number(arg.(uint32)))
		case int32:
			C.lua_pushnumber(L, C.lua_Number(arg.(int32)))
		case uint16:
			C.lua_pushnumber(L, C.lua_Number(arg.(uint16)))
		case int16:
			C.lua_pushnumber(L, C.lua_Number(arg.(int16)))
		case uint8:
			C.lua_pushnumber(L, C.lua_Number(arg.(uint8)))
		case int8:
			C.lua_pushnumber(L, C.lua_Number(arg.(int8)))
		case uint:
			C.lua_pushnumber(L, C.lua_Number(arg.(uint)))
		case int:
			C.lua_pushnumber(L, C.lua_Number(arg.(int)))
		case map[string]interface{}:
			pushMapToLua(L, arg.(map[string]interface{}))
		case []interface{}:
			pushArrayToLua(L, arg.([]interface{}))
		default:
			{
				ptr := pushDummy(L, arg)
				C.glua_pushlightuserdata(L, ptr)
			}
		}
	}
}

func pushArrayToLua(L *C.struct_lua_State, data []interface{}) {
	C.lua_createtable(L, 0, 0)
	if len(data) == 0 {
		return
	}
	for index, value := range data {
		C.lua_pushnumber(L, C.lua_Number(index))
		pushToLua(L, value)
		C.lua_settable(L, -3)
	}
}

func pushMapToLua(L *C.struct_lua_State, data map[string]interface{}) {
	C.lua_createtable(L, 0, 0)
	if len(data) == 0 {
		return
	}
	for key, value := range data {
		cStr := C.CString(key)
		defer C.free(unsafe.Pointer(cStr))
		C.lua_pushstring(L, cStr)
		pushToLua(L, value)
		C.lua_settable(L, -3)
	}
}

func pullLuaTable(_L *C.struct_lua_State) interface{} {
	keys := make([]interface{}, 0)
	values := make([]interface{}, 0)

	numKeyCount := 0
	var (
		key   interface{}
		value interface{}
	)
	C.lua_pushnil(_L)
	for C.lua_next(_L, -2) != 0 {
		kType := C.lua_type(_L, -2)
		if kType == 4 {
			key = C.GoString(C.glua_tostring(_L, -2))
		} else {
			key = int(C.lua_tointeger(_L, -2))
			numKeyCount = numKeyCount + 1
		}
		vType := C.lua_type(_L, -1)
		switch vType {
		case 0:
			C.glua_pop(_L, 1)
			continue
		case 1:
			temp := C.lua_toboolean(_L, -1)
			if temp == 1 {
				value = true
			} else {
				value = false
			}
		case 2:
			ptr := C.glua_touserdata(_L, -1)
			target, err := findDummy(_L, ptr)
			if err != nil {
				C.glua_pop(_L, 1)
				continue
			}
			value = target
		case 3:
			temp := float64(C.glua_tonumber(_L, -1))
			if math.Ceil(temp) > temp {
				value = temp
			} else {
				value = int64(temp)
			}
		case 4:
			value = C.GoString(C.glua_tostring(_L, -1))
		case 5:
			value = pullLuaTable(_L)
		}
		keys = append(keys, key)
		values = append(values, value)
		C.glua_pop(_L, 1)
	}
	if numKeyCount == len(keys) {
		return values
	}
	if numKeyCount == 0 {
		result := make(map[string]interface{})
		for index, key := range keys {
			result[key.(string)] = values[index]
		}
		return result
	} else {
		result := make(map[interface{}]interface{})
		for index, key := range keys {
			result[key] = values[index]
		}
		return result
	}
}

func pullFromLua(L *C.struct_lua_State, index int) interface{} {
	vType := C.lua_type(L, C.int(index))
	switch vType {
	case C.LUA_TBOOLEAN:
		res := C.lua_toboolean(L, C.int(index))
		if res == 0 {
			return false
		}
		return true
	case C.LUA_TNUMBER:
		temp := float64(C.glua_tonumber(L, -1))
		if math.Ceil(temp) > temp {
			return temp
		} else {
			return int64(temp)
		}
	case C.LUA_TSTRING:
		return C.GoString(C.glua_tostring(L, C.int(index)))
	case C.LUA_TTABLE:
		return pullLuaTable(L)
	case C.LUA_TLIGHTUSERDATA:
		ptr := C.glua_touserdata(L, C.int(index))
		target, err := findDummy(L, ptr)
		if err != nil {
			return nil
		} else {
			return target
		}
	case C.LUA_TNIL:
		return nil
		/*
			default:
				panic(errors.New(fmt.Sprintf("Unsupport Type %d", vType)))
		*/
	}
	return nil
}
