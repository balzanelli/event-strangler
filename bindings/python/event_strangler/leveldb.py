from contextlib import contextmanager

import bindings
from .store import Store


@contextmanager
def leveldb_store(file: str) -> Store:
    result = bindings.LevelDBStoreNew(file)
    if result.r1:
        raise ValueError(result.r1)
    try:
        yield Store(result.r0)
    finally:
        bindings.LevelDBStoreFree(result.r0)
