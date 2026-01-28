package lexer

// Template keywords used by both lexer and parser.
const (
	KeywordIf       = "if"
	KeywordElse     = "else"
	KeywordEnd      = "end"
	KeywordRange    = "range"
	KeywordDefine   = "define"
	KeywordTemplate = "template"
	KeywordBlock    = "block"
	KeywordWith     = "with"
	KeywordContinue = "continue"
	KeywordBreak    = "break"
)

// KeywordPattern is the regex pattern matching all template keywords.
// Used by the lexer to identify keywords in template blocks.
// The word boundary \b ensures keywords don't match when they're part of a longer identifier
// (e.g., "breakAll" should not match "break" as a keyword).
const KeywordPattern = `(?:` + KeywordIf + "|" + KeywordElse + "|" + KeywordEnd + "|" +
	KeywordRange + "|" + KeywordDefine + "|" + KeywordTemplate + "|" +
	KeywordBlock + "|" + KeywordWith + "|" + KeywordContinue + "|" + KeywordBreak + `)\b`
