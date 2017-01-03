# go-withcontext-checker

## What

Checks that the result of a call to \*net/http.Request.WithContext() is assigned to a value

## Why

According to [the docs](https://golang.org/pkg/net/http/#Request.WithContext):

> WithContext returns a shallow copy of r with its context changed to ctx.

Returning a shallow copy instead of modifying the context in-place was
surprising behavior to me.  This tool checks that all calls to
\*net/http.Request.WithContext() happens as part of an assignment.

## How

```
./go-withcontext-checker -t github.com/bubaflub/go-withcontext-checker/examples
net/http request.WithContext() called without lvalue at /Users/bob/gohome/src/github.com/bubaflub/go-withcontext-checker/examples/simple_test.go:12:2: "WithContext"`
```

## Todo

* [ ] Have tests that run as part of `go test`
* [ ] Print the full line where the call is made
* [ ] Provide a fixit hint like clang-tidy does
* [ ] Accept a list of functions that we should flag like GCC's `warn_unused_result`
* [ ] Can we find all "pure" functions in stdlib to generate that list?
* [ ] Open and issue to propose patch to extend `go vet`
