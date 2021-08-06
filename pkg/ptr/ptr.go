package ptr

import "reflect"

// IsPtr return true if given value v is a pointer
func IsPtr(v interface{}) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}

// Indirect returns the value of given v
func Indirect(v interface{}) interface{} {
	if !IsPtr(v) {
		return v
	}

	return reflect.Indirect(reflect.ValueOf(v)).Interface()
}

// String returns a string pointer
func String(v string) *string {
	return &v
}

// Float32 returns a float32 pointer
func Float32(v float32) *float32 {
	return &v
}

// Float64 returns a float64 pointer
func Float64(v float64) *float64 {
	return &v
}

// Bool returns a bool pointer
func Bool(v bool) *bool {
	return &v
}

// Int returns an int pointer
func Int(v int) *int {
	return &v
}

// Int8 returns an int8 pointer
func Int8(v int8) *int8 {
	return &v
}

// Int16 returns an int16 pointer
func Int16(v int16) *int16 {
	return &v
}

// Int32 returns an int32 pointer
func Int32(v int32) *int32 {
	return &v
}

// Int64 returns an int64 pointer
func Int64(v int64) *int64 {
	return &v
}

// UInt returns an uint pointer
func UInt(v uint) *uint {
	return &v
}

// UInt8 returns an uint8 pointer
func UInt8(v uint8) *uint8 {
	return &v
}

// UInt16 returns an uint16 pointer
func UInt16(v uint16) *uint16 {
	return &v
}

// UInt32 returns an uint32 pointer
func UInt32(v uint32) *uint32 {
	return &v
}

// UInt64 returns an uint64 pointer
func UInt64(v uint64) *uint64 {
	return &v
}
