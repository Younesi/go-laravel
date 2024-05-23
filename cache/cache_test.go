package cache

import (
	"testing"
)

func TestCache_EncodeDecode(t *testing.T) {
	entry := Entry{}
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
