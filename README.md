# martini-etag

## Usage

```go
package main

import . "github.com/cgarvis/martini-etag"
import "github.com/codegangsta/martini"

func main() {
  m := martini.Classic()
  m.Use(ETag())
  m.Get("/", func() string {
    return "Hello World!"
  })
  m.Run()
}
```
