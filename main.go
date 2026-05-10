package main

import (
	"fmt"

	"github.com/jjsiv/idp/internal/utils"
)

func main() {

	set := make(utils.Set)
	set.Insert("admin")

	fmt.Println(set.Has("def"))
}
