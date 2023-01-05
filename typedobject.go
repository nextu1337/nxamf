package nxamf

type TypedObject struct {
	amfClassName string
	amfData interface{}
}

func (a *TypedObject) GetAMFClassName() string {
	return a.amfClassName
}

func (a *TypedObject) GetAMFData() interface{} {
	return a.amfData
}

func (a *TypedObject) SetAMFClassName(classname string) {
	a.amfClassName = classname
}

func (a *TypedObject) SetAMFData(data interface{}) {
	a.amfData = data
}


func NewTypedObject(classname string, data interface{}) *TypedObject {
	a := new(TypedObject)
	a.amfClassName = classname
	a.amfData = data
	return a
}