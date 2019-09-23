/**
 * Copyright (c) XDE, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

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

// RPC服务器地址配置
const (
	RPCAddr = "127.0.0.1:8082" // rpc address
)

// 事件代号代码
const (
	LoginEvent = 1 // 登录事件
	chatEvent  = 2 // 聊天事件
)

//RPCServer 用于调用rpc服务
type RPCServer struct{}

// rpc服务，根据事件调用不同的rpc
func (s *RPCServer) rpc(event int, i interface{}) (*miliao.RPCResult, error) {
	conn, err := grpc.Dial(RPCAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	t := miliao.NewChatorClient(conn)

	var result miliao.RPCResult
	var rpcErr error

	switch event {
	case LoginEvent:
		user, _ := i.(miliao.User)
		res, _ := t.Login(context.Background(), &user)
		fmt.Printf("LOGIN:%s %s\n", user.GetName(), user.GetPassword())
		result = *res
	case chatEvent:
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

// http请求处理，根据event进行rpc路由
func httpRequest(w *http.ResponseWriter, event int, i interface{}) {
	var res miliao.RPCResult

	var server = &RPCServer{}
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

// 登录handle
func login(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var user miliao.User
	if err := json.Unmarshal(body, &user); err == nil {
		httpRequest(&w, LoginEvent, user)
	}
}

// 聊天handle
func chat(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var msg miliao.ChatMsg
	if err := json.Unmarshal(body, &msg); err == nil {
		httpRequest(&w, chatEvent, msg)
	}
}

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/chat", chat)

	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
