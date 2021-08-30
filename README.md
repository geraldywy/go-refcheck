# go-refcheck

    go get github.com/geraldywy/go-refcheck

`go-refcheck` reports function receivers where use of large structs are passed by value.

## Why?

Passing large structs as function receivers generates a copy in memory to work off of. This can create potential issues with the progam's memory heap.

Also, see [`copyfighter`](https://github.com/jmhodges/copyfighter) for a similar linter.

## Example
    type MyLargeStruct struct {
        A string
        B int64
        C float64
        D string
        E *bool
    }

    func (l *MyLargeStruct) Okay() bool {
        return *l.E
    }

    func (l MyLargeStruct) NotOkay() bool {
        return *l.E
    }

The function `NotOkay` is flagged as follows:
    
    go-refcheck ./...

    $ .../testdata/src/p/p.go:28:9: large struct "MyLargeStruct" passed by value to function receiver

More examples can be found in the [testdata](https://github.com/geraldywy/go-refcheck/blob/master/testdata/src/p/p.go) folder.

## Flags
By default, large structs are assumed to be larger than 32 bytes, this value can be toggled with a `max` flag.

Eg: To only flag for large structs larger than 64 bytes:

`go-refcheck -max 64 ./...`
## False positives
Kindly open an issue for such cases.

## Contributions
Yes please, create an issue/PR for any suggestions.

