package classic

const (
	EOF = iota
	_NUM_CHAR
	_ESCAPED_CHAR
	_TERM_START_CHAR
	_TERM_CHAR
	_WHITESPACE
	_QUOTED_CHAR
	_UNUSED
	AND
	OR
	NOT
	PLUS
	MINUS
	BAREOPER
	LPAREN
	RPAREN
	COLON
	STAR
	CARAT
	QUOTED
	TERM
	FUZZY_SLOP
	PREFIXTERM
	WILDTERM
	REGEXPTERM
	RANGEIN_START
	RANGEEX_START
	NUMBER
	RANGE_TO
	RANGEiN_END
	RAGEEX_END
	RANGE_QUOTED
	RANGE_GOOP
)
