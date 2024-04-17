package cracking

import (
	"math"
	"math/big"
	"strings"
	"utils/subscriber"
)


type Eva struct {
	PublicKey		[]*big.Int
	EncryptedMsg 	[]*big.Int
	KeyLen			int
	MsgLen			int

	KeySub *subscriber.Subscriber
	MsgSub *subscriber.Subscriber
}

func CreateEva(url string, keyTopic, msgTopic string) (*Eva, error) {
	eva := &Eva{}

	var err error
	eva.KeySub, err = subscriber.NewSubscriber(url, keyTopic)
	if err != nil {
		return nil, err
	}

	eva.MsgSub, err = subscriber.NewSubscriber(url, msgTopic)
	if err != nil {
		return nil, err
	}

	return eva, nil
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

func hackHelper(encryptedWord *big.Int, Eva *Eva) []byte {
	var limit *big.Int = new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2),
		big.NewInt(int64(Eva.KeyLen)), nil), big.NewInt(1))
	var possibleSol *big.Int = big.NewInt(0)
	for ; possibleSol.Cmp(limit) <= 0; possibleSol.Add(possibleSol, big.NewInt(1)) {
		tmp := bigIntToBinary(new(big.Int).Set(possibleSol), Eva.KeyLen)
		weight := getWeight(Eva.PublicKey, tmp)
		if weight.Cmp(encryptedWord) == 0 {
			return tmp
		}
	}
	return nil
}

func (Eva *Eva) Hacking() string {
	var b strings.Builder
	for i := 0; i < len(Eva.EncryptedMsg); i++ {
		b.Write(hackHelper(Eva.EncryptedMsg[i], Eva))
	}
	binary := b.String()
	binary = binary[len(binary)-Eva.MsgLen*7:]

	var oldBinary []string = make([]string, Eva.MsgLen)
	for i := 0; i < Eva.MsgLen; i++ {
		oldBinary[i] = binary[i*7 : (i+1)*7]
	}

	var tmp []byte = make([]byte, Eva.MsgLen)
	for i := 0; i < Eva.MsgLen; i++ {
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
