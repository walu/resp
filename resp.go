package resp

import (
	"io"
	"errors"
	"bytes"
	"strconv"
)

const (
	T_SimpleString = '+'
	T_Error	   = '-'
	T_Integer	   = ':'
	T_BulkString  = '$'
	T_Array	   = '*'
)

//Command
//
//Command 格式：Inline Command 与 Array With BulkString
type Command struct {
	//根据惯例，Args[0] 是Name本身
	Args []string
}

//返回Command的名称，如GET\SET
func (c Command) Name() string {
	if len(c.Args)==0 {
		return ""
	} else {
		return c.Args[0]
	}
}

//以String形式获取Command[index]
func (c Command) String(index int) (ret string) {
	if len(c.Args) > index {
		ret = c.Args[index]
	}
	return ret
}

//以int64的形式返回Command.Args[index]
func (c Command) Integer(index int) (ret int64) {
	if len(c.Args) > index {
		ret, _ = strconv.ParseInt(c.Args[index], 10, 64)
	}
	return ret
}

func NewCommand(args ...string) (*Command, error) {
	if len(args) == 0 {
		return nil, errors.New("err_new_cmd")
	}
	return &Command{Args:args}, nil
}

//从Reader中读取Command
func ReadCommand(r io.Reader) (*Command, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r, buf)
	if nil != err {
		return nil, err
	}

	if T_Array != buf[0] {
		return nil, errors.New("Unexpected Command Type")
	}

	var ret *Data
	ret = new(Data)

	ret, err = readDataForSpecType(r, buf[0])
	if nil != err {
		return nil, err
	}

	commandArgs := make([]string, len(ret.Array))
	for index := range ret.Array {
		if ret.Array[index].T != T_BulkString {
			return nil, errors.New("Unexpected Command Type")
		}
		commandArgs[index] = string(ret.Array[index].String)
	}

	return NewCommand(commandArgs...)
}

type Data struct {
	T byte
	String []byte
	Integer int64
	Array []*Data
	IsNil bool
}

func ReadData(r io.Reader) (*Data, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r, buf)
	if nil != err {
		return nil, err
	}

	return readDataForSpecType(r, buf[0])
}

func readDataForSpecType(r io.Reader, t byte) (*Data, error) {

	var err error
	var ret *Data

	ret = new(Data)
	switch t {
		case '+':
			ret.T = T_SimpleString
			ret.String, err = readRespLine(r)

		case '-':
			ret.T = T_Error
			ret.String, err = readRespLine(r)

		case ':':
			ret.T = T_Integer
			ret.Integer, err = readRespIntLine(r)

		case '$':
			var lenBulkString int64
			lenBulkString, err = readRespIntLine(r)

			ret.T = T_BulkString
			if -1 != lenBulkString {
				ret.String, err = readRespN(r, lenBulkString)
				_, err = readRespN(r, 2)
			} else {
				ret.IsNil = true
			}

		case '*':
			var lenArray int64
			var i int64
			lenArray, err = readRespIntLine(r)

			ret.T = T_Array
			if nil==err {
				if -1 != lenArray {
					ret.Array = make([]*Data, lenArray)
					for i=0; i<lenArray; i++ {
						ret.Array[i], err = ReadData(r)
					}
				} else {
					ret.IsNil = true
				}
			}

		default: //Maybe you are Inline Command
			err = errors.New("Unexpected type")

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
