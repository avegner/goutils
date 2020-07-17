package main

import (
	"fmt"
	"os"

	"github.com/avegner/utils/json"
)

const configPath = "config.json"

type Config struct {
	Field1 string `json:"field-1"`
	Field2 bool   `json:"field-2"`
	Field3 int    `json:"field-3"`
}

func main() {
	cfg := Config{}
	if err := json.UnmarshalFile(configPath, &cfg); err != nil {
		fmt.Printf("unmarshal failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("config: %+v\n", cfg)
}
