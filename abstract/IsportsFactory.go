package main

import "errors"

type ISportsFactory interface {
	MakeShoe() IShoe
	MakeShirt() IShirt
}

func GetShoeFactory(brand string) (ISportsFactory, error) {
	if brand == "nike" {
		return &Nike{}, nil
	}
	if brand == "adidas" {
		return &Adidas{}, nil
	}
	return nil, errors.New("Invalid brand")
}
