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

type NSWIndex struct {
	EfSearch    int
	EfConstruct int
	Odegree     int
}

func (h *NSWIndex) CreateContext() unsafe.Pointer {
	icontext := C.NSWContext{
		ef_search:    C.int(h.EfSearch),
		ef_construct: C.int(h.EfConstruct),
		odegree:      C.int(h.Odegree),
	}
	ptr := C.malloc(C.size_t(unsafe.Sizeof(icontext)))
	C.memcpy(ptr, unsafe.Pointer(&icontext), C.size_t(unsafe.Sizeof(icontext)))
	return ptr
}

func (h *NSWIndex) ReleaseContext(ptr unsafe.Pointer) {
	C.free(ptr)
}
