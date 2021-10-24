#include <gtest/gtest.h>
#include "libstrangler.h"
#include "include.h"

TEST(Strangler, New) {
  const auto store = LRUCacheStore_New();

  auto hash_key_config = Strangler_HashKeyOptions{
      .Name       = kHashKeyName,
      .Expression = kHashKeyExpression,
  };
  auto config = Strangler_Config{
      .HashKey = &hash_key_config,
      .Store   = store,
  };
  const auto result = Strangler_New(&config);
  if (result.r1) {
    LRUCacheStore_Free(store);
    ASSERT_EQ(result.r1, nullptr);
  }
  LRUCacheStore_Free(store);
  ASSERT_NE(result.r0, 0);
}

TEST(Strangler, Complete) {
  const auto store = LRUCacheStore_New();

  auto hash_key_config = Strangler_HashKeyOptions{
      .Name       = kHashKeyName,
      .Expression = kHashKeyExpression,
  };
  auto config = Strangler_Config{
      .HashKey = &hash_key_config,
      .Store   = store,
  };
  const auto result = Strangler_New(&config);
  if (result.r1) {
    LRUCacheStore_Free(store);
    ASSERT_EQ(result.r1, nullptr);
  }

  auto record = Strangler_Record{
      .HashKey   = kHashKey,
      .Status    = kRecordStatus,
      .CreatedAt = "2021-01-01T00:00:00Z",
      .ExpiresAt = "2021-01-01T00:00:00Z",
  };
  auto err = Store_Put(store, kHashKey, &record, 0);
  if (err) {
    Strangler_Free(result.r0);
    LRUCacheStore_Free(store);
    ASSERT_STREQ(err, "");
  }
  err = Strangler_Complete(result.r0, kHashKey);
  if (result.r1) {
    Strangler_Free(result.r0);
    LRUCacheStore_Free(store);
    ASSERT_EQ(err, nullptr);
  }
  const auto r = Store_Get(store, kHashKey);
  if (r.r1) {
    Strangler_Free(result.r0);
    LRUCacheStore_Free(store);
    ASSERT_EQ(err, nullptr);
  }
  ASSERT_STREQ(r.r0->HashKey, kHashKey);
  ASSERT_STREQ(r.r0->Status, "COMPLETE");
}
