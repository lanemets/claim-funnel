package main

import (
	"github.com/lanemets/claim-funnel/cmd"
	"github.com/lanemets/claim-funnel/interfaces/benerest"
)

func main() {
	cmd.Execute(
		benerest.NewServer(),
	)
}
