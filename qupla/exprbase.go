package qupla

type QuplaExprBase struct {
	source string
}

func NewQuplaExprBase(source string) QuplaExprBase {
	return QuplaExprBase{
		source: source,
	}
}

func (e *QuplaExprBase) GetSource() string {
	return e.source
}
