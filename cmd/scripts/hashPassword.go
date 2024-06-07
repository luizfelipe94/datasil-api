package main

import (
	"fmt"

	"github.com/luizfelipe94/datasil/modules/auth"
)

func main() {
	password, _ := auth.HashPassword("123")
	fmt.Println(password)
}
