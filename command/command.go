package command

import (
	ck "chess/command/commandkind"
	"chess/game"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Command struct {
	Kind     ck.CommandKind
	Operands []*Operand
}

func (this *Command) String() string {
	output := this.Kind.String() + " "
	for _, op := range this.Operands {
		output += op.String() + " "
	}
	return output
}

type Operand struct {
	Label    *string
	Number   *int64
	Position *game.Position
}

func (this *Operand) String() string {
	if this.IsLabel() {
		return *this.Label
	}
	if this.IsPosition() {
		return this.Position.String()
	}
	if this.IsNumber() {
		return strconv.FormatInt(*this.Number, 10)
	}
	return "???"
}

func (this *Operand) IsPosition() bool {
	return this.Position != nil
}
func (this *Operand) IsNumber() bool {
	return this.Number != nil
}
func (this *Operand) IsLabel() bool {
	return this.Label != nil
}

func Parse(cmdstr string) (*Command, *Error) {
	l := &lexer{
		Word:  nil,
		Input: cmdstr,
	}
	l.Next()
	cmd, err := parsecmd(l)
	if err != nil {
		return nil, err
	}
	err = checkCmd(cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

type textColumn int

func (this textColumn) String() string {
	return strconv.FormatInt(int64(this), 10)
}

type _range struct {
	Begin textColumn
	End   textColumn
}

func (this _range) String() string {
	if this.Begin >= this.End {
		return this.Begin.String()
	}
	return this.Begin.String() + " to " + this.End.String()
}

type lexemeKind int

func (this lexemeKind) String() string {
	switch this {
	case _label:
		return "Label"
	case _int:
		return "Int"
	case _pos:
		return "Pos"
	case _cmd:
		return "Cmd"
	case _EOF:
		return "EOF"
	}
	return "???"
}

const (
	invalidLexemeKind lexemeKind = iota

	_label
	_int
	_pos
	_cmd
	_EOF
)

type location struct {
	Input string
	Range _range
}

func (this location) Source() string {
	output := "    \u001b[36m"
	curr := textColumn(0)
	for _, r := range this.Input {
		if curr == this.Range.Begin {
			output += "\u001b[31m"
		}
		if curr == this.Range.End {
			output += "\u001b[36m"
		}
		output += string(r)
		curr++
	}
	output += "\u001b[0m"
	return output
}

type Error struct {
	message  string
	location *location
}

func (this *Error) String() string {
	if this.location == nil {
		return this.message
	}
	source := this.location.Source()
	message := this.location.Range.String() + " error: " + this.message
	if source != "" {
		return message + "\n" + source
	}
	return message
}

type lexeme struct {
	Kind lexemeKind

	CommandKind ck.CommandKind
	Text        string
	Value       int64
	Position    game.Position

	Range _range
}

func (this *lexeme) String() string {
	return this.Text + ":" + this.Kind.String()
}

type lexer struct {
	Word *lexeme

	Start, End   int
	LastRuneSize int
	Input        string
}

func (this *lexer) Next() *Error {
	symbol, err := any(this)
	if err != nil {
		return err
	}
	this.Start = this.End // this shouldn't be here
	this.Word = symbol
	return nil
}

func (this *lexer) Selected() string {
	return this.Input[this.Start:this.End]
}

func (this *lexer) Location() location {
	return location{
		Input: this.Input,
		Range: this.Range(),
	}
}

func (this *lexer) Range() _range {
	return _range{
		Begin: textColumn(this.Start),
		End:   textColumn(this.End),
	}
}

func (this *lexer) ReadAll() ([]*lexeme, *Error) {
	e := this.Next()
	if e != nil {
		return nil, e
	}
	output := []*lexeme{}
	for this.Word.Kind != _EOF {
		output = append(output, this.Word)
		e = this.Next()
		if e != nil {
			return nil, e
		}
	}
	return output, nil
}

func nextRune(l *lexer) rune {
	r, size := utf8.DecodeRuneInString(l.Input[l.End:])
	if r == utf8.RuneError && size == 1 {
		panic("Invalid UTF8 rune in string")
	}
	l.End += size
	l.LastRuneSize = size

	return r
}

func peekRune(l *lexer) rune {
	r, size := utf8.DecodeRuneInString(l.Input[l.End:])
	if r == utf8.RuneError && size == 1 {
		panic("Invalid UTF8 rune in string")
	}

	return r
}

/*ignore ignores the text previously read*/
func ignore(l *lexer) {
	l.Start = l.End
	l.LastRuneSize = 0
}

func acceptRun(l *lexer, s string) {
	r := peekRune(l)
	for strings.ContainsRune(s, r) {
		nextRune(l)
		r = peekRune(l)
	}
}

func acceptUntil(l *lexer, s string) {
	r := peekRune(l)
	for !strings.ContainsRune(s, r) {
		nextRune(l)
		r = peekRune(l)
	}
}

const (
	/*eof is equivalent to RuneError, but in this package it only shows up in EoFs
	If the rune is invalid, it panics instead*/
	eof rune = utf8.RuneError
)

const (
	digits  = "0123456789"
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
)

func isNumber(r rune) bool {
	return strings.ContainsRune(digits, r)
}

func isLetter(r rune) bool {
	return strings.ContainsRune(letters, r)
}

func ignoreWhitespace(st *lexer) {
	r := peekRune(st)
loop:
	for {
		switch r {
		case ' ', '\t', '\n':
			nextRune(st)
		default:
			break loop
		}
		r = peekRune(st)
	}
	ignore(st)
}

func any(st *lexer) (*lexeme, *Error) {
	var r rune

	ignoreWhitespace(st)

	r = peekRune(st)

	if isNumber(r) {
		return number(st), nil
	}
	if isLetter(r) {
		nextRune(st)
		r2 := peekRune(st)
		if isNumber(r2) {
			nextRune(st)
			return position(st, r, r2)
		}
		return identifier(st), nil
	}
	if r == eof {
		nextRune(st)
		return &lexeme{Kind: _EOF}, nil
	}

	nextRune(st)
	return nil, lexError(st, "invalid character: '"+string(r)+"'")
}

func lexError(st *lexer, message string) *Error {
	loc := st.Location()
	return &Error{
		message:  message,
		location: &loc,
	}
}

func number(st *lexer) *lexeme {
	acceptRun(st, digits)
	value := parseNormal(st.Selected())
	return &lexeme{
		Text:  st.Selected(),
		Kind:  _int,
		Value: value,
		Range: st.Range(),
	}
}

func parseNormal(text string) int64 {
	var output int64 = 0
	for i := range text {
		output *= 10
		char := text[i]
		if char >= '0' || char <= '9' {
			output += int64(char - '0')
		} else {
			panic(text)
		}
	}
	return output
}

func position(st *lexer, col, row rune) (*lexeme, *Error) {
	if col >= 'a' && col <= 'h' &&
		row >= '1' && row <= '8' {
		pos := game.Position{
			Column: int(col - 'a'),
			Row:    7 - int(row-'1'),
		}
		return &lexeme{
			Kind:     _pos,
			Text:     st.Selected(),
			Position: pos,
			Range:    st.Range(),
		}, nil
	}

	if col >= 'A' && col <= 'H' &&
		row >= '1' && row <= '8' {
		pos := game.Position{
			Column: int(col - 'A'),
			Row:    7 - int(row-'1'),
		}
		return &lexeme{
			Kind:     _pos,
			Text:     st.Selected(),
			Position: pos,
			Range:    st.Range(),
		}, nil
	}
	return nil, lexError(st, "invalid position: "+st.Selected())
}

func identifier(st *lexer) *lexeme {
	acceptRun(st, letters)
	selected := st.Selected()
	tp := _label
	cmdKind := ck.InvalidCommandKind

	switch selected {
	case "move":
		tp = _cmd
		cmdKind = ck.Move
	case "save":
		tp = _cmd
		cmdKind = ck.Save
	case "restore":
		tp = _cmd
		cmdKind = ck.Restore
	case "show":
		tp = _cmd
		cmdKind = ck.Show
	case "quit", "exit":
		tp = _cmd
		cmdKind = ck.Quit
	case "clear":
		tp = _cmd
		cmdKind = ck.Clear
	case "profile":
		tp = _cmd
		cmdKind = ck.Profile
	case "stopprofile":
		tp = _cmd
		cmdKind = ck.StopProfile
	case "selfplay":
		tp = _cmd
		cmdKind = ck.SelfPlay
	case "compare":
		tp = _cmd
		cmdKind = ck.Compare
	case "no", "NO":
		tp = _cmd
		cmdKind = ck.NO
	}

	return &lexeme{
		Range:       st.Range(),
		Text:        st.Selected(),
		CommandKind: cmdKind,
		Kind:        tp,
	}
}

func parsecmd(st *lexer) (*Command, *Error) {
	word, err := consume(st)
	if err != nil {
		return nil, err
	}
	if word == nil {
		panic("word is nil")
	}
	if word.Kind != _cmd {
		return nil, parseError(st, word, "invalid command")
	}
	ops, err := multipledata(st)
	if err != nil {
		return nil, err
	}
	return &Command{Kind: word.CommandKind, Operands: ops}, nil
}

func multipledata(st *lexer) ([]*Operand, *Error) {
	out := []*Operand{}
	n, err := data(st)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, nil
	}
	for n != nil {
		out = append(out, n)
		n, err = data(st)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func data(st *lexer) (*Operand, *Error) {
	word, err := consume(st)
	if err != nil {
		return nil, err
	}
	switch word.Kind {
	case _label:
		return &Operand{Label: &word.Text}, nil
	case _int:
		return &Operand{Number: &word.Value}, nil
	case _pos:
		return &Operand{Position: &word.Position}, nil
	}
	return nil, nil
}

func parseError(st *lexer, w *lexeme, message string) *Error {
	return &Error{
		message: message,
		location: &location{
			Input: st.Input,
			Range: w.Range,
		},
	}
}

func consume(st *lexer) (*lexeme, *Error) {
	n := st.Word
	err := st.Next()
	return n, err
}

func checkCmd(cmd *Command) *Error {
	switch cmd.Kind {
	case ck.Save:
		return checkCmdSave(cmd)
	case ck.Restore:
		return checkCmdRestore(cmd)
	case ck.Show:
		return checkShow(cmd)
	case ck.Move:
		return checkMove(cmd)
	case ck.Profile:
		return checkCmdProfile(cmd)
	case ck.Compare:
		return checkCmdCompare(cmd)
	case ck.Quit, ck.Clear, ck.NO, ck.StopProfile, ck.SelfPlay:
		return nil
	}
	panic("invalid command")
}

func checkMove(cmd *Command) *Error {
	if len(cmd.Operands) == 2 {
		if cmd.Operands[0].IsPosition() &&
			cmd.Operands[1].IsPosition() {
			return nil
		}
	}
	return checkErr("move <pos> <pos>")
}

func checkCmdCompare(cmd *Command) *Error {
	if len(cmd.Operands) == 2 &&
		cmd.Operands[0].IsLabel() &&
		cmd.Operands[1].IsLabel() {
		return nil
	}
	return checkErr(cmd.Kind.String() + " <label> <label>")
}

func checkCmdSave(cmd *Command) *Error {
	if len(cmd.Operands) == 1 && cmd.Operands[0].IsLabel() {
		return nil
	}
	return checkErr(cmd.Kind.String() + " <label>")
}

func checkCmdProfile(cmd *Command) *Error {
	if len(cmd.Operands) == 1 && cmd.Operands[0].IsLabel() {
		return nil
	}
	return checkErr(cmd.Kind.String() + " <label>")
}

func checkCmdRestore(cmd *Command) *Error {
	if len(cmd.Operands) == 0 {
		return nil
	}
	if len(cmd.Operands) == 1 && cmd.Operands[0].IsLabel() {
		return nil
	}
	return checkErr(cmd.Kind.String() + " <label>")
}

func checkShow(cmd *Command) *Error {
	if len(cmd.Operands) == 0 {
		return nil
	}
	if len(cmd.Operands) > 1 ||
		!cmd.Operands[0].IsLabel() ||
		!isValidShow(*cmd.Operands[0].Label) {
		return checkErr("show [moves]")
	}
	return nil
}

func checkErr(layout string) *Error {
	return &Error{message: "invalid operands, expected: " + layout}
}

func isValidShow(s string) bool {
	switch s {
	case "moves":
		return true
	}
	return false
}
