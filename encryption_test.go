package atlas

import (
	"crypto/rand"
	"testing"
)

func TestEncryptDecryptSuccess(t *testing.T) {
	// Generate valid 32-byte key (AES-256)
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		t.Fatal(err)
	}

	e := &Encryption{Key: key}

	tests := []struct {
		name      string
		plaintext string
	}{
		{"Simple Text", "Hello Atlas!"},
		{"Empty String", ""},
		{"Special Chars", "!@#$%^&*()"},
		{"Long Text", `Lorem ipsum dolor sit amet, consectetur adipiscing elit...`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := e.Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			decrypted, err := e.Decrypt(encrypted)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if decrypted != tt.plaintext {
				t.Errorf("Expected %q, got %q", tt.plaintext, decrypted)
			}
		})
	}
}
