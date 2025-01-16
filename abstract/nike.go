package main

type Nike struct{}

func (n *Nike) MakeShoe() IShoe {
	return NikeShoe{
		&Shoe{
			Logo: "Nike",
			Size: 36,
		},
	}
}

func (n *Nike) MakeShirt() IShirt {
	return NikeShirt{
		&Shirt{
			Logo: "Nike",
			Size: 36,
		},
	}
}
