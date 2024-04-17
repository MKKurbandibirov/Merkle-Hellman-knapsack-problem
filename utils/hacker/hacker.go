package hacker

import (
	"math"
	"math/big"
	"strings"
	"utils/subscriber"
)

type Hacker struct {
	PublicKey    []*big.Int
	KeyLen       int
	MsgLen       int

	KeySub *subscriber.Subscriber
	MsgSub *subscriber.Subscriber
}

func NewHacker(url string, keyTopic, msgTopic string) (*Hacker, error) {
	h := &Hacker{}

	var err error
	h.KeySub, err = subscriber.NewSubscriber(url, keyTopic)
	if err != nil {
		return nil, err
	}

	h.MsgSub, err = subscriber.NewSubscriber(url, msgTopic)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func bigIntToBinary(num *big.Int, keyLen int) []byte {
	var possibleSol []byte = make([]byte, keyLen)
	for i := 0; i < keyLen; i++ {
		possibleSol[i] = '0'
	}
	for i := keyLen - 1; i >= 0; i-- {
		if new(big.Int).Mod(num, big.NewInt(2)).Cmp(big.NewInt(1)) == 0 {
			possibleSol[i] = '1'
		}
		num.Quo(num, big.NewInt(2))
	}
	return possibleSol
}

func getWeight(publicKey []*big.Int, tmp []byte) *big.Int {
	var keyLen int = len(publicKey)
	var weight = big.NewInt(0)
	for i := 0; i < keyLen; i++ {
		if tmp[i] == '1' {
			weight.Add(weight, publicKey[i])
		}
	}
	return weight
}

func (h *Hacker) hackHelper(encryptedWord *big.Int) []byte {
	var limit *big.Int = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2),
		big.NewInt(int64(h.KeyLen)), nil), big.NewInt(1))
	var possibleSol *big.Int = big.NewInt(0)
	for ; possibleSol.Cmp(limit) <= 0; possibleSol.Add(possibleSol, big.NewInt(1)) {
		tmp := bigIntToBinary(new(big.Int).Set(possibleSol), h.KeyLen)
		weight := getWeight(h.PublicKey, tmp)
		if weight.Cmp(encryptedWord) == 0 {
			return tmp
		}
	}
	return nil
}

func (h *Hacker) Hacking(encryptedMsg []*big.Int) string {
	var b strings.Builder
	for i := 0; i < len(encryptedMsg); i++ {
		b.Write(h.hackHelper(encryptedMsg[i]))
	}
	binary := b.String()
	binary = binary[len(binary)-h.MsgLen*7:]

	var oldBinary []string = make([]string, h.MsgLen)
	for i := 0; i < h.MsgLen; i++ {
		oldBinary[i] = binary[i*7 : (i+1)*7]
	}

	var tmp []byte = make([]byte, h.MsgLen)
	for i := 0; i < h.MsgLen; i++ {
		a := 0
		for j := 6; j >= 0; j-- {
			if oldBinary[i][j] == '1' {
				a += int(math.Pow(2, float64(6-j)))
			}
		}
		tmp[i] = byte(a)
	}
	decryptedMsg := string(tmp)
	return decryptedMsg
}
