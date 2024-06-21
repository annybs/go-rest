# Go REST

Some handy functions for developing JSON-based REST APIs. In particular, it simplifies reading HTTP request bodies, writing HTTP response bodies, and handling errors.

## Error handling

You can use `errors.Is()` to ascertain the type of errors thrown by validation functions, but for the most part, this isn't necessary because the write functions already do that.

## Example

```go
package main

import (
	"errors"
	"math/rand"
	"net/http"

	"github.com/annybs/go/rest"
)

type Handler struct{}

func (*Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	n := rand.Intn(3)
	if n == 0 {
		rest.WriteResponseJSON(w, http.StatusOK, map[string]string{"status": "OK"})
	} else if n == 1 {
		rest.WriteErrorJSON(w, errors.New("the original error message is added to data.error"))
	} else {
		rest.WriteErrorJSON(w, rest.ErrNotFound)
	}
}

func main() {
	http.ListenAndServe("localhost:8000", &Handler{})
}
```

Open <http://localhost:8000> in your browser and refresh a bunch of times to see the different possible responses.

## License

See [LICENSE.md](../LICENSE.md)
