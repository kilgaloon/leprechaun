package main

import (
	"fmt"
	"os"

	"github.com/kilgaloon/leprechaun/cmd/lepretools/validate"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("HELP")

		fmt.Println("- lepretools validate:recipe {path_to_recipe.yml}")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			fmt.Println("Something went wrong, if you aren't sure how to use lepretools")
			fmt.Println("use command: 'lepretools' to get list of available commands")
		}
	}()

	switch os.Args[1] {
	case "recipe:validate":
		validate.CheckRecipe(os.Args[2])
	case "recipe:create":
		fmt.Println("Not supported yet!")
	default:
		panic("Command does not exist")
	}

}
