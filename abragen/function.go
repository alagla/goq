package abragen

import (
	"github.com/lunfardo314/goq/abra"
	"github.com/lunfardo314/goq/qupla"
)

func GenAbraBranch(function *qupla.Function, branch *abra.Branch, codeUnit *abra.CodeUnit) {
	for _, vi := range function.Sites {
		if !vi.IsParam {
			continue
		}
		branch.AddInputSite(vi.Size)
	}
}
