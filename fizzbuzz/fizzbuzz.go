package fizzbuzz

import "strconv"

// FizzBuzz contains infos for a fizz buzz game
type FizzBuzz struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

// New create a Fizzbuzz service
func New(int1, int2, limit int, str1, str2 string) *FizzBuzz {
	return &FizzBuzz{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  str1,
		Str2:  str2,
	}
}

// FizzBuzz compute all numbers until limit with fizzbuzz rule
func (s *FizzBuzz) FizzBuzz() string {

	if s.Limit < 1 || s.Int1 == 0 || s.Int2 == 0 {
		return ""
	}
	res := ""

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
