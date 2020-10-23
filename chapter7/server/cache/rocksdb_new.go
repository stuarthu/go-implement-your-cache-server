package cache

// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import "runtime"

func newRocksdbCache() *rocksdbCache {
	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, C.int(runtime.NumCPU()))
	C.rocksdb_options_set_create_if_missing(options, 1)
	var e *C.char
	db := C.rocksdb_open(options, C.CString("/tmp/rocksdb"), &e)
	if e != nil {
		panic(C.GoString(e))
	}
	C.rocksdb_options_destroy(options)
	c := make(chan *pair, 5000)
	wo := C.rocksdb_writeoptions_create()
	go write_func(db, c, wo)
	return &rocksdbCache{db, C.rocksdb_readoptions_create(), wo, e, c}
}
