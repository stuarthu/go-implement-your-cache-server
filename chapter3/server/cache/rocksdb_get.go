package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"errors"
	"unsafe"
)

func (c *rocksdbCache) Get(key string) ([]byte, error) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var length C.size_t
	v := C.rocksdb_get(c.db, c.ro, k, C.size_t(len(key)), &length, &c.e)
	if c.e != nil {
		return nil, errors.New(C.GoString(c.e))
	}
	defer C.free(unsafe.Pointer(v))
	return C.GoBytes(unsafe.Pointer(v), C.int(length)), nil
}
