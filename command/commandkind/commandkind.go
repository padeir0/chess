package commandkind

type CommandKind int

func (this CommandKind) String() string {
	switch this {
	case Next:
		return "next"
	case Move:
		return "move"
	case Undo:
		return "undo"
	case Save:
		return "save"
	case Restore:
		return "restore"
	case Show:
		return "show"
	case Quit:
		return "quit"
	case Clear:
		return "clear"
	}
	return "???"
}

const (
	InvalidCommandKind CommandKind = iota
	Next
	Move
	Undo
	Save
	Restore
	Show
	Quit
	Clear
	NO
)
