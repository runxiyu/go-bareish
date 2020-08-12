# Examples

## Generate code
First, generate necessary structures using
```shell
go generate ./...
```

## Basic examples

1 `Person` per file

```shell
go run ./basic < customer.bin
go run ./basic < exmployee.bin
go run ./basic < terminated.bin
```

## Stream example

Multiple people in a single file or stream

```shell
go run ./stream < people.bin
```
