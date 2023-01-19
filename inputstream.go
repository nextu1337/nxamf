package nxamf

import (
	"errors"
	"fmt"
	"math"
	"encoding/binary"
)

type InputStream struct {
	cursor int
	rawData string
}

func (a *InputStream) ReadBuffer(length int) (string,error) {
	if length+a.cursor > len(a.rawData) {
		return "", errors.New(fmt.Sprintf("Buffer underrun at position %d. Trying to fetch %d bytes",a.cursor,length))
	}
	data := a.rawData[a.cursor:a.cursor+length]
	a.cursor+=length
	return data,nil
}

func (a *InputStream) ReadByte() (int,error) {
	str, err := a.ReadBuffer(1)
	if err != nil {
		return -1,err
	}
	return int(str[0]),nil // ?
}

func (a *InputStream) ReadInt() (uint16,error) {
	str, err := a.ReadBuffer(2)
	if err != nil {
		return 0,err
	}
    return binary.BigEndian.Uint16([]byte(str)),nil
}

func (a *InputStream) ReadDouble() (float64,error) {
	double,err := a.ReadBuffer(8);
	if err != nil {
		return -1,err
	}
	bits := binary.BigEndian.Uint64([]byte(double))
    float := math.Float64frombits(bits)
	return float,nil
}

func (a *InputStream) ReadLong() (uint32,error) {
	long,err := a.ReadBuffer(4);
	if err != nil {
		return 0,err
	}
	return binary.BigEndian.Uint32([]byte(long)),nil
}

func (a *InputStream) ReadInt24() (uint32,error) {
	int24,err := a.ReadBuffer(3);
	if err != nil {
		return 0,err
	}
	return binary.BigEndian.Uint32([]byte("\x00"+int24)),nil
}

func NewInputStream(data string) *InputStream {
	a:= new(InputStream)
	a.rawData = data
	a.cursor = 0
	return a
}