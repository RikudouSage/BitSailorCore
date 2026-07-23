package bitwarden

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidPassphraseNumWords = errors.New("invalid passphrase word count")

type PassphraseGeneratorRequest struct {
	NumWords      *int
	WordSeparator *string
	Capitalize    *bool
	IncludeNumber *bool
}

func (receiver *PassphraseGeneratorRequest) ProvideDefaults() {
	if receiver.NumWords == nil {
		receiver.NumWords = new(3)
	}
	if receiver.WordSeparator == nil {
		receiver.WordSeparator = new("-")
	}
	if receiver.Capitalize == nil {
		receiver.Capitalize = new(false)
	}
	if receiver.IncludeNumber == nil {
		receiver.IncludeNumber = new(false)
	}
}

func (*client) GeneratePassphrase(request *PassphraseGeneratorRequest) (string, error) {
	if request == nil {
		request = &PassphraseGeneratorRequest{}
	}
	request.ProvideDefaults()

	options, err := request.validate()
	if err != nil {
		return "", err
	}

	return generatePassphrase(options, cryptoRandomizer{})
}

type passphraseGeneratorOptions struct {
	numWords      int
	wordSeparator string
	capitalize    bool
	includeNumber bool
	words         []string
}

func (receiver *PassphraseGeneratorRequest) validate() (*passphraseGeneratorOptions, error) {
	receiver.ProvideDefaults()

	if *receiver.NumWords < 3 || *receiver.NumWords > 20 {
		return nil, fmt.Errorf("%w: 'num_words' must be between 3 and 20", ErrInvalidPassphraseNumWords)
	}

	return &passphraseGeneratorOptions{
		numWords:      *receiver.NumWords,
		wordSeparator: *receiver.WordSeparator,
		capitalize:    *receiver.Capitalize,
		includeNumber: *receiver.IncludeNumber,
		words:         passphraseWordList(),
	}, nil
}

func generatePassphrase(options *passphraseGeneratorOptions, randomizer passwordRandomizer) (string, error) {
	words, err := generatePassphraseWords(options, randomizer)
	if err != nil {
		return "", err
	}

	if options.includeNumber {
		if err = includeNumberInPassphraseWords(words, randomizer); err != nil {
			return "", err
		}
	}
	if options.capitalize {
		capitalizePassphraseWords(words)
	}

	return strings.Join(words, options.wordSeparator), nil
}

func generatePassphraseWords(options *passphraseGeneratorOptions, randomizer passwordRandomizer) ([]string, error) {
	words := make([]string, 0, options.numWords)
	for range options.numWords {
		index, err := randomizer.Intn(len(options.words))
		if err != nil {
			return nil, err
		}

		words = append(words, options.words[index])
	}

	return words, nil
}

func includeNumberInPassphraseWords(words []string, randomizer passwordRandomizer) error {
	wordIndex, err := randomizer.Intn(len(words))
	if err != nil {
		return err
	}
	digit, err := randomizer.Intn(10)
	if err != nil {
		return err
	}

	words[wordIndex] = fmt.Sprintf("%s%d", words[wordIndex], digit)
	return nil
}

func capitalizePassphraseWords(words []string) {
	for i, word := range words {
		words[i] = capitalizeFirstLetter(word)
	}
}

func capitalizeFirstLetter(value string) string {
	for _, char := range value {
		return string(unicode.ToUpper(char)) + value[len(string(char)):]
	}

	return value
}

func passphraseWordList() []string {
	words := make([]string, 0, len(effLongWordList))
	for _, word := range effLongWordList {
		if strings.Contains(word, "-") {
			continue
		}
		words = append(words, word)
	}

	return words
}
