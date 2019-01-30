package ipfsutils

import (
	shell "github.com/ipfs/go-ipfs-api"
)

var sh *shell.Shell

func init () {
	sh = shell.NewShell("localhost:5001")
}

func Init (url string) {
	sh = shell.NewShell(url)
}

// string, []byte, io.Reader
func PutDag(data interface{}) (key string, err error) {
	return sh.DagPut(data, "json", "protobuf")
	//return sh.DagPut(data, "json", "cbor")
}

func GetDag(ref string) (out interface{}, err error) {
	err = sh.DagGet(ref, &out)
	return
}

