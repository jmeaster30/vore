package bytecode

import (
	"fmt"

	"github.com/jmeaster30/vore/libvore/ast"
)

type Command interface {
	// execute(string, *files.Reader, ReplaceMode) Matches
	IsCommand()
	String() string
}

type FindCommand struct {
	All  bool
	Skip int
	Take int
	Last int
	Body []SearchInstruction
}

func (f FindCommand) IsCommand() {}

func (f FindCommand) String() string {
	return fmt.Sprintf("(find (all %t) (min %d max %d) (last %d) %s)", f.All, f.Skip, f.Take, f.Last, f.Body)
}

type ReplaceCommand struct {
	All      bool
	Skip     int
	Take     int
	Last     int
	Body     []SearchInstruction
	Replacer []ReplaceInstruction
}

func (r ReplaceCommand) IsCommand() {}

func (r ReplaceCommand) String() string {
	return fmt.Sprintf("(replace (%t %d %d %d))", r.All, r.Skip, r.Take, r.Last)
}

type SetCommand struct {
	Body SetCommandBody
	Id   string
}

func (s SetCommand) IsCommand() {}

func (s SetCommand) String() string {
	return fmt.Sprintf("(set (id %s) (%v)", s.Id, s.Body)
}

type SetCommandBody interface {
	IsSetCommandBody()
}

type SetCommandExpression struct {
	Instructions []SearchInstruction
	Validate     []ProcInstruction
}

func (s SetCommandExpression) IsSetCommandBody() {}

type SetCommandMatches struct {
	Command Command
}

func (s SetCommandMatches) IsSetCommandBody() {}

type SetCommandTransform struct {
	Instructions []ProcInstruction
}

func (s SetCommandTransform) IsSetCommandBody() {}

type SearchInstruction interface {
	// execute(*SearchEngineState) *SearchEngineState
	IsSearchInstruction()
	adjust(offset int, state *GenState) SearchInstruction
	String() string
}

type ReplaceInstruction interface {
	IsReplaceInstruction()
	// execute(*ReplacerState) *ReplacerState
}

type MatchLiteral struct {
	Not      bool
	ToFind   string
	Caseless bool
}

func (i MatchLiteral) IsSearchInstruction() {}

func (i MatchLiteral) String() string {
	return fmt.Sprintf("(literal (not %t) (caseless %t) '%s')", i.Not, i.Caseless, i.ToFind)
}

func (i MatchLiteral) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

type MatchCharClass struct {
	Not   bool
	Class ast.AstCharacterClassType
}

func (i MatchCharClass) IsSearchInstruction() {}

func (i MatchCharClass) String() string {
	return fmt.Sprintf("(class (not %t) %s)", i.Not, i.Class)
}

func (i MatchCharClass) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

type MatchVariable struct {
	Name string
}

func (i MatchVariable) IsSearchInstruction() {}

func (i MatchVariable) String() string {
	return fmt.Sprintf("(var '%s')", i.Name)
}

func (i MatchVariable) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

type MatchRange struct {
	Not  bool
	From string
	To   string
}

func (i MatchRange) IsSearchInstruction() {}

func (i MatchRange) String() string {
	return fmt.Sprintf("(range (not %t) (from '%s') (to '%s'))", i.Not, i.From, i.To)
}

func (i MatchRange) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

type CallSubroutine struct {
	Name string
	ToPC int
}

func (i CallSubroutine) IsSearchInstruction() {}

func (i CallSubroutine) String() string {
	return fmt.Sprintf("(call '%s' %d)", i.Name, i.ToPC)
}

func (i CallSubroutine) adjust(offset int, state *GenState) SearchInstruction {
	i.ToPC += offset
	return i
}

type Branch struct {
	Branches []int
}

func (i Branch) IsSearchInstruction() {}

func (i Branch) String() string {
	return fmt.Sprintf("(branch %v)", i.Branches)
}

func (i Branch) adjust(offset int, state *GenState) SearchInstruction {
	for idx := range i.Branches {
		i.Branches[idx] += offset
	}
	return i
}

type StartNotIn struct {
	NextCheckpointPC int
}

func (i StartNotIn) IsSearchInstruction() {}

func (i StartNotIn) String() string {
	return fmt.Sprintf("(startNotIn %d)", i.NextCheckpointPC)
}

func (i StartNotIn) adjust(offset int, state *GenState) SearchInstruction {
	i.NextCheckpointPC += offset
	return i
}

type FailNotIn struct{}

func (i FailNotIn) IsSearchInstruction() {}

func (i FailNotIn) String() string {
	return "(failNotIn)"
}

func (i FailNotIn) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

type EndNotIn struct {
	MaxSize int
}

func (i EndNotIn) IsSearchInstruction() {}

func (i EndNotIn) String() string {
	return fmt.Sprintf("(endNotIn %d)", i.MaxSize)
}

func (i EndNotIn) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

type StartLoop struct {
	Id       int64
	MinLoops int
	MaxLoops int
	Fewest   bool
	ExitLoop int
	Name     string
}

func (i StartLoop) IsSearchInstruction() {}

func (i StartLoop) String() string {
	return fmt.Sprintf("(startLoop '%s' (min %d max %d) (lazy %t) %d %d)", i.Name, i.MinLoops, i.MaxLoops, i.Fewest, i.Id, i.ExitLoop)
}

func (i StartLoop) adjust(offset int, state *GenState) SearchInstruction {
	i.ExitLoop += offset
	return i
}

type StopLoop struct {
	Id        int64
	MinLoops  int
	MaxLoops  int
	Fewest    bool
	StartLoop int
	Name      string
}

func (i StopLoop) IsSearchInstruction() {}

func (i StopLoop) String() string {
	return fmt.Sprintf("(stopLoop '%s' (min %d max %d) (lazy %t) %d %d)", i.Name, i.MinLoops, i.MaxLoops, i.Fewest, i.Id, i.StartLoop)
}

func (i StopLoop) adjust(offset int, state *GenState) SearchInstruction {
	i.StartLoop += offset
	return i
}

type StartVarDec struct {
	Name string
}

func (i StartVarDec) IsSearchInstruction() {}

func (i StartVarDec) String() string {
	return fmt.Sprintf("(startVarDec '%s')", i.Name)
}

func (i StartVarDec) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

type EndVarDec struct {
	Name string
}

func (i EndVarDec) IsSearchInstruction() {}

func (i EndVarDec) String() string {
	return fmt.Sprintf("(endVarDec '%s')", i.Name)
}

func (i EndVarDec) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

type StartSubroutine struct {
	Id        int
	Name      string
	EndOffset int
}

func (i StartSubroutine) IsSearchInstruction() {}

func (i StartSubroutine) String() string {
	return fmt.Sprintf("(startSub '%s' %d %d)", i.Name, i.Id, i.EndOffset)
}

func (i StartSubroutine) adjust(offset int, state *GenState) SearchInstruction {
	i.EndOffset += offset
	return i
}

type EndSubroutine struct {
	Name     string
	Validate []ProcInstruction
}

func (i EndSubroutine) IsSearchInstruction() {}

func (i EndSubroutine) String() string {
	return fmt.Sprintf("(endSub '%s')", i.Name)
}

func (i EndSubroutine) adjust(offset int, state *GenState) SearchInstruction {
	return i
}

type Jump struct {
	NewProgramCounter int
}

func (i Jump) IsSearchInstruction() {}
func (i Jump) isProcInstruction()   {}

func (i Jump) String() string {
	return fmt.Sprintf("(jump %d)", i.NewProgramCounter)
}

func (i Jump) adjust(offset int, state *GenState) SearchInstruction {
	i.NewProgramCounter += offset
	return i
}

type ReplaceString struct {
	Value string
}

func (i ReplaceString) IsReplaceInstruction() {}

type ReplaceVariable struct {
	Name string
}

func (i ReplaceVariable) IsReplaceInstruction() {}

type ReplaceProcess struct {
	Process []ProcInstruction
}

func (i ReplaceProcess) IsReplaceInstruction() {}

type ProcInstruction interface {
	isProcInstruction()
	// StackChange() (int, int)
}

type Store struct {
	VariableName string
}

func (s Store) isProcInstruction() {}

type Load struct {
	VariableName string
}

func (l Load) isProcInstruction() {}

type Push struct {
	Value Value
}

func (p Push) isProcInstruction() {}

type ConditionalJump struct {
	NewProgramCounter int
}

func (c ConditionalJump) isProcInstruction() {}

type Debug struct{}

func (d Debug) isProcInstruction() {}

type Return struct{}

func (r Return) isProcInstruction() {}

type Not struct{}

func (n Not) isProcInstruction() {}

type Head struct{}

func (h Head) isProcInstruction() {}

type Tail struct{}

func (t Tail) isProcInstruction() {}

type And struct{}

func (a And) isProcInstruction() {}

type Or struct{}

func (o Or) isProcInstruction() {}

type Add struct{}

func (a Add) isProcInstruction() {}

type Subtract struct{}

func (s Subtract) isProcInstruction() {}

type Multiply struct{}

func (m Multiply) isProcInstruction() {}

type Divide struct{}

func (d Divide) isProcInstruction() {}

type Modulo struct{}

func (m Modulo) isProcInstruction() {}

type Equal struct{}

func (e Equal) isProcInstruction() {}

type NotEqual struct{}

func (n NotEqual) isProcInstruction() {}

type GreaterThan struct{}

func (g GreaterThan) isProcInstruction() {}

type GreaterThanEqual struct{}

func (g GreaterThanEqual) isProcInstruction() {}

type LessThan struct{}

func (l LessThan) isProcInstruction() {}

type LessThanEqual struct{}

func (l LessThanEqual) isProcInstruction() {}
