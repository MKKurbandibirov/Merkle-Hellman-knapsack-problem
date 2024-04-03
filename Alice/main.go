package main

import (
	"Alice/encrypt"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func main() {
	Alice, err := encrypt.CreateAlice("nats://0.0.0.0:4222", "public_key")
	if err != nil {
		log.Fatal(err)
	}

	buf, err := Alice.KeySub.GetMessage()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Reading error!")
		os.Exit(1)
	}

	// if string(buf) == "Eva" {
	// 	var b1 strings.Builder
	// 	for i := 0; i < Alice.KeyLen; i++ {
	// 		b1.WriteString(Alice.PublicKey[i].String())
	// 		b1.WriteByte(' ')
	// 	}
	// 	b1.WriteString("\n")
	// 	var MsgLen = len(Alice.EncryptedMsg)
	// 	for i := 0; i < MsgLen; i++ {
	// 		b1.WriteString(Alice.EncryptedMsg[i].String())
	// 		b1.WriteByte(' ')
	// 	}
	// 	b1.WriteString(strconv.FormatInt(int64(Alice.MsgLen), 10))

	// 	Alice.MsgPub.Publish("messages", b1.String())
	// } else {
	file, err := os.Create("publicKey.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Couldn't create file!")
		os.Exit(1)
	}
	defer file.Close()
	fileKey := strings.ReplaceAll(string(buf), " ", "\n")
	fmt.Fprint(file, fileKey)

	fmt.Println("[INFO] Public key recieved (saved in publicKey.txt)")
	var tmp []string = strings.Split(string(buf), " ")
	var Keylen int = len(tmp) - 1
	var PublicKey []*big.Int = make([]*big.Int, Keylen)
	for i := 0; i < Keylen; i++ {
		PublicKey[i] = new(big.Int)
		PublicKey[i].SetString(tmp[i], 10)
	}
	Alice.KeyLen = Keylen
	Alice.PublicKey = PublicKey

	for {
		Alice.Encrypting()
		var MsgLen = len(Alice.EncryptedMsg)
		var b strings.Builder
		for i := 0; i < MsgLen; i++ {
			b.WriteString(Alice.EncryptedMsg[i].String())
			b.WriteByte(' ')
		}
		b.WriteString(strconv.FormatInt(int64(Alice.MsgLen), 10))
	
		if err := Alice.MsgPub.Publish("messages", b.String()); err != nil {
			log.Fatal(err)
		}
	}
}
