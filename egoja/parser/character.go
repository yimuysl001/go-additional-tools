package eparser

import "unicode/utf8"

type CharacterStream struct {
	source string
	pos    int
}

func NewCharacterStream(source string) *CharacterStream {
	return &CharacterStream{source: source, pos: 0}
}

func (cs *CharacterStream) hasMore() bool {
	return cs.pos < len(cs.source)
}

func (cs *CharacterStream) match(s string, consume bool) bool {
	if cs.pos+len(s) > len(cs.source) {
		return false
	}
	if cs.source[cs.pos:cs.pos+len(s)] == s {
		if consume {
			cs.pos += len(s)
		}
		return true
	}
	return false
}

func (cs *CharacterStream) consume() rune {
	if !cs.hasMore() {
		return 0
	}
	r, size := utf8.DecodeRuneInString(cs.source[cs.pos:])
	cs.pos += size
	return r
}

func (cs *CharacterStream) getPosition() int {
	return cs.pos
}

func (cs *CharacterStream) substring(start, end int) string {
	if start < 0 || end > len(cs.source) || start > end {
		return ""
	}
	return cs.source[start:end]
}
