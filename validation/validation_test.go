package validation

import (
	"regexp"
	"testing"
)

// Test the value of the validation_expression field for the
// secondaryTid data element. NB this does not access the database,
// and so it is not testing the value stored in the script directly;
// the tested regex is simply copied from the database scripts,
// and so this serves as a convenient way of testing the regex rather
// than being a true unit test per se.
func TestDualCurrencySecondaryTidValidationExpression(t *testing.T) {
	validationExpression := "^\\d{8}$"
	tests := []struct {
		valueToValidate string
		expectedResult  bool
	}{
		/*
		 * Behaviour needs to match the validation description:
		 * "Must be numeric, not 00000000, and 8 digits long".
		 *
		 * HOWEVER: ensuring against the 00000000 case cannot be done by
		 * a regex alone since golang's regex handling does not support negative lookahead.
		 * Therefore that particular case will handled by the ValidateDataElement function.
		 */
		{valueToValidate: "12345678", expectedResult: true},
		{valueToValidate: "10000000", expectedResult: true},
		{valueToValidate: "02000000", expectedResult: true},
		{valueToValidate: "00300000", expectedResult: true},
		{valueToValidate: "00040000", expectedResult: true},
		{valueToValidate: "00005000", expectedResult: true},
		{valueToValidate: "00000600", expectedResult: true},
		{valueToValidate: "00000070", expectedResult: true},
		{valueToValidate: "00000008", expectedResult: true},

		{valueToValidate: "", expectedResult: false},
		{valueToValidate: "1", expectedResult: false},
		{valueToValidate: "12", expectedResult: false},
		{valueToValidate: "1234567", expectedResult: false},
		{valueToValidate: "123456789", expectedResult: false},
		{valueToValidate: "abcdefgh", expectedResult: false},
		{valueToValidate: "1234567x", expectedResult: false},
		{valueToValidate: "x2345678", expectedResult: false},
	}
	for _, tt := range tests {
		t.Run(tt.valueToValidate, func(t *testing.T) {
			result, err := regexp.MatchString(validationExpression, tt.valueToValidate)
			if err != nil {
				t.Fatal(err.Error())
			}

			if result != tt.expectedResult {
				t.Errorf("%v: expected %v, but got %v", tt.valueToValidate, tt.expectedResult, result)
			}
		})
	}
}

func Test_isValidJson(t *testing.T) {
	tests := []struct {
		name string
		jsonString string
		want bool
	}{
		{
			name: "test1 - empty json string",
			jsonString: "",
			want: false,
		},
		{
			name: "test2 - invalid string",
			jsonString: "test string",
			want: false,
		},
		{
			name: "test3 - invalid string",
			jsonString: "{{test string}",
			want: false,
		},
		{
			name: "test4 - valid json object",
			jsonString: "{ \"merchantNo\": \"111111111111\" }",
			want: true,
		},
		{
			name: "test5 - valid json array",
			jsonString: "[{ \"merchantNo\": \"111111111111\" }]",
			want: true,
		},
		{
			name: "test6 - empty json array",
			jsonString: "[]",
			want: true,
		},
		{
			name: "test7 - empty json object",
			jsonString: "{}",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidJson(tt.jsonString); got != tt.want {
				t.Errorf("isValidJsonObject() = %v, want %v", got, tt.want)
			}
		})
	}
}