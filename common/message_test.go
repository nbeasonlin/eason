package common

import (
	"fmt"
	"testing"
)

func TestMessageA(t *testing.T) {

	m := NewJSONMessage(1, 2, []byte("Hello World"))
	fmt.Println(string(m.Bytes()))
	rm := NewJSONMessageByBytes(m.Bytes())
	fmt.Println(rm.From(), rm.To(), rm.Timestamp(), string(rm.Data()))

}
