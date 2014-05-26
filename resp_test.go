package resp

import (
	"bytes"
	"testing"
)

var (
	respSimpleString = Data{T:T_SimpleString, String:[]byte("OK")}
	respSimpleStringText = "+OK\r\n"

	respError = Data{T:T_Error, String:[]byte("Error message")}
	respErrorText = "-Error message\r\n"

	respBulkString = Data{T:T_BulkString, String:[]byte("foobar")}
	respBulkStringText = "$6\r\nfoobar\r\n"

	respNilBulkString = Data{T:T_BulkString, IsNil:true}
	respNilBulkStringText = "$-1\r\n"

	respInteger = Data{T:T_Integer, Integer:1000}
	respIntegerText = ":1000\r\n"

	respArray = Data{T:T_Array, Array:[]*Data{&respSimpleString, &respInteger}}
	respArrayText = "*2\r\n" + respSimpleStringText + respIntegerText
)


var validCommand map[string]string
var validData map[string]Data

func TestValidData(t *testing.T) {
	for text, data := range validData {
		buf := bytes.NewReader([]byte(text))
		//test read
		d, err := ReadData(buf)
		if nil!=err || d.T != data.T {
			t.Error(err)
		}

		if false == eqData(*d, data) {
			t.Error(text, *d, data)
		}

		//test format
		if text != string(data.Format()) {
			t.Error(text, data)
		}
	}
}

func BenchmarkDataFormat(b *testing.B) {
	for i:=0; i<b.N; i++ {
		for _, data := range validData {
			data.Format()
		}
	}
}

func BenchmarkReadData(b *testing.B) {
	for i:=0; i<b.N; i++ {
		for text, _ := range validData {
			buf := bytes.NewReader([]byte(text))
			ReadData(buf)
		}
	}
}



func eqData(d1, d2 Data) bool {
	eqType := d1.T == d2.T
	eqString := 0==bytes.Compare(d1.String, d2.String)
	eqInteger := d1.Integer == d2.Integer
	eqNil := d1.IsNil == d2.IsNil
	eqArrayLen := len(d1.Array) == len(d2.Array)
	eqArray := true
	if len(d1.Array) > 0 && eqArrayLen {
		for index := range d1.Array {
			if false == eqData(*d1.Array[index], *d2.Array[index]) {
				eqArray = false
				break
			}
		}
	}
	return eqType && eqString && eqInteger && eqNil && eqArrayLen && eqArray
}

func TestValidCommand(t *testing.T) {
	for input, cmd := range validCommand {
		reader := bytes.NewReader([]byte(input))
		c, err := ReadCommand(reader)
		if nil != err {
			t.Error("read command error", err)
		} else if c.Name() != cmd {
			t.Error("read command error", c.Name(), cmd)
		}
	}
}

func _validCommand(tb testing.TB) {
	for input, cmd := range validCommand {
		reader := bytes.NewReader([]byte(input))
		c, err := ReadCommand(reader)
		if nil != err {
			tb.Error("read command error", err)
		} else if c.Name() != cmd {
			tb.Error("read command error", c.Name(), cmd)
		}
	}

}

func BenchmarkValidCommand(b *testing.B) {
	for i:=0; i<b.N; i++ {
		_validCommand(b)
	}
}

func testCommandFormat(t *testing.T) {
	cmd, _ := NewCommand("LLEN", "walu.cc")
	if "*2\r\n$7\r\nwalu.cc\r\n" != string(cmd.Format()) {
		t.Error("cmd format error")
	}

}

func BenchmarkCommandFormat(b *testing.B) {
	cmd, _ := NewCommand("LLEN", "walu.cc")
	for i:=0; i< b.N; i++ {
		cmd.Format()
	}
}

func init() {
	validCommand = map[string]string{
		"PING" : "PING",
		"PING\n" : "PING",
		"PING\r" : "PING",
		"  PING ": "PING",
		"*2\r\n$4\r\nLLEN\r\n$6\r\nmysist\r\n" : "LLEN",
	}

	validData = map[string]Data {
		respSimpleStringText : respSimpleString,
		respErrorText : respError,
		respBulkStringText : respBulkString,
		respNilBulkStringText : respNilBulkString,
		respIntegerText : respInteger,
		respArrayText : respArray,
	}
}
