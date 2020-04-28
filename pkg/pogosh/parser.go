// Copyright 2020 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pogosh

type tokenizer struct {
	ts []token
}

// The remainder of this file parses and evaluates an LL(1) grammar using a
// predictive parse. This grammar was found in the POSIX.1-2017 spec and
// converted to LL(1).
// TODO: explain the Backus-Naur
// When a parse function is called, can assume FIRST conditions are met.

// Program ::= LineBreak Program2
// Program2 ::= CompleteCommands |
// CompleteCommands ::= CompleteCommand CompleteCommands2
// CompleteCommands2 ::= NewLineList CompleteCommands3 |
// CompleteCommands3 ::= CompleteCommands |
func parseProgram(s *State, t *tokenizer) command {
	cmd := compoundList{}

	parseLineBreak(s, t)
	for t.ts[0].ttype != ttEOF {
		cmd.cmds = append(cmd.cmds, parseCompleteCommand(s, t))
		switch t.ts[0].ttype {
		case ttNewLine:
			parseLineBreak(s, t)
		case ttEOF:
		default:
			panic("Parse error")
		}
		parseLineBreak(s, t)
	}
	return &cmd
}

// CompleteCommand ::= List CompleteCommand2
// CompleteCommand2 ::= SeparatorOp |
// List ::= AndOr List2
// List2 ::= SeparatorOp List |
func parseCompleteCommand(s *State, t *tokenizer) command {
	cmd := compoundList{}

	for {
		cmd.cmds = append(cmd.cmds, parseAndOr(s, t))
		if t.ts[0].ttype == ttNewLine || t.ts[0].ttype == ttEOF {
			break
		}
		switch t.ts[0].value {
		case "&":
			parseSeparatorOp(s, t)
			cmd.cmds[len(cmd.cmds)-1] = &async{cmd.cmds[len(cmd.cmds)-1]}
		case ";":
			parseSeparatorOp(s, t)
		}
		if t.ts[0].ttype == ttNewLine || t.ts[0].ttype == ttEOF {
			break
		}
	}
	return &cmd
}

// AndOr ::= Pipeline AndOr2
// AndOr2 ::= '&&' LineBreak AndOr | '||' LineBreak AndOr |
func parseAndOr(s *State, t *tokenizer) command {
	cmd := parsePipeline(s, t)

	for {
		switch t.ts[0].ttype {
		case ttAndIf: // FIRST['&&' LineBreak AndOr]
			t.ts = t.ts[1:]
			parseLineBreak(s, t)
			cmd = &and{cmd, parsePipeline(s, t)}
		case ttOrIf: // FIRST['||' LineBreak AndOr]
			t.ts = t.ts[1:]
			parseLineBreak(s, t)
			cmd = &or{cmd, parsePipeline(s, t)}
		default: // TODO: FOLLOW[AndOr2]
			return cmd
		}
	}
}

// Pipeline ::= '!' PipeSequence | PipeSequence
func parsePipeline(s *State, t *tokenizer) command {
	switch t.ts[0] {
	case token{"!", ttWord}:
		t.ts = t.ts[1:]
		return &not{parsePipeSequence(s, t)}
	default:
		return parsePipeSequence(s, t)
	}
}

// PipeSequence ::= Command PipeSequence2
// PipeSequence2 ::= '|' LineBreak PipeSequence |
func parsePipeSequence(s *State, t *tokenizer) command {
	cmd := pipeline{}
	cmd.cmds = append(cmd.cmds, parseCommand(s, t))

	switch t.ts[0] {
	case token{"|", ttWord}:
		t.ts = t.ts[1:]
		parseLineBreak(s, t)
		// TODO: possibly wrong associativity
		cmd.cmds = append(cmd.cmds, parsePipeSequence(s, t))
	default: // TODO: FOLLOW[PipeSequence]
	}
	return &cmd
}

// TODO: make LL(0)
// Command ::= SimpleCommand | CompoundCommand | CompoundCommand RedirectList | FunctionDefinition
func parseCommand(s *State, t *tokenizer) command {
	// TODO: support more than simple command
	return parseSimpleCommand(s, t)
}

// compound_command  : brace_group
//                   | subshell
//                   | for_clause
//                   | case_clause
//                   | if_clause
//                   | while_clause
//                   | until_clause
//                   ;
func parseCompoundCommand(s *State, t *tokenizer) {

}

// subshell          : '(' compound_list ')
//                   ;
func parseSubshell(s *State, t *tokenizer) {

}

// compound_list     : LineBreak term
//                   | LineBreak term separator
//                   ;
func parseCompoundList(s *State, t *tokenizer) {

}

// term              : term separator '&&'
//                   | '&&'
//                   ;
func parseTerm(s *State, t *tokenizer) {

}

// for_clause        : 'for' name do_group
//                   | 'for' name sequential_sep do_group
//                   | 'for' name LineBreak in sequential_sep do_group
//                   | 'for' name LineBreak in wordlist sequential_sep do_group
//                   ;
func parseForClause(s *State, t *tokenizer) {

}

// in                : 'in'
//                   ;
func parseIn(s *State, t *tokenizer) {

}

// wordlist          : wordlist WORD
//                   | WORD
//                   ;
func parseWordList(s *State, t *tokenizer) {

}

// case_clause       : 'case' WORD LineBreak 'in' LineBreak case_list 'esac'
//                   | 'case' WORD LineBreak 'in' LineBreak case_list_ns 'esac'
//                   | 'case' WORD LineBreak 'in' LineBreak 'esac'
//                   ;
func parseCaseClause(s *State, t *tokenizer) {

}

// case_list_ns      : case_list case_item_ns
//                   | case_item_ns
//                   ;
func parseCaseListNS(s *State, t *tokenizer) {

}

// case_list         : case_list case_item
//                   | case_item
//                   ;
func parseCaseList(s *State, t *tokenizer) {

}

// case_item_ns      : pattern ')' LineBreak
//                   | pattern ')' compound_list
//                   | '(' pattern ')' LineBreak
//                   | '(' pattern ')' compound_list
//                   ;
func parseCaseItemNS(s *State, t *tokenizer) {

}

// case_item         : pattern ')' ';;' LineBreak
//                   | pattern ')' ';;' compound_list
//                   | '(' pattern ')' ';;' LineBreak
//                   | '(' pattern ')' ';;' compound_list
//                   ;
func parseCaseItem(s *State, t *tokenizer) {

}

// pattern           : WORD
//                   | pattern '|' WORD
//                   ;
func parsePattern(s *State, t *tokenizer) {

}

// if_clause         : 'if' compound_list 'then' compound_list else_part 'fi'
//                   | 'if' compound_list 'then' compound_list 'fi'
func parseIfClause(s *State, t *tokenizer) {

}

// else_part         : 'elif' compound_list 'then' compound_list
//                   | 'elif' compound_list 'then' compound_list else_part
//                   | 'else' compound_list
//                   ;
func parseElsePart(s *State, t *tokenizer) {

}

// while_clause      : 'while' compound_list do_group
//                   ;
func parseWhileClause(s *State, t *tokenizer) {

}

// until_clause      : 'until' compound_list do_group
//                   ;
func parseUntilClause(s *State, t *tokenizer) {

}

// function_definition : fname '(' ')' LineBreak function_body
//                   ;
func parseFunctionDefinition(s *State, t *tokenizer) {

}

// function_body     : compound_command
//                   | compound_command redirect_list
//                   ;
func parseFunctionBody(s *State, t *tokenizer) {

}

// fname             : NAME
//                   ;
func parseFName(s *State, t *tokenizer) {

}

// brace_group       : '{' compound_list '}'
//                   ;
func parseBraceGroup(s *State, t *tokenizer) {

}

// do_group          : 'do' compound_list 'done'
//                   ;
func parseDoGroup(s *State, t *tokenizer) {

}

// SimpleCommand ::= CmdPrefix SimpleCommand2 | CmdName CmdSuffix
// SimpleCommand2 ::= CmdWord CmdSuffix |
func parseSimpleCommand(s *State, t *tokenizer) command {
	cmd := simpleCommand{}
	parseCmdPrefix(s, t, &cmd)
	parseCmdName(s, t, &cmd)
	parseCmdSuffix(s, t, &cmd)
	return &cmd
}

// CmdName ::= WORD
func parseCmdName(s *State, t *tokenizer, cmd *simpleCommand) {
	if t.ts[0].ttype == ttWord {
		cmd.name = []byte(t.ts[0].value)
		cmd.args = [][]byte{cmd.name}
		t.ts = t.ts[1:]
	} else {
		panic("Bad parse") // TODO: better error handling
	}
}

// TODO: generalize to parseWord ???
// CmdWord ::= WORD
func parseCmdWord(s *State, t *tokenizer) []byte {
	if t.ts[0].ttype == ttWord {
		cmdWord := t.ts[0].value
		t.ts = t.ts[1:]
		return []byte(cmdWord)
	}
	panic("Bad parse") // TODO: better error handling
}

// CmdPrefix ::= IORedirect CmdPrefix | Assignment_WORD CmdPrefix |
func parseCmdPrefix(s *State, t *tokenizer, cmd *simpleCommand) {
	// TODO
}

// CmdSuffix ::= IORedirect CmdSuffix | WORD CmdSuffix |
func parseCmdSuffix(s *State, t *tokenizer, cmd *simpleCommand) {
	for {
		switch t.ts[0].value {
		case "<", "<&", ">", ">&", ">>", "<>", ">|": // TODO: IO_NUMBER
			parseIORedirect(s, t, cmd)
		case "&&", "||", ";", "&", "|", "\n", "": // TODO: follow set
			return
		default:
			cmd.args = append(cmd.args, []byte(t.ts[0].value))
			t.ts = t.ts[1:]
		}
	}
}

// redirect_list     : io_redirect
//                   | redirect_list io_redirect
//                   ;
func parseRedirectList(s *State, t *tokenizer) {

}

// IORedirect ::= IORedirect2 | IO_NUMBER IORedirect2
// IORedirect2 ::= IOFile | io_here
func parseIORedirect(s *State, t *tokenizer, cmd *simpleCommand) {
	// TODO: IO_NUMBER io_here
	parseIOFile(s, t, cmd)
}

// IOFile ::= IOOp Filename
// IOOp ::= '<' | '<&' | '>' | '>&' | '>>' | '<>' | '>|'
func parseIOFile(s *State, t *tokenizer, cmd *simpleCommand) {
	cmd.redirects = append(cmd.redirects, redirect{
		ioOp:     parseCmdWord(s, t),
		filename: parseFilename(s, t),
	})
}

// TODO: might be able to replace by parseWord
// Filename ::= WORD
func parseFilename(s *State, t *tokenizer) []byte {
	if t.ts[0].ttype == ttWord {
		filename := t.ts[0].value
		t.ts = t.ts[1:]
		return []byte(filename)
	}
	panic("Bad parse") // TODO: better error handling
}

// io_here           : DLESS here_end
//                   | DLESSDASH here_end
//                   ;
func parseIOHere(s *State, t *tokenizer) {

}

// here_end          : WORD
//                   ;
func parseHereEnd(s *State, t *tokenizer) {

}

// NewLineList ::= NEWLINE NewLineList | NEWLINE
func parseNewLineList(s *State, t *tokenizer) {
	if t.ts[0].ttype != ttNewLine {
		panic("Parse error") // TODO: better error message
	}
	for t.ts[0].ttype == ttNewLine {
		t.ts = t.ts[1:]
	}
	// TODO: follow set?
}

// LineBreak ::= NEWLINE LineBreak |
func parseLineBreak(s *State, t *tokenizer) {
	for t.ts[0].ttype == ttNewLine {
		t.ts = t.ts[1:]
	}
	// TODO: follow set?
}

// SeparatorOp ::= '&' | ';'
func parseSeparatorOp(s *State, t *tokenizer) {
	switch t.ts[0].value {
	case "&", ";":
		t.ts = t.ts[1:]
	default:
		panic("Parse error")
	}
}

// separator         : separator_op LineBreak
//                   | NewLineList
//                   ;
func parseSeparator(s *State, t *tokenizer) {

}

// sequential_sep    : ';' LineBreak
//                   | NewLineList
//                   ;
func parseSequentialSep(s *State, t *tokenizer) {

}
