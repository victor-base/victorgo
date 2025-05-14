package index

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lvictor
#include <victor/victor.h>
#include <stdlib.h>
#include <string.h>     // memcpy
*/
import "C"

import "unsafe"

type HNSWIndex struct {
	EfSearch    int
	EfConstruct int
	M0          int
}

func (h *HNSWIndex) CreateContext() unsafe.Pointer {
	icontext := C.HNSWContext{
		ef_search:    C.int(h.EfSearch),
		ef_construct: C.int(h.EfConstruct),
		M0:           C.int(h.M0),
	}
	ptr := C.malloc(C.size_t(unsafe.Sizeof(icontext)))
	C.memcpy(ptr, unsafe.Pointer(&icontext), C.size_t(unsafe.Sizeof(icontext)))
	return ptr
}

func (h *HNSWIndex) ReleaseContext(ptr unsafe.Pointer) {
	C.free(ptr)
}
