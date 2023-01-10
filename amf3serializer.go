package nxamf

import (
	"time"
)

type AMF3_Serializer struct {
	Stream *OutputStream
}

func (a *AMF3_Serializer) WriteAMFData(data interface{}, forcetype int) {
	types := AMF3Const()
	// Recognize the type if it's not forced
	var typ int = -1
	if forcetype==-1 {
		switch data.(type) {
		case nil:
			typ = types["DT_NULL"]
			break
		case bool:
			if data.(bool) {
				typ = types["DT_BOOL_TRUE"]
			} else {
				typ = types["DT_BOOL_FALSE"]
			}
			break
		case int:
			if data.(int) > 0xFFFFFFF || data.(int) < -268435456 {
				typ = types["DT_NUMBER"]
			} else {
				typ = types["DT_INTEGER"]
			}
			break
		case float64:
			typ = types["DT_NUMBER"]
			break
		case string:
			typ = types["DT_STRING"]
			break
		case []interface{},MixedArray:
			typ = types["DT_ARRAY"]
			break
		case *ByteArray:
			typ = types["DT_BYTEARRAY"]
			break
		case time.Time:
			typ = types["DT_DATE"]
			break
		case map[string]interface{}:
			typ = types["DT_OBJECT"]
			break
		}
		if typ == types["DT_INTEGER"] && (data.(int) > 268435455||data.(int) < -268435456) {
			typ = types["DT_NUMBER"]
		}
	} else {
		typ = forcetype
	}

	a.Stream.WriteByte(typ)

	switch typ {
	case types["DT_NULL"],types["DT_BOOL_FALSE"],types["DT_BOOL_TRUE"]:
		break
	case types["DT_INTEGER"]:
		a.WriteInt(data.(int))
		break
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
	case types["DT_STRING"]:
		a.WriteString(data.(string))
		break
	case types["DT_DATE"]:
		a.WriteDate(data.(time.Time))
		break
	case types["DT_ARRAY"]:
		switch i := data.(type) {
		case MixedArray:
			a.WriteMixedArray(i)
			break
		case []interface{}:
			a.WriteArray(i)
			break
		}
		break
	case types["DT_OBJECT"]:
		a.WriteObject(data.(map[string]interface{}))
		break
	case types["DT_BYTEARRAY"]:
		a.WriteByteArray(data.(*ByteArray))
		break
	default:
		panic("Type unsupported")
	}
}

func (a *AMF3_Serializer) WriteObject(data map[string]interface{}) {
	encoding := AMF3Const()["ET_PROPLIST"]
	// code below was not ported (wouldn't make sense if it was anyways)
	/*
	if ($data instanceof SabreAMF_ITypedObject) {

                $classname = $data->getAMFClassName();
                $data = $data->getAMFData();

            } else if (!$classname = $this->getRemoteClassName(get_class($data))) {

                
                $classname = '';

            } else {

                if ($data instanceof SabreAMF_Externalized) {

                    $encodingType = SabreAMF_AMF3_Const::ET_EXTERNALIZED;

                }

            }
	*/
	objectInfo := 0x03
	objectInfo = objectInfo | (encoding << 2)

	switch encoding {
	case AMF3Const()["ET_PROPLIST"]:
		propertyCount := 0
		for range data {
			propertyCount+=1
		}

		objectInfo = objectInfo | (propertyCount << 4)

		a.WriteInt(objectInfo)
		a.WriteString("") // should be classname but i don't support that
		for k,_ := range data {
			a.WriteString(k)
		}
		for _,v := range data {
			a.WriteAMFData(v,-1)
		}
		break
	}
}

func (a *AMF3_Serializer) WriteArray(arr []interface{}) {
	length := len(arr)
	id := (length << 1) | 0x01

	a.WriteInt(id)
	a.WriteString("")

	for _, v := range arr {
		a.WriteAMFData(v,-1)
	}
}

func (a *AMF3_Serializer) WriteMixedArray(arr MixedArray) {
	a.WriteInt(1)
	for k, v := range arr {
		a.WriteString(k)
		a.WriteAMFData(v,-1)
	}
	a.WriteString("")
}

func (a *AMF3_Serializer) WriteInt(i int) {
	if (i & 0xffffff80) == 0 {
		a.Stream.WriteByte((i & 0x7f))
		return
	}
	if (i & 0xffffc000) == 0 {
		a.Stream.WriteByte(((i>>7)|0x80))
		a.Stream.WriteByte((i & 0x7f))
		return
	}
	if (i & 0xffe00000) == 0 {
		a.Stream.WriteByte(((i>>14)|0x80))
		a.Stream.WriteByte(((i>>7)|0x80))
		a.Stream.WriteByte((i & 0x7f))
		return
	}
	a.Stream.WriteByte(((i>>22)|0x80))
	a.Stream.WriteByte(((i>>15)|0x80))
	a.Stream.WriteByte(((i>>8)|0x80))
	a.Stream.WriteByte((i & 0xff))
	return
}

func (a *AMF3_Serializer) WriteByteArray(bytearr *ByteArray) {
	a.WriteString(bytearr.GetData())
}

func (a *AMF3_Serializer) WriteString(str string) {
	strref := len(str) << 1 | 0x01
	a.WriteInt(strref)
	a.Stream.WriteBuffer(str)
}

func (a *AMF3_Serializer) WriteDate(data time.Time) {
	a.WriteInt(0x01)
	a.Stream.WriteDouble(float64(data.Unix()*1000))
}

func (a *AMF3_Serializer) GetStream() *OutputStream {
	return a.Stream
}

func NewAMF3_Serializer(str *OutputStream) *AMF3_Serializer {
	a := new(AMF3_Serializer)
	a.Stream = str
	return a
}