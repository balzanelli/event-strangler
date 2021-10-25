from contextlib import contextmanager

import event_strangler_binding
from .config import Config, _get_config


class Strangler:
    def __init__(self, c_uintptr):
        self._c_uintptr = c_uintptr

    def complete(self, hash_key: str) -> None:
        err = event_strangler_binding.strangler_complete(self._c_uintptr, hash_key)
        if err:
            raise ValueError(err)

    def purge(self, hash_key: str) -> None:
        err = event_strangler_binding.strangler_purge(self._c_uintptr, hash_key)
        if err:
            raise ValueError(err)


@contextmanager
def build(config: Config) -> Strangler:
    with _get_config(config) as strangler_config:
        result = event_strangler_binding.strangler_new(strangler_config)
        if result.r1:
            raise ValueError(result.r1)
        try:
            yield Strangler(result.r0)
        finally:
            event_strangler_binding.strangler_free(result.r0)
