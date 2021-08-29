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

func (l MyLargeStruct) NotOkay() bool { // want "large struct MyLargeStruct passed as value to function receiver"
	return *l.E
}

func (l MyLargeStruct) NotOkayAlso() string { // want "large struct MyLargeStruct passed as value to function receiver"
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

func (m MyOtherLargeStruct) NotOkay() int { // want "large struct MyOtherLargeStruct passed as value to function receiver"
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
