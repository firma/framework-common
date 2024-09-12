package errno

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoDataError(t *testing.T) {
	err := NewError(-1, "TEST", "参数错误")

	assert.Equal(t, err.ToString(), `{"code":-1,"err":"参数错误","data":{}}`)
}

func TestDataError(t *testing.T) {
	err := NewError(0, "TEST", "成功").WithData(
		struct {
			Foo string
		}{Foo: "bar"},
	)

	assert.Equal(t, err.ToString(), `{"code":0,"err":"成功","data":{"Foo":"bar"}}`)
}
