package lua

//#include "glua.h"
import "C"

import (
	"errors"
	"unsafe"
)

// DoFile -- dofile on all core vm
func DoFile(scriptPath string) error {
	vms := getCore().vms
	count := vms.Len()
	e := vms.Front()
	target, err := LoadScript(scriptPath)
	if err != nil {
		return err
	}
	script := C.CString(target)
	defer C.free(unsafe.Pointer(script))
	for i := 0; i < count && e != nil; i++ {
		vm := e.Value.(*gLuaVM)
		ret := C.gluaL_dostring(vm.vm, script)
		if ret != C.LUA_OK {
			ExpireScript(scriptPath)
			return errors.New("luaL_loadfile failed")
		}
		// ignore returns of luaL_dostring
		num := C.lua_gettop(vm.vm)
		C.glua_pop(vm.vm, num)
		e = e.Next()
		vms.PushBack(vm)
	}
	return nil
}
