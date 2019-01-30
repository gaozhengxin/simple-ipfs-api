package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
	url := "http://127.0.0.1:8888"

	//post := `{"method":"HttpHandler.PutDag","params":[[66,66,66,66]],"id":1}}`
	post := `{"method":"HttpHandler.GetDag","params":["QmV76vLiJCR79Rg8rsrPMQt5keok8NWuB2xdisNVknw3aP"],"id":1}`

	var jsonStr = []byte(post)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
