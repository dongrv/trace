package trace

import (
	"database/sql"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestTrace(t *testing.T) {
	t.Run("new-trace", func(t *testing.T) {
		ctx := New("server")
		defer ctx.Clear()

		// 模拟数据库操作
		query, args := "SELECT * FROM users WHERE UID=%d", []interface{}{1}
		_ = mockDB(ctx.New("run-sql"), query, args)

		// 模拟HTTP GET请求
		_ = mockHttpGet(ctx.New("run-http"), "https://github.com/dongrv/trace")

		// 模拟RPC请求
		var reply uint64
		_ = mockGRPC(ctx.New("run-rpc"), "/count/sum", nil, uint64(1), &reply)

		// 模拟执行Redis命令
		_, _ = mockRedisCMD(ctx.New("run-redis-cmd"), "LPOP", "user:list")

		t.Logf("%s", ctx.Stop().String())
	})
}

func mockDB(trace *Context, query string, args ...interface{}) error {
	defer trace.Set(query, args).Stop()
	time.Sleep(200 * time.Millisecond) // 模拟运行时间
	db := func(query string, args ...interface{}) (sql.Result, error) {
		return nil, errors.New("connection fail")
	}
	_, err := db(query, args)
	if err != nil {
		trace.SetKV("err", err.Error())
		return err
	}
	return nil
}

func mockHttpGet(ctx *Context, url string) error {
	defer ctx.Set(url, nil).Stop()
	_, err := http.Get(url)
	if err != nil {
		ctx.SetKV("err", err.Error())
		return err
	}
	time.Sleep(2 * time.Second)
	return nil
}

func mockGRPC(ctx *Context, fullMethod string, headers map[string]string, args interface{}, reply interface{}) error {
	defer ctx.Set(fullMethod, []interface{}{args}).Stop()
	rpc := func(fullMethod string, headers map[string]string, args interface{}, reply interface{}) error {
		time.Sleep(time.Second)
		p := reply.(*uint64)
		*p += args.(uint64)
		return nil
	}
	if err := rpc(fullMethod, headers, args, reply); err != nil {
		ctx.SetKV("err", err.Error())
		return err
	}
	ctx.SetKV("reply", reply)
	return nil
}

func mockRedisCMD(ctx *Context, cmd string, args ...interface{}) (interface{}, error) {
	defer ctx.Set(cmd, args).Stop()
	redis := func(cmd string, args ...interface{}) (interface{}, error) {
		time.Sleep(10 * time.Millisecond)
		return 100, nil
	}
	if v, err := redis(cmd, args); err != nil {
		ctx.SetKV("err", err.Error())
		return nil, err
	} else {
		return v, nil
	}
}
