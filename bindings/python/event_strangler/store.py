from typing import Optional

import bindings
from .record import Record


class Store:
    def __init__(self, c_uintptr):
        self.c_uintptr = c_uintptr

    def exists(self, hash_key: str) -> bool:
        result = bindings.StoreExists(self.c_uintptr, hash_key)
        if result.r1:
            raise ValueError(result.r1)
        return result.r0

    def get(self, hash_key: str) -> Record:
        result = bindings.StoreGet(self.c_uintptr, hash_key)
        if result.r1:
            raise ValueError(result.r1)
        return Record.from_binding(result.r0)

    def put(self, hash_key: str, record: Record, time_to_live: Optional[int] = 0) -> None:
        item = Record.to_binding(record)
        err = bindings.StorePut(self.c_uintptr, hash_key, item, time_to_live)
        if err:
            raise ValueError(err)

    def delete(self, hash_key: str) -> None:
        err = bindings.StoreDelete(self.c_uintptr, hash_key)
        if err:
            raise ValueError(err)
