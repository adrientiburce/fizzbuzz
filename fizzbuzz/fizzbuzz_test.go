package fizzbuzz

import (
	"testing"
)

func TestFizzBuzz_FizzBuzz(t *testing.T) {
	type fields struct {
		Int1  int
		Int2  int
		Limit int
		Str1  string
		Str2  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "fizzBuzz basic",
			fields: fields{
				Int1:  3,
				Int2:  5,
				Limit: 16,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want: "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16",
		},
		{
			name: "fooBar until 15",
			fields: fields{
				Int1:  3,
				Int2:  5,
				Limit: 15,
				Str1:  "foo",
				Str2:  "bar",
			},
			want: "1,2,foo,4,bar,foo,7,8,foo,bar,11,foo,13,14,foobar",
		},
		{
			name: "fizzBuzz with 2 and 4",
			fields: fields{
				Int1:  2,
				Int2:  4,
				Limit: 10,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want: "1,fizz,3,fizz,5,fizz,7,fizzbuzz,9,fizz",
		},
		{
			name: "fizzBuzz with wrong params",
			fields: fields{
				Int1:  2,
				Int2:  4,
				Limit: 0,
				Str1:  "fizz",
				Str2:  "buzz",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &FizzBuzz{
				Int1:  tt.fields.Int1,
				Int2:  tt.fields.Int2,
				Limit: tt.fields.Limit,
				Str1:  tt.fields.Str1,
				Str2:  tt.fields.Str2,
			}
			if got := s.computeSuite(); got != tt.want {
				t.Errorf("FizzBuzz.FizzBuzz() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateStruct(t *testing.T) {
	fizzbuzz := FizzBuzz{
		Int1:  3,
		Int2:  5,
		Limit: 15,
		Str1:  "foo",
		Str2:  "bar",
	}

	err := validateStruct(fizzbuzz)
	if err != nil {
		t.Error("fizzBuzz shoudl be valid")
	}
}

func TestValidateStruct_WithError(t *testing.T) {
	fizzbuzz := FizzBuzz{
		Int1:  3,
		Limit: 10,
		Str1:  "foo",
		Str2:  "bar",
	}

	err := validateStruct(fizzbuzz)
	expected := "Key: 'FizzBuzz.Int2' Error:Field validation for 'Int2' failed on the 'required' tag"
	if got := err.Error(); got != expected {
		t.Errorf("error returned: %v, want %v", got, expected)
	}
}
