package utils

import (
	"math/rand"
	"time"
)

var SeededRand = rand.New(rand.NewSource(time.Now().UnixNano()))
