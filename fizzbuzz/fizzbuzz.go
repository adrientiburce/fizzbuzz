package fizzbuzz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/prologic/bitcask"
)

// Service to handle a fizzbuzz server
type Service struct {
	Db *bitcask.Bitcask
}

type Stat struct {
	Int1 map[int]int
	Int2 map[int]int
}

// FizzBuzz contains infos for a fizz buzz game
type FizzBuzz struct {
	Int1  int    `json:"int1"`
	Int2  int    `json:"int2"`
	Limit int    `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

// New create a Fizzbuzz service
func New(dbPath string) (*Service, error) {
	db, err := bitcask.Open(dbPath, bitcask.WithSync(true))
	if err != nil {
		log.Fatal("can't open databse")
		return nil, err
	}

	return &Service{
		Db: db,
	}, nil
}

// FizzBuzz compute all numbers until limit with fizzbuzz rule
func (s *FizzBuzz) computeSuite() string {

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

// Statistics endpoint to show server stats
func (s *Service) Statistics(w http.ResponseWriter, r *http.Request) {
	res := ""
	for key := range s.Db.Keys() {
		val, _ := s.Db.Get(key)
		res += fmt.Sprintf("%s -> %v \n", key, string(val))
	}

	fmt.Fprint(w, res)
}

// DeleteAll remove old db keys, used mainly for debug purpose
func (s *Service) DeleteAll(w http.ResponseWriter, r *http.Request) {
	if err := s.Db.DeleteAll(); err != nil {

		http.Error(w, fmt.Sprintf("failed delete database (%s)", err), http.StatusInternalServerError)
	}
	fmt.Fprint(w, "Database reset")
}

// FizzBuzzEndpoint show fizzbuzz suite
func (s *Service) FizzBuzzEndpoint(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)

	var fizzBuzz FizzBuzz
	err := json.Unmarshal(reqBody, &fizzBuzz)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request (%s)", err), http.StatusInternalServerError)
		return
	}
	if fizzBuzz.Limit == 0 {
		http.Error(w, "fizzbuzz limit empty", http.StatusBadRequest)
	} else if fizzBuzz.Int1 == 0 || fizzBuzz.Int2 == 0 {
		http.Error(w, "fizzbuzz int1 or int2 empty", http.StatusBadRequest)
	} else if fizzBuzz.Str1 == "" || fizzBuzz.Str2 == "" {
		http.Error(w, "fizzbuzz str1 or str2 empty", http.StatusBadRequest)
	}

	s.incrementDbValue("str1:" + fizzBuzz.Str1)
	s.incrementDbValue("str2:" + fizzBuzz.Str2)
	s.incrementDbValue("int1:" + strconv.Itoa(fizzBuzz.Int1))
	s.incrementDbValue("int2:" + strconv.Itoa(fizzBuzz.Int2))
	s.incrementDbValue("limit:" + strconv.Itoa(fizzBuzz.Limit))

	res := fizzBuzz.computeSuite()
	fmt.Fprint(w, res)
}

func (s *Service) incrementDbValue(key string) {
	val, _ := s.Db.Get([]byte(key))
	if len(val) == 0 {
		if err := s.Db.Put([]byte(key), []byte(strconv.Itoa(1))); err != nil {
			log.Print(err)
		}
	} else {
		oldCount, err := strconv.Atoi(string(val))
		if err != nil {
			log.Print(err)
		}
		if err := s.Db.Put([]byte(key), []byte(strconv.Itoa(oldCount+1))); err != nil {
			log.Print(err)
		}
	}
}
