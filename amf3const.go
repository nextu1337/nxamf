package nxamf

func AMF3Const() map[string]int {
	return map[string]int {"DT_UNDEFINED":0x00,
	"DT_NULL":0x01,
	"DT_BOOL_FALSE":0x02,
	"DT_BOOL_TRUE":0x03,
	"DT_INTEGER":0x04,
	"DT_NUMBER":0x05,
	"DT_STRING":0x06,
	"DT_XML":0x07,
	"DT_DATE":0x08,
	"DT_ARRAY":0x09,
	"DT_OBJECT":0x0a,
	"DT_XMLSTRING":0x0b,
	"DT_BYTEARRAY":0x0c,
	"ET_PROPLIST":0x00,
	"ET_EXTERNALIZED":0x01,
	"ET_SERIAL":0x02}
}