package commandkind

type CommandKind int

func (this CommandKind) String() string {
	switch this {
	case Move:
		return "move"
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
	case Compare:
		return "compare"
	case Profile:
		return "profile"
	case StopProfile:
		return "stopprofile"
	}
	return "???"
}

const (
	InvalidCommandKind CommandKind = iota
	Move
	Save
	Restore
	Show
	Quit
	Clear
	NO
	SelfPlay
	Compare
	Championship

	Profile
	StopProfile
)
