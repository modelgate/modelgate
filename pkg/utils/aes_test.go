package utils

import (
	"testing"

	"github.com/samber/lo"
)

func TestEncryptAESGCM(t *testing.T) {
	txt := "hello world"
	key := []byte("0123456789abcdef")
	encrypted, err := EncryptAESGCM([]byte(txt), key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(encrypted)
}

func TestDecryptAESGCM(t *testing.T) {
	txt := "hello world"
	key := []byte("0123456789abcdef")
	encrypted, err := EncryptAESGCM([]byte(txt), key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(encrypted)

	decrypted, err := DecryptAESGCM(encrypted, key)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(decrypted))
}

func TestSha256Hex(t *testing.T) {
	txt := lo.RandomString(48, append(lo.LowerCaseLettersCharset, lo.NumbersCharset...))
	txt = "sk-mg-api01-" + txt

	t.Log(txt)
	t.Log(txt[0:16])
	t.Log(txt[len(txt)-4:])
	t.Log(Sha256Hex(txt))
}
