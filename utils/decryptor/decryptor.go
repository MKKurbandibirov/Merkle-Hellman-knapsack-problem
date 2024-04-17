package decryptor

import (
	"flag"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	"utils/publisher"
	"utils/subscriber"
)

type Decryptor struct {
	PrivateKey []*big.Int
	PublicKey  []*big.Int
	KeyLen     int
	MsgLen     int
	Q          *big.Int
	R          *big.Int

	MsgSub *subscriber.Subscriber
	Pub    *publisher.Publisher
}

func (dec *Decryptor) KeyGen(keyLen int) {
	dec.KeyLen = keyLen
	dec.PrivateKey = make([]*big.Int, keyLen)
	var sum, rnd *big.Int = new(big.Int), new(big.Int)
	for i := 0; i < keyLen; i++ {
		rnd.Rand(rand.New(rand.NewSource(time.Now().UnixNano())), big.NewInt(int64(100)))
		sum.Mul(sum, big.NewInt(int64(2)))
		sum.Add(sum, rnd)
		dec.PrivateKey[i] = new(big.Int)
		dec.PrivateKey[i].Set(sum)
	}
}

func modexp(a, t, n *big.Int) *big.Int {
	if t.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(1)
	}
	z := modexp(a, new(big.Int).Quo(t, big.NewInt(2)), n)
	if new(big.Int).Mod(t, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return new(big.Int).Mod(new(big.Int).Mul(z, z), n)
	} else {
		return new(big.Int).Mod(new(big.Int).Mul(new(big.Int).Mul(z, z), a), n)
	}
}

func millerRabinTest(num *big.Int, k int) bool {
	if num.Cmp(big.NewInt(2)) == 0 || num.Cmp(big.NewInt(3)) == 0 {
		return true
	}
	if num.Cmp(big.NewInt(2)) == -1 || new(big.Int).Mod(num, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}
	
	var t *big.Int = new(big.Int).Sub(num, big.NewInt(1))
	var s int = 0
	for new(big.Int).Mod(t, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		t.Quo(t, big.NewInt(2))
		s++
	}

	var rnd *big.Int = new(big.Int)
	for i := 0; i < k; i++ {
		rnd.Rand(rand.New(rand.NewSource(time.Now().UnixNano())), new(big.Int).Sub(num, big.NewInt(4)))
		rnd.Add(rnd, big.NewInt(2))
		var x *big.Int = modexp(rnd, t, num)
		if x.Cmp(big.NewInt(1)) == 0 || x.Cmp(new(big.Int).Sub(num, big.NewInt(1))) == 0 {
			continue
		}
		for j := 1; j < s; j++ {
			x = modexp(x, big.NewInt(2), num)
			if x.Cmp(big.NewInt(1)) == 0 {
				return false
			}
			if x.Cmp(new(big.Int).Sub(num, big.NewInt(1))) == 0 {
				break
			}
		}
		if x.Cmp(new(big.Int).Sub(num, big.NewInt(1))) != 0 {
			return false
		}
	}
	return true
}

func evclideGCD(a, b *big.Int) *big.Int {
	for a.Cmp(big.NewInt(0)) != 0 && b.Cmp(big.NewInt(0)) != 0 {
		if a.Cmp(b) == 1 {
			a.Mod(a, b)
		} else {
			b.Mod(b, a)
		}
	}
	return (new(big.Int).Add(a, b))
}

func NewDecryptor(url, msgTopic string) (*Decryptor, error) {
	dec := &Decryptor{}
	if len(os.Args) == 1 {
		return nil, fmt.Errorf("[ERROR] Not enough arguments")
	} else if os.Args[1] == "-g" {
		keyLen := flag.Int("g", 0, "Enter the key length!")
		flag.Parse()
		dec.KeyGen(*keyLen)
	} else {
		dec.KeyLen = len(os.Args) - 1
		dec.PrivateKey = make([]*big.Int, dec.KeyLen)
		for i := 0; i < dec.KeyLen; i++ {
			dec.PrivateKey[i] = big.NewInt(int64(0))
			dec.PrivateKey[i].SetString(os.Args[i+1], 10)
		}
	}
	dec.Q = big.NewInt(0)
	for i := 0; i < dec.KeyLen; i++ {
		dec.Q.Add(dec.Q, dec.PrivateKey[i])
	}
	for !millerRabinTest(dec.Q, 10000) {
		dec.Q.Add(dec.Q, big.NewInt(int64(1)))
	}
	dec.R = new(big.Int)
	dec.R.Rand(rand.New(rand.NewSource(time.Now().UnixNano())), dec.Q)
	for evclideGCD(new(big.Int).Set(dec.Q), new(big.Int).Set(dec.R)).Cmp(big.NewInt(1)) != 0 {
		dec.R.Rand(rand.New(rand.NewSource(time.Now().UnixNano())), dec.Q)
	}
	dec.PublicKey = make([]*big.Int, dec.KeyLen)
	for i := 0; i < dec.KeyLen; i++ {
		dec.PublicKey[i] = big.NewInt(int64(0))
		mul := new(big.Int)
		mod := new(big.Int)
		dec.PublicKey[i].Set(mod.Mod(mul.Mul(dec.PrivateKey[i], dec.R), dec.Q))
	}

	var err error
	dec.MsgSub, err = subscriber.NewSubscriber(url, msgTopic)
	if err != nil {
		return nil, err
	}

	dec.Pub, err = publisher.NewPublisher(url)
	if err != nil {
		return nil, err
	}

	return dec, nil
}

func (dec *Decryptor) decryptToBinary(encryptedMsg []*big.Int) []string {
	modInv := new(big.Int)
	modInv.ModInverse(dec.R, dec.Q)
	var newLen = len(encryptedMsg)
	var tmp []*big.Int = make([]*big.Int, newLen)
	for i := 0; i < newLen; i++ {
		tmp[i] = new(big.Int)
		mul := new(big.Int)
		mod := new(big.Int)
		tmp[i].Set(mod.Mod((mul.Mul(encryptedMsg[i], modInv)), dec.Q))
	}
	var tmpBinary [][]byte = make([][]byte, newLen)
	var binary []string = make([]string, newLen)
	for i := 0; i < newLen; i++ {
		tmpBinary[i] = make([]byte, dec.KeyLen)
		for j := dec.KeyLen - 1; j >= 0; j-- {
			if dec.PrivateKey[j].Cmp(tmp[i]) == 1 {
				tmpBinary[i][j] = '0'
			} else {
				tmpBinary[i][j] = '1'
				neg := big.NewInt(int64(0))
				tmp[i].Add(tmp[i], neg.Neg(dec.PrivateKey[j]))
			}
		}
		binary[i] = string(tmpBinary[i])
	}
	return binary
}

func (dec *Decryptor) oldBinary(encryptedMsg []*big.Int) []string {
	var newBinary []string = dec.decryptToBinary(encryptedMsg)
	var tmp string
	var builder strings.Builder
	newLen := len(newBinary)
	builder.Grow(newLen * dec.KeyLen)
	for i := 0; i < newLen; i++ {
		builder.WriteString(newBinary[i])
	}
	tmp = builder.String()
	zeros := len(tmp) % (dec.MsgLen * 7)
	n := len(tmp) / (dec.MsgLen * 7)
	if n > 1 {
		zeros += (n - 1) * (dec.MsgLen * 7)
	}
	tmp = tmp[zeros:]
	var oldBinary []string = make([]string, dec.MsgLen)
	for i := 0; i < dec.MsgLen; i++ {
		oldBinary[i] = tmp[i*7 : (i+1)*7]
	}
	return oldBinary
}

func (dec *Decryptor) Decrypting(encryptedMsg []*big.Int) string {
	var binary []string = dec.oldBinary(encryptedMsg)
	var tmp []byte = make([]byte, dec.MsgLen)
	for i := 0; i < dec.MsgLen; i++ {
		a := 0
		for j := 6; j >= 0; j-- {
			if binary[i][j] == '1' {
				a += int(math.Pow(2, float64(6-j)))
			}
		}
		tmp[i] = byte(a)
	}

	return string(tmp)

}
