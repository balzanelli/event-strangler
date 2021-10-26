#include <gtest/gtest.h>
#include "event_strangler.h"
#include "include.h"

TEST(EventStranglerLRUCacheStore, EventStranglerLRUCacheStoreNew) {
  const auto store = EventStranglerLRUCacheStoreNew();
  EventStranglerLRUCacheStoreFree(store);
  ASSERT_NE(store, 0);
}

TEST(EventStranglerLRUCacheStore, EventStranglerLRUCacheStorePut) {
  const auto store = EventStranglerLRUCacheStoreNew();

  auto record = EventStranglerRecord{
      .HashKey = kHashKey,
      .Status = kRecordStatus,
      .CreatedAt = const_cast<char*>("2021-01-01T00:00:00Z"),
      .ExpiresAt = const_cast<char*>("2021-01-01T00:00:00Z"),
  };
  const auto err = EventStranglerStorePut(store, kHashKey, &record, 0);
  EventStranglerLRUCacheStoreFree(store);

  ASSERT_EQ(err, nullptr);
}

TEST(EventStranglerLRUCacheStore, EventStranglerLRUCacheStoreGet) {
  const auto store = EventStranglerLRUCacheStoreNew();

  auto record = EventStranglerRecord{
      .HashKey = kHashKey,
      .Status = kRecordStatus,
      .CreatedAt = const_cast<char*>("2021-01-01T00:00:00Z"),
      .ExpiresAt = const_cast<char*>("2021-01-01T00:00:00Z"),
  };
  const auto err = EventStranglerStorePut(store, kHashKey, &record, 0);
  if (err) {
    EventStranglerLRUCacheStoreFree(store);
    ASSERT_EQ(err, nullptr);
  }

  const auto result = EventStranglerStoreGet(store, kHashKey);
  if (result.r1) {
    EventStranglerLRUCacheStoreFree(store);
    ASSERT_EQ(result.r1, nullptr);
  }
  EventStranglerLRUCacheStoreFree(store);
  ASSERT_NE(result.r0, nullptr);
  ASSERT_STREQ(result.r0->HashKey, kHashKey);
  ASSERT_STREQ(result.r0->Status, kRecordStatus);
}

TEST(EventStranglerLRUCacheStore, EventStranglerLRUCacheStoreDelete) {
  const auto store = EventStranglerLRUCacheStoreNew();

  auto record = EventStranglerRecord{
      .HashKey = kHashKey,
      .Status = kRecordStatus,
      .CreatedAt = const_cast<char*>("2021-01-01T00:00:00Z"),
      .ExpiresAt = const_cast<char*>("2021-01-01T00:00:00Z"),
  };
  auto err = EventStranglerStorePut(store, kHashKey, &record, 0);
  if (err) {
    EventStranglerLRUCacheStoreFree(store);
    ASSERT_EQ(err, nullptr);
  }

  err = EventStranglerStoreDelete(store, kHashKey);
  if (err) {
    EventStranglerLRUCacheStoreFree(store);
    ASSERT_EQ(err, nullptr);
  }
  EventStranglerLRUCacheStoreFree(store);
}

TEST(EventStranglerLRUCacheStore, EventStranglerLRUCacheStoreClose) {
  const auto store = EventStranglerLRUCacheStoreNew();

  const auto err = EventStranglerStoreClose(store);
  if (err) {
    EventStranglerLRUCacheStoreFree(store);
    ASSERT_EQ(err, nullptr);
  }
  EventStranglerLRUCacheStoreFree(store);
}
