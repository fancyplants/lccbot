package main

import (
	// "../main"

	"testing"
	"fmt"
)

func TestRolls(t *testing.T) {
	r, err := NewRoll("3d6")
	if err != nil {
		t.Error(err)
		return
	}


	m := make(map[int]int)

	for i := 0; i < 50000; i++ {
		val := r.Calc()
		_, ok := m[val]
		if !ok {
			m[val] = 1
		} else {
			m[val] += 1
		}
	}
	
	fmt.Println(m)
}