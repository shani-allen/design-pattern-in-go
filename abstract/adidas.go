// concrete type
package main

type Adidas struct{}

func (a *Adidas) MakeShoe() IShoe {
	return AdidasShoe{
		&Shoe{
			Logo: "Adidas",
			Size: 36,
		},
	}
}

func (a *Adidas) MakeShirt() IShirt {
	return AdidasShirt{
		&Shirt{
			Logo: "Adidas",
			Size: 36,
		},
	}
}
