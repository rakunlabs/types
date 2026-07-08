# types

[![License](https://img.shields.io/github/license/worldline-go/types?color=red&style=flat-square)](https://raw.githubusercontent.com/worldline-go/types/main/LICENSE)
[![Coverage](https://img.shields.io/sonar/coverage/worldline-go_types?logo=sonarcloud&server=https%3A%2F%2Fsonarcloud.io&style=flat-square)](https://sonarcloud.io/summary/overall?id=worldline-go_types)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/worldline-go/types/test.yml?branch=main&logo=github&style=flat-square&label=ci)](https://github.com/worldline-go/types/actions)
[![Go PKG](https://raw.githubusercontent.com/worldline-go/guide/main/badge/custom/reference.svg)](https://pkg.go.dev/github.com/worldline-go/types)

Library for types conversion.

```sh
go get github.com/worldline-go/types
```

## Usage

### types.Map

Map based on `map[string]any` and null values stay as nil.

> This type not convert to base64 when marshaling, it is directly a json string.

```go
type Train struct {
	Details types.Map `db:"details"`
}
```

### types.RawJSON

`[]byte` type behind, same as `json.RawMessage` with scan and value methods.

```go
type Train struct {
	Details types.RawJSON `db:"details"`
}
```

### types.Slice[T]

Slice based on `[]T` and null values stay as nil.

```go
type Train struct {
	Slice types.Slice[string] `db:"slice"`
}
```

### types.JSON[T]

For any type of json nullable value. Useful for struct values.

```go
type Details struct {
	Name string `json:"name"`
}

type Train struct {
	Details types.JSON[Details] `db:"details"`
}
```

### types.Null[T]

Wrapper of `sql.Null[T]` with additional json marshal and unmarshal methods.

### types.Time and types.Null[types.Time]

Wrapper of `time.Time` with additional json marshal and unmarshal methods with database scan and value methods.

---

### string, json.Number OR decimal.Decimal

> Use `decimal.Decimal` always for calculations.

To use `decimal.Decimal` type from `github.com/greatcloak/decimal` package.  
Use `json.Number` type for json number values and after that convert to `decimal.Decimal`.  
Use `string` type to direct get numeric values as string and convert to `decimal.Decimal`.

Use with nullable `sql.Null[decimal.Decimal]` package or pointer or `decimal.NullDecimal` type.

In struct use like this:

```go
type Train struct {
	ID          int64                     `db:"id"          goqu:"skipinsert"`
	Details     types.Map                 `db:"details"     json:"details,omitempty"`
	Additionals types.RawJSON             `db:"additionals" json:"additionals,omitempty"`
	Price       sql.Null[json.Number]     `db:"price"`
	LastPrice   decimal.NullDecimal       `db:"last_price"`
}
```

</details>

## Development

Go to `example` folder and run `make` command and fallow usage.
