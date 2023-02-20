# nxamf
nxamf is a port of SabreAMF in Go. It's an AMF client written in Go. Supports AMF0 and AMF3. <br>
##### Please make sure all the issues are library related! <br>
###### If any error is being encountered please create an entry in the `Issues` tab or make a pull request if you can fix it yourself
###### This project in no way is representing my coding skills. Every problem was solved the same way it is in `github.com/evert/SabreAMF` 

## Installation
Do `go install github.com/nextu1337/nxamf@latest`


### AMF
Action Message Format (AMF) is a binary format used to serialize object graphs such as ActionScript objects and XML, or send messages between an Adobe Flash client and a remote service, usually a Flash Media Server or third party alternatives.

## More on the project
This project was started purely out of boredom. I once looked around GitHub for AMF serialization libraries for Go and couldn't find any that would work in the way I wanted them to. Having experienced with evert's SabreAMF library in the past I decided that I should port it over to Go. For some reason I decided that porting has to mean exact copy in different language therefore the code looks pretty much the same (apart from the fact that of course it's been written in Go)

## AMF0
- Number `int,float64`
- Bool `bool`
- String `string`
- Object `map[string]interface{}`
- Null `nil`
- MixedArray `MixedArray` (which is an "alias" of map[string]interface{} but if you use it, it will be serialized as an object)
- Array `[]interface{}`
- Date `time.Time`
- AMF3 `AMF3_Wrapper`
- TypedObject `isn't supported when serializing but is supported when deserializing so TypedObject`
## AMF3
- Integer `int`
- Number `float64`
- Bool `bool`
- String `string`
- Date `time.Time`
- Array `[]interface{}`
- Object `map[string]interface{}`
- Associative Array `MixedArray`
- Null `nil`

### What is unsupported?
Things that haven't been ported from SabreAMF or not made in SabreAMF (that of course weren't made here either) include:
- DT_TYPEDOBJECT isn't fully supported, can't be serialized due to Go having weird way of recognizing private methods/fields from public.
- DT_MOVIECLIP
blah blah blah I don't think there is anything more missing (which there probably is but it's not as relevant)

## Example code
Due to Go being Go I decided to create functions for instantiation of the structs
That is `NewOutputStream()`, `NewMessage()`, `NewAMF3_Wrapper(data)` etc
```go
package main

import (
  "fmt"
  "github.com/nextu1337/nxamf" // Import the library
  "encoding/hex" // Needed only in this example for the hex dump
)

func main() {
  opt := nxamf.NewOutputStream() // Create new OutputStream
  msg := nxamf.NewMessage() // Create new Message
  data := []interface{}{"Hello!",123.456,false}
  msg.AddHeader(nxamf.NewHeader("name",false,"data")) // Add new header, false is the "required" field
  msg.AddBody(nxamf.NewBody("Target","Response",nxamf.NewAMF3_Wrapper(data))) // Add new body
  msg.Serialize(opt) // Serialize using the OutputStream
  // Code below is not needed
  fmt.Println(hex.EncodeToString([]byte(opt.GetRawData()))) // hex dump
}
```

```go
package main

import (
  nx "github.com/nextu1337/nxamf" // Import the library
  "fmt"
)

func main() {
  body := []interface{}{"anything can be put here"}
  c := nx.NewClient("url")
  // c.SetHttpProxy("http://127.0.0.1:8888") // allows setting a http proxy
  c.AddHeader("name",false,"data")
  resp := c.SendRequest("target",nx.NewAMF3_Wrapper(body)) // Response can be anything that is listed in the supported tab
  fmt.Println(resp)
}
```
