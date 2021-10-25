from contextlib import contextmanager

import event_strangler_binding
from .store import Store


@contextmanager
def get_leveldb_store(filepath: str) -> Store:
    result = event_strangler_binding.strangler_leveldb_store_new(filepath)
    if result.r1:
        raise ValueError(result.r1)
    try:
        yield Store(result.r0)
    finally:
        event_strangler_binding.strangler_leveldb_store_free(result.r0)
