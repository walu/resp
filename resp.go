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

type Data struct {
	T byte
	str []byte
	num int64
	array []*Data

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

func ReadData(r io.Reader) (*Data, error) {

	var buf []byte
	var n int
	var err error

	buf = make([]byte, 1)
	n, err = io.ReadFull(r, buf)
	if nil != err {
		return nil, errors.New("err_first_byte")
	}

	if n==1 {
	//
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
			ret.str, err = readRespN(r, lenBulkString)

			//read the followed \r\n
			_, err = readRespN(r, 2)

		case '*':
			var lenArray int64
			var i int64
			lenArray, err = readRespIntLine(r)

			ret.T = T_Array
			if nil==err {
				ret.array = make([]*Data, lenArray)
				for i=0; i<lenArray; i++ {
					ret.array[i], err = ReadData(r)
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
