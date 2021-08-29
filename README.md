# go-refcheck

    go get github.com/geraldywy/go-refcheck

`go-refcheck` reports function receivers where use of large structs are passed by value. 

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

    $ .../testdata/src/p/p.go:28:9: large struct MyLargeStruct passed as value to function receiver

More examples can be found in the [testdata](https://github.com/geraldywy/go-refcheck/blob/master/testdata/src/p/p.go) folder.


## Why?

Q: I'm still a little confused about when to choose value or pointer
receivers. Can you provide any concrete/real-world examples of when we
would choose one over the other?

A: When you want to modify the state of the receiver, you have to use
pointer receivers.  If the struct is very big, you probably want to
use a pointer receiver because value receivers operate on a copy.  If
neither applies, you can use a value receiver.  However, be careful
with value receivers; e.g., if you have a mutex in a struct, you
cannot make it a value receiver, because the mutex would be copied,
defeating its purpose.

[Source](https://pdos.csail.mit.edu/6.824/papers/tour-faq.txt)

Also, see [`copyfighter`](https://github.com/jmhodges/copyfighter) for a similar linter.

## Flags
By default, large structs are assumed to be larger than 32 bytes, this value can be toggled with a `max` flag.

Eg: To only flag for large structs larger than 64 bytes:

`go-refcheck -max 64 ./...`
## False positives
Kindly open an issue for such cases.

## Contributions
Yes please, create an issue/PR for any suggestions.

