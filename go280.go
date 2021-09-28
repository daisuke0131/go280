package go280

import (
	"fmt"
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
	reportPanic(pass)
	reportNotRecover(pass)
	return nil, nil
}

func reportPanic(pass *analysis.Pass){
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	positives := map[*ssa.Function][]*ssa.Function{}
	for _, f := range s.SrcFuncs {
		if isPanicFunc(pass,f){
			if f.Object() != nil {
				pass.ExportObjectFact(f.Object(), new(isPanic))
			}
		} else {
			recordCallee(pass,f,positives)
		}
	}
	for k, v := range positives{
		if pass.ImportObjectFact(k.Object(),new(isPanic)){
			for _, v2 := range v {
				exportFact(pass, v2, positives)
			}
		}
	}
}

func reportNotRecover(pass *analysis.Pass){
	s := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA)
	for _, f := range s.SrcFuncs {
		if isRecover(pass,f){
			continue
		}
		for _, b := range f.Blocks {
			for _, instr := range b.Instrs {
				call, _ := instr.(*ssa.Call)
				if call == nil {
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
					pass.Reportf(call.Pos(),"panic")
				}
			}
		}
	}
}

func isRecover(pass *analysis.Pass,f *ssa.Function) bool{
	for _, b := range f.Blocks {
		for _, instr := range b.Instrs {
			call, _ := instr.(*ssa.Defer)
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
			for _, v := range callee.Blocks{
				for _, instr2 := range v.Instrs {
					call2, _ := instr2.(*ssa.Call)
					if call2 == nil{
						continue
					}
					fmt.Printf("%T\n",call2.Common().Value)
					callee2, _ := call2.Common().Value.(*ssa.Builtin)
					if callee2 == nil{
						continue
					}
					if callee2.Object() == nil{
						continue
					}
					if callee2.Name() == "recover"{
						return true
					}
				}
			}
		}
	}
	return false
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

func recordCallee(pass *analysis.Pass,f *ssa.Function, m map[*ssa.Function][]*ssa.Function){
	if f.Object() == nil{
		return
	}
	if pass.ImportObjectFact(f.Object(),new(isPanic)){
		return
	}
	for _, b := range f.Blocks {
		for _, instr := range b.Instrs {
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
			m[callee] = append(m[callee], f)
		}
	}
}

func exportFact(pass *analysis.Pass,f *ssa.Function,m map[*ssa.Function][]*ssa.Function){
	if pass.ImportObjectFact(f.Object(),new(isPanic)){
		return
	}

	pass.ExportObjectFact(f.Object(),new(isPanic))
	for _, v2 := range m[f] {
		exportFact(pass,v2,m)
	}
}

