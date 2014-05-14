# pushd-client-go

Client library for pushd, written in golang.

###### Usage
For more details see: http://godoc.org/github.com/catalyst-zero/pushd-client-go
```go
import pushdClientPkg "github.com/catalyst-zero/pushd-client-go"

pushdClient :=  pushdClientPkg.NewHttpClient("127.0.0.1:8080")

subscriber, err := pushdClient.V1.SubscribeDevice("gcm", "weiufb34dhinf...", "de")
```
