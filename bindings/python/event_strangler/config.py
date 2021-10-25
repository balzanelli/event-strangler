from contextlib import contextmanager

import event_strangler_binding
from .hash_key import HashKeyOptions, _get_hash_key_options
from .store import Store


class Config:
    def __init__(self, hash_key_options: HashKeyOptions, store: Store):
        self.hash_key_options = hash_key_options
        self.store = store


@contextmanager
def _get_config(config: Config) -> event_strangler_binding.strangler_config:
    result = event_strangler_binding.strangler_config()
    try:
        with _get_hash_key_options(config.hash_key_options) as hash_key_options:
            result.hash_key = hash_key_options
            result.store = config.store.c_uintptr
            yield result
    finally:
        pass
