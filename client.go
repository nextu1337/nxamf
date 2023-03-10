package nxamf

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

type Client struct {
	endPoint string
	httpProxy string
	amfInputStream *InputStream
	amfOutputStream *OutputStream
	amfRequest *Message
	amfResponse *Message
	encoding int
	httpHeaders []interface{}
	proxySet bool
}

func (this *Client) SendRequest(servicePath string, data interface{}) (interface{},error) {
	// if this.encoding & Const()["FLEXMSG"] == 1 {
		/*
		if($this->encoding & SabreAMF_Const::FLEXMSG) {

				XD we're not doing this LOLOLO LLMFAOOOO AHAHAHAH XDDD
                // Setting up the message
                $message = new SabreAMF_AMF3_RemotingMessage();
                $message->body = $data;

                // We need to split serviceName.methodName into separate variables
                $service = explode('.',$servicePath);
                $method = array_pop($service);
                $service = implode('.',$service);
                $message->operation = $method; 
                $message->source = $service;

                $data = $message;
            }
			*/
	// }
	// target := "null"
	// if this.encoding & Const()["FLEXMSG"] == 1 {
	// 	target = servicePath
	// }
	target := servicePath
	this.amfRequest.AddBody(NewBody(target,"/1",data))

	this.amfRequest.Serialize(this.amfOutputStream)

	headers := append(this.httpHeaders, "Content-Type: "+GetMIMETYPE())

	if len(this.httpProxy) > 0 {
		os.Setenv("HTTP_PROXY", this.httpProxy)
		this.proxySet = true
	} else {
		if this.proxySet {
			os.Unsetenv("HTTP_PROXY")
			this.proxySet = false
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST",this.endPoint,strings.NewReader(this.amfOutputStream.GetRawData()))
	if err != nil {
		return nil,err
	}
	
	for _,v := range headers {
		s := strings.Split(v.(string),":")
		req.Header.Add(s[0],strings.Join(s[1:],":"))
	}
	resp,err := client.Do(req)
	if err != nil {
		return nil,err
	}

	switch resp.StatusCode {
	case 200:
		break;
	default:
		return nil,errors.New(resp.Status)
	}

	bodyBytes,_ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	this.amfInputStream = NewInputStream(string(bodyBytes))
	this.amfResponse = NewMessage()
	this.amfResponse.Deserialize(this.amfInputStream)

	this.parseHeaders()

	

	for _, body := range this.amfResponse.GetBodies().([]*Body) {
		if body.Target[:2] == "/1" {
			return body.Data,nil
		}
	}
	return nil,errors.New("Something went wrong, this should not happen")
}

func (this *Client) AddHTTPHeader(header string) {
	this.httpHeaders = append(this.httpHeaders, header)
}

func (this *Client) AddHeader(name string, required bool, data interface{}) {
	this.amfRequest.AddHeader(NewHeader(name,required,data))
}

func (this *Client) SetCredentials(username string, password string) {
	this.AddHeader("Credentials",false,map[string]interface{}{"userid":username,"password":password})
}

func (this *Client) SetHttpProxy(httpProxy string) {
	this.httpProxy = httpProxy
}

func (this *Client) SetEncoding(encoding int) {
	this.encoding = encoding
	this.amfRequest.SetEncoding((encoding & Const()["AMF3"]))
}

func (this *Client) parseHeaders() {
	for _, header := range this.amfResponse.GetHeaders().([]*Header) {
		switch header.Name {
		case "ReplaceGatewayUrl":
			switch v := header.Data.(type) {
			case string:
				this.endPoint = v
				break
			}
			break
		}
	}
}

func NewClient(endPoint string) *Client {
	client := new(Client)
	client.endPoint = endPoint
	client.amfRequest = NewMessage()
	client.amfOutputStream = NewOutputStream()
	client.proxySet = false
	return client
}
