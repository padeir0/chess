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
	case Profile:
		return "profile"
	case StopProfile:
		return "stopprofile"
	case SelfPlay:
		return "selfplay"
	case Compare:
		return "compare"
	case Championship:
		return "championship"
	case Test:
		return "test"
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

	Test

	Profile
	StopProfile
)
