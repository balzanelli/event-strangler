from contextlib import contextmanager

import bindings
from .config import Config


class Strangler:
    def __init__(self, c_uintptr):
        self.c_uintptr = c_uintptr

    def complete(self, hash_key: str) -> None:
        err = bindings.Complete(self.c_uintptr, hash_key)
        if err:
            raise ValueError(err)

    def purge(self, hash_key: str) -> None:
        err = bindings.Purge(self.c_uintptr, hash_key)
        if err:
            raise ValueError(err)


@contextmanager
def build(config: Config) -> Strangler:
    result = bindings.New(Config.to_binding(config))
    if result.r1:
        raise ValueError(result.r1)
    try:
        yield Strangler(result.r0)
    finally:
        bindings.Free(result.r0)
