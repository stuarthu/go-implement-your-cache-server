package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"runtime"
	"time"
	"unsafe"
)

type writeTask struct {
	key   string
	value []byte
}

type rocksdbCache struct {
	db           *C.rocksdb_t
	readoptions  *C.rocksdb_readoptions_t
	writeoptions *C.rocksdb_writeoptions_t
	err          *C.char
	writeChan    chan *writeTask
}

const BATCH_SIZE = 100

func flush_batch(db *C.rocksdb_t, batch *C.rocksdb_writebatch_t, wo *C.rocksdb_writeoptions_t) {
	var err *C.char
	C.rocksdb_write(db, wo, batch, &err)
	if err != nil {
		panic(C.GoString(err))
	}
	C.rocksdb_writebatch_clear(batch)
}

func write_func(db *C.rocksdb_t, c chan *writeTask, wo *C.rocksdb_writeoptions_t) {
	count := 0
	timer := time.NewTimer(time.Second)
	batch := C.rocksdb_writebatch_create()
	for {
		select {
		case t := <-c:
			count++
			key := C.CString(t.key)
			value := C.CBytes(t.value)
			C.rocksdb_writebatch_put(batch, key, C.size_t(len(t.key)), (*C.char)(value), C.size_t(len(t.value)))
			C.free(unsafe.Pointer(key))
			C.free(value)
			if count == BATCH_SIZE {
				flush_batch(db, batch, wo)
				count = 0
			}
		case <-timer.C:
			if count != 0 {
				flush_batch(db, batch, wo)
				count = 0
			}
			timer.Reset(time.Second)
		}
	}
}

func (r *rocksdbCache) Set(key string, value []byte) error {
	r.writeChan <- &writeTask{key, value}
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
	writeChan := make(chan *writeTask, 5000)

	writeoptions := C.rocksdb_writeoptions_create()
	go write_func(db, writeChan, writeoptions)
	return &rocksdbCache{db, C.rocksdb_readoptions_create(), writeoptions, nil, writeChan}
}
