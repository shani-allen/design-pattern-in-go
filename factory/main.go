package main

import "fmt"

type Factory interface {
	Create() IShoe
}

func NewFactory(brand string) Factory {
	if brand == "nike" {
		return &Nike{
			Shoe: Shoe{},
		}
	}
	return nil
}

type Shoe struct {
}

func (s *Shoe) Create() IShoe {
	fmt.Println("create")
	return &Nike{
		Shoe: Shoe{},
	}
}

type IShoe interface {
	Check()
}

func (s *Nike) Check() {
	fmt.Println("check")
}

type Nike struct {
	Shoe
}

func main() {
	nike := NewFactory("nike")
	check := nike.Create()
	check.Check()
}
