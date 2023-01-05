package nxamf

type AMF3_Wrapper struct {
	data interface{}
}

func (a *AMF3_Wrapper) SetData(data interface{}) {
	a.data = data
}

func (a *AMF3_Wrapper) GetData() interface{} {
	return a.data
}

func NewAMF3_Wrapper(data interface{}) *AMF3_Wrapper {
	a := new(AMF3_Wrapper)
	a.SetData(data)
	return a
}