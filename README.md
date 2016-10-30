# Go Package: srcutil

  [About](#about) | [Usage](#usage) | [Go Doc](https://godoc.org/github.com/cstockton/go-srcutil)

  > Get:
  > ```bash
  > go get -u github.com/cstockton/go-srcutil
  > ```
  >
  > Example:
  > ```Go
  > pkg, err := srcutil.Import("io")
  > pkg.Doc
  > > Package io provides basic interfaces to I/O primitives.
  > ```


## About

Package srcutil provides utilities for working with Go source code. The Go
standard library provides a powerful suite of packages "go/{ast,doc,...}"
which are used by the Go tool chain to compile Go programs. As you initially
try to find your way around you hit a small dependency barrier and have to
learn a small portion of each package. There is a fantastic write up and
collection of examples that I used to learn (or shamelessly copy pasta'd)
while creating this package, currently maintained by:

```
  Alan Donovan (https://github.com/golang/example/tree/master/gotypes)
```

In the mean time this package can help you get started with some common use
cases.


## Bugs and Patches

  Feel free to report bugs and submit pull requests.

  * bugs:
    <https://github.com/cstockton/go-srcutil/issues>
  * patches:
    <https://github.com/cstockton/go-srcutil/pulls>



[Go Doc]: https://godoc.org/github.com/cstockton/go-srcutil
