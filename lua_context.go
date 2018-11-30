package lua

type gLuaContext struct {
	id         int64
	vmId       int64
	threadId   int64
	scriptPath string
	methodName string
	args       []interface{}
	callback   chan interface{}
}
