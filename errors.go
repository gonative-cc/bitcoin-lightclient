package main

import (
	"fmt"
)

type InvalidHeaderErr struct {
	HeaderHex string
	Id        int
}

func NewInvalidHeaderErr(headerHex string, id int) InvalidHeaderErr {
	return InvalidHeaderErr{
		HeaderHex: headerHex,
		Id:        id,
	}
}

func (e InvalidHeaderErr) Error() string {
	return fmt.Sprintf("invalid header %s at index %d", e.HeaderHex, e.Id)
}
