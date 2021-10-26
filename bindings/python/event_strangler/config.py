from typing import Any, Optional

import bindings
from .hash_key import HashKeyOptions


class Config:
    def __init__(self,
                 hash_key: Optional[HashKeyOptions] = None,
                 store: Optional[Any] = None):
        self.hash_key: Optional[HashKeyOptions] = hash_key
        self.store: Optional[Any] = store

    @staticmethod
    def to_binding(model) -> bindings.Config:
        result = bindings.Config()
        result.HashKey = HashKeyOptions.to_binding(model.hash_key)
        result.Store = model.store.c_uintptr
        return result
