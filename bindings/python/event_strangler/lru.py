from contextlib import contextmanager

import bindings
from .store import Store


@contextmanager
def get_lru_cache_store() -> Store:
    result = bindings.LRUCacheStoreNew()
    try:
        yield Store(result)
    finally:
        bindings.LRUCacheStoreFree(result)
