package nxamf

type ByteArray struct {
	data string
}

func (a *ByteArray) SetData(data string) {
	a.data = data
}

func (a *ByteArray) GetData() string {
	return a.data
}

func NewByteArray(data string) *ByteArray {
	a := new(ByteArray)
	a.SetData(data)
	return a
}