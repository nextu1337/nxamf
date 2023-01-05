package nxamf

func AMF0Const() map[string]int {
	return map[string]int {"DT_NUMBER":0x00,
	"DT_BOOL":0x01,
	"DT_STRING":0x02,
	"DT_OBJECT":0x03,
	"DT_MOVIECLIP":0x04,
	"DT_NULL":0x05,
	"DT_UNDEFINED":0x06,
	"DT_REFERENCE":0x07,
	"DT_MIXEDARRAY":0x08,
	"DT_OBJECTTERM":0x09,
	"DT_ARRAY":0x0a,
	"DT_DATE":0x0b,
	"DT_LONGSTRING":0x0c,
	"DT_UNSUPPORTED":0x0e,
	"DT_XML":0x0f,
	"DT_TYPEDOBJECT":0x10,
	"DT_AMF3":0x11}
}