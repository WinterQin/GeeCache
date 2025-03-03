package singleflight

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	var g Group
	v, err := g.Do("key", func() (interface{}, error) {
		return "bar", nil
	})

	if v != "bar" || err != nil {
		t.Errorf("Do v = %v, error = %v", v, err)
	}
}

func TestDoConcurrent(t *testing.T) {
	var g Group
	var wg sync.WaitGroup
	var calls int32 // 使用atomic来安全地计数

	// 创建10个并发请求
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, err := g.Do("key", func() (interface{}, error) {
				// 记录函数被调用的次数
				atomic.AddInt32(&calls, 1)
				// 模拟耗时操作
				time.Sleep(time.Millisecond)
				return "bar", nil
			})

			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			if v != "bar" {
				t.Errorf("got %v, want %v", v, "bar")
			}
		}()
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 验证函数只被调用一次
	if calls != 1 {
		t.Errorf("call function %d times, want 1", calls)
	}
}

// 添加一个新的测试用例
func TestDoWithDifferentKeys(t *testing.T) {
	var g Group
	var wg sync.WaitGroup

	// 使用不同的key测试
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		wg.Add(1)
		go func(k string) {
			defer wg.Done()
			v, err := g.Do(k, func() (interface{}, error) {
				return k + "_value", nil
			})

			if err != nil {
				t.Errorf("Do error: %v", err)
			}
			expected := k + "_value"
			if v != expected {
				t.Errorf("got %v, want %v", v, expected)
			}
		}(key)
	}

	wg.Wait()
}
