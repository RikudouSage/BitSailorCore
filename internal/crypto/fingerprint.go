package crypto

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"strings"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/helper"
)

const fingerprintMinimumEntropy = 64

func FingerprintPhrase(email string, publicKey []byte) ([]string, error) {
	if len(publicKey) == 0 {
		return nil, fmt.Errorf("public key is required to generate a fingerprint")
	}

	keyFingerprint := sha256.Sum256(publicKey)
	userFingerprint, err := hkdfSha256(keyFingerprint[:], strings.ToLower(email), 32)
	if err != nil {
		return nil, fmt.Errorf("failed deriving fingerprint: %w", err)
	}

	return fingerprintHashPhrase(userFingerprint, fingerprintMinimumEntropy)
}

func fingerprintHashPhrase(hash []byte, minimumEntropy int) ([]string, error) {
	entropyPerWord := math.Log2(float64(len(helper.EffLongWordList)))
	numWords := int(math.Ceil(float64(minimumEntropy) / entropyPerWord))

	entropyAvailable := len(hash) * 4
	if float64(numWords)*entropyPerWord > float64(entropyAvailable) {
		return nil, fmt.Errorf("output entropy of hash function is too small")
	}

	phrase := make([]string, 0, numWords)
	hashNumber := new(big.Int).SetBytes(hash)
	wordListLength := big.NewInt(int64(len(helper.EffLongWordList)))
	remainder := new(big.Int)

	for len(phrase) < numWords {
		hashNumber.QuoRem(hashNumber, wordListLength, remainder)
		phrase = append(phrase, helper.EffLongWordList[remainder.Int64()])
	}

	return phrase, nil
}
