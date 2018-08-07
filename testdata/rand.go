package main

import(
	"math/rand"
)

func insecureRand() int {
	// there are many possible uses for math/rand
	// it's impractical to check for every possible use
	
	return rand.Int();
}
