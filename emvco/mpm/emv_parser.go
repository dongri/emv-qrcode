package mpm

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

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
		p.current += p.ValueLength() + IDWordCount + ValueLengthWordCount
		if p.err != nil {
			return false
		}
	}
	if p.current >= p.max {
		return false
	}
	return true
}

func (p *Parser) ID() int64 {
	start := p.current
	end := start + IDWordCount
	if p.current < 0 || p.max < end {
		p.err = fmt.Errorf("Next() is not called, or the length of slice bounds out of range")
		return -1
	}
	id, err := strconv.ParseInt(string(p.source[start:end]), 10, 64)
	if err != nil {
		p.err = err
		return -1
	}
	return id
}

func (p *Parser) ValueLength() int64 {
	start := p.current + IDWordCount
	end := start + ValueLengthWordCount
	if p.current < 0 || p.max < end {
		p.err = fmt.Errorf("Next() is not called, or the length of slice bounds out of range")
		return 0
	}
	len, err := strconv.ParseInt(string(p.source[start:end]), 10, 64)
	if err != nil {
		p.err = err
		return 0
	}
	return len
}

func (p *Parser) Value() string {
	start := p.current + IDWordCount + ValueLengthWordCount
	end := start + p.ValueLength()
	if p.current < 0 || p.max < end {
		p.err = fmt.Errorf("Next() is not called, or the length of slice bounds out of range")
		return ""
	}
	return string(p.source[start:end])
}

func (p *Parser) Err() error {
	return p.err
}
