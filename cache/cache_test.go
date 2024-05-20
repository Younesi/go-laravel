package cache

import (
	"testing"
)

func TestRedisCache_Set(t *testing.T) {
	key := "foo"
	err := testRedisCache.Set(key, "bar")
	if err != nil {
		t.Errorf("error setting redis value: %s", err)
	}

	has, err := testRedisCache.Has(key)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Errorf("could not find the key in cache, key: %s", key)
	}
}

func TestRedisCache_Forget(t *testing.T) {
	key := "foo"
	err := testRedisCache.Set(key, "bar")
	if err != nil {
		t.Errorf("error setting redis value: %s", err)
	}

	err = testRedisCache.Forget(key)
	if err != nil {
		t.Error(err)
	}

	has, err := testRedisCache.Has(key)
	if err != nil {
		t.Error(err)
	}
	if has {
		t.Errorf("did not expect key in cache, key: %s", key)
	}
}

func TestRedisCache_Get(t *testing.T) {
	key := "foo"
	value := "bar"
	err := testRedisCache.Set(key, value)
	if err != nil {
		t.Errorf("error setting redis value: %s", err)
	}

	v, err := testRedisCache.Get(key)
	if err != nil {
		t.Error(err)
	}
	if v != value {
		t.Errorf("did not get correct value from cache, value: %s", value)
	}
}

func TestRedisCache_Has(t *testing.T) {
	key := "foo"
	err := testRedisCache.Forget(key)
	if err != nil {
		t.Error(err)
	}

	has, err := testRedisCache.Has(key)
	if err != nil {
		t.Error(err)
	}
	if has {
		t.Errorf("unexpected key found in cache, key: %s", key)
	}

	testRedisCache.Set(key, "foo")

	has, err = testRedisCache.Has(key)
	if err != nil {
		t.Error(err)
	}
	if !has {
		t.Errorf("could not find the key in cache, key: %s", key)
	}
}

//todo: add tests for encode and decode

func TestRedisCache_EncodeDecode(t *testing.T) {
	entry := entry{}
	entry["foo"] = "bar"
	bytes, err := encode(entry)
	if err != nil {
		t.Error(err)
	}

	decoded, err := decode(string(bytes))
	if err != nil {
		t.Error(err)
	}
	if decoded["foo"] != entry["foo"] {
		t.Error("expected to get entry when decoded")
	}
}
