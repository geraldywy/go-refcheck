package analyzer

import (
	"flag"
	"go/ast"
	"go/build"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var flagSet flag.FlagSet

var maxStructWidth = flag.Int64("max", 32, "maximum size in bytes a struct can be before by-value uses are flagged, accounts for padding")

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "gorefcheck",
		Doc:      "reports function receivers where use of large structs are passed by value",
		Run:      run,
		Flags:    flagSet,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.StarExpr)(nil),
		(*ast.TypeSpec)(nil),
	}

	// maps to process after
	starReceivers := make(map[token.Pos]bool)
	receivers := make(map[token.Pos]string)
	largeStructsCache := make(map[string]bool)

	for range pass.Files {
		inspector.Preorder(nodeFilter, func(node ast.Node) {
			switch node := node.(type) {
			case *ast.FuncDecl:
				list := node.Recv
				if list == nil || len(list.List) < 1 {
					return
				}
				rec := list.List[0]
				st, ok := pass.TypesInfo.Types[rec.Type]
				if !ok {
					return
				}
				if arr := strings.Split(st.Type.String(), "."); len(arr) != 0 {
					receivers[rec.Type.Pos()] = arr[len(arr)-1]
				}
			case *ast.StarExpr:
				starReceivers[node.Pos()] = true
			case *ast.TypeSpec:
				checkWideStruct(pass, node, *maxStructWidth, largeStructsCache)
			}
		})
	}
	for pos, structName := range receivers {
		if starReceivers[pos] {
			continue
		}
		isLargeStruct := largeStructsCache[structName]
		if isLargeStruct {
			pass.Reportf(pos, "large struct %s passed as value to function receiver", structName)
		}
	}
	return nil, nil
}

// returns whether a node is a large struct
func checkWideStruct(pass *analysis.Pass, node *ast.TypeSpec, maxSize int64, largeStructsCache map[string]bool) {
	structName := node.Name.Name
	if _, ok := largeStructsCache[structName]; ok { // check cache
		return
	}
	st, ok := node.Type.(*ast.StructType)
	if ok {
		wordSize := int64(8)
		maxAlign := int64(8)
		switch build.Default.GOARCH {
		case "386", "arm":
			wordSize, maxAlign = 4, 4
		case "amd64p32":
			wordSize = 4
		}
		str := pass.TypesInfo.Types[st].Type.(*types.Struct)
		s := gcSizes{wordSize, maxAlign}
		sz := s.Sizeof(str)
		if sz > maxSize {
			largeStructsCache[structName] = true
			return
		}
	}
	largeStructsCache[structName] = false
}

// Code below based on go/types.StdSizes.

type gcSizes struct {
	WordSize int64
	MaxAlign int64
}

func (s *gcSizes) Alignof(T types.Type) int64 {
	// NOTE: On amd64, complex64 is 8 byte aligned,
	// even though float32 is only 4 byte aligned.

	// For arrays and structs, alignment is defined in terms
	// of alignment of the elements and fields, respectively.
	switch t := T.Underlying().(type) {
	case *types.Array:
		// spec: "For a variable x of array type: unsafe.Alignof(x)
		// is the same as unsafe.Alignof(x[0]), but at least 1."
		return s.Alignof(t.Elem())
	case *types.Struct:
		// spec: "For a variable x of struct type: unsafe.Alignof(x)
		// is the largest of the values unsafe.Alignof(x.f) for each
		// field f of x, but at least 1."
		max := int64(1)
		for i, nf := 0, t.NumFields(); i < nf; i++ {
			if a := s.Alignof(t.Field(i).Type()); a > max {
				max = a
			}
		}
		return max
	}
	a := s.Sizeof(T) // may be 0
	// spec: "For a variable x of any type: unsafe.Alignof(x) is at least 1."
	if a < 1 {
		return 1
	}
	if a > s.MaxAlign {
		return s.MaxAlign
	}
	return a
}

var basicSizes = [...]byte{
	types.Bool:       1,
	types.Int8:       1,
	types.Int16:      2,
	types.Int32:      4,
	types.Int64:      8,
	types.Uint8:      1,
	types.Uint16:     2,
	types.Uint32:     4,
	types.Uint64:     8,
	types.Float32:    4,
	types.Float64:    8,
	types.Complex64:  8,
	types.Complex128: 16,
}

func (s *gcSizes) Sizeof(T types.Type) int64 {
	switch t := T.Underlying().(type) {
	case *types.Basic:
		k := t.Kind()
		if int(k) < len(basicSizes) {
			if s := basicSizes[k]; s > 0 {
				return int64(s)
			}
		}
		if k == types.String {
			return s.WordSize * 2
		}
	case *types.Array:
		n := t.Len()
		if n == 0 {
			return 0
		}
		a := s.Alignof(t.Elem())
		z := s.Sizeof(t.Elem())
		return align(z, a)*(n-1) + z
	case *types.Slice:
		return s.WordSize * 3
	case *types.Struct:
		nf := t.NumFields()
		if nf == 0 {
			return 0
		}

		var o int64
		max := int64(1)
		for i := 0; i < nf; i++ {
			ft := t.Field(i).Type()
			a, sz := s.Alignof(ft), s.Sizeof(ft)
			if a > max {
				max = a
			}
			if i == nf-1 && sz == 0 && o != 0 {
				sz = 1
			}
			o = align(o, a) + sz
		}
		return align(o, max)
	case *types.Interface:
		return s.WordSize * 2
	}
	return s.WordSize // catch-all
}

// align returns the smallest y >= x such that y % a == 0.
func align(x, a int64) int64 {
	y := x + a - 1
	return y - y%a
}
