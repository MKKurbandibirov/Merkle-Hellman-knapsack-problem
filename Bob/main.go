package main

import (
	"cursach/bob"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:4045")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Couldn't connect to server!")
		os.Exit(1)
	}
	defer conn.Close()

	// Отправка ключа
	var Bob *bob.T_Bob = bob.CreateBob()
	file, err := os.Create("privateKey.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Couldn't create file!")
		os.Exit(1)
	}
	defer file.Close()
	for i := 0; i < Bob.KeyLen; i++ {
		fmt.Fprintln(file, Bob.PrivateKey[i].String())
	}

	var b strings.Builder
	for i := 0; i < Bob.KeyLen; i++ {
		b.WriteString(Bob.PublicKey[i].String())
		b.WriteByte(' ')
	}
	conn.Write([]byte(b.String()))

	buf := make([]byte, 8192)
	n, err := conn.Read(buf)
	if n == 0 || err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Reading error!")
		os.Exit(1)
	}
	var tmp []string = strings.Split(string(buf[:n]), " ")
	var MsgLen int = len(tmp) - 1
	Len, _ := strconv.ParseInt(tmp[MsgLen], 10, 64)
	Bob.MsgLen = int(Len)
	tmp = tmp[:MsgLen]

	fmt.Println("[OK] Recieve encrypted message:")
	for i := 0; i < len(tmp); i++ {
		fmt.Println(tmp[i])
	}
	fmt.Println()

	var CryptedMsg []*big.Int = make([]*big.Int, MsgLen)
	for i := 0; i < MsgLen; i++ {
		CryptedMsg[i] = new(big.Int)
		CryptedMsg[i].SetString(tmp[i], 10)
	}
	Bob.CryptedMsg = CryptedMsg
	decryptedMsg := Bob.Decrypting()
	fmt.Println(decryptedMsg)
}