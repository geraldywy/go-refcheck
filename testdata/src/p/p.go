package p

type SmallStruct struct {
	A bool
	B int
}

func (s SmallStruct) Okay() bool {
	return s.A
}

func (s *SmallStruct) OkayAlso() bool {
	return s.A
}

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

func (l MyLargeStruct) NotOkay() bool { // want "function with receiver passed by value to large struct MyLargeStruct"
	return *l.E
}

func (l MyLargeStruct) NotOkayAlso() string { // want "function with receiver passed by value to large struct MyLargeStruct"
	l.A = "assigning to a copy"
	return l.A
}

func (l *MyLargeStruct) OkayAlso() string {
	l.A = "assigning to a ref"
	return l.A
}

type MyOtherLargeStruct struct {
	InnerLargeStruct MyLargeStruct
	A                int
}

func (m MyOtherLargeStruct) NotOkay() int { // want "function with receiver passed by value to large struct MyOtherLargeStruct"
	return m.A
}

func (m *MyOtherLargeStruct) Okay() int {
	return m.A
}

type NotSoLargeStruct struct {
	A *MyLargeStruct
	B *MyOtherLargeStruct
	C string
}

func (n NotSoLargeStruct) Okay() string {
	return n.C
}

func (n *NotSoLargeStruct) OkayAlso() string {
	return n.C
}
