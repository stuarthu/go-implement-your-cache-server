package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"

import "unsafe"

type rocksdbScanner struct {
	i *C.rocksdb_iterator_t
}

func (s *rocksdbScanner) Close() {
	C.rocksdb_iter_destroy(s.i)
}

func (s *rocksdbScanner) Scan() bool {
	C.rocksdb_iter_next(s.i)
	valid := C.rocksdb_iter_valid(s.i)
	return valid == 0
}

func (s *rocksdbScanner) Key() string {
	var length C.size_t
	k := C.rocksdb_iter_key(s.i, &length)
	defer C.free(unsafe.Pointer(k))
	return C.GoString(k)
}

func (s *rocksdbScanner) Value() []byte {
	var length C.size_t
	v := C.rocksdb_iter_value(s.i, &length)
	defer C.free(unsafe.Pointer(v))
	return C.GoBytes(unsafe.Pointer(v), C.int(length))
}

func (c *rocksdbCache) NewScanner() Scanner {
	return &rocksdbScanner{C.rocksdb_create_iterator(c.db, c.ro)}
}
