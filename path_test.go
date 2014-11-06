package main

import (
	"fmt"
	"testing"
)

func TestPath(_ *testing.T) {
	l := lex("name", "some {{.in}} put", "{{", "}}")
	for item := range l.items {
		fmt.Println("Item: ", item)
	}
}
