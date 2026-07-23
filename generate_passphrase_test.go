package bitwarden

import (
	"reflect"
	"testing"
)

func TestPassphraseGenWords(t *testing.T) {
	options := mustValidatePassphraseRequest(t, &PassphraseGeneratorRequest{})

	words := assertGeneratedPassphraseWords(t, options, 4,
		wordIndexes("crust", "substance", "undertook", "protector"),
	)
	assertStringsEqual(t, words, []string{"crust", "substance", "undertook", "protector"})

	words = assertGeneratedPassphraseWords(t, options, 1, wordIndexes("sighing"))
	assertStringsEqual(t, words, []string{"sighing"})

	words = assertGeneratedPassphraseWords(t, options, 2, wordIndexes("dinghy", "numbing"))
	assertStringsEqual(t, words, []string{"dinghy", "numbing"})
}

func TestPassphraseCapitalize(t *testing.T) {
	tests := map[string]string{
		"hello": "Hello",
		"1ello": "1ello",
		"Hello": "Hello",
		"h":     "H",
		"":      "",
		// Also supports non-ascii, though the EFF list doesn't have any.
		"áéíóú": "Áéíóú",
	}

	for value, expected := range tests {
		if actual := capitalizeFirstLetter(value); actual != expected {
			t.Fatalf("capitalizeFirstLetter(%q) = %q, want %q", value, actual, expected)
		}
	}
}

func TestPassphraseCapitalizeWords(t *testing.T) {
	words := []string{"hello", "world"}
	capitalizePassphraseWords(words)
	assertStringsEqual(t, words, []string{"Hello", "World"})
}

func TestPassphraseIncludeNumber(t *testing.T) {
	words := []string{"hello", "world"}
	err := includeNumberInPassphraseWords(words, &passwordTestRandomizer{
		t:      t,
		values: []int{0, 8},
	})
	if err != nil {
		t.Fatalf("includeNumberInPassphraseWords() returned error: %v", err)
	}
	assertStringsEqual(t, words, []string{"hello8", "world"})

	words = []string{"This", "is", "a", "test"}
	err = includeNumberInPassphraseWords(words, &passwordTestRandomizer{
		t:      t,
		values: []int{3, 6},
	})
	if err != nil {
		t.Fatalf("includeNumberInPassphraseWords() returned error: %v", err)
	}
	assertStringsEqual(t, words, []string{"This", "is", "a", "test6"})
}

func TestPassphraseSeparator(t *testing.T) {
	options := mustValidatePassphraseRequest(t, &PassphraseGeneratorRequest{
		NumWords: new(4),
		// This emoji is 35 bytes long, but represented as a single character.
		WordSeparator: new("👨🏻‍❤️‍💋‍👨🏻"),
		Capitalize:    new(false),
		IncludeNumber: new(true),
	})

	assertGeneratedPassphrase(t, options,
		append(wordIndexes("crust", "substance", "undertook", "protector"), 3, 2),
		"crust👨🏻‍❤️‍💋‍👨🏻substance👨🏻‍❤️‍💋‍👨🏻undertook👨🏻‍❤️‍💋‍👨🏻protector2",
	)
}

func TestPassphrase(t *testing.T) {
	options := mustValidatePassphraseRequest(t, &PassphraseGeneratorRequest{
		NumWords:      new(4),
		WordSeparator: new("-"),
		Capitalize:    new(true),
		IncludeNumber: new(true),
	})
	assertGeneratedPassphrase(t, options,
		append(wordIndexes("crust", "substance", "undertook", "protector"), 3, 2),
		"Crust-Substance-Undertook-Protector2",
	)

	options = mustValidatePassphraseRequest(t, &PassphraseGeneratorRequest{
		NumWords:      new(3),
		WordSeparator: new(" "),
		Capitalize:    new(false),
		IncludeNumber: new(true),
	})
	assertGeneratedPassphrase(t, options,
		append(wordIndexes("numbing", "catnap", "jokester"), 0, 4),
		"numbing4 catnap jokester",
	)

	options = mustValidatePassphraseRequest(t, &PassphraseGeneratorRequest{
		NumWords:      new(5),
		WordSeparator: new(";"),
		Capitalize:    new(false),
		IncludeNumber: new(false),
	})
	assertGeneratedPassphrase(t, options,
		wordIndexes("cabana", "pungent", "acts", "sappy", "duller"),
		"cabana;pungent;acts;sappy;duller",
	)
}

func mustValidatePassphraseRequest(t *testing.T, request *PassphraseGeneratorRequest) *passphraseGeneratorOptions {
	t.Helper()

	request.ProvideDefaults()
	options, err := request.validate()
	if err != nil {
		t.Fatalf("validate() returned error: %v", err)
	}

	return options
}

func assertGeneratedPassphraseWords(t *testing.T, options *passphraseGeneratorOptions, numWords int, randomValues []int) []string {
	t.Helper()

	options = &passphraseGeneratorOptions{
		numWords:      numWords,
		wordSeparator: options.wordSeparator,
		capitalize:    options.capitalize,
		includeNumber: options.includeNumber,
		words:         options.words,
	}

	words, err := generatePassphraseWords(options, &passwordTestRandomizer{
		t:      t,
		values: append([]int(nil), randomValues...),
	})
	if err != nil {
		t.Fatalf("generatePassphraseWords() returned error: %v", err)
	}

	return words
}

func assertGeneratedPassphrase(t *testing.T, options *passphraseGeneratorOptions, randomValues []int, expected string) {
	t.Helper()

	randomizer := &passwordTestRandomizer{
		t:      t,
		values: append([]int(nil), randomValues...),
	}
	actual, err := generatePassphrase(options, randomizer)
	if err != nil {
		t.Fatalf("generatePassphrase() returned error: %v", err)
	}
	if actual != expected {
		t.Fatalf("generatePassphrase() = %q, want %q", actual, expected)
	}
}

func wordIndexes(words ...string) []int {
	wordList := passphraseWordList()
	indexes := make([]int, 0, len(words))
	for _, word := range words {
		for index, candidate := range wordList {
			if candidate == word {
				indexes = append(indexes, index)
				break
			}
		}
	}

	return indexes
}

func assertStringsEqual(t *testing.T, actual []string, expected []string) {
	t.Helper()

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("strings = %v, want %v", actual, expected)
	}
}
