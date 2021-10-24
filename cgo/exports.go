package main

/*
#include <stdint.h>
#include <stdlib.h>

typedef const char const_char;

struct Strangler_Record {
    const char* HashKey;
    const char* Status;
    const char* CreatedAt;
    const char* ExpiresAt;
};
typedef struct Strangler_Record Strangler_Record;

struct Strangler_HashKeyOptions {
	const char* Name;
	const char* Expression;
};
typedef struct Strangler_HashKeyOptions Strangler_HashKeyOptions;

struct Strangler_Config {
	Strangler_HashKeyOptions* HashKey;
	uintptr_t			  	  Store;
};
typedef struct Strangler_Config Strangler_Config;
*/
import "C"
import (
	eventstrangler "github.com/balzanelli/event-strangler"
	"runtime/cgo"
	"time"
	"unsafe"
)

//export Strangler_Error_Free
func Strangler_Error_Free(err *C.char) {
	C.free(unsafe.Pointer(err))
}

func RecordToCType(record *eventstrangler.Record) *C.struct_Strangler_Record {
	result := (*C.struct_Strangler_Record)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_Strangler_Record{}))))
	result.HashKey = C.CString(record.HashKey)
	result.Status = C.CString(string(record.Status))
	result.CreatedAt = C.CString(record.CreatedAt.String())
	result.ExpiresAt = C.CString(record.ExpiresAt.String())
	return result
}

func RecordToGoType(record *C.struct_Strangler_Record) (*eventstrangler.Record, error) {
	createdAt, err := time.Parse(time.RFC3339, C.GoString(record.CreatedAt))
	if err != nil {
		return nil, err
	}
	expiresAt, err := time.Parse(time.RFC3339, C.GoString(record.ExpiresAt))
	if err != nil {
		return nil, err
	}
	return &eventstrangler.Record{
		HashKey:   C.GoString(record.HashKey),
		Status:    eventstrangler.RecordState(C.GoString(record.Status)),
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}, nil
}

//export Record_Free
func Record_Free(record *C.struct_Strangler_Record) {
	C.free(unsafe.Pointer(record))
}

//export LevelDBStore_New
func LevelDBStore_New(filepath *C.const_char) (C.uintptr_t, *C.char) {
	f := C.GoString(filepath)
	leveldb, err := eventstrangler.NewLevelDBStore(f)
	if err != nil {
		return C.uintptr_t(0), C.CString(err.Error())
	}
	return C.uintptr_t(cgo.NewHandle(leveldb)), nil
}

//export LevelDBStore_Free
func LevelDBStore_Free(leveldb C.uintptr_t) {
	cgo.Handle(leveldb).Delete()
}

//export LRUCacheStore_New
func LRUCacheStore_New() C.uintptr_t {
	lru := eventstrangler.NewLRUCacheStore()
	return C.uintptr_t(cgo.NewHandle(lru))
}

//export LRUCacheStore_Free
func LRUCacheStore_Free(lru C.uintptr_t) {
	cgo.Handle(lru).Delete()
}

//export Store_Exists
func Store_Exists(store C.uintptr_t, hash_key *C.const_char) (bool, *C.char) {
	k := C.GoString(hash_key)
	exists, err := cgo.Handle(store).Value().(eventstrangler.Store).
		Exists(k)
	if err != nil {
		return false, C.CString(err.Error())
	}
	return exists, nil
}

//export Store_Get
func Store_Get(store C.uintptr_t, hash_key *C.const_char) (*C.struct_Strangler_Record, *C.char) {
	k := C.GoString(hash_key)
	record, err := cgo.Handle(store).Value().(eventstrangler.Store).
		Get(k)
	if err != nil {
		return nil, C.CString(err.Error())
	}
	return RecordToCType(record), nil
}

//export Store_Put
func Store_Put(store C.uintptr_t, hash_key *C.const_char, record *C.struct_Strangler_Record, time_to_live int) *C.char {
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

//export Store_Delete
func Store_Delete(store C.uintptr_t, hash_key *C.const_char) *C.char {
	k := C.GoString(hash_key)
	err := cgo.Handle(store).Value().(eventstrangler.Store).
		Delete(k)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export Store_Close
func Store_Close(store C.uintptr_t) *C.char {
	err := cgo.Handle(store).Value().(eventstrangler.Store).
		Close()
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

func HashKeyOptionsToCType(opt *eventstrangler.HashKeyOptions) *C.struct_Strangler_HashKeyOptions {
	result := (*C.struct_Strangler_HashKeyOptions)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_Strangler_HashKeyOptions{}))))
	result.Name = C.CString(opt.Name)
	result.Expression = C.CString(opt.Expression)
	return result
}

func HashKeyOptionsToGoType(opt *C.struct_Strangler_HashKeyOptions) *eventstrangler.HashKeyOptions {
	return &eventstrangler.HashKeyOptions{
		Name:       C.GoString(opt.Name),
		Expression: C.GoString(opt.Expression),
	}
}

//export HashKeyOptions_Free
func HashKeyOptions_Free(opt *C.struct_Strangler_HashKeyOptions) {
	C.free(unsafe.Pointer(opt))
}

func ConfigToCType(opt *eventstrangler.Config) *C.struct_Strangler_Config {
	if opt == nil {
		return nil
	}
	result := (*C.struct_Strangler_Config)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_Strangler_Config{}))))
	if opt.HashKey != nil {
		result.HashKey = HashKeyOptionsToCType(opt.HashKey)
	}
	return result
}

func ConfigToGoType(opt *C.struct_Strangler_Config) *eventstrangler.Config {
	if opt == nil {
		return nil
	}
	result := &eventstrangler.Config{}
	if opt.HashKey != nil {
		result.HashKey = HashKeyOptionsToGoType(opt.HashKey)
	}
	if opt.Store != C.uintptr_t(0) {
		result.Store = cgo.Handle(opt.Store).Value().(eventstrangler.Store)
	}
	return result
}

//export Strangler_New
func Strangler_New(config *C.struct_Strangler_Config) (C.uintptr_t, *C.char) {
	c := ConfigToGoType(config)
	result, err := eventstrangler.New(c)
	if err != nil {
		return C.uintptr_t(0), C.CString(err.Error())
	}
	return C.uintptr_t(cgo.NewHandle(result)), nil
}

//export Strangler_Free
func Strangler_Free(strangler C.uintptr_t) {
	cgo.Handle(strangler).Delete()
}

//export Strangler_Complete
func Strangler_Complete(strangler C.uintptr_t, hash_key *C.const_char) *C.char {
	k := C.GoString(hash_key)
	err := cgo.Handle(strangler).Value().(*eventstrangler.Strangler).
		Complete(k)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export Strangler_Purge
func Strangler_Purge(strangler C.uintptr_t, hash_key *C.const_char) *C.char {
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
