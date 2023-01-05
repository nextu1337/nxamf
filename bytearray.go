package nxamf

type ByteArray struct {
	data interface{}
}

func (a *ByteArray) SetData(data interface{}) {
	a.data = data
}

func (a *ByteArray) GetData() interface{} {
	return a.data
}

func NewByteArray(data interface{}) *ByteArray {
	a := new(ByteArray)
	a.SetData(data)
	return a
}