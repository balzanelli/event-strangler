from typing import Optional

import event_strangler_binding
from .record import Record


class Store:
    def __init__(self, c_uintptr):
        self.c_uintptr = c_uintptr

    def exists(self, hash_key: str) -> bool:
        result = event_strangler_binding.strangler_store_exists(self.c_uintptr, hash_key)
        if result.r1:
            raise ValueError(result.r1)
        return result.r0

    def get(self, hash_key: str) -> Record:
        result = event_strangler_binding.strangler_store_get(self.c_uintptr, hash_key)
        if result.r1:
            raise ValueError(result.r1)
        return Record(
            hash_key=result.r0.hash_key,
            status=result.r0.status,
            created_at=result.r0.created_at,
            expires_at=result.r0.expires_at
        )

    def put(self, hash_key: str, record: Record, time_to_live: Optional[int] = 0) -> None:
        item = event_strangler_binding.strangler_record()
        item.hash_key = record.hash_key
        item.status = record.status
        item.created_at = record.created_at
        item.expires_at = record.expires_at

        err = event_strangler_binding.strangler_store_put(self.c_uintptr, hash_key, item, time_to_live)
        if err:
            raise ValueError(err)

    def delete(self, hash_key: str) -> None:
        err = event_strangler_binding.strangler_store_delete(self.c_uintptr, hash_key)
        if err:
            raise ValueError(err)
