package entity

import (
	"reflect"
	"testing"
)

func TestReplies(t *testing.T) {
	res := Replies{}
	if res == nil {
		t.Error("new Replies should not be nil")
	}
	if len(res) != 0 {
		t.Errorf("new replies should have 0 length, but was %d", len(res))
	}
}

func TestRepliesAdd(t *testing.T) {
	tests := []struct {
		input Replies
		arg   bool
		want  Replies
	}{
		{Replies{true}, false, Replies{false, true}},
		{Replies{false, false}, true, Replies{true, false, false}},
		{Replies{false, true, true}, false, Replies{false, false, true, true}},
		{Replies{true, true, true, false}, false, Replies{false, true, true, true, false}},
		{Replies{false, false, false, false, true}, false, Replies{false, false, false, false, false}},
	}
	for i, v := range tests {
		res := v.input.Add(v.arg)
		if !reflect.DeepEqual(res, v.want) {
			t.Errorf("TestRepliesAdd #%d failed", i)
		}
	}
}

func TestRepliesLastGoodAnsw(t *testing.T) {
	tests := []struct {
		input Replies
		want  int
	}{
		{Replies{true}, 1},
		{Replies{false, false, true}, 0},
		{Replies{true, true, true}, 3},
		{Replies{true, false, true, true, false}, 1},
		{Replies{true, true, false, true}, 2},
	}
	for i, v := range tests {
		res := v.input.LastGoodAnsw()
		if res != v.want {
			t.Errorf("TestRepliesLastGoodAnsw #%d failed", i)
		}
	}
}
