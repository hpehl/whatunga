package main

import (
	"fmt"
	"testing"
)

func TestPath(_ *testing.T) {
	l := lex("name", "some {{.in}} put", "{{", "}}")
	var items []item
	for {
		item := l.nextItem()
		items = append(items, item)
		if item.typ == itemEOF || item.typ == itemError {
			break
		}
	}
	fmt.Println("Items: ", items)
}
