package fizzbuzz

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputeSuite(t *testing.T) {
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
			name: "fizzBuzz basic suite",
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
			name: "fizzbuzz with foo & bar until 15",
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
			name: "fizzBuzz with wrong parameters",
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
	require.NoError(t, err)
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
	require.Equal(t, expected, err.Error())
}

var fizzbuzzStr = `{
	"int1": 3,
	"int2": 5,
	"limit": 16,
	"str1": "fizz",
	"str2": "buzz"
}`

var fooBarStr = `{
	"int1": 3,
	"int2": 5,
	"limit": 10,
	"str1": "foo",
	"str2": "bar"
}`

func TestFizzBuzzEndpoint(t *testing.T) {
	service := New()
	req, err := http.NewRequest("GET", "/fizzbuzz", strings.NewReader(fizzbuzzStr))
	require.NoError(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.FizzBuzzEndpoint)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16", rr.Body.String())

	// first request to statistic endpoint
	statReq, err := http.NewRequest("GET", "/statistics", strings.NewReader(fizzbuzzStr))
	require.NoError(t, err)

	rr = httptest.NewRecorder()
	handlerStat := http.HandlerFunc(service.Statistics)
	handlerStat.ServeHTTP(rr, statReq)

	file, err := ioutil.ReadFile("../testdata/stat/stat1.json")
	require.NoError(t, err)
	require.Equal(t, string(file), rr.Body.String())

	// second fizzbuzz call with other parameters
	req, err = http.NewRequest("GET", "/fizzbuzz", strings.NewReader(fooBarStr))
	require.NoError(t, err)
	handler.ServeHTTP(rr, req)

	rr = httptest.NewRecorder()
	handlerStat.ServeHTTP(rr, statReq)
	file, err = ioutil.ReadFile("../testdata/stat/stat2.json")
	require.NoError(t, err)
	var expectedStat, resStat Statistics
	err = json.Unmarshal(file, &expectedStat)
	require.NoError(t, err)

	err = json.NewDecoder(rr.Body).Decode(&resStat)
	require.NoError(t, err)
	require.Equal(t, expectedStat, resStat)
}

func TestFizzBuzzEndpoint_WithError(t *testing.T) {
	tests := []struct {
		name         string
		inputString  string
		expectedCode int
	}{
		{
			name: "with missing limit",
			inputString: `{
				"int1": 3,
				"int2": 5,
				"str1": "fizz",
				"str2": "buzz"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "with bad json",
			inputString: `{
				"int1": 3,
				"int2": 5,
				"str1": "fizz"
				"str2": "buzz"
			}`,
			expectedCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := New()
			req, err := http.NewRequest("GET", "/fizzbuzz", strings.NewReader(tt.inputString))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(service.FizzBuzzEndpoint)
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.expectedCode, rr.Code)
		})
	}
}
