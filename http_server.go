package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"rpc_demo/miliao"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	RPC_ADDR = "127.0.0.1:8082"
)

const (
	LOING_EVENT = 1
	CHAT_EVENT  = 2
)

type RpcServer struct{}

func (s *RpcServer) rpc(event int, i interface{}) (*miliao.RPCResult, error) {
	conn, err := grpc.Dial(RPC_ADDR, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	t := miliao.NewChatorClient(conn)

	var result miliao.RPCResult
	var rpcErr error

	switch event {
	case LOING_EVENT:
		user, _ := i.(miliao.User)
		res, _ := t.Login(context.Background(), &user)
		fmt.Printf("LOGIN:%s %s\n", user.GetName(), user.GetPassword())
		result = *res
	case CHAT_EVENT:
		msg, _ := i.(miliao.ChatMsg)
		res, _ := t.Chat(context.Background(), &msg)
		result = *res
		fmt.Printf("CHAT:%s %s\n", msg.GetUserId(), msg.GetMsg())
	default:
		break
	}

	if rpcErr != nil {
		return nil, err
	}

	fmt.Printf("RPC: code:%d err:%s token:%s", result.ErrorCode, result.Err, result.Result)

	return &result, nil
}

func httpRequest(w *http.ResponseWriter, event int, i interface{}) {
	var res miliao.RPCResult

	var server = &RpcServer{}
	result, err := server.rpc(event, i)

	if err == nil {
		res.ErrorCode = result.GetErrorCode()
		res.Err = result.GetErr()
		res.Result = result.GetResult()
	} else {
		res.Err = err.Error()
		res.ErrorCode = -1
		res.Result = ""
	}

	ret, _ := json.Marshal(res)
	fmt.Fprint(*w, string(ret))
}

// http login handle
func login(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var user miliao.User
	if err := json.Unmarshal(body, &user); err == nil {
		httpRequest(&w, LOING_EVENT, user)
	}
}

// http chat handle
func chat(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var msg miliao.ChatMsg
	if err := json.Unmarshal(body, &msg); err == nil {
		httpRequest(&w, CHAT_EVENT, msg)
	}
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/chat", chat)

	fmt.Println("------ http json server --------")

	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
