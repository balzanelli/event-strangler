package main

/*
#include <stdint.h>
#include <stdlib.h>

struct EventStranglerRecord {
    char* HashKey;
    char* Status;
    char* CreatedAt;
    char* ExpiresAt;
};
typedef struct EventStranglerRecord EventStranglerRecord;

struct EventStranglerHashKeyOptions {
	char* Name;
	char* Expression;
};
typedef struct EventStranglerHashKeyOptions EventStranglerHashKeyOptions;

struct EventStranglerConfig {
	EventStranglerHashKeyOptions* HashKey;
	uintptr_t			  	      Store;
};
typedef struct EventStranglerConfig EventStranglerConfig;
*/
import "C"
import (
	eventstrangler "github.com/balzanelli/event-strangler"
	"runtime/cgo"
	"time"
	"unsafe"
)

func RecordToCType(record *eventstrangler.Record) *C.struct_EventStranglerRecord {
	result := (*C.struct_EventStranglerRecord)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_EventStranglerRecord{}))))
	result.HashKey = C.CString(record.HashKey)
	result.Status = C.CString(string(record.Status))
	result.CreatedAt = C.CString(record.CreatedAt.String())
	result.ExpiresAt = C.CString(record.ExpiresAt.String())
	return result
}

func RecordToGoType(record *C.struct_EventStranglerRecord) (*eventstrangler.Record, error) {
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

//export EventStranglerDynamoDBStoreNew
func EventStranglerDynamoDBStoreNew(tableName *C.char, profile *C.char) (C.uintptr_t, *C.char) {
	t := C.GoString(tableName)
	p := C.GoString(profile)
	dynamoDB, err := eventstrangler.NewDynamoDBStore(t, p)
	if err != nil {
		return C.uintptr_t(0), C.CString(err.Error())
	}
	return C.uintptr_t(cgo.NewHandle(dynamoDB)), nil
}

//export EventStranglerDynamoDBStoreFree
func EventStranglerDynamoDBStoreFree(dynamoDB C.uintptr_t) {
	cgo.Handle(dynamoDB).Delete()
}

//export EventStranglerLevelDBStoreNew
func EventStranglerLevelDBStoreNew(filepath *C.char) (C.uintptr_t, *C.char) {
	f := C.GoString(filepath)
	levelDB, err := eventstrangler.NewLevelDBStore(f)
	if err != nil {
		return C.uintptr_t(0), C.CString(err.Error())
	}
	return C.uintptr_t(cgo.NewHandle(levelDB)), nil
}

//export EventStranglerLevelDBStoreFree
func EventStranglerLevelDBStoreFree(levelDB C.uintptr_t) {
	cgo.Handle(levelDB).Delete()
}

//export EventStranglerLRUCacheStoreNew
func EventStranglerLRUCacheStoreNew() C.uintptr_t {
	lru := eventstrangler.NewLRUCacheStore()
	return C.uintptr_t(cgo.NewHandle(lru))
}

//export EventStranglerLRUCacheStoreFree
func EventStranglerLRUCacheStoreFree(lru C.uintptr_t) {
	cgo.Handle(lru).Delete()
}

//export EventStranglerStoreExists
func EventStranglerStoreExists(store C.uintptr_t, hashKey *C.char) (bool, *C.char) {
	k := C.GoString(hashKey)
	exists, err := cgo.Handle(store).Value().(eventstrangler.Store).
		Exists(k)
	if err != nil {
		return false, C.CString(err.Error())
	}
	return exists, nil
}

//export EventStranglerStoreGet
func EventStranglerStoreGet(store C.uintptr_t, hashKey *C.char) (*C.struct_EventStranglerRecord, *C.char) {
	k := C.GoString(hashKey)
	record, err := cgo.Handle(store).Value().(eventstrangler.Store).
		Get(k)
	if err != nil {
		return nil, C.CString(err.Error())
	}
	return RecordToCType(record), nil
}

//export EventStranglerStorePut
func EventStranglerStorePut(store C.uintptr_t, hashKey *C.char, record *C.struct_EventStranglerRecord,
	timeToLive int) *C.char {
	k := C.GoString(hashKey)
	r, err := RecordToGoType(record)
	if err != nil {
		return C.CString(err.Error())
	}
	err = cgo.Handle(store).Value().(eventstrangler.Store).
		Put(k, r, timeToLive)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export EventStranglerStoreDelete
func EventStranglerStoreDelete(store C.uintptr_t, hashKey *C.char) *C.char {
	k := C.GoString(hashKey)
	err := cgo.Handle(store).Value().(eventstrangler.Store).
		Delete(k)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export EventStranglerStoreClose
func EventStranglerStoreClose(store C.uintptr_t) *C.char {
	err := cgo.Handle(store).Value().(eventstrangler.Store).
		Close()
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

func HashKeyOptionsToCType(hashKeyOptions *eventstrangler.HashKeyOptions) *C.struct_EventStranglerHashKeyOptions {
	result := (*C.struct_EventStranglerHashKeyOptions)(
		C.malloc(C.size_t(unsafe.Sizeof(C.struct_EventStranglerHashKeyOptions{}))))
	result.Name = C.CString(hashKeyOptions.Name)
	result.Expression = C.CString(hashKeyOptions.Expression)
	return result
}

func HashKeyOptionsToGoType(hashKeyOptions *C.struct_EventStranglerHashKeyOptions) *eventstrangler.HashKeyOptions {
	return &eventstrangler.HashKeyOptions{
		Name:       C.GoString(hashKeyOptions.Name),
		Expression: C.GoString(hashKeyOptions.Expression),
	}
}

func ConfigToCType(config *eventstrangler.Config) *C.struct_EventStranglerConfig {
	if config == nil {
		return nil
	}
	result := (*C.struct_EventStranglerConfig)(C.malloc(C.size_t(unsafe.Sizeof(C.struct_EventStranglerConfig{}))))
	if config.HashKey != nil {
		result.HashKey = HashKeyOptionsToCType(config.HashKey)
	}
	return result
}

func ConfigToGoType(config *C.struct_EventStranglerConfig) *eventstrangler.Config {
	if config == nil {
		return nil
	}
	result := &eventstrangler.Config{}
	if config.HashKey != nil {
		result.HashKey = HashKeyOptionsToGoType(config.HashKey)
	}
	if config.Store != C.uintptr_t(0) {
		result.Store = cgo.Handle(config.Store).Value().(eventstrangler.Store)
	}
	return result
}

//export EventStranglerNew
func EventStranglerNew(config *C.struct_EventStranglerConfig) (C.uintptr_t, *C.char) {
	c := ConfigToGoType(config)
	result, err := eventstrangler.New(c)
	if err != nil {
		return C.uintptr_t(0), C.CString(err.Error())
	}
	return C.uintptr_t(cgo.NewHandle(result)), nil
}

//export EventStranglerFree
func EventStranglerFree(strangler C.uintptr_t) {
	cgo.Handle(strangler).Delete()
}

//export EventStranglerComplete
func EventStranglerComplete(strangler C.uintptr_t, hashKey *C.char) *C.char {
	k := C.GoString(hashKey)
	err := cgo.Handle(strangler).Value().(*eventstrangler.Strangler).
		Complete(k)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

//export EventStranglerPurge
func EventStranglerPurge(strangler C.uintptr_t, hashKey *C.char) *C.char {
	k := C.GoString(hashKey)
	err := cgo.Handle(strangler).Value().(*eventstrangler.Strangler).
		Purge(k)
	if err != nil {
		return C.CString(err.Error())
	}
	return nil
}

func main() {
}
