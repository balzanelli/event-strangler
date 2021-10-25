class Record:
    def __init__(self, hash_key: str, status: str, created_at: str, expires_at: str):
        self.hash_key = hash_key
        self.status = status
        self.created_at = created_at
        self.expires_at = expires_at
