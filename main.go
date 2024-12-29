package main

import (
	"github.com/etcd-io/auger/pkg/encoding"

	"github.com/etcd-io/auger/pkg/scheme"

	"encoding/base64"

	"syscall/js"
)

func stripNewline(d []byte) []byte {
	if len(d) > 0 && d[len(d)-1] == '\n' {
		return d[:len(d)-1]
	}
	return d
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("decodeEtcdData", js.FuncOf(decodeEtcdData))
	js.Global().Set("encodeEtcdData", js.FuncOf(encodeEtcdData))
	<-c
}

func decodeEtcdData(this js.Value, args []js.Value) interface{} {
	// Get the base64 string from the arguments
	var base64str = args[0].String()

	// Decode Base64 to binary (byte slice)
	binaryData, err := base64.StdEncoding.DecodeString(base64str)
	if err != nil {
		return err.Error()
	}

	// etcdctl appends a new line automatically, so strip it before processing
	binaryData = stripNewline(binaryData)

	// Pass the binary data to auger to detect it's type
	inMediaType, in, err := encoding.DetectAndExtract(binaryData)
	if err != nil {
		return err.Error()
	}

	// if metaOnly {
	// 	return encoding.DecodeSummary(inMediaType, in, out)
	// }

	// Convert the data to yaml
	buf, _, err := encoding.Convert(scheme.Codecs, inMediaType, "application/yaml", in)
	if err != nil {
		return err.Error()
	}

	return js.ValueOf(string(buf))
}

func encodeEtcdData(this js.Value, args []js.Value) interface{} {
	var in = []byte(args[0].String())
	inMediaType, err := encoding.ToMediaType("yaml")
	if err != nil {
		return err.Error()
	}
	buf, _, err := encoding.Convert(scheme.Codecs, inMediaType, encoding.StorageBinaryMediaType, in)
	if err != nil {
		return err.Error()
	}

	var base64data = base64.StdEncoding.EncodeToString(buf)

	return js.ValueOf(base64data)
}