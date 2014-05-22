package resp

import (
	"bytes"
	"testing"
)

type testResult struct {
	T byte
	Value interface{}
}

var validBody map[string]testResult

func TestValidData(t *testing.T) {
	for body, result := range validBody {
		buf := bytes.NewReader([]byte(body))
		ret, err := ReadData(buf)
		if nil != err {
			t.Error(err)
		}

		switch ret.T {
			case T_SimpleString, T_BulkString:
				if 0 != bytes.Compare(result.Value.([]byte), ret.String) {
					t.Error("not eq")
				}
			case T_Integer:
				if result.Value.(int64) != ret.Integer {
					t.Error("not eq")
				}
		}
	}
}

func BenchmarkReadDataBulkString(b *testing.B) {
}

func init() {
	validBody = map[string]testResult {
		"+OK\r\n" : {T_SimpleString, []byte("OK")},
		"-Errors\r\n" : {T_Error, []byte("Errors")},
		":100\r\n" : {T_Integer, int64(100)},
		"$7\r\nwalu.cc\r\n" : {T_BulkString, []byte("walu.cc")},
	}
}
