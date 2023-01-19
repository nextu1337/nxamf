package nxamf

import "time"
import "math"

type AMF0_Deserializer struct {
	Stream *InputStream
	SavedReferences []interface{}
	RefList []interface{}
	SavedRefs []interface{}
	Amf3Deserializer *AMF3_Deserializer
}

func (a *AMF0_Deserializer) ReadAMFData(settype int, newscope bool) interface{} {
	types := AMF0Const()
	if newscope {
		a.RefList = []interface{}{}
	}
	if settype == -1 {
		settype, _ = a.Stream.ReadByte()
	}
	switch settype {
	case types["DT_NUMBER"]:
		val,_ := a.Stream.ReadDouble()
		return val
	case types["DT_BOOL"]:
		val,_ := a.Stream.ReadByte()
		if val > 0 {
			return true
		} else {
			return false
		}
	case types["DT_STRING"]:
		return a.ReadString()
	case types["DT_OBJECT"]:
		return a.ReadObject()
	case types["DT_NULL"]:
	case types["DT_UNDEFINED"]:
	case types["DT_UNSUPPORTED"]:
		return nil
	case types["DT_REFERENCE"]:
		return a.ReadReference()
	case types["DT_MIXEDARRAY"]:
		return a.ReadMixedArray()
	case types["DT_ARRAY"]:
		return a.ReadArray()
	case types["DT_DATE"]:
		return a.ReadDate()
	case types["DT_LONGSTRING"]:
	case types["DT_XML"]:
		return a.ReadLongString()
	case types["DT_TYPEDOBJECT"]:
		return a.ReadTypedObject()
	case types["DT_AMF3"]:
		return a.ReadAMF3Data()
	default:
		return []interface{}{"Unsupported type",settype}
	}
	return nil
}

func (a *AMF0_Deserializer) ReadObject() map[string]interface{} {
	object := map[string]interface{}{}
	a.RefList = append(a.RefList,&object)
	for true {
		key := a.ReadString()
		vartype,_ := a.Stream.ReadByte()
		if vartype==AMF0Const()["DT_OBJECTTERM"] {
			break
		}
		object[key] = a.ReadAMFData(vartype,false)
	}
	/* //No clue how to implement it
		if (defined('SABREAMF_OBJECT_AS_ARRAY')) {
                $object = (object)$object;
            }
	*/
	return object
}

func (a *AMF0_Deserializer) ReadReference() interface{} {
	refId,_ := a.Stream.ReadInt()
	if len(a.RefList)>int(refId) {
		return a.RefList[refId]
	} else {
		return []interface{}{"Invalid reference offset",refId}
	}
	return nil
}

func (a *AMF0_Deserializer) ReadArray() []interface{} {
	length,_ := a.Stream.ReadLong()
	arr := []interface{}{}
	a.RefList = append(a.RefList,&arr)
	for length>0 {
		arr = append(arr,a.ReadAMFData(-1,false))
		length-=1
	}
	return arr
}

func (a *AMF0_Deserializer) ReadMixedArray() map[string]interface{} {
	a.Stream.ReadLong()
	return a.ReadObject() 
}

func (a *AMF0_Deserializer) ReadString() string {
	strLen,_ := a.Stream.ReadInt()
	b,_ := a.Stream.ReadBuffer(int(strLen))
	return b
}

func (a *AMF0_Deserializer) ReadLongString() string {
	strLen,_ := a.Stream.ReadLong()
	b,_ := a.Stream.ReadBuffer(int(strLen))
	return b
}

func (a *AMF0_Deserializer) ReadDate() time.Time {
	timestamp,_ := a.Stream.ReadDouble()
	timestamp = math.Floor(timestamp/1000)
	a.Stream.ReadInt() // timezone offset
	dateTime := time.Unix(int64(timestamp), 0)
	return dateTime
}
// TODO: those two

func (a *AMF0_Deserializer) ReadAMF3Data() interface{} {
	amf3Deserializer := NewAMF3_Deserializer(a.Stream,a.SavedRefs)
	data := amf3Deserializer.ReadAMFData(-1)
	a.SavedRefs = amf3Deserializer.GetReferences()
	return NewAMF3_Wrapper(data)
}

func (a *AMF0_Deserializer) ReadTypedObject() interface{} {
	classname := a.ReadString()
	// isMapped := false
	rObject := NewTypedObject(classname,nil)
	a.RefList = append(a.RefList,&rObject)
	props := map[string]interface{}{}
	for ;; {
		key := a.ReadString()
		vartype, _ := a.Stream.ReadByte()
		if vartype == AMF0Const()["DT_OBJECTTERM"] {
			break
		}
		props[key] = a.ReadAMFData(vartype,false)
	}
	rObject.SetAMFData(props)
	return rObject
}

func NewAMF0_Deserializer(stream *InputStream, savedRefs []interface{}) *AMF0_Deserializer {
	a := new(AMF0_Deserializer)
	a.Stream = stream
	a.SavedReferences = savedRefs
	return a
}