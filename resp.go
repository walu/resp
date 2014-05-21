//package resp provides methods to parse and format resp(redis protocal) data
package resp

import (
	"io"
	"fmt"
	"errors"
	"bytes"
	"strconv"
)

const (
	T_SimpleString	= '+'
	T_Error		= '-'
	T_Integer	= ':'
	T_BulkString	= '$'
	T_Array		= '*'
)

type Data struct {
	T byte
	str []byte
	num int64
	array []*Data
	isNil bool
}

//string\bulkString
func (d *Data) String() string {
	return string(d.str)
}

func (d *Data) Byte() []byte {
	return d.str
}

func (d *Data) Error() string {
	return string(d.str)
}

func (d *Data) Integer() int64 {
	return d.num
}

func (d *Data) Array() []*Data {
	return d.array
}

func (d *Data) IsNil() bool {
	return d.isNil == true
}

/*
valid type: string, []byte, error, int64, []interface{}
string -> T_SimpleString
[]byte -> T_BulkString
error -> T_Error
int* -> T_Integer
[]interface{} -> array

*/
func NewData(val interface{}) (ret *Data, err error) {
	ret = new(Data)
	switch val.(type) {
		case string:
			ret.T = T_SimpleString
			ret.str = []byte(val.(string))
		case error:
			ret.T = T_Error
			ret.str = []byte(val.(error).Error())
		case int64, int8, int16, int32:
			ret.T = T_Integer
			ret.num = val.(int64)
		case []interface{}:
			ret.T = T_Array
			ret.array = make([]*Data, len(val.([]interface{})))
			for index := range val.([]interface{}) {
				ret.array[index], err = NewData(val.([]interface{}))
				if nil!=err {
					goto end
				}
			}
		default:
			return nil, errors.New("unsupported type")
	}
end:
	if nil!=err {
		ret = nil
	}
	return ret, err
}

func NewSimpleString() {}
func NewError() {}
func NewInteger() {}
func NewArray() {}

//format *Data to []byte
func FormatData(d *Data) []byte {
	ret := new(bytes.Buffer)
	ret.WriteByte(d.T)
	switch d.T {
	case T_SimpleString, T_Error:
		fmt.Fprintf(ret, "%s\r\n", d.str)
	case T_Integer:
		fmt.Fprintf(ret, "%d\r\n", d.num)
	case T_BulkString:
		fmt.Fprintf(ret, "%d\r\n%s\r\n", len(d.str), string(d.str))
	case T_Array:
		fmt.Fprintf(ret, "%d\r\n", len(d.array))
		for index := range d.array {
			ret.Write(FormatData(d.array[index]))
		}
	}
	return ret.Bytes()
}

//read from io.Reader, and parse into *Data
func ReadData(r io.Reader) (*Data, error) {

	var buf []byte
	var err error

	buf = make([]byte, 1)
	_, err = io.ReadFull(r, buf)
	if nil != err {
		return nil, err
	}

	ret := &Data{}
	switch buf[0] {
		case '+':
			ret.T = T_SimpleString
			ret.str, err = readRespLine(r)

		case '-':
			ret.T = T_Error
			ret.str, err = readRespLine(r)

		case ':':
			ret.T = T_Integer
			ret.num, err = readRespIntLine(r)

		case '$':
			var lenBulkString int64
			lenBulkString, err = readRespIntLine(r)

			ret.T = T_BulkString
			if -1 == lenBulkString {
				ret.isNil = true
			} else {
				ret.str, err = readRespN(r, lenBulkString)
				//read the followed \r\n
				_, err = readRespN(r, 2)
			}

		case '*':
			var lenArray int64
			var i int64
			lenArray, err = readRespIntLine(r)

			ret.T = T_Array
			if -1 == lenArray {
				ret.isNil = true
			} else if nil==err {
				ret.array = make([]*Data, lenArray)
				for i=0; i<lenArray; i++ {
					ret.array[i], err = ReadData(r)
				}
			}
		default:
			//Inline Commands
			tmp, err := readRespLine(r)
			if nil==err {
				tmpSlice := bytes.Fields(tmp)

				ret.T = T_Array
				ret.array = make([]*Data, len(tmpSlice))
				for index := range tmpSlice {
					t := &Data{}
					t.str = tmpSlice[index]
					t.T = T_SimpleString
					ret.array[index] = t
				}
			}

	}
	return ret, err
}

//读取当前行，并去掉最后的\r\n
func readRespLine(r io.Reader) ([]byte, error) {

	var n, i int
	var err error
	var buf []byte
	var ret *bytes.Buffer

	buf = make([]byte, 1)
	ret = &bytes.Buffer{}

	for {
		n, err = io.ReadFull(r, buf)
		if nil != err {
			return nil, err
		}

		if n==0 {
			continue
		}

		i++
		ret.WriteByte(buf[0])
		if '\n' == buf[0] {
			break
		}
	}

	return ret.Next(i-2), nil
}

//读取N个字节，并去掉最后的\r\n
func readRespN(r io.Reader, n int64) ([]byte, error) {
	var err error
	var ret []byte

	ret = make([]byte, n)
	_, err = io.ReadFull(r, ret)
	if nil!=err {
		ret = nil
	}
	return ret, err
}

//读取当前行的数字，并去掉最后的\r\n
func readRespIntLine(r io.Reader) (int64, error) {
	line, err := readRespLine(r)
	if nil!=err {
		return 0, err
	}
	return strconv.ParseInt(string(line), 10, 64)
}
