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

const BATCH_SIZE = 100

func flush_batch(db *C.rocksdb_t, b *C.rocksdb_writebatch_t, o *C.rocksdb_writeoptions_t) {
	var e *C.char
	C.rocksdb_write(db, o, b, &e)
	if e != nil {
		panic(C.GoString(e))
	}
	C.rocksdb_writebatch_clear(b)
}

func write_func(db *C.rocksdb_t, c chan *pair, o *C.rocksdb_writeoptions_t) {
	count := 0
	t := time.NewTimer(time.Second)
	b := C.rocksdb_writebatch_create()
	for {
		select {
		case p := <-c:
			count++
			key := C.CString(p.k)
			value := C.CBytes(p.v)
			C.rocksdb_writebatch_put(b, key, C.size_t(len(p.k)), (*C.char)(value), C.size_t(len(p.v)))
			C.free(unsafe.Pointer(key))
			C.free(value)
			if count == BATCH_SIZE {
				flush_batch(db, b, o)
				count = 0
			}
			if !t.Stop() {
				<-t.C
			}
			t.Reset(time.Second)
		case <-t.C:
			if count != 0 {
				flush_batch(db, b, o)
				count = 0
			}
			t.Reset(time.Second)
		}
	}
}

func (c *rocksdbCache) Set(key string, value []byte) error {
	c.ch <- &pair{key, value}
	return nil
}
