package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	c := GetConfig()

	fmt.Printf("%v", c)
}
