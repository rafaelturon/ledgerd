### Download protoc
http://google.github.io/proto-lens/installing-protoc.html

### Install protoc-gen-gi
https://grpc.io/docs/languages/go/quickstart/
1. Install the protocol compiler plugin for Go (protoc-gen-go) using the following command:

        $ go get github.com/golang/protobuf/protoc-gen-go
2. Update your PATH so that the protoc compiler can find the plugin:

        $ export PATH="$PATH:$(go env GOPATH)/bin"

 * Alternative for MacOS using Homebrew:
   - brew install protobuf
   - brew install protoc-gen-go

### Compile proto files 
protoc --proto_path=pb pb/*.proto --go_out=plugins=grpc:pb --go_opt=paths=source_relative
