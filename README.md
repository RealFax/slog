## structured logging
_All usage of this package is the same as `log/slog`_

### install
```shell
go get github.com/RealFax/slog@latest
```

### performance
you can use build tags to declare which JSON serialization method to use

e.g., use the jsoniter flag to declare the use of `json-iterator/go`

### environment variables

#### LOG_PURE
type: boolean

This variable declares whether colorful output is disabled.

e.g., LOG_PURE=true

#### LOG_LEVEL
type: string

- debug
- info
- warn
- error

default: info
e.g., LOG_LEVEL=debug || LOG_LEVEL=DEBUG