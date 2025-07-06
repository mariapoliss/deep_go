package main

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type UInteger interface {
	uint32 | uint64 | uint16
}

func ToLittleEndian[T UInteger](number T) T {
	numSize := (int)(unsafe.Sizeof(number))
	for i := 0; i < (int)(numSize>>1); i++ {
		first := *(*uint8)(unsafe.Add(unsafe.Pointer(&number), i))
		last := *(*uint8)(unsafe.Add(unsafe.Pointer(&number), numSize-1-i))
		// switch
		*(*uint8)(unsafe.Add(unsafe.Pointer(&number), i)) = last
		*(*uint8)(unsafe.Add(unsafe.Pointer(&number), numSize-1-i)) = first
	}
	return number
}

func TestConversionUint32(t *testing.T) {
	tests := map[string]struct {
		number uint32
		result uint32
	}{
		"test case #1": {
			number: 0x00000000,
			result: 0x00000000,
		},
		"test case #2": {
			number: 0xFFFFFFFF,
			result: 0xFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF,
			result: 0xFF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF,
			result: 0xFFFF0000,
		},
		"test case #5": {
			number: 0x01020304,
			result: 0x04030201,
		},
		"test case #6": {
			number: 0x01020314,
			result: 0x14030201,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestConversionUint16(t *testing.T) {
	tests := map[string]struct {
		number uint16
		result uint16
	}{
		"test case #1": {
			number: 0x0000,
			result: 0x0000,
		},
		"test case #2": {
			number: 0xFFFF,
			result: 0xFFFF,
		},
		"test case #3": {
			number: 0x00FF,
			result: 0xFF00,
		},
		"test case #4": {
			number: 0x14FE,
			result: 0xFE14,
		},
		"test case #5": {
			number: 0x0102,
			result: 0x0201,
		},
		"test case #6": {
			number: 0x0314,
			result: 0x1403,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}

func TestConversionUint64(t *testing.T) {
	tests := map[string]struct {
		number uint64
		result uint64
	}{
		"test case #1": {
			number: 0x0000000000000000,
			result: 0x0000000000000000,
		},
		"test case #2": {
			number: 0xFFFFFFFFFFFFFFFF,
			result: 0xFFFFFFFFFFFFFFFF,
		},
		"test case #3": {
			number: 0x00FF00FF00FF00FF,
			result: 0xFF00FF00FF00FF00,
		},
		"test case #4": {
			number: 0x0000FFFF0000FFFF,
			result: 0xFFFF0000FFFF0000,
		},
		"test case #5": {
			number: 0x0102030401020304,
			result: 0x0403020104030201,
		},
		"test case #6": {
			number: 0xFFFFFFFF00000000,
			result: 0x00000000FFFFFFFF,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := ToLittleEndian(test.number)
			assert.Equal(t, test.result, result)
		})
	}
}
