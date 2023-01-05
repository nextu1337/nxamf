package nxamf

import (
	"math"
	"encoding/binary"
)

type OutputStream struct {
	RawData string
}

func (a *OutputStream) WriteBuffer(buf string) {
	a.RawData+=buf
}

func (a *OutputStream) WriteByte(bt int) {
	a.RawData+=string(bt)
}

func (a *OutputStream) WriteInt(integer uint16) {
	bytes := make([]byte, 4) 
    binary.BigEndian.PutUint16(bytes, integer)
	a.RawData+=string(bytes[:2])
}

func (a *OutputStream) WriteDouble(double float64) {
	var buf [8]byte
    binary.BigEndian.PutUint64(buf[:], math.Float64bits(double))
    a.RawData+=string(buf[:])
	// bits := math.Float64bits(double)
    // bytes := make([]byte, 8)
    // binary.LittleEndian.PutUint64(bytes, bits)
	// a.RawData+=strrev(string(bytes))
}

func (a *OutputStream) WriteLongMinus1() {
	bytes := make([]byte, 8) 
    binary.BigEndian.PutUint64(bytes, 4294967295)
	a.RawData+=string(bytes[4:])
}

func (a *OutputStream) WriteLong(long uint64) {
	bytes := make([]byte, 8) 
    binary.BigEndian.PutUint64(bytes, long)
	a.RawData+=string(bytes[4:])
}


func (a *OutputStream) GetRawData() string {
	return a.RawData
}

func NewOutputStream() *OutputStream {
	a := new(OutputStream)
	a.RawData = ""
	return a
}