package resp

import (
	"bytes"
	"testing"
)

func TestReadDataString(t *testing.T) {

	stringData := map[string]string {
		"+Hello\r\n" : "Hello",
	}

	for key, target := range stringData {
		body := bytes.NewBufferString(key)
		ret, err := ReadData(body)

		if nil!=err || target != ret.String() || ret.T != T_SimpleString {
			t.Errorf("target:%s result:%s, err:%s", target, ret.String(), err.Error())
		}
	}
}

func TestReadDataError(t *testing.T) {
	stringData := map[string]string {
		"-Redis Server Went Away\r\n" : "Redis Server Went Away",
	}

	for key, target := range stringData {
		body := bytes.NewBufferString(key)
		ret, err := ReadData(body)

		if nil!=err || target != ret.Error() || ret.T != T_Error {
			t.Errorf("target:%s result:%s, err:%s", target, ret.Error(), err.Error())
		}
	}

}


func TestReadDataInteger(t *testing.T) {
	stringData := map[string]int64 {
		":1000\r\n" : 1000,
	}

	for key, target := range stringData {
		body := bytes.NewBufferString(key)
		ret, err := ReadData(body)

		if nil!=err || target != ret.Integer() || ret.T != T_Integer {
			t.Errorf("target:%s result:%d, err:%s", target, ret.Integer(), err.Error())
		}
	}

}


func TestReadDataBulkString(t *testing.T) {
	stringData := map[string]string {
		"$6\r\nfoobar\r\n" : "foobar",
	}

	for key, target := range stringData {
		body := bytes.NewBufferString(key)
		ret, err := ReadData(body)

		if nil!=err || target != ret.String() || ret.T != T_BulkString {
			t.Errorf("target:%s result:%d, err:%s", target, ret.String(), err.Error())
		}
	}

}

func TestReadDataArray(t *testing.T) {
	key := "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"
	target := []string{"foo", "bar"}

	body := bytes.NewBufferString(key)
	ret, err := ReadData(body)

	if nil!=err ||  ret.T != T_Array {
		t.Error(ret.T)
	}

	arr := ret.Array()
	if len(arr) != 2 || arr[0].String() != target[0] || arr[1].String() != target[1] {
		t.Error(arr[0], arr[1])
	}



}
