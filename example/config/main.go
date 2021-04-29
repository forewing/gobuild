package main

import (
	"github.com/forewing/gobuild"
)

func main() {
	target, err := gobuild.GetTargetFromJson("./config.json")
	if err != nil {
		panic(err)
	}
	err = target.Build()
	if err != nil {
		panic(err)
	}
}
