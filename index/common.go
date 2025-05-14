package index

import "unsafe"

type IndexType int

const (
	FlatIndexType = IndexType(0) //define FLAT_INDEX    0x00
	NSWIndexType  = IndexType(2) // define NSW_INDEX     0x02
	HNSWIndexType = IndexType(3) // define HNSW_INDEX    0x03
)

type MethodType int

const (
	L2NORM = MethodType(iota) //define L2NORM 0x00
	COSINE                    //define COSINE 0x01
)

type IndexContext interface {
	CreateContext() unsafe.Pointer
	ReleaseContext(ptr unsafe.Pointer)
}
