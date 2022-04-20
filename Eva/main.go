package main

import (
	"cursach/eva"
	"fmt"
	"net"
	"os"
	"strings"
	"strconv"
	"math/big"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:4045")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Couldn't connect to server!")
		os.Exit(1)
	}
	defer conn.Close()

	var Eva *eva.T_Eva = eva.CreateEva()
	conn.Write([]byte("Eva"))

	buf := make([]byte, 8192)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Reading error!")
		os.Exit(1)
	}
	var KeyAndMsg []string = strings.Split(string(buf[:n]), "\n")

	// fmt.Println(KeyAndMsg[1])

	fmt.Println("[OK] Public key recieved!")
	var tmp1 []string = strings.Split(KeyAndMsg[0], " ")
	var Keylen int = len(tmp1) - 1
	var PublicKey []*big.Int = make([]*big.Int, Keylen)
	for i := 0; i < Keylen; i++ {
		PublicKey[i] = new(big.Int)
		PublicKey[i].SetString(tmp1[i], 10)
		fmt.Println(PublicKey[i])
	}
	fmt.Println()
	Eva.KeyLen = Keylen
	Eva.PublicKey = PublicKey

	var tmp []string = strings.Split(KeyAndMsg[1], " ")
	var MsgLen int = len(tmp) - 1
	Len, _ := strconv.ParseInt(tmp[MsgLen], 10, 64)
	Eva.MsgLen = int(Len)
	tmp = tmp[:MsgLen]

	fmt.Println("[OK] Recieve encrypted message:")
	for i := 0; i < len(tmp); i++ {
		fmt.Println(tmp[i])
	}
	fmt.Println()

	var EncryptedMsg []*big.Int = make([]*big.Int, MsgLen)
	for i := 0; i < MsgLen; i++ {
		EncryptedMsg[i] = new(big.Int)
		EncryptedMsg[i].SetString(tmp[i], 10)
	}
	Eva.EncryptedMsg = EncryptedMsg

	fmt.Println(Eva.Hacking())
}