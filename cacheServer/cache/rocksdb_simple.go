package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"unsafe"
)

type simpleRocksdbCache struct {
	db           *C.rocksdb_t
	readoptions  *C.rocksdb_readoptions_t
	writeoptions *C.rocksdb_writeoptions_t
}

func (r *simpleRocksdbCache) read(key string) []byte {
	var length C.size_t
	var err *C.char
	ckey := C.CString(key)
	value := C.rocksdb_get(r.db, r.readoptions, ckey, C.size_t(len(key)), &length, &err)
	if err != nil {
		panic(C.GoString(err))
	}
	b := C.GoBytes(unsafe.Pointer(value), C.int(length))
	C.free(unsafe.Pointer(ckey))
	C.free(unsafe.Pointer(value))
	return b
}

func (r *simpleRocksdbCache) write(key string, value string) {
	var err *C.char
	ckey := C.CString(key)
	cvalue := C.CString(value)
	C.rocksdb_put(r.db, r.writeoptions, ckey, C.size_t(len(key)), cvalue, C.size_t(len(value)), &err)
	if err != nil {
		panic(C.GoString(err))
	}
	C.free(unsafe.Pointer(ckey))
	C.free(unsafe.Pointer(cvalue))
}

func (r *simpleRocksdbCache) set(key string, value []byte) {
	r.write(key, string(value))
}

func (r *simpleRocksdbCache) get(key string) []byte {
	return r.read(key)
}

func NewSimpleRocksdbCache() *simpleRocksdbCache {
	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, 16)
	C.rocksdb_options_optimize_level_style_compaction(options, 512*1024*1024)
	C.rocksdb_options_set_create_if_missing(options, 1)
	readoptions := C.rocksdb_readoptions_create()
	writeoptions := C.rocksdb_writeoptions_create()
	var err *C.char
	db := C.rocksdb_open(options, C.CString("/mnt/rocksdb"), &err)
	if err != nil {
		panic(C.GoString(err))
	}
	C.rocksdb_options_destroy(options)
	return &simpleRocksdbCache{db, readoptions, writeoptions}
}
