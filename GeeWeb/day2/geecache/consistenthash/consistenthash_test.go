package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {

	// 定义一个hash算法，只处理数字
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key)) // 数字转字符串
		return uint32(i)
	})

	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
	hash.Add("8")
	testCases["27"] = "8"
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

}
