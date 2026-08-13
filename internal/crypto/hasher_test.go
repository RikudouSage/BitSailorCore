package crypto

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func TestDeriveMasterKeyArgon2ID(t *testing.T) {
	memory := 32
	parallelism := 2

	key, err := DeriveMasterKey("test_key", "67t9b5g67$%Dh89n", KDFTypeArgon2ID, &KDFConfig{
		Iterations:  4,
		Memory:      &memory,
		Parallelism: &parallelism,
	})
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte{
		207, 240, 225, 177, 162, 19, 163, 76, 98, 106, 179, 175, 224, 9, 17, 240,
		20, 147, 237, 47, 246, 150, 141, 184, 62, 225, 131, 242, 51, 53, 225, 242,
	}
	if !bytes.Equal(key, expected) {
		t.Fatalf("expected %s, got %s", base64.StdEncoding.EncodeToString(expected), base64.StdEncoding.EncodeToString(key))
	}
}

func TestDeriveMasterKeyHashArgon2ID(t *testing.T) {
	memory := 32
	parallelism := 2
	password := "asdfasdf"

	key, err := DeriveMasterKey("test_salt", password, KDFTypeArgon2ID, &KDFConfig{
		Iterations:  4,
		Memory:      &memory,
		Parallelism: &parallelism,
	})
	if err != nil {
		t.Fatal(err)
	}

	if hash := DeriveMasterKeyHash(key, password); hash != "PR6UjYmjmppTYcdyTiNbAhPJuQQOmynKbdEl1oyi/iQ=" {
		t.Fatalf("unexpected hash: %s", hash)
	}
}
