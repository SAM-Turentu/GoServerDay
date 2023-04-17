package lru

import (
	"reflect"
	"testing"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("name", String("12344"))
	if v, ok := lru.Get("name"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit name=1234 failed")
	}
	if _, ok := lru.Get("phone"); ok {
		t.Fatalf("cache miss phone failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value Value) {
		keys = append(keys, key)
	}
	lru := New(int64(10), callback)
	lru.Add("key1", String("1234"))
	lru.Add("k2", String("k2"))
	lru.Add("k3", String("k3"))
	lru.Add("k4", String("k4"))
	lru.Add("k5", String("k5"))

	expect := []string{"key1", "k2"}
	if !reflect.DeepEqual(expect, keys) {
		t.Fatalf("Call OnEvicted failed, expect keys equals to %s", expect)
	}

}
