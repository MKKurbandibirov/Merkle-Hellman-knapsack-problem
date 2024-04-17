package main

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"utils/encryptor"
)

func main() {
	Alice, err := encryptor.NewEncryptor("nats://0.0.0.0:4222", "public_key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[INFO] Alice is waiting a key!")

	buf, err := Alice.KeySub.GetMessage()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Reading error!")
		os.Exit(1)
	}

	file, err := os.Create("publicKey.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Couldn't create file!")
		os.Exit(1)
	}
	defer file.Close()

	fileKey := strings.ReplaceAll(string(buf), " ", "\n")
	fmt.Fprint(file, fileKey)

	fmt.Println("[INFO] Public key recieved (saved in publicKey.txt)")

	tmp := strings.Split(string(buf), " ")
	keylen := len(tmp) - 1
	publicKey := make([]*big.Int, keylen)
	for i := 0; i < keylen; i++ {
		publicKey[i] = new(big.Int)
		publicKey[i].SetString(tmp[i], 10)
	}
	Alice.KeyLen = keylen
	Alice.PublicKey = publicKey

	for {
		var (
			encryptedMsg = Alice.Encrypting()

			b strings.Builder
		)

		for i := 0; i < len(encryptedMsg); i++ {
			b.WriteString(encryptedMsg[i].String())
			b.WriteByte(' ')
		}
		b.WriteString(strconv.FormatInt(int64(Alice.MsgLen), 10))

		if err := Alice.MsgPub.Publish("messages", b.String()); err != nil {
			log.Fatal(err)
		}
	}
}
