package main

import (
	"Eva/cracking"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"strings"
)

func main() {
	Eva, err := cracking.CreateEva("nats://0.0.0.0:4222", "public_key", "messages")

	key, err := Eva.KeySub.GetMessage()
	if err != nil {
		log.Fatal(err)
	}

	msg, err := Eva.MsgSub.GetMessage()
	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(KeyAndMsg[1])

	fmt.Println("[OK] Public key recieved!")
	var tmp1 []string = strings.Split(key, " ")
	var Keylen int = len(tmp1) - 1
	var PublicKey []*big.Int = make([]*big.Int, Keylen)
	for i := 0; i < Keylen; i++ {
		PublicKey[i] = new(big.Int)
		PublicKey[i].SetString(tmp1[i], 10)
		fmt.Println(PublicKey[i])
	}
	fmt.Println(key, msg)
	Eva.KeyLen = Keylen
	Eva.PublicKey = PublicKey

	var tmp []string = strings.Split(msg, " ")
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
