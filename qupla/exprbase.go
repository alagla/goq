package qupla

type QuplaExprBase struct {
	source   string
	hasState bool
}

func NewQuplaExprBase(source string) QuplaExprBase {
	return QuplaExprBase{
		source: source,
	}
}

func (e *QuplaExprBase) GetSource() string {
	return e.source
}

func (e *QuplaExprBase) HasState() bool {
	return e.hasState
}
