package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../../../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../../../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -O3
import "C"
import (
	"errors"
	"regexp"
	"runtime"
	"strconv"
	"unsafe"
)

type rocksdbCache struct {
	db           *C.rocksdb_t
	readoptions  *C.rocksdb_readoptions_t
	writeoptions *C.rocksdb_writeoptions_t
	err          *C.char
}

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
	pattern := regexp.MustCompile(`(?P<key>[^;]+)=(?P<value>[^;]+);`)
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

func NewRocksdbCache() *rocksdbCache {
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
