package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"errors"
	"runtime"
	"unsafe"
)

type rocksdbCache struct {
	db           *C.rocksdb_t
	readoptions  *C.rocksdb_readoptions_t
	writeoptions *C.rocksdb_writeoptions_t
	err          *C.char
}

func (r *rocksdbCache) Set(key string, value []byte) error {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	cvalue := C.CBytes(value)
	defer C.free(cvalue)
	C.rocksdb_put(r.db, r.writeoptions, ckey, C.size_t(len(key)), (*C.char)(cvalue), C.size_t(len(value)), &r.err)
	if r.err != nil {
		return errors.New(C.GoString(r.err))
	}
	return nil
}

func newRocksdbCache() *rocksdbCache {
	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	C.rocksdb_options_set_create_if_missing(options, 1)
	var err *C.char
	db := C.rocksdb_open(options, C.CString("/mnt/rocksdb"), &err)
	if err != nil {
		panic(C.GoString(err))
	}
	C.rocksdb_options_destroy(options)
	return &rocksdbCache{db, C.rocksdb_readoptions_create(), C.rocksdb_writeoptions_create(), err}
}
