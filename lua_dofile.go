package lua

// #cgo LDFLAGS:  -lluajit-5.1 -ldl -lm
//#include "glua.h"
import "C"
import "errors"

// DoFile -- dofile on all core vm
func DoFile(scriptFile string) error {
	vms := getCore().vms
	count := vms.Len()
	e := vms.Front()
	target, err := LoadScript(scriptFile)
	if err != nil {
		return err
	}
	for i := 0; i < count && e != nil; i++ {
		vm := e.Value.(*gLuaVM)
		ret := C.gluaL_dostring(vm.vm, C.CString(target))
		if ret != C.LUA_OK {
			return errors.New("luaL_loadfile failed")
		}
		e = e.Next()
		vms.PushBack(vm)
	}
	return nil
}
