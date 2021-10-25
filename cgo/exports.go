package main

/*
#include <stdint.h>
#include <stdlib.h>

struct strangler_record {
    char* hash_key;
    char* status;
    char* created_at;
    char* expires_at;
};
typedef struct strangler_record strangler_record;

struct strangler_hash_key_options {
	char* name;
	char* expression;
};
typedef struct strangler_hash_key_options strangler_hash_key_options;

struct strangler_config {
	strangler_hash_key_options* hash_key;
	uintptr_t			  	    store;
};
typedef struct strangler_config strangler_config;
*/
import "C"
import (
	eventstrangler "github.com/balzanelli/event-strangler"
	"runtime/cgo"
	"time"
	"unsafe"
)

func RecordToCType(record *eventstrangler.Record) *C.struct_strangler_record {
	result := (*C.struct_strangler_record)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_strangler_record{}))))
	result.hash_key = C.CString(record.HashKey)
	result.status = C.CString(string(record.Status))
	result.created_at = C.CString(record.CreatedAt.String())
	result.expires_at = C.CString(record.ExpiresAt.String())
	return result
}

func RecordToGoType(record *C.struct_strangler_record) (*eventstrangler.Record, error) {
	createdAt, err := time.Parse(time.RFC3339, C.GoString(record.created_at))
	if err != nil {
		return nil, err
	}
	expiresAt, err := time.Parse(time.RFC3339, C.GoString(record.expires_at))
	if err != nil {
		return nil, err
	}
	return &eventstrangler.Record{
		HashKey:   C.GoString(record.hash_key),
		Status:    eventstrangler.RecordState(C.GoString(record.status)),
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

//export strangler_record_free
func strangler_record_free(record *C.struct_strangler_record) {
	C.free(unsafe.Pointer(record))
}

//export strangler_leveldb_store_new
func strangler_leveldb_store_new(filepath *C.char) (C.uintptr_t, *C.char) {
	f := C.GoString(filepath)
	leveldb, err := eventstrangler.NewLevelDBStore(f)
	if err != nil {
		return C.uintptr_t(0), C.CString(err.Error())
	}
	return C.uintptr_t(cgo.NewHandle(leveldb)), nil
}

//export strangler_leveldb_store_free
func strangler_leveldb_store_free(leveldb C.uintptr_t) {
	cgo.Handle(leveldb).Delete()
}

//export strangler_lru_cache_store_new
func strangler_lru_cache_store_new() C.uintptr_t {
	lru := eventstrangler.NewLRUCacheStore()
	return C.uintptr_t(cgo.NewHandle(lru))
}

//export strangler_lru_cache_store_free
func strangler_lru_cache_store_free(lru C.uintptr_t) {
	cgo.Handle(lru).Delete()
}

//export strangler_store_exists
func strangler_store_exists(store C.uintptr_t, hash_key *C.char) (bool, *C.char) {
	k := C.GoString(hash_key)
	exists, err := cgo.Handle(store).Value().(eventstrangler.Store).
		Exists(k)
	if err != nil {
		return false, C.CString(err.Error())
	}
	return exists, nil
}

//export strangler_store_get
func strangler_store_get(store C.uintptr_t, hash_key *C.char) (*C.struct_strangler_record, *C.char) {
	k := C.GoString(hash_key)
	record, err := cgo.Handle(store).Value().(eventstrangler.Store).
		Get(k)
	if err != nil {
		return nil, C.CString(err.Error())
	}
	return RecordToCType(record), nil
}

//export strangler_store_put
func strangler_store_put(store C.uintptr_t, hash_key *C.char, record *C.struct_strangler_record, time_to_live int) *C.char {
	k := C.GoString(hash_key)
	r, err := RecordToGoType(record)
	if err != nil {
		return C.CString(err.Error())
	}
	err = cgo.Handle(store).Value().(eventstrangler.Store).
		Put(k, r, time_to_live)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export strangler_store_delete
func strangler_store_delete(store C.uintptr_t, hash_key *C.char) *C.char {
	k := C.GoString(hash_key)
	err := cgo.Handle(store).Value().(eventstrangler.Store).
		Delete(k)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export strangler_store_close
func strangler_store_close(store C.uintptr_t) *C.char {
	err := cgo.Handle(store).Value().(eventstrangler.Store).
		Close()
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

func HashKeyOptionsToCType(opt *eventstrangler.HashKeyOptions) *C.struct_strangler_hash_key_options {
	result := (*C.struct_strangler_hash_key_options)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_strangler_hash_key_options{}))))
	result.name = C.CString(opt.Name)
	result.expression = C.CString(opt.Expression)
	return result
}

func HashKeyOptionsToGoType(opt *C.struct_strangler_hash_key_options) *eventstrangler.HashKeyOptions {
	return &eventstrangler.HashKeyOptions{
		Name:       C.GoString(opt.name),
		Expression: C.GoString(opt.expression),
	}
}

//export strangler_hash_key_options_free
func strangler_hash_key_options_free(opt *C.struct_strangler_hash_key_options) {
	C.free(unsafe.Pointer(opt))
}

func ConfigToCType(opt *eventstrangler.Config) *C.struct_strangler_config {
	if opt == nil {
		return nil
	}
	result := (*C.struct_strangler_config)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_strangler_config{}))))
	if opt.HashKey != nil {
		result.hash_key = HashKeyOptionsToCType(opt.HashKey)
	}
	return result
}

func ConfigToGoType(opt *C.struct_strangler_config) *eventstrangler.Config {
	if opt == nil {
		return nil
	}
	result := &eventstrangler.Config{}
	if opt.hash_key != nil {
		result.HashKey = HashKeyOptionsToGoType(opt.hash_key)
	}
	if opt.store != C.uintptr_t(0) {
		result.Store = cgo.Handle(opt.store).Value().(eventstrangler.Store)
	}
	return result
}

//export strangler_new
func strangler_new(config *C.struct_strangler_config) (C.uintptr_t, *C.char) {
	c := ConfigToGoType(config)
	result, err := eventstrangler.New(c)
	if err != nil {
		return C.uintptr_t(0), C.CString(err.Error())
	}
	return C.uintptr_t(cgo.NewHandle(result)), nil
}

//export strangler_free
func strangler_free(strangler C.uintptr_t) {
	cgo.Handle(strangler).Delete()
}

//export strangler_complete
func strangler_complete(strangler C.uintptr_t, hash_key *C.char) *C.char {
	k := C.GoString(hash_key)
	err := cgo.Handle(strangler).Value().(*eventstrangler.Strangler).
		Complete(k)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export strangler_purge
func strangler_purge(strangler C.uintptr_t, hash_key *C.char) *C.char {
	k := C.GoString(hash_key)
	err := cgo.Handle(strangler).Value().(*eventstrangler.Strangler).
		Purge(k)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

func main() {
}
