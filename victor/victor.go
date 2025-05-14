package victor

/*
#cgo CFLAGS: -I${SRCDIR}/../include
#cgo LDFLAGS: -L${SRCDIR}/../lib -lvictor
#include <victor/victor.h>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
	idx "victorgo/index"
)

// ErrorCode maps C error codes to Go
type ErrorCode int

const (
	SUCCESS ErrorCode = iota
	INVALID_INIT
	INVALID_INDEX
	INVALID_VECTOR
	INVALID_RESULT
	INVALID_DIMENSIONS
	INVALID_ARGUMENT
	INVALID_ID
	INVALID_REF
	INVALID_METHOD
	DUPLICATED_ENTRY
	NOT_FOUND_ID
	INDEX_EMPTY
	THREAD_ERROR
	SYSTEM_ERROR
	FILEIO_ERROR
	NOT_IMPLEMENTED
	INVALID_FILE
)

// errorMessages maps error codes to human-readable messages
var errorMessages = map[ErrorCode]string{
	SUCCESS:            "Success",
	INVALID_INIT:       "Invalid initialization",
	INVALID_INDEX:      "Invalid index",
	INVALID_VECTOR:     "Invalid vector",
	INVALID_RESULT:     "Invalid result",
	INVALID_DIMENSIONS: "Invalid dimensions",
	INVALID_ARGUMENT:   "Invalid argument",
	INVALID_ID:         "Invalid ID",
	INVALID_REF:        "Invalid reference",
	DUPLICATED_ENTRY:   "Duplicated entry",
	NOT_FOUND_ID:       "ID not found",
	INDEX_EMPTY:        "Index is empty",
	THREAD_ERROR:       "Thread error",
	SYSTEM_ERROR:       "System error",
	NOT_IMPLEMENTED:    "Not implemented",
}

var (
	ErrIndexNotInitialized = fmt.Errorf("index not initialized")
	ErrEmptyVector         = fmt.Errorf("empty vector")
)

// toError converts a C error code to a Go error
func toError(code C.int) error {
	if code == C.int(SUCCESS) {
		return nil
	}
	if msg, exists := errorMessages[ErrorCode(code)]; exists {
		return fmt.Errorf(msg)
	}
	return fmt.Errorf("unknown error code: %d", code)
}

// Index represents an index structure in Go
type Index struct {
	ptr *C.Index
}

// AllocIndex creates a new index
func AllocIndex(indexType, method int, dims uint16, icontext idx.IndexContext) (*Index, error) {
	var ptrCtx unsafe.Pointer
	if icontext != nil {
		ptrCtx = icontext.CreateContext()
		defer func() {
			icontext.ReleaseContext(ptrCtx)
		}()
	}

	idx := C.alloc_index(C.int(indexType), C.int(method), C.uint16_t(dims), ptrCtx)
	if idx == nil {
		return nil, fmt.Errorf("failed to allocate index")
	}
	return &Index{ptr: idx}, nil
}

// Insert adds a vector to the index with a given ID
func (idx *Index) Insert(id uint64, vector []float32) error {
	if idx.ptr == nil {
		return ErrIndexNotInitialized
	}
	if len(vector) == 0 {
		return ErrEmptyVector
	}

	cVector := (*C.float)(unsafe.Pointer(&vector[0]))
	return toError(C.insert(idx.ptr, C.uint64_t(id), cVector, C.uint16_t(len(vector))))
}

// Search finds the closest match for a given vector
func (idx *Index) Search(vector []float32, dims int) (*MatchResult, error) {
	if idx.ptr == nil {
		return nil, ErrIndexNotInitialized
	}

	if len(vector) == 0 {
		return nil, ErrEmptyVector
	}

	var cResult C.MatchResult
	cVector := (*C.float)(unsafe.Pointer(&vector[0]))
	err := C.search(idx.ptr, cVector, C.uint16_t(dims), &cResult)
	if e := toError(err); e != nil {
		return nil, e
	}

	return &MatchResult{
		ID:       int(cResult.id),
		Distance: float32(cResult.distance),
	}, nil
}

// SearchN finds the n closest matches for a given vector
func (idx *Index) SearchN(vector []float32, n int) ([]MatchResult, error) {
	if idx.ptr == nil {
		return nil, ErrIndexNotInitialized
	}

	if len(vector) == 0 {
		return nil, ErrEmptyVector
	}

	// Allocate memory for results
	cResults := make([]C.MatchResult, n)
	cVector := (*C.float)(unsafe.Pointer(&vector[0]))

	// Call the C function
	err := C.search_n(idx.ptr, cVector, C.uint16_t(len(vector)),
		(*C.MatchResult)(unsafe.Pointer(&cResults[0])), C.int(n))

	if e := toError(err); e != nil {
		return nil, e
	}

	// Convert C results to Go results
	results := make([]MatchResult, n)
	for i := 0; i < n; i++ {
		results[i] = MatchResult{
			ID:       int(cResults[i].id),
			Distance: float32(cResults[i].distance),
		}
	}

	return results, nil
}

// Delete removes a vector from the index by its ID
func (idx *Index) Delete(id uint64) error {
	if idx.ptr == nil {
		return ErrIndexNotInitialized
	}
	return toError(C.delete(idx.ptr, C.uint64_t(id)))
}

// DestroyIndex releases index memory
func (idx *Index) DestroyIndex() {
	if idx.ptr != nil {
		C.destroy_index(&idx.ptr)
		idx.ptr = nil
	}
}

// Size returns the current number of elements in the index
func (idx *Index) Size() (uint64, error) {
	if idx.ptr == nil {
		return 0, ErrIndexNotInitialized
	}
	var size C.uint64_t

	err := C.size(idx.ptr, &size)
	if e := toError(err); e != nil {
		return 0, e
	}

	return uint64(size), nil
}

// Contains checks whether a given vector ID exists in the index
func (idx *Index) Contains(id uint64) (bool, error) {
	if idx.ptr == nil {
		return false, ErrIndexNotInitialized
	}

	result := C.contains(idx.ptr, C.uint64_t(id))

	return result == 1, nil
}

// GetStats retrieves the internal statistics of the index
func (idx *Index) GetStats() (*IndexStats, error) {
	if idx.ptr == nil {
		return nil, ErrIndexNotInitialized
	}

	var cStats C.IndexStats

	err := C.stats(idx.ptr, &cStats)
	if e := toError(err); e != nil {
		return nil, e
	}

	// Convert C stats to Go stats
	stats := &IndexStats{
		Insert: TimeStat{
			Count: uint64(cStats.insert.count),
			Total: float64(cStats.insert.total),
			Last:  float64(cStats.insert.last),
			Min:   float64(cStats.insert.min),
			Max:   float64(cStats.insert.max),
		},
		Delete: TimeStat{
			Count: uint64(cStats.delete.count),
			Total: float64(cStats.delete.total),
			Last:  float64(cStats.delete.last),
			Min:   float64(cStats.delete.min),
			Max:   float64(cStats.delete.max),
		},
		Dump: TimeStat{
			Count: uint64(cStats.dump.count),
			Total: float64(cStats.dump.total),
			Last:  float64(cStats.dump.last),
			Min:   float64(cStats.dump.min),
			Max:   float64(cStats.dump.max),
		},
		Search: TimeStat{
			Count: uint64(cStats.search.count),
			Total: float64(cStats.search.total),
			Last:  float64(cStats.search.last),
			Min:   float64(cStats.search.min),
			Max:   float64(cStats.search.max),
		},
		SearchN: TimeStat{
			Count: uint64(cStats.search_n.count),
			Total: float64(cStats.search_n.total),
			Last:  float64(cStats.search_n.last),
			Min:   float64(cStats.search_n.min),
			Max:   float64(cStats.search_n.max),
		},
	}
	return stats, nil
}
