package trace

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
)

type Context struct {
	context     context.Context
	Id          string        `json:"id"`
	Start       int64         `json:"start"`
	End         int64         `json:"end"`
	Format      string        `json:"format,omitempty"`
	Args        []interface{} `json:"args,omitempty"`
	Cost        float64       `json:"cost,omitempty"` // unit: second
	Description KV            `json:"description,omitempty"`
	Traces      []*Context    `json:"traces,omitempty"`
}

func clear(t *Context) {
	t.Format = ""
	t.Args = nil
	t.Cost = 0
	t.Description = nil
	for _, m := range t.Traces {
		clear(m)
	}
	metaPool.Put(t)
}

func (ctx *Context) Set(text string, args []interface{}) *Context {
	ctx.Format, ctx.Args = text, args
	return ctx
}

func (ctx *Context) SetKV(key string, values ...interface{}) *Context {
	if ctx.Description == nil {
		ctx.Description = make(KV)
	}
	for _, value := range values {
		ctx.Description.Set(key, value)
	}
	return ctx
}

func (ctx *Context) WithContext(parent context.Context) *Context {
	ctx.context = parent
	return ctx
}

func (ctx *Context) New(id string) *Context {
	trace := New(id)
	ctx.Traces = append(ctx.Traces, trace)
	return trace
}

func (ctx *Context) Stop() *Context {
	ctx.End = time.Now().UnixMilli()
	ctx.Cost = float64(ctx.End-ctx.Start) / 1000
	return ctx
}

func (ctx *Context) String() string {
	bytes, _ := json.Marshal(ctx)
	return string(bytes)
}

func (ctx *Context) Clear() {
	clear(ctx)
}

func (ctx *Context) Upload(fn func(stream []byte) error) error {
	bytes, err := json.Marshal(ctx)
	if err != nil {
		return err
	}
	if fn == nil {
		return errors.New("no upload target")
	}
	return fn(bytes)
}

var metaPool = sync.Pool{New: func() interface{} { return &Context{} }}

func New(id string) *Context {
	ctx := metaPool.Get().(*Context)
	ctx.Id, ctx.Start = strings.TrimSpace(id), time.Now().UnixMilli()
	return ctx
}

func End(t *Context) {
	clear(t)
}
