// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package httpserver

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonE34310f8DecodeSpOfficeLookuperInternalHttpserver(in *jlexer.Lexer, out *errorResponse) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "status":
			out.Status = int(in.Int())
		case "code":
			out.Code = int(in.Int())
		case "message":
			out.Message = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE34310f8EncodeSpOfficeLookuperInternalHttpserver(out *jwriter.Writer, in errorResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Status))
	}
	{
		const prefix string = ",\"code\":"
		out.RawString(prefix)
		out.Int(int(in.Code))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v errorResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE34310f8EncodeSpOfficeLookuperInternalHttpserver(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v errorResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE34310f8EncodeSpOfficeLookuperInternalHttpserver(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *errorResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE34310f8DecodeSpOfficeLookuperInternalHttpserver(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *errorResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE34310f8DecodeSpOfficeLookuperInternalHttpserver(l, v)
}