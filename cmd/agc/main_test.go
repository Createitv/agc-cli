package main

import (
	"os"
	"testing"
)

func TestMainFunctionVersionPath(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = []string{"agc", "version"}
	main()
}
