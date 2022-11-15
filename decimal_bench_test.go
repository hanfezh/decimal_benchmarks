package main

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	edecimal "github.com/ericlagergren/decimal"
	mydecimal "github.com/pingcap/tidb/types"
	sdecimal "github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

type pair struct {
	a float64
	b float64
}

const (
	numTestPairs = 1000
)

var (
	testPairs   []pair
	initialized = false
)

func init() {
	initTestPairs()
}

func genRandFloat() float64 {
	f := rand.Float64()
	digits := rand.Int()%12 + 1
	offset := rand.Int()%digits + 1
	f = math.Round(f*math.Pow10(digits)) / math.Pow10(digits-offset)
	if f != 0.0 {
		return f
	}
	return genRandFloat()
}

func initTestPairs() {
	if initialized {
		return
	}
	for i := 0; i < numTestPairs; i++ {
		a := genRandFloat()
		b := genRandFloat()
		testPairs = append(testPairs, pair{a, b})
	}
	initialized = true
}

func shopDecimalAdd(b *testing.B, p pair) {
	d1 := sdecimal.NewFromFloat(p.a)
	d2 := sdecimal.NewFromFloat(p.b)
	result := d1.Add(d2)
	_, _ = result.Float64()
	_ = result.String()
}

func shopDecimalSub(b *testing.B, p pair) {
	d1 := sdecimal.NewFromFloat(p.a)
	d2 := sdecimal.NewFromFloat(p.b)
	result := d1.Sub(d2)
	_, _ = result.Float64()
	_ = result.String()
}

func shopDecimalMul(b *testing.B, p pair) {
	d1 := sdecimal.NewFromFloat(p.a)
	d2 := sdecimal.NewFromFloat(p.b)
	result := d1.Mul(d2)
	_, _ = result.Float64()
	_ = result.String()
}

func shopDecimalDiv(b *testing.B, p pair) {
	d1 := sdecimal.NewFromFloat(p.a)
	d2 := sdecimal.NewFromFloat(p.b)
	result := d1.Div(d2)
	_, _ = result.Float64()
	_ = result.String()
}

func ericDecimalAdd(b *testing.B, p pair) {
	d1 := new(edecimal.Big).SetFloat64(p.a)
	d2 := new(edecimal.Big).SetFloat64(p.b)
	_, _ = d1.Add(d1, d2).Float64()
	_ = d1.String()
}

func ericDecimalSub(b *testing.B, p pair) {
	d1 := new(edecimal.Big).SetFloat64(p.a)
	d2 := new(edecimal.Big).SetFloat64(p.b)
	_, _ = d1.Sub(d1, d2).Float64()
	_ = d1.String()
}

func ericDecimalMul(b *testing.B, p pair) {
	d1 := new(edecimal.Big).SetFloat64(p.a)
	d2 := new(edecimal.Big).SetFloat64(p.b)
	_, _ = d1.Mul(d1, d2).Float64()
	_ = d1.String()
}

func ericDecimalDiv(b *testing.B, p pair) {
	d1 := new(edecimal.Big).SetFloat64(p.a)
	d2 := new(edecimal.Big).SetFloat64(p.b)
	_, _ = d1.Quo(d1, d2).Float64()
	_ = d1.String()
}

func myDecimalAdd(b *testing.B, p pair) {
	d1 := mydecimal.NewDecFromFloatForTest(p.a)
	d2 := mydecimal.NewDecFromFloatForTest(p.b)
	var result mydecimal.MyDecimal
	err := mydecimal.Add(d1, d2, &result)
	require.Equal(b, err, nil)
	_, _ = result.ToFloat64()
	_ = result.String()
}

func myDecimalSub(b *testing.B, p pair) {
	d1 := mydecimal.NewDecFromFloatForTest(p.a)
	d2 := mydecimal.NewDecFromFloatForTest(p.b)
	var result mydecimal.MyDecimal
	err := mydecimal.Sub(d1, d2, &result)
	require.Equal(b, err, nil)
	_, _ = result.ToFloat64()
	_ = result.String()
}

func myDecimalMul(b *testing.B, p pair) {
	d1 := mydecimal.NewDecFromFloatForTest(p.a)
	d2 := mydecimal.NewDecFromFloatForTest(p.b)
	var result mydecimal.MyDecimal
	err := mydecimal.Mul(d1, d2, &result)
	require.Equal(b, err, nil)
	_, _ = result.ToFloat64()
	_ = result.String()
}

func myDecimalDiv(b *testing.B, p pair) {
	d1 := mydecimal.NewDecFromFloatForTest(p.a)
	d2 := mydecimal.NewDecFromFloatForTest(p.b)
	var result mydecimal.MyDecimal
	err := mydecimal.Div(d1, d2, &result, 5)
	require.Equal(b, err, nil)
	_, _ = result.ToFloat64()
	_ = result.String()
}

func BenchmarkDecimal(b *testing.B) {
	type op struct {
		name string
		fn   func(b *testing.B, p pair)
	}
	pkgs := []struct {
		name string
		ops  []op
	}{
		{"shopspring/decimal", []op{
			{"add", shopDecimalAdd},
			{"sub", shopDecimalSub},
			{"mul", shopDecimalMul},
			{"div", shopDecimalDiv}},
		},
		{"ericlagergren/decimal", []op{
			{"add", ericDecimalAdd},
			{"sub", ericDecimalSub},
			{"mul", ericDecimalMul},
			{"div", ericDecimalDiv}},
		},
		{"tidb/types/mydecimal", []op{
			{"add", myDecimalAdd},
			{"sub", myDecimalSub},
			{"mul", myDecimalMul},
			{"div", myDecimalDiv}},
		},
	}
	for _, pkg := range pkgs {
		for _, op := range pkg.ops {
			b.Run(fmt.Sprintf("pkg=%s/op=%s", pkg.name, op.name), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					for _, p := range testPairs {
						op.fn(b, p)
					}
				}
			})
		}
	}
}
