package go280

import (
	"golang.org/x/tools/go/ssa"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

const doc = "go280 is ..."

// Analyzer is ...
var Analyzer = &analysis.Analyzer{
	Name: "go280",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		buildssa.Analyzer,
	},
	FactTypes:  []analysis.Fact{new(isPanic)},
}

type isPanic struct{
	analysis.Fact
}

func (*isPanic) String() string{
	return "isPanic"
}

func run(pass *analysis.Pass) (interface{}, error) {
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	positives := map[*ssa.Function]
	for _, f := range s.SrcFuncs {
		if isPanicFunc(pass,f){
			if f.Object() != nil {
				pass.ExportObjectFact(f.Object(), new(isPanic))
			}
		}
	}
	return nil, nil
}

func isPanicFunc(pass *analysis.Pass,f *ssa.Function) bool{
	for _, b := range f.Blocks {
		for _, instr := range b.Instrs {
			p, _ := instr.(*ssa.Panic)
			if p != nil {
				return true
			}

			call, _ := instr.(*ssa.Call)
			if call == nil{
				continue
			}
			callee, _ := call.Common().Value.(*ssa.Function)
			if callee == nil{
				continue
			}

			if callee.Object() == nil{
				continue
			}
			if pass.ImportObjectFact(callee.Object(),new(isPanic)){
				return true
			}
		}
	}
	return false
}
