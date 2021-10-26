#include <gtest/gtest.h>
#include "event_strangler.h"
#include "include.h"

TEST(EventStrangler, EventStranglerNew) {
  const auto store = EventStranglerLRUCacheStoreNew();

  auto hash_key_options = EventStranglerHashKeyOptions{
      .Name = kHashKeyName,
      .Expression = kHashKeyExpression,
  };
  auto config = EventStranglerConfig{
      .HashKey = &hash_key_options,
      .Store = store,
  };
  const auto result = EventStranglerNew(&config);
  if (result.r1) {
    EventStranglerLRUCacheStoreFree(store);
    ASSERT_EQ(result.r1, nullptr);
  }
  EventStranglerLRUCacheStoreFree(store);
  ASSERT_NE(result.r0, 0);
}

TEST(EventStrangler, EventStranglerComplete) {
  const auto store = EventStranglerLRUCacheStoreNew();

  auto hash_key_options = EventStranglerHashKeyOptions{
      .Name = kHashKeyName,
      .Expression = kHashKeyExpression,
  };
  auto config = EventStranglerConfig{
      .HashKey = &hash_key_options,
      .Store = store,
  };
  const auto result = EventStranglerNew(&config);
  if (result.r1) {
    EventStranglerLRUCacheStoreFree(store);
    ASSERT_EQ(result.r1, nullptr);
  }

  auto record = EventStranglerRecord{
      .HashKey = kHashKey,
      .Status = kRecordStatus,
      .CreatedAt = const_cast<char*>("2021-01-01T00:00:00Z"),
      .ExpiresAt = const_cast<char*>("2021-01-01T00:00:00Z"),
  };
  auto err = EventStranglerStorePut(store, kHashKey, &record, 0);
  if (err) {
    EventStranglerFree(result.r0);
    EventStranglerLRUCacheStoreFree(store);
    ASSERT_STREQ(err, "");
  }
  err = EventStranglerComplete(result.r0, kHashKey);
  if (result.r1) {
    EventStranglerFree(result.r0);
    EventStranglerLRUCacheStoreFree(store);
    ASSERT_EQ(err, nullptr);
  }
  const auto r = EventStranglerStoreGet(store, kHashKey);
  if (r.r1) {
    EventStranglerFree(result.r0);
    EventStranglerLRUCacheStoreFree(store);
    ASSERT_EQ(err, nullptr);
  }
  ASSERT_STREQ(r.r0->HashKey, kHashKey);
  ASSERT_STREQ(r.r0->Status, "COMPLETE");
}
