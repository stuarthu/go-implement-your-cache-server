package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"errors"
	"regexp"
	"strconv"
	"unsafe"
)

func (r *rocksdbCache) Get(key string) ([]byte, error) {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	var length C.size_t
	value := C.rocksdb_get(r.db, r.readoptions, ckey, C.size_t(len(key)), &length, &r.err)
	if r.err != nil {
		return nil, errors.New(C.GoString(r.err))
	}
	b := C.GoBytes(unsafe.Pointer(value), C.int(length))
	C.free(unsafe.Pointer(value))
	return b, nil
}

func (r *rocksdbCache) Del(key string) error {
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	C.rocksdb_delete(r.db, r.writeoptions, ckey, C.size_t(len(key)), &r.err)
	if r.err != nil {
		return errors.New(C.GoString(r.err))
	}
	return nil
}

func (r *rocksdbCache) GetStat() Stat {
	prop := C.CString("rocksdb.aggregated-table-properties")
	defer C.free(unsafe.Pointer(prop))
	value := C.rocksdb_property_value(r.db, prop)
	defer C.free(unsafe.Pointer(value))
	stat := C.GoString(value)
	pattern := regexp.MustCompile(`([^;]+)=([^;]+);`)
	s := Stat{}
	for _, submatches := range pattern.FindAllStringSubmatch(stat, -1) {
		if submatches[1] == " # entries" {
			s.Count, _ = strconv.ParseInt(submatches[2], 10, 64)
		} else if submatches[1] == " raw key size" {
			s.KeySize, _ = strconv.ParseInt(submatches[2], 10, 64)
		} else if submatches[1] == " raw value size" {
			s.ValueSize, _ = strconv.ParseInt(submatches[2], 10, 64)
		}
	}
	return s
}
