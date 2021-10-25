from contextlib import contextmanager

import event_strangler_binding
from .store import Store


@contextmanager
def get_lru_cache_store() -> Store:
    result = event_strangler_binding.strangler_lru_cache_store_new()
    try:
        yield Store(result)
    finally:
        event_strangler_binding.strangler_lru_cache_store_free(result)
