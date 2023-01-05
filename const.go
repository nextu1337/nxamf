package nxamf

func Const() map[string]int {
	return map[string]int {"AC_Flash":0,
	"AMF0":0,
	"R_DEBUG":3,
	"R_STATUS":2,
	"R_RESULT":1,
	"AC_Flash9":3,
	"AC_FlashCom":1,
	"AMF3":3,
	"FLEXMSG":16}
}

func GetMIMETYPE() string {
	return "application/x-amf"
}