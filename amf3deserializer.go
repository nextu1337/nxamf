package nxamf

import (
	"fmt"
	"time"
	"math"
)

type AMF3_Deserializer struct {
	objectCount int
	storedStrings []string
	storedObjects []interface{}
	storedClasses []map[string]interface{}
	refsLoaded bool
	Stream *InputStream
	savedReferences []interface{}
}

func (this *AMF3_Deserializer) ReadAMFData(settype int) interface{} {
	types := AMF3Const()
	if !this.refsLoaded && len(this.savedReferences)>0 {
		this.storedStrings = this.savedReferences[0].([]string)
		this.storedClasses = this.savedReferences[2].([]map[string]interface{})
		this.storedObjects = this.savedReferences[1].([]interface{})

		this.refsLoaded = true
	}
	if settype == -1 {
		settype,_ = this.Stream.ReadByte()
	}
	switch settype {
	case types["DT_UNDEFINED"],types["DT_NULL"]:
		return nil
	case types["DT_BOOL_FALSE"]:
		return false
	case types["DT_BOOL_TRUE"]:
		return true
	case types["DT_INTEGER"]:
		return this.ReadInt()
	case types["DT_NUMBER"]:
		num, _ := this.Stream.ReadDouble()
		return num
	case types["DT_STRING"],types["DT_XML"]:
		return this.ReadString()
	case types["DT_DATE"]:
		return this.ReadDate()
	case types["DT_ARRAY"]:
		return this.ReadArray()
	case types["DT_OBJECT"]:
		return this.ReadObject()
	case types["DT_XMLSTRING"]:
		return this.ReadXMLString()
	case types["DT_BYTEARRAY"]:
		return this.ReadByteArray()
	default:
		panic(fmt.Sprintf("Unsupported type: 0x%s",settype))
		return nil
	}
}

func (this *AMF3_Deserializer) GetReferences() []interface{} {
	return []interface{}{this.storedStrings,this.storedObjects,this.storedClasses}
}

func (this *AMF3_Deserializer) ReadObject() map[string]interface{} {
	objInfo := this.ReadU29()
	storedObject := false
	if (objInfo & 0x01) == 0 {
		storedObject = true;
	}
	objInfo = objInfo >> 1
	var rObject map[string]interface{}
	var encodingType int = -1
	propertyNames := []interface{}{}
	var className string = ""
	var properties map[string]interface{}
	if storedObject {
		// objectReference := objInfo (pointless)
		if objInfo > len(this.storedObjects) {
			panic(fmt.Sprintf("Object reference %v not found", objInfo))
			return map[string]interface{}{}
		} else {
			rObject = this.storedObjects[objInfo].(map[string]interface{})
		}
	} else {
		storedClass := false
		if (objInfo & 0x01) == 0 {
			storedClass = true
		}
		objInfo = objInfo >> 1
		if storedClass {
			// classReference := objInfo (pointless)
			if objInfo > len(this.storedClasses) {
				panic(fmt.Sprintf("Class reference %v not found", objInfo))
				return map[string]interface{}{}
			} else {
				encodingType = this.storedClasses[objInfo]["encodingType"].(int)
				propertyNames = this.storedClasses[objInfo]["propertyNames"].([]interface{})
				className = this.storedClasses[objInfo]["className"].(string)
			}
		} else {
			className = this.ReadString()
			encodingType = (objInfo & 0x03)
			propertyNames = []interface{}{}
			objInfo = objInfo >> 2
		}

		if className != "" {
			/* i dont like it
			if ($localClassName = $this->getLocalClassName($className)) {

                        $rObject = new $localClassName();

                    } else {

                        $rObject = new SabreAMF_TypedObject($className,array());

                    }
			*/
			rObject = map[string]interface{}{} // cba
		} else {
			rObject = map[string]interface{}{}
		}

		this.storedObjects = append(this.storedObjects,&rObject)
		if encodingType == AMF3Const()["ET_EXTERNALIZED"] {
			if !storedClass {
				this.storedClasses = append(this.storedClasses,map[string]interface{}{"className":className,"encodingType":encodingType,"propertyNames":propertyNames})
			}

			/*
			if ($rObject instanceof SabreAMF_Externalized) {
                        $rObject->readExternal($this->readAMFData());
                    } elseif ($rObject instanceof SabreAMF_TypedObject) {
                        $rObject->setAMFData(array('externalizedData'=>$this->readAMFData()));
			*/
			rObject["externalizedData"] = this.ReadAMFData(-1)
		} else {
			if encodingType == AMF3Const()["ET_SERIAL"] {
				if !storedClass {
					this.storedClasses = append(this.storedClasses,map[string]interface{}{"className":className,"encodingType":encodingType,"propertyNames":propertyNames})
				}
				properties = map[string]interface{}{}
				var propertyName string = " "
				for propertyName != "" {
					propertyName = this.ReadString()
					if propertyName != "" {
						propertyNames = append(propertyNames,propertyName)
						properties[propertyName] = this.ReadAMFData(-1)
					}
				}
			} else {
				if !storedClass {
					propertyCount := objInfo
					for i:=0;i<propertyCount;i++ {
						propertyNames = append(propertyNames,this.ReadString())
					}
					this.storedClasses = append(this.storedClasses,map[string]interface{}{"className":className,"encodingType":encodingType,"propertyNames":propertyNames})
				}
				properties = map[string]interface{}{}
				for _,propertyName := range propertyNames {
					properties[propertyName.(string)] = this.ReadAMFData(-1)
				}
				rObject = properties // idk
			}
			/* lol
			if ($rObject instanceof SabreAMF_TypedObject) {
                        $rObject->setAMFData($properties);
                    } else {
                        foreach($properties as $k=>$v) if ($k) $rObject->$k = $v;
                    }
					*/ 
		}

	}
	return rObject
}

func (this *AMF3_Deserializer) ReadString() string {
	strref := this.ReadU29()
	if (strref & 0x01) == 0 {
		strref = strref >> 1
		if strref>=len(this.storedStrings) {
			panic("Undefined string reference: "+fmt.Sprintf("%s",strref))
			return ""
		}
		return this.storedStrings[strref]
	} else {
		strlen := strref >> 1
		str,_ := this.Stream.ReadBuffer(strlen)
		if str != "" {
			this.storedStrings = append(this.storedStrings,str)
		}
		return str
	}
}

func (this *AMF3_Deserializer) ReadXMLString() string {
	strref := this.ReadU29()
	strlen := strref >> 1
	str,_ := this.Stream.ReadBuffer(strlen)
	return str
}

func (this *AMF3_Deserializer) ReadArray() map[string]interface{} {
	arrId := this.ReadU29()
	if (arrId & 0x01) == 0 {
		arrId = arrId >> 1
		if arrId >= len(this.storedObjects) {
			panic("Undefined array reference: "+fmt.Sprintf("%s",arrId))
			return map[string]interface{}{}
		}
		return this.storedObjects[arrId].(map[string]interface{})
	}
	arrId = arrId >> 1

	data := map[string]interface{}{}
	this.storedObjects = append(this.storedObjects,&data)
	key := this.ReadString()

	for key != "" {
		data[key] = this.ReadAMFData(-1)
		key = this.ReadString()
	}

	for i:=0;i<arrId;i++ {
		data[fmt.Sprintf("%v",len(data))] = this.ReadAMFData(-1)
	}
	
	return data
}

func (this *AMF3_Deserializer) ReadByteArray() *ByteArray {
	strref := this.ReadU29()
	strlen := strref >> 1
	str,_ := this.Stream.ReadBuffer(strlen)
	return NewByteArray(str)
}

func (this *AMF3_Deserializer) ReadU29() int {
	var count int = 0
	var u29 int = 0
	bytee,_ := this.Stream.ReadByte()
	for ((bytee & 0x80) != 0) && count < 4 {
		u29 = u29 << 7
		u29 = u29 | (bytee & 0x7f)
		bytee,_ = this.Stream.ReadByte()
		count+=1
	}

	if count < 4 {
		u29 = u29 << 7
		u29 = u29 | bytee
	} else {
		u29 = u29 << 8
		u29 = u29 | bytee
	}
	return u29
}

func (this *AMF3_Deserializer) ReadInt() int {
	integer := this.ReadU29()
	if (integer & 0x18000000) == 0x18000000 {
		integer = integer ^ 0x1fffffff
		integer *= -1
		integer -= 1
	} else if (integer & 0x10000000) == 0x10000000 {
		// remove the signed flag
		integer = integer & 0x0fffffff;
	}
	return integer
}



func (this *AMF3_Deserializer) ReadDate() time.Time {
	dateref := this.ReadU29()
	if (dateref & 0x01) == 0 {
		dateref = dateref >> 1
		if dateref>=len(this.storedObjects) {
			panic("Undefined date reference: "+fmt.Sprintf("%s",dateref))
			return time.Unix(0,0)
		}
		return (this.storedObjects[dateref]).(time.Time)
	}
	tme,_ := this.Stream.ReadDouble()
	timestamp := math.Floor(tme/1000)
	dateTime := time.Unix(int64(timestamp),0)
	return dateTime
}

func NewAMF3_Deserializer(stream *InputStream, savedRefs []interface{}) *AMF3_Deserializer {
	a := new(AMF3_Deserializer)
	a.Stream = stream
	a.savedReferences = savedRefs
	return a
}