package bitwarden

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"slices"
)

var ErrNoPasswordCharacterSetEnabled = errors.New("no password character set enabled")
var ErrInvalidPasswordLength = errors.New("invalid password length")

type PasswordGeneratorRequest struct {
	Lowercase *bool
	Uppercase *bool
	Numbers   *bool
	Special   *bool

	Length *int

	AvoidAmbiguous *bool
	MinLowercase   *int
	MinUppercase   *int
	MinNumber      *int
	MinSpecial     *int
}

func (receiver *PasswordGeneratorRequest) ProvideDefaults() {
	if receiver.Lowercase == nil {
		receiver.Lowercase = new(true)
	}
	if receiver.Uppercase == nil {
		receiver.Uppercase = new(true)
	}
	if receiver.Numbers == nil {
		receiver.Numbers = new(true)
	}
	if receiver.Special == nil {
		receiver.Special = new(false)
	}
	if receiver.Length == nil {
		receiver.Length = new(16)
	}
	if receiver.AvoidAmbiguous == nil {
		receiver.AvoidAmbiguous = new(false)
	}
}

func (*client) GeneratePassword(request *PasswordGeneratorRequest) (string, error) {
	if request == nil {
		request = &PasswordGeneratorRequest{}
	}
	request.ProvideDefaults()

	options, err := request.validate()
	if err != nil {
		return "", err
	}

	return generatePassword(options, cryptoRandomizer{})
}

type passwordGeneratorOptions struct {
	lower   characterSet
	upper   characterSet
	number  characterSet
	special characterSet
	all     characterSet

	minLower   int
	minUpper   int
	minNumber  int
	minSpecial int
	remaining  int
	length     int
}

type characterSet []rune

type passwordRandomizer interface {
	Intn(maxExclusive int) (int, error)
}

type cryptoRandomizer struct{}

func (receiver *PasswordGeneratorRequest) validate() (*passwordGeneratorOptions, error) {
	receiver.ProvideDefaults()

	if !*receiver.Lowercase && !*receiver.Uppercase && !*receiver.Numbers && !*receiver.Special {
		return nil, ErrNoPasswordCharacterSetEnabled
	}

	if *receiver.Length < 4 {
		return nil, fmt.Errorf("%w: the password cannot have fewer than 4 characters", ErrInvalidPasswordLength)
	}

	minLowercase := passwordMinimum(receiver.MinLowercase, *receiver.Lowercase)
	minUppercase := passwordMinimum(receiver.MinUppercase, *receiver.Uppercase)
	minNumber := passwordMinimum(receiver.MinNumber, *receiver.Numbers)
	minSpecial := passwordMinimum(receiver.MinSpecial, *receiver.Special)

	minimumLength := minLowercase + minUppercase + minNumber + minSpecial
	if minimumLength > *receiver.Length {
		return nil, fmt.Errorf("%w: due to settings the password must be at least %d characters long but you requested length of %d", ErrInvalidPasswordLength, minimumLength, *receiver.Length)
	}

	lower := newCharacterSet(*receiver.Lowercase, []rune("abcdefghijklmnopqrstuvwxyz"), *receiver.AvoidAmbiguous, []rune{'l'})
	upper := newCharacterSet(*receiver.Uppercase, []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"), *receiver.AvoidAmbiguous, []rune{'I', 'O'})
	number := newCharacterSet(*receiver.Numbers, []rune("0123456789"), *receiver.AvoidAmbiguous, []rune{'0', '1'})
	special := newCharacterSet(*receiver.Special, []rune{'!', '@', '#', '$', '%', '^', '&', '*'}, false, nil)
	all := unionCharacterSets(lower, upper, number, special)

	return &passwordGeneratorOptions{
		lower:      lower,
		upper:      upper,
		number:     number,
		special:    special,
		all:        all,
		minLower:   minLowercase,
		minUpper:   minUppercase,
		minNumber:  minNumber,
		minSpecial: minSpecial,
		remaining:  *receiver.Length - minimumLength,
		length:     *receiver.Length,
	}, nil
}

func generatePassword(options *passwordGeneratorOptions, randomizer passwordRandomizer) (string, error) {
	buf := make([]rune, 0, options.length)

	var err error
	buf, err = appendRandomRunes(buf, options.all, options.remaining, randomizer)
	if err != nil {
		return "", err
	}
	buf, err = appendRandomRunes(buf, options.upper, options.minUpper, randomizer)
	if err != nil {
		return "", err
	}
	buf, err = appendRandomRunes(buf, options.lower, options.minLower, randomizer)
	if err != nil {
		return "", err
	}
	buf, err = appendRandomRunes(buf, options.number, options.minNumber, randomizer)
	if err != nil {
		return "", err
	}
	buf, err = appendRandomRunes(buf, options.special, options.minSpecial, randomizer)
	if err != nil {
		return "", err
	}

	if err = shuffleRunes(buf, randomizer); err != nil {
		return "", err
	}

	return string(buf), nil
}

func passwordMinimum(min *int, enabled bool) int {
	if !enabled {
		return 0
	}
	if min == nil || *min < 1 {
		return 1
	}

	return *min
}

func newCharacterSet(enabled bool, chars []rune, exclude bool, excluded []rune) characterSet {
	if !enabled {
		return nil
	}

	excludedSet := make(map[rune]struct{}, len(excluded))
	if exclude {
		for _, char := range excluded {
			excludedSet[char] = struct{}{}
		}
	}

	result := make([]rune, 0, len(chars))
	for _, char := range chars {
		if _, ok := excludedSet[char]; ok {
			continue
		}
		result = append(result, char)
	}

	return dedupeAndSortRunes(result)
}

func unionCharacterSets(sets ...characterSet) characterSet {
	var result []rune
	for _, set := range sets {
		result = append(result, set...)
	}

	return dedupeAndSortRunes(result)
}

func dedupeAndSortRunes(chars []rune) characterSet {
	seen := make(map[rune]struct{}, len(chars))
	result := make([]rune, 0, len(chars))
	for _, char := range chars {
		if _, ok := seen[char]; ok {
			continue
		}
		seen[char] = struct{}{}
		result = append(result, char)
	}

	slices.Sort(result)

	return result
}

func appendRandomRunes(buf []rune, set characterSet, count int, randomizer passwordRandomizer) ([]rune, error) {
	for range count {
		char, err := randomRune(set, randomizer)
		if err != nil {
			return nil, err
		}
		buf = append(buf, char)
	}

	return buf, nil
}

func randomRune(set characterSet, randomizer passwordRandomizer) (rune, error) {
	if len(set) == 0 {
		return 0, errors.New("cannot sample from an empty character set")
	}

	index, err := randomizer.Intn(len(set))
	if err != nil {
		return 0, err
	}

	return set[index], nil
}

func shuffleRunes(value []rune, randomizer passwordRandomizer) error {
	for i := len(value) - 1; i > 0; i-- {
		j, err := randomizer.Intn(i + 1)
		if err != nil {
			return err
		}
		value[i], value[j] = value[j], value[i]
	}

	return nil
}

func (cryptoRandomizer) Intn(maxExclusive int) (int, error) {
	if maxExclusive <= 0 {
		return 0, fmt.Errorf("invalid random range: %d", maxExclusive)
	}

	value, err := rand.Int(rand.Reader, big.NewInt(int64(maxExclusive)))
	if err != nil {
		return 0, fmt.Errorf("failed generating random number: %w", err)
	}

	return int(value.Int64()), nil
}
