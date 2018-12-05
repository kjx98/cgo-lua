package lua

import (
	"encoding/json"
	"testing"
	"time"
)

func test_sum(args ...interface{}) (interface{}, error) {
	sum := 0
	for _, arg := range args {
		sum = sum + int(arg.(int64))
	}
	return sum, nil
}

func json_decode(args ...interface{}) (interface{}, error) {
	raw := args[0].(string)

	var res map[string]interface{}
	err := json.Unmarshal([]byte(raw), &res)
	return res, err
}

func TestLuaRun(t *testing.T) {
	RegisterExternMethod("json_decode", json_decode)
	RegisterExternMethod("test_sum", test_sum)

	start_t := time.Now()
	res, err := Call("script.lua", "async_json_encode", nil)
	end_t := time.Now()
	if err != nil {
		t.Error("async_json_encode", err)
	} else {
		delta_t := end_t.Sub(start_t)
		t.Log(res, "cost", delta_t.Seconds(), "seconds")
	}

	start_t = time.Now()
	res, err = Call("", "test_args", 69)
	end_t = time.Now()
	if err != nil {
		t.Log("ok!! test_args", err)
	} else {
		t.Error("shouldnot found function")
	}
	err = DoFile("script.lua")
	if err != nil {
		t.Error("Loadfile", err)
	}
	start_t = time.Now()
	res, err = Call("", "test_args", 69)
	end_t = time.Now()
	if err != nil {
		t.Error("test_args", err)
	} else {
		delta_t := end_t.Sub(start_t)
		t.Log(res, "cost", delta_t.Seconds(), "seconds")
	}

	start_t = time.Now()
	res, err = Call("", "test_pull_table", "ok.?ok")
	end_t = time.Now()
	if err != nil {
		t.Error("pull_table", err)
	} else {
		delta_t := end_t.Sub(start_t)
		t.Log(res, "cost", delta_t.Seconds(), "seconds")
	}

	start_t = time.Now()
	res, err = Call("", "fib", 35)
	end_t = time.Now()
	if err != nil {
		t.Error("fib", err)
	} else {
		delta_t := end_t.Sub(start_t)
		t.Log(res, "cost", delta_t.Seconds(), "seconds")
	}
}
