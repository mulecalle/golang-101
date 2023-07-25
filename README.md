# go-generics


```go
// Final K is the constraint comparable, V is the constraint
func Final[K comparable, V Number](m map[K]V) V {
	var s V
	for _, v := range m {
		s += v
	}
	return s
}

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

fmt.Printf("Generic Sums with Constraint: %v and %v\n", Final(ints), Final(floats))
```