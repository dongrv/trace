package trace

import "encoding/json"

type KV map[string][]interface{}

func (kv KV) Bytes() ([]byte, error) {
	if len(kv) == 0 {
		return nil, nil
	}
	return json.Marshal(kv)
}

func (kv KV) String() (string, error) {
	bytes, err := kv.Bytes()
	return string(bytes), err
}

func (kv KV) Set(key string, value interface{}) {
	_, ok := kv[key]
	if !ok {
		kv[key] = make([]interface{}, 0, 1)
	}
	kv[key] = append(kv[key], value)
}
