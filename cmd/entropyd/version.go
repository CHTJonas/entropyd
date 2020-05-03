package main

import "fmt"

type ver struct {
	Major int
	Minor int
	Patch int
}

func (v *ver) getString() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func getVer() *ver {
	return &ver{
		Major: 1,
		Minor: 3,
		Patch: 0,
	}
}
