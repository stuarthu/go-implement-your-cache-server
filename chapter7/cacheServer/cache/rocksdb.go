package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"time"
	"unsafe"
)

type readTask struct {
	key string
	ret chan []byte
}

type writeTask struct {
	key   string
	value []byte
}

type rocksdbCache struct {
	db          *C.rocksdb_t
	readoptions *C.rocksdb_readoptions_t
	writeChan   chan *writeTask
}

func (r *rocksdbCache) Get(key string) []byte {
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

const BATCH_SIZE = 100

func flush_batch(db *C.rocksdb_t, batch *C.rocksdb_writebatch_t) {
	writeoptions := C.rocksdb_writeoptions_create()
	var err *C.char
	C.rocksdb_write(db, writeoptions, batch, &err)
	if err != nil {
		panic(C.GoString(err))
	}
	C.rocksdb_writebatch_clear(batch)
	C.rocksdb_writeoptions_destroy(writeoptions)
}

func write_func(db *C.rocksdb_t, c chan *writeTask) {
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
				flush_batch(db, batch)
				count = 0
			}
		case <-timer.C:
			if count != 0 {
				flush_batch(db, batch)
				count = 0
			}
			timer.Reset(time.Second)
		}
	}
}

func (r *rocksdbCache) Set(key string, value []byte) {
	r.writeChan <- &writeTask{key, value}
}

func NewRocksdbCache() *rocksdbCache {
	options := C.rocksdb_options_create()
	C.rocksdb_options_increase_parallelism(options, 16)
	C.rocksdb_options_optimize_level_style_compaction(options, 512*1024*1024)
	C.rocksdb_options_set_create_if_missing(options, 1)
	var err *C.char
	db := C.rocksdb_open(options, C.CString("/mnt/rocksdb"), &err)
	if err != nil {
		panic(C.GoString(err))
	}
	C.rocksdb_options_destroy(options)
	writeChan := make(chan *writeTask, 5000)

	go write_func(db, writeChan)
	return &rocksdbCache{db, C.rocksdb_readoptions_create(), writeChan}
}
