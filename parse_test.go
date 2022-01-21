package main

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	dat, err := os.ReadFile("env.toml")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(dat))

	var env Env
	err = toml.Unmarshal([]byte(string(dat)), &env)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(env.FilterInclude)
	fmt.Println(env.FilterExclude)
}
