package main
import (
	ipfsutils "./ipfs_utils"
	"encoding/base64"
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strconv"
	"strings"
)

type ServerHandler struct{
}

// adapt HTTP connection to ReadWriteCloser
type HttpConn struct {
    in  io.Reader
    out io.Writer
}

func (c *HttpConn) Read(p []byte) (n int, err error)  { return c.in.Read(p) }
func (c *HttpConn) Write(d []byte) (n int, err error) { return c.out.Write(d) }
func (c *HttpConn) Close() error                      { return nil }

// our service
type HttpHandler struct{}

type PutDagResult struct {
	S string `json:"key"`
}

type GetDagResult struct {
	//Dag DagStruct `json:"dag"`
	Dag interface{} `json:"dag"`
}

type DagStruct struct {
	Data []uint `json:"data"`
}

func (h *HttpHandler) PutDag(data []byte, ret *PutDagResult) error {
	log.Printf("data is %v", data)
	encoded := base64.StdEncoding.EncodeToString(data)
	r := strings.NewReader(`{"data":"` + encoded + `"}`)
	key, err := ipfsutils.PutDag(r)
	ret.S = key
	return err
}

func (h *HttpHandler) PutDagBase64(data string, ret *PutDagResult) error {
	log.Printf("data is %v", data)
	r := strings.NewReader(`{"data":"` + data + `"}`)
	key, err := ipfsutils.PutDag(r)
	ret.S = key
	return err
}

func (h *HttpHandler) GetDag(key string, ret *GetDagResult) error {
	log.Printf("key is %+v\n\n", key)
	dag, err := ipfsutils.GetDag(key)
	if err != nil {
		return err
	}
	ret.Dag = dag
	defer func() {
		if r := recover(); r != nil {
			log.Printf("r is %v", r)
			return
		}
	}()
	if data := dag.(map[string]interface{})["data"]; data != nil {
		decodeBytes, err := base64.StdEncoding.DecodeString(data.(string))
		log.Printf("decodeBytes is %v", decodeBytes)
		if err != nil {
			return err
		}
		var ints []uint
		for _, b := range decodeBytes {
			ints = append(ints, uint(b))
		}
		dag.(map[string]interface{})["data"] = ints
	}
	return nil
}

func (h *HttpHandler) GetDagBase64(key string, ret *GetDagResult) error {
	log.Printf("key is %+v\n\n", key)
	dag, err := ipfsutils.GetDag(key)
	if err != nil {
		return err
	}
	ret.Dag = dag
	return nil
}

func main() {
	port := flag.Int("port", 8888, "listening port")
	ipfs_url := flag.String("ipfsurl", "127.0.0.1:5001", "ipfs http api url")
	ipfsutils.Init(*ipfs_url)
	flag.Parse()
	server := rpc.NewServer()
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal("server\t-", "listen error:", err.Error())
	}
	defer listener.Close()
	log.Println("server\t-", "start listion on port "+strconv.Itoa(*port))
	// 新建处理器
	httpHandler := &HttpHandler{}
	// 注册处理器
	server.Register(httpHandler)
	// 等待并处理链接
	go http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCodec := jsonrpc.NewServerCodec(&HttpConn{in: r.Body, out: w})
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(200)
		err := server.ServeRequest(serverCodec)
		if err != nil {
			log.Printf("Error while serving JSON request: %v", err)
			http.Error(w, "Error while serving JSON request, details have been logged.", 500)
			return
		}
	}))
	select {}
}
