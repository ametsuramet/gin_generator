# GIN GENERATOR

INSTALL PACKAGE:
```bash
go get -u github.com/ametsuramet/gin_generator
```

Example:
Create file ex: generator.go
```Go
package main

import (
	// "fmt"
	"github.com/ametsuramet/gin_generator"
	"path/filepath"
)

func main() {
	jsonFile, _ := filepath.Abs("builder.json")
	path, _ := filepath.Abs("")
	gen := gin_generator.Set(jsonFile, path, nil)

	gen.Generate()
}

```
Create file json: generator.go
```json
[
	{
		"name": "AboutUs",
		"schema": [
			{
				"field": "banner",
				"type": "string",
			},
			{
				"field": "image",
				"type": "string",
			},
			{
				"field": "who_we_are_image",
				"type": "string",
			},
			{
				"field": "active_flag",
				"type": "boolean",
			},
			{
				"field": "created_by",
				"type": "integer::unsigned",
			},
			{
				"field": "updated_by",
			}
		]
	}
]
```

RUN:
```bash
go run generator.go
```