package resp

import (
	"io/ioutil"
)

const (
	var T_SimpleString = '+'
	var T_Errors	   = '-'
	var T_Integers	   = ':'
	var T_BulkStrings  = '$'
	var T_Arrays	   = '*'
)

type Data struct {
	T byte
	str []byte
	num int64
	arrays []Data

}

//string\bulkString
func (d Data) String() []byte {

}

func (d Data) Errors() string {

}

func (d Data) Array() []Data {

}

func ReadData(r io.Reader) (Data, error) {

	var buf []byte
	var n int
	var err error

	var requestType []byte
	requestType = make([]byte, 1)
	n, err = io.ReadAtLeast(r, buf, 1)
	if nil != err {
		return errors.New("err_first_byte")
	}

	ret := Data{}
	switch buf[0] {
		case '+':
			ret.T = T_SimgpleString
			ret.common, err = readRespLine(r)

		case '-':
			ret.T = T_Errors
			ret.common, err = readRespLine(r)

		case ':':
			ret.T = T_Integers
			ret.num, err = readRespIntLine(r)

		case '$':
			var lenBulkString int64
			lenBulkString, err = readRespIntLine(r)

			ret.T = T_BulkString
			ret.common, err = readRespN(r, lenBulkString)

		case '*':
			var lenArray int64
			lenArray, err = readRespIntLine(r)
			data.array = make(Data, lenArray)
			for i:=0; i<lenArray; i++ {
				data.array[i], err = readData(r)
			}

	}
	return data, err
}

//读取当前行，并去掉最后的\r\n
func readRespLine(r io.Reader, d byte) ([]byte, error) {

}

//读取N个字节，并去掉最后的\r\n
func readRespN(r io.Reader, n int) ([]byte, error) {

}

//读取当前行的数字，并去掉最后的\r\n
func readRespIntLine(r io.Reader) (int64, error) {

}
