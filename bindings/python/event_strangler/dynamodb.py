from contextlib import contextmanager
from typing import Optional

import bindings
from .store import Store


@contextmanager
def get_dynamodb_store(table_name: str, profile: Optional[str] = None) -> Store:
    result = bindings.DynamoDBStoreNew(table_name, profile)
    if result.r1:
        raise ValueError(result.r1)
    try:
        yield Store(result.r0)
    finally:
        bindings.DynamoDBStoreFree(result.r0)
