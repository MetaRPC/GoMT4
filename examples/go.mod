module github.com/MetaRPC/GoMT4

go 1.23.6

require (
	git.mtapi.io/root/mrpc-proto.git/mt4/libraries/go v0.0.0-20250801133633-34bb1da6e4e5
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.74.2
	google.golang.org/protobuf v1.36.6   
)

require (
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250728155136-f173205681a0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250728155136-f173205681a0 // indirect
)

replace github.com/MetaRPC/GoMT4 => ../
