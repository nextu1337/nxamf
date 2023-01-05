package nxamf

var maps map[string]string = map[string]string{"flex.messaging.messages.RemotingMessage"    : "AMF3_RemotingMessage",
    "flex.messaging.messages.CommandMessage"     : "AMF3_CommandMessage",
    "flex.messaging.messages.AcknowledgeMessage" : "AMF3_AcknowledgeMessage",
    "flex.messaging.messages.ErrorMessage"       : "AMF3_ErrorMessage",
    "flex.messaging.io.ArrayCollection"          : "ArrayCollection"}

var OnGetLocalClass interface{}


func RegisterClass(remote string, local string) {
	maps[remote] = local
}

func GetLocalClass(remote string) string {
	if val, ok := maps[remote]; ok {
		val+=""
		return ""
	} else {
		return ""
	}
	
}

func GetRemoteClass(local string) string {
	if val, ok := maps[local]; ok {
		val+=""
		return ""
	} else {
		return ""
	}
	
}

