package nxamf

import "time"
import "fmt"
import "reflect"


type AMF0_Serializer struct {
	Stream *OutputStream
}


func (a *AMF0_Serializer) WriteAMFData(data interface{}, forcetype int) {
	types := AMF0Const()
	// Recognize the type if it's not forced
	var typ int = -1
	if forcetype==-1 {
		switch data.(type) {
		case nil:
			typ = types["DT_NULL"]
			break
		case bool:
			typ = types["DT_BOOL"]
			break
		case float64,int:
			typ = types["DT_NUMBER"]
			break
		case string:
			if len(data.(string))>65536 {
				typ = types["DT_LONGSTRING"]
			} else {
				typ = types["DT_STRING"]
			}
			break
		case map[string]interface{}:
			typ = types["DT_OBJECT"]
			break
		case []interface{}:
			typ = types["DT_ARRAY"]
			break
		case map[int]interface{}:
			if !a.IsPureArray(data.(map[int]interface{})) {
				mS := make(map[string]interface{})
				for k,v := range data.(map[int]interface{}) {
					mS[fmt.Sprintf("%v",k)] = v
				}
				data = mS
				typ = types["DT_MIXEDARRAY"]
			} else {
				mS := []interface{}{}
				for _,v := range data.(map[int]interface{}) {
					mS = append(mS,v)
				}
				data = mS
				typ = types["DT_ARRAY"]
			}
			break
		case *AMF3_Wrapper:
			typ = types["DT_AMF3"]
			break
		case time.Time:
			typ = types["DT_DATE"]
			break
		default:
			if reflect.ValueOf(data).Kind() == reflect.Struct {
				//if a.GetRemoteClassName(data.(string)) != "" { // temporary
				if false == true {
					typ = types["DT_TYPEDOBJECT"]
					break
				} 
				// else {
				// 	typ = types["DT_OBJECT"]
				// 	break
				// }
			}
			
		}
	} else {
		typ = forcetype
	}
	a.Stream.WriteByte(typ)
	
	switch typ {
	case types["DT_NUMBER"]:
		switch i := data.(type) {
		case float64:
			a.Stream.WriteDouble(i)	
		case float32:
			a.Stream.WriteDouble(float64(i))
		case int:
			a.Stream.WriteDouble(float64(i))
		}
		
		break
	case types["DT_BOOL"]:
		c := 0
		if data.(bool) {
			c = 1
		}
		a.Stream.WriteByte(c)
		break
	case types["DT_STRING"]:
		a.WriteString(data.(string))
		break
	case types["DT_OBJECT"]:
		a.WriteObject(data.(map[string]interface{}))
		break
	case types["DT_NULL"]:
		break
	case types["DT_MIXEDARRAY"]:
		a.WriteMixedArray(data.(map[string]interface{}))
		break
	case types["DT_ARRAY"]:
		a.WriteArray(data.([]interface{}))
		break
	case types["DT_DATE"]:
		a.WriteDate(data.(time.Time))
		break
	case types["DT_LONGSTRING"]:
		a.WriteLongString(data.(string))
		break
	case types["DT_TYPEDOBJECT"]:
		a.WriteTypedObject(data)
		break
	case types["DT_AMF3"]:
		a.WriteAMF3Data(data.(*AMF3_Wrapper))
		break
	}
}

func (a *AMF0_Serializer) WriteMixedArray(data map[string]interface{}) {
	a.Stream.WriteLong(0)
	for k,v := range data {
		a.WriteString(string(k))
		a.WriteAMFData(v,-1)
	}
	a.WriteString("");
	a.Stream.WriteByte(AMF0Const()["DT_OBJECTTERM"])
}

func (a *AMF0_Serializer) WriteArray(data []interface{}) {
	if len(data)<1 {
		a.Stream.WriteLong(0)
	} else {
		a.Stream.WriteLong(uint64(len(data)))
		for i := range data {
			if data[i] != nil {
				a.WriteAMFData(data[i],-1)
			} else {
				a.Stream.WriteByte(AMF0Const()["DT_UNDEFINED"])
			}
		}
	}
}

func (a *AMF0_Serializer) WriteObject(data map[string]interface{}) {
    for i,v := range data {
        a.WriteString(i)
		a.WriteAMFData(v,-1)
    }
	a.WriteString("")
	a.Stream.WriteByte(AMF0Const()["DT_OBJECTTERM"])
}

func (a *AMF0_Serializer) WriteString(data string) {
	a.Stream.WriteInt(uint16(len(data)))
	a.Stream.WriteBuffer(data)
}

func (a *AMF0_Serializer) WriteLongString(data string) {
	a.Stream.WriteLong(uint64(len(data)))
	a.Stream.WriteBuffer(data)
}

func (a *AMF0_Serializer) WriteTypedObject(data interface{}) {
	// im lost
}

func (a *AMF0_Serializer) WriteAMF3Data(data *AMF3_Wrapper) {
	// TODO: Change interface to AMF3_Wrapper // done
	//  NewAMF3Serializer(a.Stream).WriteAMFData(data->getData())
}

func (a *AMF0_Serializer) WriteDate(data time.Time) {
	a.Stream.WriteDouble(float64(data.Unix()*1000))
	a.Stream.WriteInt(0)
}

// Checks whether array (or in this case map[int]interface{}) is sparse or not, for some reason doesn't work even if indexes are in a good order..
func (a *AMF0_Serializer) IsPureArray(array map[int]interface{}) bool {
	var i int = 0
	for k,_ := range array {
		if k!=i {
			return false
		}
		i+=1
	}
	return true
}

func (a *AMF0_Serializer) GetRemoteClassName(localClass string) string {
	return GetRemoteClass(localClass)
}

// public function to get OutputStream
func (a *AMF0_Serializer) GetStream() *OutputStream {
	return a.Stream
}

// Create object of AMF0 Serializer
func NewAMF0_Serializer(str *OutputStream) AMF0_Serializer {
	var a AMF0_Serializer
	a.Stream = str
	return a
}