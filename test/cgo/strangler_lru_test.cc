#include <gtest/gtest.h>
#include "libstrangler.h"
#include "include.h"

TEST(LRUCacheStore, New) {
  const auto store = LRUCacheStore_New();
  LRUCacheStore_Free(store);
  ASSERT_NE(store, 0);
}

TEST(LRUCacheStore, Put) {
  const auto store = LRUCacheStore_New();

  auto record = Strangler_Record{
      .HashKey   = kHashKey,
      .Status    = kRecordStatus,
      .CreatedAt = "2021-01-01T00:00:00Z",
      .ExpiresAt = "2021-01-01T00:00:00Z",
  };
  const auto err = Store_Put(store, kHashKey, &record, 0);
  LRUCacheStore_Free(store);

  ASSERT_EQ(err, nullptr);
}

TEST(LRUCacheStore, Get) {
  const auto store = LRUCacheStore_New();

  auto record = Strangler_Record{
      .HashKey   = kHashKey,
      .Status    = kRecordStatus,
      .CreatedAt = "2021-01-01T00:00:00Z",
      .ExpiresAt = "2021-01-01T00:00:00Z",
  };
  const auto err = Store_Put(store, kHashKey, &record, 0);
  if (err) {
    LRUCacheStore_Free(store);
    ASSERT_EQ(err, nullptr);
  }

  const auto result = Store_Get(store, kHashKey);
  if (result.r1) {
    LRUCacheStore_Free(store);
    ASSERT_EQ(result.r1, nullptr);
  }
  LRUCacheStore_Free(store);
  ASSERT_NE(result.r0, nullptr);
  ASSERT_STREQ(result.r0->HashKey, kHashKey);
  ASSERT_STREQ(result.r0->Status, kRecordStatus);
}

TEST(LRUCacheStore, Delete) {
  const auto store = LRUCacheStore_New();

  auto record = Strangler_Record{
      .HashKey   = kHashKey,
      .Status    = kRecordStatus,
      .CreatedAt = "2021-01-01T00:00:00Z",
      .ExpiresAt = "2021-01-01T00:00:00Z",
  };
  auto err = Store_Put(store, kHashKey, &record, 0);
  if (err) {
    LRUCacheStore_Free(store);
    ASSERT_EQ(err, nullptr);
  }

  err = Store_Delete(store, kHashKey);
  if (err) {
    LRUCacheStore_Free(store);
    ASSERT_EQ(err, nullptr);
  }
  LRUCacheStore_Free(store);
}

TEST(LRUCacheStore, Close) {
  const auto store = LRUCacheStore_New();

  const auto err = Store_Close(store);
  if (err) {
    LRUCacheStore_Free(store);
    ASSERT_EQ(err, nullptr);
  }
  LRUCacheStore_Free(store);
}
