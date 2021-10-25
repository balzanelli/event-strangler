from contextlib import contextmanager
from typing import Optional

import event_strangler_binding


class HashKeyOptions:
    def __init__(self, name: Optional[str] = None, expression: Optional[str] = None):
        self.name = name
        self.expression = expression


@contextmanager
def _get_hash_key_options(hash_key_options: HashKeyOptions) -> event_strangler_binding.strangler_hash_key_options:
    result = event_strangler_binding.strangler_hash_key_options()
    try:
        result.name = hash_key_options.name
        result.expression = hash_key_options.expression
        yield result
    finally:
        event_strangler_binding.strangler_hash_key_options_free(result)
