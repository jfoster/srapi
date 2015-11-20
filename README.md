speedrun.com API Client for Go
==============================

This Go package implements a client for the
[speedrun.com API](https://github.com/speedruncom/api). It's not 100% complete
and a relatively direct mapping of API structures to Go ``struct`` values.

Installation
------------

```
go get github.com/sgt-kabukiman/srapi
```

Usage
-----

```go
package main

import (
	"fmt"

	"github.com/sgt-kabukiman/srapi"
)

func main() {
	// optional sorting
	sort := &srapi.Sorting{"name", srapi.Descending}

	// pagination
	cursor := &srapi.Cursor{2, 2} // offset, max

	regions, err := srapi.Regions(sort, cursor)
	if err != nil {
		panic(err) // err is an srapi.*ApiError struct, containing more information
	}

	fmt.Printf("regions = %+v\n", regions)
}
```

License
-------

This code is licensed under the MIT license.