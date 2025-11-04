package main

import "fmt"

// Number type constraint as an interface. The constraint allows any type implementing the interface. For example,
// if you declare a type constraint interface with three methods,
// then use it with a type parameter in a generic function, type arguments used to call the function must have all of those methods.
type Number interface {
	int64 | float64
}

// SumInts adds together the values of m.
func SumInts(m map[string]int64) int64 {
	var s int64
	for _, v := range m {
		s += v
	}
	return s
}

// SumFloats adds together the values of m.
func SumFloats(m map[string]float64) float64 {
	var s float64
	for _, v := range m {
		s += v
	}
	return s
}

// SumIntsOrFloats sums the values of map m. It supports both int64 and float64
// as types for map values.
func SumIntsOrFloats[K comparable, V int64 | float64](m map[K]V) V {
	// 	- Specify for the K type parameter the type constraint comparable. Intended specifically for cases like these,
	//	the comparable constraint is predeclared in Go. It allows any type whose values may be used as an operand of
	//	the comparison operators == and !=. Go requires that map keys be comparable.
	//	So declaring K as comparable is necessary so you can use K as the key in the map variable.
	//	It also ensures that calling code uses an allowable type for map keys.
	// 	- Specify for the V type parameter a constraint that is a union of two types: int64 and float64.
	//	Using | specifies a union of the two types, meaning that this constraint allows either type.
	//	Either type will be permitted by the compiler as an argument in the calling code.
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

// SumNumbers sums the values of map m. It supports both integers
// and floats as map values.
func SumNumbers[K comparable, V Number](m map[K]V) V {
	// - Essentially, you’re moving the union from the function declaration into a new type constraint.
	// That way, when you want to constrain a type parameter to either int64 or float64,
	// you can use this Number type constraint instead of writing out int64 | float64.
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

func main() {
	// Initialize a map for the integer values
	ints := map[string]int64{
		"first":  34,
		"second": 12,
	}

	// Initialize a map for the float values
	floats := map[string]float64{
		"first":  35.98,
		"second": 26.99,
	}

	fmt.Printf("Non-Generic Sums: %v and %v\n",
		SumInts(ints),
		SumFloats(floats))

	// Specify type arguments – the type names in square brackets – to be clear about the types that should replace type parameters in the function you’re calling.
	fmt.Printf("Generic Sums: %v and %v\n",
		SumIntsOrFloats[string, int64](ints),
		SumIntsOrFloats[string, float64](floats))

	// You can omit type arguments in calling code when the Go compiler can infer the types you want to use. The compiler infers type arguments from the types of function arguments.
	fmt.Printf("Generic infered Sums: %v and %v\n",
		SumIntsOrFloats(ints),
		SumIntsOrFloats(floats))

	fmt.Printf("Generic Sums with Constraint: %v and %v\n",
		SumNumbers(ints),
		SumNumbers(floats))
}
