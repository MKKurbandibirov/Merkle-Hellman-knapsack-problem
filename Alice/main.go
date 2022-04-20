package main

import (
	"cursach/alice"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	ln, err := net.Listen("tcp", ":4045")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Couldn't create the server!")
		os.Exit(1)
	}
	defer ln.Close()

	var Alissa *alice.T_Alissa = alice.CreateAlissa(0)

	for {
		conn, err := ln.Accept()
		if err != nil {
			conn.Close()
			fmt.Fprintln(os.Stderr, "[ERROR] Couldn't create connection!")
			os.Exit(1)
		}

		buf := make([]byte, 8192)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[ERROR] Reading error!")
			os.Exit(1)
		}

		if (string(buf[:n]) == "Eva") {
			var b1 strings.Builder
			for i := 0; i < Alissa.KeyLen; i++ {
				b1.WriteString(Alissa.PublicKey[i].String())
				b1.WriteByte(' ')
			}
			b1.WriteString("\n")
			var MsgLen = len(Alissa.EncryptedMsg)
			for i := 0; i < MsgLen; i++ {
				b1.WriteString(Alissa.EncryptedMsg[i].String())
				b1.WriteByte(' ')
			}
			b1.WriteString(strconv.FormatInt(int64(Alissa.MsgLen), 10))
			conn.Write([]byte(b1.String()))
		} else {
			file, err := os.Create("publicKey.txt")
			if err != nil {
				fmt.Fprintln(os.Stderr, "[ERROR] Couldn't create file!")
				os.Exit(1)
			}
			defer file.Close()
			fileKey := strings.ReplaceAll(string(buf[:n]), " ", "\n")
			fmt.Fprint(file, fileKey)
	
			fmt.Println("[OK] Public key recieved (saved in publicKey.txt)")
			var tmp []string = strings.Split(string(buf[:n]), " ")
			var Keylen int = len(tmp) - 1
			var PublicKey []*big.Int = make([]*big.Int, Keylen)
			for i := 0; i < Keylen; i++ {
				PublicKey[i] = new(big.Int)
				PublicKey[i].SetString(tmp[i], 10)
			}
			Alissa.KeyLen = Keylen
			Alissa.PublicKey = PublicKey

			Alissa.Encrypting()
			var MsgLen = len(Alissa.EncryptedMsg)
			var b strings.Builder
			for i := 0; i < MsgLen; i++ {
				b.WriteString(Alissa.EncryptedMsg[i].String())
				b.WriteByte(' ')
			}
			b.WriteString(strconv.FormatInt(int64(Alissa.MsgLen), 10))
			fmt.Println()
			conn.Write([]byte(b.String()))
		}
	}
}
