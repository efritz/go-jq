# Go-JQ

[![GoDoc](https://godoc.org/github.com/efritz/go-jq?status.svg)](https://godoc.org/github.com/efritz/go-jq)
[![Build Status](https://secure.travis-ci.org/efritz/go-jq.png)](http://travis-ci.org/efritz/go-jq)
[![Maintainability](https://api.codeclimate.com/v1/badges/9aab8d8dce9e96f2ab9a/maintainability)](https://codeclimate.com/github/efritz/go-jq/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/9aab8d8dce9e96f2ab9a/test_coverage)](https://codeclimate.com/github/efritz/go-jq/test_coverage)

Go bindings for JQ.

## Example

This package defines only one function, `Run`. This function takes the JQ
expression and the input on which to evaluate it.

```go
results, err := jq.Run(".[] | .id", values)
if err != nil {
    // Handle error
}

for _, result := range results {
    fmt.Printf("Matched: %#v\n", result)
}
```

## License

Copyright (c) 2019 Eric Fritz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
