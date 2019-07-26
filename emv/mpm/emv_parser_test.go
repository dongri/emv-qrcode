package mpm

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestParserError_Error(t *testing.T) {
	type fields struct {
		Func string
		Err  error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				Func: "Sample",
				Err:  errors.New("sample error"),
			},
			want: "parser.Sample: sample error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ParserError{
				Func: tt.fields.Func,
				Err:  tt.fields.Err,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("ParserError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewParser(t *testing.T) {
	type args struct {
		payload string
	}
	tests := []struct {
		name string
		args args
		want *Parser
	}{
		{
			name: "ok",
			args: args{
				payload: "1234",
			},
			want: &Parser{
				current: -1,
				max:     4,
				source:  []rune("1234"),
				err:     nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewParser(tt.args.payload)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Next(t *testing.T) {
	const fnValueLength = "ValueLength"

	tests := []struct {
		name       string
		parser     *Parser
		want       bool
		wantParser *Parser
	}{
		{
			name: "first time",
			parser: &Parser{
				current: -1,
				max:     12,
				source:  []rune("000201010211"),
				err:     nil,
			},
			want: true,
			wantParser: &Parser{
				current: 0,
				max:     12,
				source:  []rune("000201010211"),
				err:     nil,
			},
		},
		{
			name: "second time",
			parser: &Parser{
				current: 0,
				max:     12,
				source:  []rune("000201010211"),
				err:     nil,
			},
			want: true,
			wantParser: &Parser{
				current: 6,
				max:     12,
				source:  []rune("000201010211"),
				err:     nil,
			},
		},
		{
			name: "third time (finish success)",
			parser: &Parser{
				current: 6,
				max:     12,
				source:  []rune("000201010211"),
				err:     nil,
			},
			want: false,
			wantParser: &Parser{
				current: 12,
				max:     12,
				source:  []rune("000201010211"),
				err:     nil,
			},
		},
		{
			name: "first time irregular value length",
			parser: &Parser{
				current: -1,
				max:     12,
				source:  []rune("00ab0101cd11"),
				err:     nil,
			},
			want: true,
			wantParser: &Parser{
				current: 0,
				max:     12,
				source:  []rune("00ab0101cd11"),
				err:     nil,
			},
		},
		{
			name: "second time irregular value length",
			parser: &Parser{
				current: 0,
				max:     12,
				source:  []rune("00ab0101cd11"),
				err:     nil,
			},
			want: false,
			wantParser: &Parser{
				current: 0,
				max:     12,
				source:  []rune("00ab0101cd11"),
				err:     syntaxError(fnValueLength, "ab"),
			},
		},
		{
			name: "not enough value length",
			parser: &Parser{
				current: 0,
				max:     5,
				source:  []rune("00021"),
				err:     nil,
			},
			want: false,
			wantParser: &Parser{
				current: 6,
				max:     5,
				source:  []rune("00021"),
				err:     nil,
			},
		},
		{
			name: "empty payload",
			parser: &Parser{
				current: 0,
				max:     0,
				source:  []rune(""),
				err:     nil,
			},
			want: false,
			wantParser: &Parser{
				current: 0,
				max:     0,
				source:  []rune(""),
				err:     outOfRangeError(fnValueLength, 0, 0, IDWordCount, IDWordCount+ValueLengthWordCount),
			},
		},
		{
			name: "alread exist error",
			parser: &Parser{
				current: 0,
				max:     0,
				source:  []rune(""),
				err:     fmt.Errorf("sample error"),
			},
			want: false,
			wantParser: &Parser{
				current: 0,
				max:     0,
				source:  []rune(""),
				err:     fmt.Errorf("sample error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.parser.Next(); got != tt.want {
				t.Errorf("Parser.Next() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.parser, tt.wantParser) {
				t.Errorf("WantParser = %v, want %v", tt.parser, tt.wantParser)
			}
		})
	}
}

func TestParser_ID(t *testing.T) {
	type fields struct {
		current int64
		max     int64
		source  []rune
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		want   ID
	}{
		{
			name: "ok",
			fields: fields{
				current: 0,
				max:     6,
				source:  []rune("000201"),
				err:     nil,
			},
			want: ID("00"),
		},
		{
			name: "not called Next()",
			fields: fields{
				current: -1,
				max:     6,
				source:  []rune("000201"),
				err:     nil,
			},
			want: ID(""),
		},
		{
			name: "slice bounds out of range",
			fields: fields{
				current: 0,
				max:     0,
				source:  []rune(""),
				err:     nil,
			},
			want: ID(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				current: tt.fields.current,
				max:     tt.fields.max,
				source:  tt.fields.source,
				err:     tt.fields.err,
			}
			if got := p.ID(); got != tt.want {
				t.Errorf("Parser.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_ValueLength(t *testing.T) {
	type fields struct {
		current int64
		max     int64
		source  []rune
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name: "ok",
			fields: fields{
				current: 0,
				max:     6,
				source:  []rune("000201"),
				err:     nil,
			},
			want: 2,
		},
		{
			name: "not called Next()",
			fields: fields{
				current: -1,
				max:     6,
				source:  []rune("000201"),
				err:     nil,
			},
			want: 0,
		},
		{
			name: "slice bounds out of range",
			fields: fields{
				current: 0,
				max:     0,
				source:  []rune(""),
				err:     nil,
			},
			want: 0,
		},
		{
			name: "value length is not number",
			fields: fields{
				current: 0,
				max:     6,
				source:  []rune("00ab00"),
				err:     nil,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				current: tt.fields.current,
				max:     tt.fields.max,
				source:  tt.fields.source,
				err:     tt.fields.err,
			}
			if got := p.ValueLength(); got != tt.want {
				t.Errorf("Parser.ValueLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Value(t *testing.T) {
	type fields struct {
		current int64
		max     int64
		source  []rune
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				current: 0,
				max:     6,
				source:  []rune("000201"),
				err:     nil,
			},
			want: "01",
		},
		{
			name: "not called Next()",
			fields: fields{
				current: -1,
				max:     6,
				source:  []rune("000201"),
				err:     nil,
			},
			want: "",
		},
		{
			name: "slice bounds out of range",
			fields: fields{
				current: 0,
				max:     5,
				source:  []rune("00021"),
				err:     nil,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				current: tt.fields.current,
				max:     tt.fields.max,
				source:  tt.fields.source,
				err:     tt.fields.err,
			}
			if got := p.Value(); got != tt.want {
				t.Errorf("Parser.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParser_Err(t *testing.T) {
	type fields struct {
		current int64
		max     int64
		source  []rune
		err     error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "exist error",
			fields: fields{
				current: 0,
				max:     0,
				source:  []rune(""),
				err:     fmt.Errorf("sample error"),
			},
			wantErr: true,
		},
		{
			name: "not exist error",
			fields: fields{
				current: 0,
				max:     0,
				source:  []rune(""),
				err:     nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parser{
				current: tt.fields.current,
				max:     tt.fields.max,
				source:  tt.fields.source,
				err:     tt.fields.err,
			}
			if err := p.Err(); (err != nil) != tt.wantErr {
				t.Errorf("Parser.Err() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
