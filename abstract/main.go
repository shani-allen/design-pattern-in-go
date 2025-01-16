package main

import (
	"fmt"
)

func main() {

	factory, _ := GetShoeFactory("adidas")
	shirts := factory.MakeShirt()

	// prints the shirt Logo
	shirts.SetLogo("Adidas")
	fmt.Println(shirts.GetLogo())
}
