[![GO Report](https://goreportcard.com/badge/github.com/rck/unit)](https://goreportcard.com/report/github.com/rck/unit)
[![GoDoc](https://godoc.org/github.com/rck/unit?status.svg)](https://godoc.org/github.com/rck/unit)

# unit
Unit is a library to parse user defined units in go. It was written with sizes in mind, but everything where
the numeric part can be converted to an `int64` is possible.

It can be used standalone, but also implements the `flag` interfaces.

# Installing
```
go get -u github.com/rck/unit
```
Next, include `unit` in your application:

```golang
import "github.com/rck/unit"
```

# Example

```golang
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/rck/unit"
)

func main() {
	size := unit.MustNewUnit(unit.DefaultUnits).MustNewValue(1, unit.None)
	flag.Var(size, "s", "Parse a size")
	flag.Parse()
	fmt.Println(size.Value)

	myUnits := unit.DefaultUnits
	myUnits["kB"] = unit.DefaultUnits["KB"]

	u, err := unit.NewUnit(myUnits)
	if err != nil {
		log.Fatalf("Could not create unit based on mapping: %v\n", err)
	}
	if v, err := u.ValueFromString("+1024K"); err == nil {
		fmt.Println(v.Value, v)
		switch v.ExplicitSign {
		case unit.Positive:
			fmt.Println("Explicit positive sign")
		case unit.Negative:
			fmt.Println("Explicit negative sign")
		case unit.None:
			fmt.Println("No xxplicit sign")
		default:
			fmt.Println("This can not happen :-)")
		}
	}

	// if called with -s 20M, prints:
	// 20971520
	// 1048576 +1M
	// Explicit positive sign
}
```

# LICENSE
This project is licensed under the same terms as [u-root](https://github.com/u-root).
