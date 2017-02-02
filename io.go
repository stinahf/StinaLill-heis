package io

// #include io.h
// #cgo LDFLAGS: -lcomedi -lm

import "C"

func io_init() bool {
	return int(C.io_init()) != 0
}

func io_set_bit(int channel) {
	C.io_set_bit(C.int(channel))
}

func io_clear_bit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func io_write_analag(channel int, value int) {
	C.io_write_analag(C.int(channel), C.int(value))
}

func io_read_bit(channel int) bool {
	return int(C.io_read_bit(C.int(channel))) != 0
}

func io_read_analog(channel int) int {
	return int(C.io_read_analog(C.int(channel)))
}