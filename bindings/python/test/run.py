import event_strangler.hash_key
import event_strangler.lru
import event_strangler.record
import event_strangler.strangler

HASH_KEY = "d260350bd70e2dd5bd6b1e9b1047f58b9b58e5ed2eb912081d1f902a1370bfdc"

with event_strangler.lru.get_lru_cache_store() as store:
    hash_key_options = event_strangler.hash_key.HashKeyOptions(
        name="event-strangler/tests/python",
        expression="[subject, transaction_id]",
    )

    with event_strangler.strangler.build(event_strangler.strangler.Config(hash_key_options, store)) as strangler:
        store.put(HASH_KEY, event_strangler.record.Record(
            hash_key=HASH_KEY,
            status="PROCESSING",
            created_at="2021-01-01T00:00:00Z",
            expires_at="2021-01-01T00:00:00Z"
        ))
        record = store.get(HASH_KEY)
        print(f'record->hash_key: {record.hash_key}\n' +
              f'record->status: {record.status}')

        strangler.complete(HASH_KEY)
        record = store.get(HASH_KEY)
        print(f'record->hash_key: {record.hash_key}\n' +
              f'record->status: {record.status}')
