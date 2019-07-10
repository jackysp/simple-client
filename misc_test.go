package simple_client

import (
	"testing"
)

func TestAppend(t *testing.T) {
	a := make([]*int, 0, 10)
	t.Log(len(a))
	a = append(a, nil)
	t.Log(len(a))
	a = a[:0]
	t.Log(len(a))
}

