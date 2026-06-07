//go:generate sh -c "mkdir -p .bin && GOBIN=\"$(pwd)/.bin\" go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11 && protoc -I proto --go_out=./gen --go_opt=paths=source_relative --plugin=protoc-gen-go=\"$(pwd)/.bin/protoc-gen-go\" proto/base.proto proto/system.proto proto/disco.proto proto/report.proto proto/event.proto proto/config.proto proto/control.proto"
package main
