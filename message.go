package nxamf
// TODO: Add "Deserialize" function

type Header struct {
	Name string
	Required int
	Data interface{}
}

type Body struct {
	Target string
	Response string
	Data interface{}
}

type Message struct {
	ClientType int
	Bodies []*Body
	Headers []*Header
	Encoding int
}

func (a *Message) Serialize(stream *OutputStream) {
	stream.WriteByte(0x00)
	stream.WriteByte(a.Encoding)
	stream.WriteInt(uint16(len(a.Headers)))

	for _, header := range a.Headers {
		serializer := NewAMF0_Serializer(stream)
		serializer.WriteString(header.Name)
		stream.WriteByte(header.Required)
		stream.WriteLongMinus1()
		serializer.WriteAMFData(header.Data,-1)
	}

	stream.WriteInt(uint16(len(a.Bodies)))

	for _, body := range a.Bodies {
		serializer := NewAMF0_Serializer(stream)
		serializer.WriteString(body.Target)
		serializer.WriteString(body.Response)
		stream.WriteLongMinus1()
		
		switch a.Encoding {
		case Const()["AMF0"]:
			serializer.WriteAMFData(body.Data,-1)
			break
		case Const()["AMF3"]:
			serializer.WriteAMF3Data(NewAMF3_Wrapper(body.Data))
			break
		
		}
		
	}
}

func (a *Message) Deserialize(stream *InputStream) {
	a.Headers = []*Header{}
	a.Bodies = []*Body{}
	stream.ReadByte()
	a.ClientType,_ = stream.ReadByte()
	deserializer := NewAMF0_Deserializer(stream,[]interface{}{})
	totalHeaders,_ := stream.ReadInt()

	for i:=0;i<int(totalHeaders);i++ {
		str := deserializer.ReadString()
		odbyt,_ := stream.ReadByte()
		abc := false
		if odbyt > 0 {
			abc = true
		}
		stream.ReadLong()
		data := deserializer.ReadAMFData(-1,true)
		header := NewHeader(str,abc,data)
		a.Headers = append(a.Headers,header)
	}

	totalBodies,_ := stream.ReadInt()

	for i:=0;i<int(totalBodies);i++ {
		target := deserializer.ReadString()
		response := deserializer.ReadString()
		stream.ReadLong()
		data := deserializer.ReadAMFData(-1,true)
		body := NewBody(target,response,data)
		
		aaaaa := 0

		switch data.(type) {
		case *AMF3_Wrapper:
			aaaaa = 1
		case []interface{}:
			aaaaa = 2
		}
		
		if aaaaa == 1 {
			body.Data = (data.(*AMF3_Wrapper)).GetData()
			a.Encoding = Const()["AMF3"]
		} else if aaaaa==2 {
			/*
			if (!defined("SABREAMF_AMF3_PRESERVE_ARGUMENTS")) {
                        $body['data'] = $body['data'][0]->getData();
                    } else {
						*/
			i := 0 
			for i < len(data.([]interface{})) {
				switch data.(type) {
				case *AMF3_Wrapper:
					body.Data.([]interface{})[i] = (body.Data.([]interface{})[i].(*AMF3_Wrapper)).GetData()
				}
				i+=1
			}
			a.Encoding = Const()["AMF3"]
		}

		a.Bodies = append(a.Bodies,body)
	}
}

func (a *Message) GetClientType() int {
	return a.ClientType
}

func (a *Message) GetBodies() interface{} {
	return a.Bodies
}

func (a *Message) GetHeaders() interface{} {
	return a.Headers
}

func (a *Message) GetEncoding() int {
	return a.Encoding
}

func (a *Message) SetEncoding(encoding int) {
	a.Encoding = encoding
}

func (a *Message) AddBody(body *Body) {
	a.Bodies = append(a.Bodies,body)
}

func (a *Message) AddHeader(header *Header) {
	a.Headers = append(a.Headers,header)
}

func NewHeader(name string, required bool, data interface{}) *Header {
	a := new(Header)
	a.Name = name
	a.Required = 0
	if required { 
		a.Required = 1
	}
	a.Data = data
	return a
}

func NewBody(target string, response string, data interface{}) *Body {
	a := new(Body)
	a.Target = target
	a.Response = response
	a.Data = data
	return a
}

func NewMessage() *Message {
	a := new(Message)
	a.ClientType=0
	a.Encoding = Const()["AMF0"]
	return a
}