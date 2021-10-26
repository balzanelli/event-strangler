from typing import Optional

import bindings


class Record:
    def __init__(self, hash_key: Optional[str], status: Optional[str], created_at: Optional[str],
                 expires_at: Optional[str]):
        self.hash_key: Optional[str] = hash_key
        self.status: Optional[str] = status
        self.created_at: Optional[str] = created_at
        self.expires_at: Optional[str] = expires_at

    @staticmethod
    def to_binding(model) -> bindings.Record:
        result = bindings.Record()
        result.HashKey = model.hash_key
        result.Status = model.status
        result.CreatedAt = model.created_at
        result.ExpiresAt = model.expires_at
        return result

    @staticmethod
    def from_binding(model: bindings.Record):
        return Record(
            hash_key=model.HashKey,
            status=model.Status,
            created_at=model.CreatedAt,
            expires_at=model.ExpiresAt
        )
