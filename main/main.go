package main

import (
	"access"
	"fmt"
)

type A struct {
	AA  *string  `default:"48d4v84v8df"`
	B   *int     `default:"455555"`
	C   *int8    `default:"4"`
	D   *int16   `default:"2"`
	E   *int32   `default:"2323"`
	F   *int64   `default:"233333323"`
	G   *uint    `default:"33333323"`
	H   *uint8   `default:"36"`
	I   *uint16  `default:"666"`
	J   *uint32  `default:"662226"`
	K   *uint64  `default:"665444442226"`
	L   *bool    `default:"true"`
	M   *float32 `default:"3.555"`
	N   *float64 `default:"33.5555555"`
	ZZZ *BBBB    `default:"5555555"`
	O   []bool   `default:"true,false,false,true"`

	P []*BBBB `default:"4,6,7,8,8"`
}

type BBBB struct {
	aaa []string
}

func (bb *BBBB) Default(val string) error {
	bb.aaa = []string{val}
	return nil
}

func (bb BBBB) IsZero() bool {
	return true
}

func main() {
	aaa := A{}
	AA := "ccccccc"
	aaa.AA = &AA
	access.Set(&aaa)

	fmt.Println(*aaa.AA)
	fmt.Println(*aaa.B)
	fmt.Println(*aaa.C)
	fmt.Println(*aaa.D)
	fmt.Println(*aaa.E)
	fmt.Println(*aaa.F)
	fmt.Println(*aaa.G)
	fmt.Println(*aaa.H)
	fmt.Println(*aaa.I)
	fmt.Println(*aaa.J)
	fmt.Println(*aaa.K)
	fmt.Println(*aaa.L)
	fmt.Println(*aaa.M)
	fmt.Println(*aaa.N)

	fmt.Println(aaa.O)
	fmt.Println(aaa.P)

	fmt.Println(aaa.ZZZ)
}
