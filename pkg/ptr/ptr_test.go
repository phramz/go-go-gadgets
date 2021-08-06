package ptr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndirect(t *testing.T) {
	outs := []struct {
		Ptr      interface{}
		Expected interface{}
	}{
		{Ptr: nil, Expected: nil},
		{Ptr: String("foo"), Expected: "foo"},
		{Ptr: Bool(true), Expected: true},
		{Ptr: Bool(false), Expected: false},
		{Ptr: Int(123), Expected: 123},
		{Ptr: Int8(123), Expected: int8(123)},
		{Ptr: Int16(123), Expected: int16(123)},
		{Ptr: Int32(123), Expected: int32(123)},
		{Ptr: Int64(123), Expected: int64(123)},
		{Ptr: UInt(123), Expected: uint(123)},
		{Ptr: UInt8(123), Expected: uint8(123)},
		{Ptr: UInt16(123), Expected: uint16(123)},
		{Ptr: UInt32(123), Expected: uint32(123)},
		{Ptr: UInt64(123), Expected: uint64(123)},
		{Ptr: Float32(123.456), Expected: float32(123.456)},
		{Ptr: Float64(123.456), Expected: 123.456},
	}

	for _, tt := range outs {
		tt := tt
		t.Run(fmt.Sprintf("%T", tt.Expected), func(t *testing.T) {
			result := Indirect(tt.Ptr)

			assert.Equal(t, tt.Expected, result)
			assert.False(t, IsPtr(result))
		})
	}
}
