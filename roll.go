package main

import (
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// represents like 2d20
type Roll struct {
	Multiplier int
	Sides      int
}

var IncorrectFormatting = errors.New("Incorrect formatting")

func NewRoll(roll string) (*Roll, error) {
	r := &Roll{
		Multiplier: 1,
		Sides:      -1,
	}

	trimmed := strings.TrimSpace(roll)
	vals := strings.Split(trimmed, "d")
	if len(vals) != 2 {
		return nil, IncorrectFormatting
	}

	posMultiplier := vals[0]

	// if blank, there's no multiplier at all!
	if posMultiplier == "" {
		r.Multiplier = 1
	}

	mult, err := strconv.Atoi(posMultiplier)

	// if it parses correctly, set multiplier
	if err == nil {
		r.Multiplier = mult
	}

	posSides := vals[1]

	sides, err := strconv.Atoi(posSides)

	// must have amount of sides of dice
	if err != nil || sides < 1 {
		return nil, IncorrectFormatting
	}

	r.Sides = sides
	return r, nil
}

func (r Roll) Calc() int {
	// reset seed each time
	rand.Seed(time.Now().UnixNano())
	sum := 0

	for i := 0; i < r.Multiplier; i++ {
		// [0, n)
		sum += rand.Intn(r.Sides) + 1
	}

	return sum
}
