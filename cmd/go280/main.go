package main

import (
	"go280"

	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(go280.Analyzer) }
