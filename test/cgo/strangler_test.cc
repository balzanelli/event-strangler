#include <gtest/gtest.h>
#include "event_strangler.h"
#include "include.h"

TEST(strangler, strangler_new) {
  const auto store = strangler_lru_cache_store_new();

  auto hash_key_options = strangler_hash_key_options{
      .name = kHashKeyName,
      .expression = kHashKeyExpression,
  };
  auto config = strangler_config{
      .hash_key = &hash_key_options,
      .store = store,
  };
  const auto result = strangler_new(&config);
  if (result.r1) {
    strangler_lru_cache_store_free(store);
    ASSERT_EQ(result.r1, nullptr);
  }
  strangler_lru_cache_store_free(store);
  ASSERT_NE(result.r0, 0);
}

TEST(strangler, strangler_complete) {
  const auto store = strangler_lru_cache_store_new();

  auto hash_key_options = strangler_hash_key_options{
      .name = kHashKeyName,
      .expression = kHashKeyExpression,
  };
  auto config = strangler_config{
      .hash_key = &hash_key_options,
      .store = store,
  };
  const auto result = strangler_new(&config);
  if (result.r1) {
    strangler_lru_cache_store_free(store);
    ASSERT_EQ(result.r1, nullptr);
  }

  auto record = strangler_record{
      .hash_key = kHashKey,
      .status = kRecordStatus,
      .created_at = const_cast<char*>("2021-01-01T00:00:00Z"),
      .expires_at = const_cast<char*>("2021-01-01T00:00:00Z"),
  };
  auto err = strangler_store_put(store, kHashKey, &record, 0);
  if (err) {
    strangler_free(result.r0);
    strangler_lru_cache_store_free(store);
    ASSERT_STREQ(err, "");
  }
  err = strangler_complete(result.r0, kHashKey);
  if (result.r1) {
    strangler_free(result.r0);
    strangler_lru_cache_store_free(store);
    ASSERT_EQ(err, nullptr);
  }
  const auto r = strangler_store_get(store, kHashKey);
  if (r.r1) {
    strangler_free(result.r0);
    strangler_lru_cache_store_free(store);
    ASSERT_EQ(err, nullptr);
  }
  ASSERT_STREQ(r.r0->hash_key, kHashKey);
  ASSERT_STREQ(r.r0->status, "COMPLETE");
}
