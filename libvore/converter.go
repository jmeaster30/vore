package libvore

import (
	"fmt"
	"strconv"
)

type NodeName = uint32

const (
	NodeFindCommand NodeName = iota
	NodeReplaceCommand
	NodeSetCommand
	NodeSetCommandExpression
	NodeSetCommandMatches
	NodeSetCommandTransform
	NodeMatchLiteral
	NodeMatchCharClass
	NodeMatchVariable
	NodeMatchRange
	NodeCallSubroutine
	NodeBranch
	NodeStartNotIn
	NodeFailNotIn
	NodeEndNotIn
	NodeStartLoop
	NodeStopLoop
	NodeStartVarDec
	NodeEndVarDec
	NodeStartSubroutine
	NodeEndSubroutine
	NodeJump
	NodeReplaceString
	NodeReplaceVariable
	NodeReplaceProcess
)

func ToMap(commands []Command) []map[string]any {
	var results []map[string]any
	for _, c := range commands {
		results = append(results, c.ToMap())
	}
	return results
}

func CommandsFromMap(commands []map[string]any) ([]Command, error) {
	var results []Command
	for _, c := range commands {
		com, err := CommandFromMap(c)
		if err != nil {
			return nil, err
		}
		results = append(results, com)
	}
	return results, nil
}

func CommandFromMap(command map[string]any) (Command, error) {
	nodeTypeInt, err := strconv.ParseInt(command["node"].(string), 10, 32)
	if err != nil {
		return nil, err
	}
	nodeType := NodeName(nodeTypeInt)
	switch nodeType {
	case NodeFindCommand:
		return FindCommandFromMap(command)
	case NodeReplaceCommand:
	case NodeSetCommand:
	default:
	}
	//log.Panicf("'%v' is not a Command node", nodeType)
	return nil, fmt.Errorf("'%v' is not a Command node", nodeType)
}

func SearchInstructionsFromMap(searchInstructions []any) ([]SearchInstruction, error) {
	var results []SearchInstruction
	for _, c := range searchInstructions {
		instruction, ok := c.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("can't assert type %+v to map[string]any", c)
		}
		inst, err := SearchInstructionFromMap(instruction)
		if err != nil {
			return nil, err
		}
		results = append(results, inst)
	}
	return results, nil
}

func SearchInstructionFromMap(command map[string]any) (SearchInstruction, error) {
	nodeTypeInt, err := strconv.ParseInt(command["node"].(string), 10, 32)
	if err != nil {
		return nil, err
	}
	nodeType := NodeName(nodeTypeInt)
	switch nodeType {
	case NodeMatchLiteral:
		return MatchLiteralFromMap(command)
	case NodeMatchCharClass:
		return MatchCharClassFromMap(command)
	case NodeMatchVariable:
	case NodeMatchRange:
	case NodeCallSubroutine:
	case NodeBranch:
	case NodeStartNotIn:
	case NodeFailNotIn:
	case NodeEndNotIn:
	case NodeStartLoop:
	case NodeStopLoop:
	case NodeStartVarDec:
	case NodeEndVarDec:
	case NodeStartSubroutine:
	case NodeEndSubroutine:
	case NodeJump:
	default:
	}
	//log.Panicf("'%v' is not a SearchInstruction node", nodeType)
	return nil, fmt.Errorf("'%v' is not a SearchInstruction node", nodeType)
}

func (c FindCommand) ToMap() map[string]any {
	var body []map[string]any
	for _, b := range c.body {
		body = append(body, b.ToMap())
	}

	return map[string]any{
		"node": strconv.FormatInt(int64(NodeFindCommand), 10),
		"all":  strconv.FormatBool(c.all),
		"skip": strconv.FormatInt(int64(c.skip), 10),
		"take": strconv.FormatInt(int64(c.take), 10),
		"last": strconv.FormatInt(int64(c.last), 10),
		"body": body,
	}
}

func FindCommandFromMap(obj map[string]any) (*FindCommand, error) {
	all, err := strconv.ParseBool(obj["all"].(string))
	if err != nil {
		return nil, err
	}

	skip, err := strconv.ParseInt(obj["skip"].(string), 10, 32)
	if err != nil {
		return nil, err
	}

	take, err := strconv.ParseInt(obj["take"].(string), 10, 32)
	if err != nil {
		return nil, err
	}

	last, err := strconv.ParseInt(obj["last"].(string), 10, 32)
	if err != nil {
		return nil, err
	}

	insts, ok := obj["body"].([]any)
	if !ok {
		return nil, fmt.Errorf("bad type assertion: %+v to []map[string]any", obj["body"])
	}
	//return nil, fmt.Errorf("insts asserted %+v to any", insts)
	body, err := SearchInstructionsFromMap(insts)
	if err != nil {
		return nil, err
	}

	return &FindCommand{
		all:  all,
		skip: int(skip),
		take: int(take),
		last: int(last),
		body: body,
	}, nil
}

func (c ReplaceCommand) ToMap() map[string]any {
	var body []map[string]any
	for _, b := range c.body {
		body = append(body, b.ToMap())
	}

	var replacer []map[string]any
	for _, c := range c.replacer {
		replacer = append(replacer, c.ToMap())
	}

	return map[string]any{
		"node":     NodeReplaceCommand,
		"all":      c.all,
		"skip":     c.skip,
		"take":     c.take,
		"last":     c.last,
		"body":     body,
		"replacer": replacer,
	}
}

func (c SetCommand) ToMap() map[string]any {
	return map[string]any{
		"node": NodeSetCommand,
		"id":   c.id,
		"body": c.body.ToMap(),
	}
}

func (s SetCommandExpression) ToMap() map[string]any {
	var instructions []map[string]any
	for _, b := range s.instructions {
		instructions = append(instructions, b.ToMap())
	}
	return map[string]any{
		"node":         NodeSetCommandExpression,
		"instructions": instructions,
	}
}

func (s SetCommandMatches) ToMap() map[string]any {
	return map[string]any{
		"node":    NodeSetCommandMatches,
		"command": s.command.ToMap(),
	}
}

func (s SetCommandTransform) ToMap() map[string]any {
	return map[string]any{
		"node":       NodeSetCommandTransform,
		"statements": []map[string]any{},
	}
}

func (i MatchLiteral) ToMap() map[string]any {
	return map[string]any{
		"node":     strconv.FormatInt(int64(NodeMatchLiteral), 10),
		"not":      strconv.FormatBool(i.not),
		"toFind":   i.toFind,
		"caseless": strconv.FormatBool(i.caseless),
	}
}

func MatchLiteralFromMap(obj map[string]any) (*MatchLiteral, error) {
	notFlag, err := strconv.ParseBool(obj["not"].(string))
	if err != nil {
		return nil, err
	}

	toFind := obj["toFind"].(string)

	caseless, err := strconv.ParseBool(obj["caseless"].(string))
	if err != nil {
		return nil, err
	}

	return &MatchLiteral{
		not:      notFlag,
		toFind:   toFind,
		caseless: caseless,
	}, nil
}

func (i MatchCharClass) ToMap() map[string]any {
	return map[string]any{
		"node":  NodeMatchCharClass,
		"not":   i.not,
		"class": i.class,
	}
}

func MatchCharClassFromMap(obj map[string]any) (*MatchCharClass, error) {
	notFlag, err := strconv.ParseBool(obj["not"].(string))
	if err != nil {
		return nil, err
	}

	class, err := strconv.ParseInt(obj["class"].(string), 10, 32)
	if err != nil {
		return nil, err
	}

	return &MatchCharClass{
		not:   notFlag,
		class: AstCharacterClassType(class),
	}, nil
}

func (i MatchVariable) ToMap() map[string]any {
	return map[string]any{
		"node": NodeMatchVariable,
		"name": i.name,
	}
}

func (i MatchRange) ToMap() map[string]any {
	return map[string]any{
		"node": NodeMatchRange,
		"not":  i.not,
		"from": i.from,
		"to":   i.to,
	}
}

func (i CallSubroutine) ToMap() map[string]any {
	return map[string]any{
		"node": NodeCallSubroutine,
		"name": i.name,
		"toPC": i.toPC,
	}
}

func (i Branch) ToMap() map[string]any {
	return map[string]any{
		"node":     NodeBranch,
		"branches": i.branches,
	}
}

func (i StartNotIn) ToMap() map[string]any {
	return map[string]any{
		"node":             NodeStartNotIn,
		"nextCheckpointPC": i.nextCheckpointPC,
	}
}

func (i FailNotIn) ToMap() map[string]any {
	return map[string]any{
		"node": NodeFailNotIn,
	}
}

func (i EndNotIn) ToMap() map[string]any {
	return map[string]any{
		"node":    NodeEndNotIn,
		"maxSize": i.maxSize,
	}
}

func (i StartLoop) ToMap() map[string]any {
	return map[string]any{
		"node":     NodeStartLoop,
		"id":       i.id,
		"minLoops": i.minLoops,
		"maxLoops": i.maxLoops,
		"fewest":   i.fewest,
		"exitLoop": i.exitLoop,
		"name":     i.name,
	}
}

func (i StopLoop) ToMap() map[string]any {
	return map[string]any{
		"node":      NodeStopLoop,
		"id":        i.id,
		"minLoops":  i.minLoops,
		"maxLoops":  i.maxLoops,
		"fewest":    i.fewest,
		"startLoop": i.startLoop,
		"name":      i.name,
	}
}

func (i StartVarDec) ToMap() map[string]any {
	return map[string]any{
		"node": NodeStartVarDec,
		"name": i.name,
	}
}

func (i EndVarDec) ToMap() map[string]any {
	return map[string]any{
		"node": NodeEndVarDec,
		"name": i.name,
	}
}

func (i StartSubroutine) ToMap() map[string]any {
	return map[string]any{
		"node":      NodeStartSubroutine,
		"id":        i.id,
		"name":      i.name,
		"endOffset": i.endOffset,
	}
}

func (i EndSubroutine) ToMap() map[string]any {
	return map[string]any{
		"node":     NodeEndSubroutine,
		"name":     i.name,
		"validate": []map[string]any{},
	}
}

func (i Jump) ToMap() map[string]any {
	return map[string]any{
		"node":              NodeJump,
		"newProgramCounter": i.newProgramCounter,
	}
}

func (i ReplaceString) ToMap() map[string]any {
	return map[string]any{
		"node":  NodeReplaceString,
		"value": i.value,
	}
}

func (i ReplaceVariable) ToMap() map[string]any {
	return map[string]any{
		"node": NodeReplaceVariable,
		"name": i.name,
	}
}

func (i ReplaceProcess) ToMap() map[string]any {
	return map[string]any{
		"node":    NodeReplaceProcess,
		"process": []map[string]any{},
	}
}
