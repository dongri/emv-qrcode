package mpm

import (
	"errors"
	"fmt"
	"strconv"
	"unicode/utf8"
)

type Error struct {
	Func string
	Err  error
}

func (e *Error) Error() string {
	return "parser." + e.Func + ": " + e.Err.Error()
}

func notCallError(fn string) *Error {
	return &Error{
		Func: fn,
		Err:  errors.New("not call Next()"),
	}
}

func outOfRangeError(fn string, current, max, start, end int64) *Error {
	return &Error{
		Func: fn,
		Err:  fmt.Errorf("bounds out of range. current: %d, max: %d, start: %d, end: %d", current, max, start, end),
	}
}

type NumError struct {
	Func string
	Num  string
	Err  error
}

func (e *NumError) Error() string {
	return "parser." + e.Func + ": " + "parsing " + strconv.Quote(e.Num) + ": " + e.Err.Error()
}

func syntaxError(fn, str string) *NumError {
	return &NumError{
		Func: fn,
		Num:  str,
		Err:  strconv.ErrSyntax,
	}
}

const (
	IDWordCount          = 2
	ValueLengthWordCount = 2
)

type Parser struct {
	current int64
	max     int64
	source  []rune
	err     error
}

func NewParser(payload string) (*Parser, error) {
	return &Parser{
		current: -1,
		max:     int64(utf8.RuneCountInString(payload)),
		source:  []rune(payload),
		err:     nil,
	}, nil
}

func (p *Parser) Next() bool {
	if p.err != nil {
		return false
	}
	if p.current < 0 {
		p.current = 0
	} else {
		valueLength := p.ValueLength()
		if p.err != nil {
			return false
		}
		p.current += valueLength + IDWordCount + ValueLengthWordCount
	}
	if p.current >= p.max {
		return false
	}
	return true
}

func (p *Parser) ID() int64 {
	const fnID = "ID"
	start := p.current
	end := start + IDWordCount
	if p.current < 0 {
		p.err = notCallError(fnID)
		return -1
	}
	if p.max < end {
		p.err = outOfRangeError(fnID, p.current, p.max, start, end)
		return -1
	}
	strID := string(p.source[start:end])
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		p.err = syntaxError(fnID, strID)
		return -1
	}
	return id
}

func (p *Parser) ValueLength() int64 {
	const fnValueLength = "ValueLength"
	start := p.current + IDWordCount
	end := start + ValueLengthWordCount
	if p.current < 0 {
		p.err = notCallError(fnValueLength)
		return 0
	}
	if p.max < end {
		p.err = outOfRangeError(fnValueLength, p.current, p.max, start, end)
		return 0
	}
	strValueLength := string(p.source[start:end])
	len, err := strconv.ParseInt(strValueLength, 10, 64)
	if err != nil {
		p.err = syntaxError(fnValueLength, strValueLength)
		return 0
	}
	return len
}

func (p *Parser) Value() string {
	const fnValue = "Value"
	start := p.current + IDWordCount + ValueLengthWordCount
	end := start + p.ValueLength()
	if p.current < 0 {
		p.err = notCallError(fnValue)
		return ""
	}
	if p.max < end {
		p.err = outOfRangeError(fnValue, p.current, p.max, start, end)
		return ""
	}
	return string(p.source[start:end])
}

func (p *Parser) Err() error {
	return p.err
}
