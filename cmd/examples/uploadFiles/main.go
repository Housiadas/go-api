package main

import (
	"flag"
	"log"
	"net/http"
)

const html = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="X-UA-Compatible" content="ie=edge" />
    <title>File upload demo</title>
  </head>
  <body>
    <form
      id="form"
      enctype="multipart/form-data"
      action="http://localhost:4000/v1/users/1/documents-upload"
      method="POST"
    >
      <input class="input file-input" type="file" name="file" multiple />
      <button class="button" type="submit">Submit</button>
    </form>
  </body>
</html>`

func main() {
	addr := flag.String("addr", ":9001", "Server address")
	flag.Parse()
	log.Printf("starting server on %s", *addr)
	// Start an HTTP server listening on the given address, which responds to all
	// requests with the webpage HTML above.
	err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(html))
		if err != nil {
			return
		}
	}))
	log.Fatal(err)
}
