package main

import (
	".."
	"encoding/base64"
	"fmt"
	"strings"
	"io"
)

func main () {
	// 添加data到ipfs，获取data的hash
	data := bigData(100)
	fmt.Printf("data is %v\n\n", data)
	key, err := ipfsutils.PutDag(data)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("key: %v\n\n", key)  //QmZF9436VS6D7YYe5toC1KHm23J9DbA343mt4eaei3kMG6

	// 通过hash获取原data
	dag, err := ipfsutils.GetDag(key)
	if err != nil {
		fmt.Printf("Error: %v\n\n", err.Error())
		return
	}
	fmt.Printf("%+v\n\n", dag)
	fmt.Printf("len is %v\n\n", len(dag.(map[string]interface{})))

}

func bigData(size int) io.Reader {
	if size > 1e9 {
		panic(fmt.Errorf("%v is too large", size))
	}
	var data [1e9]byte
	for i := 0; i < size; i ++ {
		data[i] = '@'
	}
	encoded := base64.StdEncoding.EncodeToString(data[:size])
	r := strings.NewReader(`{"data":"` + encoded + `"}`)
	//r := strings.NewReader(`{"lalala":"lalala"}`)
	return r
}
