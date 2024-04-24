package main

import (
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
	"utils/hacker"
)

func main() {
	Eva, err := hacker.NewHacker("nats://0.0.0.0:4222", "public_key", "messages")
	if err != nil {
		log.Fatal(err)
	}

	key, err := Eva.KeySub.GetMessage()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[INFO] Public key recieved!")

	var tmp1 []string = strings.Split(key, " ")
	var Keylen int = len(tmp1) - 1
	var PublicKey []*big.Int = make([]*big.Int, Keylen)
	for i := 0; i < Keylen; i++ {
		PublicKey[i] = new(big.Int)
		PublicKey[i].SetString(tmp1[i], 10)
	}

	Eva.KeyLen = Keylen
	Eva.PublicKey = PublicKey

	fmt.Println("[INFO] Hacking messages:")
	for {
		msg, err := Eva.MsgSub.GetMessage()
		if err != nil {
			log.Fatal(err)
		}

		var tmp []string = strings.Split(msg, " ")
		var MsgLen int = len(tmp) - 1
		Len, _ := strconv.ParseInt(tmp[MsgLen], 10, 64)
		Eva.MsgLen = int(Len)
		tmp = tmp[:MsgLen]

		fmt.Println("[INFO] Recieve encrypted message:")

		for i := 0; i < len(tmp); i++ {
			fmt.Println(tmp[i])
		}
		fmt.Println()

		var encryptedMsg []*big.Int = make([]*big.Int, MsgLen)
		for i := 0; i < MsgLen; i++ {
			encryptedMsg[i] = new(big.Int)
			encryptedMsg[i].SetString(tmp[i], 10)
		}

		fmt.Printf("[INFO] Message after hacking: \n>>> %s\n", Eva.Hacking(encryptedMsg))
	}
}
