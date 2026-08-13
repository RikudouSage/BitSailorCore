package result

import "testing"

func TestFieldCheckboxValue(t *testing.T) {
	t.Parallel()

	trueValue := "true"
	falseValue := "false"
	tests := []struct {
		name  string
		field Field
		want  bool
	}{
		{name: "checkbox true", field: Field{Type: FieldTypeCheckbox, Value: &trueValue}, want: true},
		{name: "checkbox false", field: Field{Type: FieldTypeCheckbox, Value: &falseValue}, want: false},
		{name: "checkbox nil", field: Field{Type: FieldTypeCheckbox}, want: false},
		{name: "text true", field: Field{Type: FieldTypeText, Value: &trueValue}, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.field.CheckboxValue(); got != test.want {
				t.Fatalf("CheckboxValue() = %v, want %v", got, test.want)
			}
		})
	}
}
