package criteria

import (
	"fmt"
	"time"
)

type equalsFilter struct {
	Type   string
	Field  string
	ValStr *string
	ValInt *int
}

type equalsAnyFilter struct {
	Type  string
	Field string
	Value []string
}

type containsFilter struct {
	Type  string
	Field string
	Value string
}

type rangeFilter struct {
	Type       string
	Field      string
	Parameters []parameters
}
type parameters struct {
	Name    string
	ValInt  *int
	ValDate time.Time
}

func (f rangeFilter) Check() {
	for _, x := range f.Parameters {
		switch x.Name {
		case "gte":
			return
		case "lte":
			return
		case "gt":
			return
		case "lt":
			return
		default:
			error := fmt.Errorf("%s is not a valid range", x.Name)
			panic(error)
		}
	}
}

type FilterType struct{}
