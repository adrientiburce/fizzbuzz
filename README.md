
[![Build status](https://github.com/adrientiburce/fizzbuzz/actions/workflows/ci.yml/badge.svg)](https://github.com/adrientiburce/fizzbuzz/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/adrientiburce/fizzbuzz/branch/master/graph/badge.svg?token=UH5M5HCWJU)](https://codecov.io/gh/adrientiburce/fizzbuzz)

# Endpoints 

- /fizzbuzz (GET) with a body containing 5 parameters : displays the fizzbuzz suite

Request body example :

```json
{
	"int1": 3,
	"int2": 5,
	"limit": 100,
	"str1": "fizz",
	"str2": "buzz"
}
```

- /statistics (GET) : displays the number of times each parameter has been called

# Fizz-buzz REST server

The original fizz-buzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by "fizz", all multiples of 5 by "buzz", and all multiples of 15 by "fizzbuzz".
The output would look like this: "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,...".

- Your goal is to implement a web server that will expose a REST API endpoint that:

Accepts five parameters : three integers int1, int2 and limit, and two strings str1 and str2.
Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.

The server needs to be:
 - Ready for production
 - Easy to maintain by other developers

Add a statistics endpoint allowing users to know what the most frequent request has been.
This endpoint should:
- Accept no parameter
- Return the parameters corresponding to the most used request, as well as the number of hits for this request

