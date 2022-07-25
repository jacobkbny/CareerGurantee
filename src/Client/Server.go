package main

// go version  go 1.18.4 window/amd64
//Restful API

import (
	"crypto/ecdsa"
	_ "encoding/json"
	"fmt"
	_ "net/http"
	"net/rpc"
)

type Args struct {
	A, B int
}

// Request Alias != nil -> Wallet 생성

// Request Tx != nil -> 블록 생성

// Request Tx == nil && Address != nil -> 블록 조회

type Request struct {
	Alias   string
	Address string
	T       *Transaction
}

type Transaction struct {
	Txid      [32]byte // 거래 ID
	TimeStamp []byte   // 거래시간
	Applier   []byte   // 요청자
	Company   []byte   // 경력회사
	Career    []byte   // 경력기간
	Payment   []byte   // 결제수단
	Job       []byte   // 직종 , 업무
	Proof     []byte   // pdf 링크
}
type Response struct {
	C          int
	Address    string
	publicKey  []byte
	privateKey ecdsa.PrivateKey
}

// WBS ( Work Based Schedule ) by Excel
// QS (Quality of Service ) by Excel
func main() {
	// Client, err := rpc.Dial("tcp", "127.0.0.1:9000")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer Client.Close() // 메인이 끝나기 직전에 실행되는 함수 (연결 해지)
	// args := &Args{1, 2}
	r := &Request{}
	response := new(Response)
	// Alias 가 비어있지 않다 -> Wallet 생성 요청
	if r.Alias != "" {
		Client, err := rpc.Dial("tcp", "127.0.0.1:9000")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer Client.Close()
		err = Client.Call("Return.SendWalletAddress", r.Alias, response)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(response.Address, r.Alias, "님의 지갑의 Address 입니다 ")
	}

	// alias := &Request{"test", [32]byte{}, ""}
	// err = Client.Call("Return.SendWalletAddress", alias, reply)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(reply.Address, alias, "님의 지갑의 Address 입니다 ")
	// args.A = 4
	// args.B = 9
	// SubCall := Client.Go("Calc.Subtract", args, reply, nil)
	// <-SubCall.Done
	// fmt.Println(reply.C, "you called Subtract Function")
	// MultiCall := Client.Go("Calc.Multiply", args, reply, nil)
	// <-MultiCall.Done
	// fmt.Println(reply.C, "you called Multiply Function")
	// DivideCall := Client.Go("Calc.Divide", args, reply, nil)
	// fmt.Println(reply.C, "you called Divide Function")
	// <-DivideCall.Done
	// fmt.Println(reply.C)
}
