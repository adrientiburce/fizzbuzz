package fizzbuzz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-playground/validator/v10"
)

// Service to handle a fizzbuzz server
type Service struct {
	fizzBuzzStat Statistics
	sync.Mutex
}

// Statistics store satisticts on used parameters for fizzbuzz endpoint
type Statistics struct {
	Int1Stat  map[int]int    `json:"int1"`
	Int2Stat  map[int]int    `json:"int2"`
	LimitStat map[int]int    `json:"limit"`
	Str1Stat  map[string]int `json:"str1"`
	Str2Stat  map[string]int `json:"str2"`
}

// FizzBuzz contains all parameters needed for a fizz buzz game
type FizzBuzz struct {
	Int1  int    `json:"int1" validate:"required"`
	Int2  int    `json:"int2" validate:"required"`
	Limit int    `json:"limit" validate:"required"`
	Str1  string `json:"str1" validate:"required"`
	Str2  string `json:"str2" validate:"required"`
}

// New create a Fizzbuzz service
func New() *Service {
	return &Service{
		fizzBuzzStat: Statistics{
			Int1Stat:  make(map[int]int),
			Int2Stat:  make(map[int]int),
			LimitStat: make(map[int]int),
			Str1Stat:  make(map[string]int),
			Str2Stat:  make(map[string]int),
		},
	}
}

// computeSuite compute all numbers until limit with fizzbuzz rule
func (s *FizzBuzz) computeSuite() (res string) {
	if s.Limit < 1 || s.Int1 == 0 || s.Int2 == 0 {
		return
	}

	for i := 1; i <= s.Limit; i++ {
		if i%(s.Int1*s.Int2) == 0 {
			res += s.Str1 + s.Str2
		} else if i%s.Int1 == 0 {
			res += s.Str1
		} else if i%s.Int2 == 0 {
			res += s.Str2
		} else {
			res += strconv.Itoa(i)
		}
		res += ","
	}

	return res[:len(res)-1]
}

// Statistics endpoint to show fizzbuzz statistics fir every parameters
func (s *Service) Statistics(w http.ResponseWriter, r *http.Request) {
	res, err := json.MarshalIndent(s.fizzBuzzStat, "", "\t")
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request (%s)", err), http.StatusInternalServerError)
	}
	fmt.Fprint(w, string(res))
}

// FizzBuzzEndpoint show fizzbuzz suite according to input parameters
func (s *Service) FizzBuzzEndpoint(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read body (%s)", err), http.StatusInternalServerError)
		return

	}

	var fizzBuzz FizzBuzz
	if err = json.Unmarshal(reqBody, &fizzBuzz); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request (%s)", err), http.StatusInternalServerError)
		return
	}

	if err := validateStruct(fizzBuzz); err != nil {
		http.Error(w, fmt.Sprintf("Missing params: %s", err), http.StatusBadRequest)
		return
	}

	s.Lock()
	s.fizzBuzzStat.Int1Stat[fizzBuzz.Int1]++
	s.fizzBuzzStat.Int2Stat[fizzBuzz.Int2]++
	s.fizzBuzzStat.LimitStat[fizzBuzz.Limit]++
	s.fizzBuzzStat.Str1Stat[fizzBuzz.Str1]++
	s.fizzBuzzStat.Str2Stat[fizzBuzz.Str2]++
	s.Unlock()

	res := fizzBuzz.computeSuite()
	fmt.Fprint(w, res)
}

func validateStruct(i interface{}) error {
	validate := validator.New()
	if err := validate.Struct(i); err != nil {
		return err
	}
	return nil
}
