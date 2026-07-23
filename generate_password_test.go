package bitwarden

import (
	"errors"
	"reflect"
	"testing"
)

func TestPasswordGenAllCharsetsEnabled(t *testing.T) {
	request := &PasswordGeneratorRequest{
		Lowercase:      new(true),
		Uppercase:      new(true),
		Numbers:        new(true),
		Special:        new(true),
		AvoidAmbiguous: new(false),
	}
	options := mustValidatePasswordRequest(t, request)

	assertRunesEqual(t, options.lower, []rune("abcdefghijklmnopqrstuvwxyz"))
	assertRunesEqual(t, options.upper, []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	assertRunesEqual(t, options.number, []rune("0123456789"))
	assertRunesEqual(t, options.special, []rune("!#$%&*@^"))

	assertGeneratedPassword(t, options, "0oA772tQjaUO$a@L")
}

func TestPasswordGenOnlyLettersEnabled(t *testing.T) {
	request := &PasswordGeneratorRequest{
		Lowercase:      new(true),
		Uppercase:      new(true),
		Numbers:        new(false),
		Special:        new(false),
		AvoidAmbiguous: new(false),
	}
	options := mustValidatePasswordRequest(t, request)

	assertRunesEqual(t, options.lower, []rune("abcdefghijklmnopqrstuvwxyz"))
	assertRunesEqual(t, options.upper, []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	assertRunesEqual(t, options.number, nil)
	assertRunesEqual(t, options.special, nil)

	assertGeneratedPassword(t, options, "FrNSJGvhnAbXggMU")
}

func TestPasswordGenOnlyNumbersAndLowerEnabledNoAmbiguous(t *testing.T) {
	request := &PasswordGeneratorRequest{
		Lowercase:      new(true),
		Uppercase:      new(false),
		Numbers:        new(true),
		Special:        new(false),
		AvoidAmbiguous: new(true),
	}
	options := mustValidatePasswordRequest(t, request)

	assertRunesEqual(t, options.lower, []rune("abcdefghijkmnopqrstuvwxyz"))
	assertRunesEqual(t, options.upper, nil)
	assertRunesEqual(t, options.number, []rune("23456789"))
	assertRunesEqual(t, options.special, nil)

	assertGeneratedPassword(t, options, "5uat85wos2jg4n9f")
}

func TestPasswordGenOnlyUpperAndSpecialEnabledNoAmbiguous(t *testing.T) {
	request := &PasswordGeneratorRequest{
		Lowercase:      new(false),
		Uppercase:      new(true),
		Numbers:        new(false),
		Special:        new(true),
		AvoidAmbiguous: new(true),
	}
	options := mustValidatePasswordRequest(t, request)

	assertRunesEqual(t, options.lower, nil)
	assertRunesEqual(t, options.upper, []rune("ABCDEFGHJKLMNPQRSTUVWXYZ"))
	assertRunesEqual(t, options.number, nil)
	assertRunesEqual(t, options.special, []rune("!#$%&*@^"))

	assertGeneratedPassword(t, options, "%VBT*%YPT!LH$PAF")
}

func TestPasswordGenMinimumLimits(t *testing.T) {
	request := &PasswordGeneratorRequest{
		Lowercase:      new(true),
		Uppercase:      new(true),
		Numbers:        new(true),
		Special:        new(true),
		AvoidAmbiguous: new(false),
		Length:         new(24),
		MinLowercase:   new(5),
		MinUppercase:   new(5),
		MinNumber:      new(5),
		MinSpecial:     new(5),
	}
	options := mustValidatePasswordRequest(t, request)

	assertRunesEqual(t, options.lower, []rune("abcdefghijklmnopqrstuvwxyz"))
	assertRunesEqual(t, options.upper, []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	assertRunesEqual(t, options.number, []rune("0123456789"))
	assertRunesEqual(t, options.special, []rune("!#$%&*@^"))

	if options.minLower != 5 || options.minUpper != 5 || options.minNumber != 5 || options.minSpecial != 5 {
		t.Fatalf("minimums = lower:%d upper:%d number:%d special:%d, want all 5",
			options.minLower, options.minUpper, options.minNumber, options.minSpecial)
	}

	assertGeneratedPassword(t, options, "t&c0L73*D*G%aak7goq!N2T4")
}

func TestGeneratePasswordRejectsNoCharacterSet(t *testing.T) {
	_, err := (&client{}).GeneratePassword(&PasswordGeneratorRequest{
		Lowercase: new(false),
		Uppercase: new(false),
		Numbers:   new(false),
		Special:   new(false),
	})
	if !errors.Is(err, ErrNoPasswordCharacterSetEnabled) {
		t.Fatalf("GeneratePassword() error = %v, want %v", err, ErrNoPasswordCharacterSetEnabled)
	}
}

func TestGeneratePasswordRejectsInvalidLength(t *testing.T) {
	_, err := (&client{}).GeneratePassword(&PasswordGeneratorRequest{Length: new(3)})
	if !errors.Is(err, ErrInvalidPasswordLength) {
		t.Fatalf("GeneratePassword() error = %v, want %v", err, ErrInvalidPasswordLength)
	}
}

func TestGeneratePasswordRejectsMinimumsLongerThanPassword(t *testing.T) {
	_, err := (&client{}).GeneratePassword(&PasswordGeneratorRequest{
		Length:       new(4),
		MinLowercase: new(2),
		MinUppercase: new(2),
		MinNumber:    new(1),
	})
	if !errors.Is(err, ErrInvalidPasswordLength) {
		t.Fatalf("GeneratePassword() error = %v, want %v", err, ErrInvalidPasswordLength)
	}
}

func mustValidatePasswordRequest(t *testing.T, request *PasswordGeneratorRequest) *passwordGeneratorOptions {
	t.Helper()

	request.ProvideDefaults()
	options, err := request.validate()
	if err != nil {
		t.Fatalf("validate() returned error: %v", err)
	}

	return options
}

func assertGeneratedPassword(t *testing.T, options *passwordGeneratorOptions, expected string) {
	t.Helper()

	randomizer := newPasswordTestRandomizer(t, options, expected)
	actual, err := generatePassword(options, randomizer)
	if err != nil {
		t.Fatalf("generatePassword() returned error: %v", err)
	}
	if actual != expected {
		t.Fatalf("generatePassword() = %q, want %q", actual, expected)
	}
}

func assertRunesEqual(t *testing.T, actual characterSet, expected []rune) {
	t.Helper()

	if !reflect.DeepEqual([]rune(actual), expected) {
		t.Fatalf("character set = %q, want %q", string(actual), string(expected))
	}
}

type passwordTestRandomizer struct {
	t      *testing.T
	values []int
}

func newPasswordTestRandomizer(t *testing.T, options *passwordGeneratorOptions, expected string) *passwordTestRandomizer {
	t.Helper()

	initial, sampleValues := scriptedInitialPassword(t, options, []rune(expected))
	shuffleValues := scriptedShuffleValues(t, initial, []rune(expected))

	return &passwordTestRandomizer{
		t:      t,
		values: append(sampleValues, shuffleValues...),
	}
}

func (receiver *passwordTestRandomizer) Intn(maxExclusive int) (int, error) {
	receiver.t.Helper()

	if len(receiver.values) == 0 {
		receiver.t.Fatalf("unexpected random request with max %d", maxExclusive)
	}

	value := receiver.values[0]
	receiver.values = receiver.values[1:]
	if value < 0 || value >= maxExclusive {
		receiver.t.Fatalf("scripted random value %d outside range [0, %d)", value, maxExclusive)
	}

	return value, nil
}

func scriptedInitialPassword(t *testing.T, options *passwordGeneratorOptions, expected []rune) ([]rune, []int) {
	t.Helper()

	remaining := append([]rune(nil), expected...)
	initial := make([]rune, 0, options.length)

	special := takeMatchingRunes(t, &remaining, options.special, options.minSpecial)
	number := takeMatchingRunes(t, &remaining, options.number, options.minNumber)
	lower := takeMatchingRunes(t, &remaining, options.lower, options.minLower)
	upper := takeMatchingRunes(t, &remaining, options.upper, options.minUpper)

	all := make([]rune, options.remaining)
	copy(all, remaining)
	initial = append(initial, all...)
	initial = append(initial, upper...)
	initial = append(initial, lower...)
	initial = append(initial, number...)
	initial = append(initial, special...)

	values := make([]int, 0, options.length)
	values = appendSampleIndexes(t, values, options.all, all)
	values = appendSampleIndexes(t, values, options.upper, upper)
	values = appendSampleIndexes(t, values, options.lower, lower)
	values = appendSampleIndexes(t, values, options.number, number)
	values = appendSampleIndexes(t, values, options.special, special)

	return initial, values
}

func takeMatchingRunes(t *testing.T, remaining *[]rune, set characterSet, count int) []rune {
	t.Helper()

	result := make([]rune, 0, count)
	for range count {
		found := false
		for i := len(*remaining) - 1; i >= 0; i-- {
			if !containsRune(set, (*remaining)[i]) {
				continue
			}

			result = append(result, (*remaining)[i])
			*remaining = append((*remaining)[:i], (*remaining)[i+1:]...)
			found = true
			break
		}
		if !found {
			t.Fatalf("could not find %d runes in set %q from %q", count, string(set), string(*remaining))
		}
	}

	return result
}

func appendSampleIndexes(t *testing.T, values []int, set characterSet, samples []rune) []int {
	t.Helper()

	for _, sample := range samples {
		index := indexRune(set, sample)
		if index < 0 {
			t.Fatalf("sample %q not found in set %q", sample, string(set))
		}
		values = append(values, index)
	}

	return values
}

func scriptedShuffleValues(t *testing.T, initial []rune, expected []rune) []int {
	t.Helper()

	current := append([]rune(nil), initial...)
	values := make([]int, 0, len(current)-1)

	for i := len(current) - 1; i > 0; i-- {
		j := -1
		for candidate := 0; candidate <= i; candidate++ {
			if current[candidate] == expected[i] {
				j = candidate
				break
			}
		}
		if j < 0 {
			t.Fatalf("could not script shuffle from %q to %q", string(initial), string(expected))
		}

		values = append(values, j)
		current[i], current[j] = current[j], current[i]
	}

	if !reflect.DeepEqual(current, expected) {
		t.Fatalf("scripted shuffle produced %q, want %q", string(current), string(expected))
	}

	return values
}

func containsRune(set characterSet, target rune) bool {
	return indexRune(set, target) >= 0
}

func indexRune(set characterSet, target rune) int {
	for i, char := range set {
		if char == target {
			return i
		}
	}

	return -1
}
