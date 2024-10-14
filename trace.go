package trace

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
)

type Trace struct {
	context     context.Context
	Id          string        `json:"id"`
	Start       int64         `json:"start"`
	End         int64         `json:"end"`
	Text        string        `json:"text,omitempty"`
	Args        []interface{} `json:"args,omitempty"`
	Cost        float64       `json:"cost,omitempty"` // unit: second
	Description KV            `json:"description,omitempty"`
	Traces      []*Trace      `json:"children,omitempty"`
}

func clear(t *Trace) {
	t.Text = ""
	t.Args = nil
	t.Cost = 0
	t.Description = nil
	for _, m := range t.Traces {
		clear(m)
	}
	metaPool.Put(t)
}

func (t *Trace) Set(text string, args []interface{}) *Trace {
	t.Text, t.Args = text, args
	return t
}

func (t *Trace) SetDescription(key string, values ...interface{}) *Trace {
	if t.Description == nil {
		t.Description = make(KV)
	}
	for _, value := range values {
		t.Description.Set(key, value)
	}
	return t
}

func (t *Trace) WithContext(parent context.Context) *Trace {
	t.context = parent
	return t
}

func (t *Trace) NewTrace(id string) *Trace {
	trace := Pop(id)
	t.Traces = append(t.Traces, trace)
	return trace
}

func (t *Trace) Stop() *Trace {
	t.End = time.Now().UnixMilli()
	t.Cost = float64(t.End-t.Start) / 1000
	return t
}

func (t *Trace) Clear() {
	clear(t)
}

func (t *Trace) Upload(fn func(stream []byte) error) error {
	bytes, err := json.Marshal(t)
	if err != nil {
		return err
	}
	if fn == nil {
		return errors.New("now upload target")
	}
	return fn(bytes)
}

var metaPool = sync.Pool{New: func() interface{} { return &Trace{} }}

func Pop(id string) *Trace {
	meta := metaPool.Get().(*Trace)
	meta.Id, meta.Start = strings.TrimSpace(id), time.Now().UnixMilli()
	return meta
}

func End(t *Trace) {
	clear(t)
}
