package main

import (
	"Bob/decrypt"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func main() {
	Bob, err := decrypt.CreateBob("nats://0.0.0.0:4222", "messages")
	if err != nil {
		log.Fatal(err)
	}

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

	if err := Bob.Pub.Publish("public_key", b.String()); err != nil {
		log.Fatal(err)
	}

	for {
		buf, err := Bob.MsgSub.GetMessage()
		fmt.Println(buf)
		if len(buf) == 0 || err != nil {
			fmt.Fprintln(os.Stderr, "[ERROR] While getting message: %w", err)
			os.Exit(1)
		}

		var tmp []string = strings.Split(string(buf), " ")
		var MsgLen int = len(tmp) - 1
		Len, _ := strconv.ParseInt(tmp[MsgLen], 10, 64)
		Bob.MsgLen = int(Len)
		tmp = tmp[:MsgLen]

		fmt.Println("[INFO] Recieve encrypted message:")
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

		fmt.Printf("[INFO] Message after decrypting: \n%s\n", decryptedMsg)
	}
}
