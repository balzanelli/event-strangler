#include <gtest/gtest.h>
#include "event_strangler.h"
#include "include.h"

TEST(strangler_lru_cache_store, strangler_lru_cache_store_new) {
  const auto store = strangler_lru_cache_store_new();
  strangler_lru_cache_store_free(store);
  ASSERT_NE(store, 0);
}

TEST(strangler_lru_cache_store, strangler_lru_cache_store_put) {
  const auto store = strangler_lru_cache_store_new();

  auto record = strangler_record{
      .hash_key = kHashKey,
      .status = kRecordStatus,
      .created_at = const_cast<char*>("2021-01-01T00:00:00Z"),
      .expires_at = const_cast<char*>("2021-01-01T00:00:00Z"),
  };
  const auto err = strangler_store_put(store, kHashKey, &record, 0);
  strangler_lru_cache_store_free(store);

  ASSERT_EQ(err, nullptr);
}

TEST(strangler_lru_cache_store, strangler_lru_cache_store_get) {
  const auto store = strangler_lru_cache_store_new();

  auto record = strangler_record{
      .hash_key = kHashKey,
      .status = kRecordStatus,
      .created_at = const_cast<char*>("2021-01-01T00:00:00Z"),
      .expires_at = const_cast<char*>("2021-01-01T00:00:00Z"),
  };
  const auto err = strangler_store_put(store, kHashKey, &record, 0);
  if (err) {
    strangler_lru_cache_store_free(store);
    ASSERT_EQ(err, nullptr);
  }

  const auto result = strangler_store_get(store, kHashKey);
  if (result.r1) {
    strangler_lru_cache_store_free(store);
    ASSERT_EQ(result.r1, nullptr);
  }
  strangler_lru_cache_store_free(store);
  ASSERT_NE(result.r0, nullptr);
  ASSERT_STREQ(result.r0->hash_key, kHashKey);
  ASSERT_STREQ(result.r0->status, kRecordStatus);
}

TEST(strangler_lru_cache_store, strangler_lru_cache_store_delete) {
  const auto store = strangler_lru_cache_store_new();

  auto record = strangler_record{
      .hash_key = kHashKey,
      .status = kRecordStatus,
      .created_at = const_cast<char*>("2021-01-01T00:00:00Z"),
      .expires_at = const_cast<char*>("2021-01-01T00:00:00Z"),
  };
  auto err = strangler_store_put(store, kHashKey, &record, 0);
  if (err) {
    strangler_lru_cache_store_free(store);
    ASSERT_EQ(err, nullptr);
  }

  err = strangler_store_delete(store, kHashKey);
  if (err) {
    strangler_lru_cache_store_free(store);
    ASSERT_EQ(err, nullptr);
  }
  strangler_lru_cache_store_free(store);
}

TEST(strangler_lru_cache_store, strangler_lru_cache_store_close) {
  const auto store = strangler_lru_cache_store_new();

  const auto err = strangler_store_close(store);
  if (err) {
    strangler_lru_cache_store_free(store);
    ASSERT_EQ(err, nullptr);
  }
  strangler_lru_cache_store_free(store);
}
