package geecachepb

import (
	"testing"

	"github.com/golang/protobuf/proto"
)

func TestProtobuf(t *testing.T) {
	req := &Request{
		Group: "scores",
		Key:   "Tom",
	}

	// 测试序列化
	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatal("marshal error:", err)
	}

	// 测试反序列化
	newReq := &Request{}
	err = proto.Unmarshal(data, newReq)
	if err != nil {
		t.Fatal("unmarshal error:", err)
	}

	// 验证数据
	if newReq.Group != req.Group || newReq.Key != req.Key {
		t.Errorf("data mismatch, got %v, want %v", newReq, req)
	}
}
