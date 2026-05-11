package main

import (
	"encoding/json"
	"fmt"
	"log"

	idpv1alpha1 "github.com/jjsiv/idp/api/v1alpha1"
)

func main() {
	params := []idpv1alpha1.ResourceParameter{
		{
			Name:  "param1",
			Value: "value2",
		},
		{
			Name:  "param2",
			Value: "value3",
		},
	}

	json, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(json))
}
