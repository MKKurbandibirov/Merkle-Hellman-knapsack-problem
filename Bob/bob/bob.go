package bob

import (
	"flag"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"
)

type T_Bob struct {
	PrivateKey	[]*big.Int
	PublicKey	[]*big.Int
	CryptedMsg	[]*big.Int
	KeyLen		int
	MsgLen		int
	Q			*big.Int
	R			*big.Int
}

func (Bob *T_Bob) KeyGen(keyLen int) {
	Bob.KeyLen = keyLen
	Bob.PrivateKey = make([]*big.Int, keyLen)
	var sum, rnd *big.Int = new(big.Int), new(big.Int)
	for i := 0; i < keyLen; i++ {
		rnd.Rand(rand.New(rand.NewSource(time.Now().UnixNano())), big.NewInt(int64(100)))
		sum.Mul(sum, big.NewInt(int64(2)))
		sum.Add(sum, rnd)
		Bob.PrivateKey[i] = new(big.Int)
		Bob.PrivateKey[i].Set(sum)
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

func MillerRabinTest(num *big.Int, k int) bool {
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

func CreateBob() *T_Bob {
	var Bob *T_Bob = new(T_Bob)
	if len(os.Args) == 1 {
		os.Exit(1)
	} else if os.Args[1] == "-g" {
		keyLen := flag.Int("g", 0, "Enter the key length!")
		flag.Parse()
		Bob.KeyGen(*keyLen)
	} else {
		Bob.KeyLen = len(os.Args) - 1
		Bob.PrivateKey = make([]*big.Int, Bob.KeyLen)
		for i := 0; i < Bob.KeyLen; i++ {
			Bob.PrivateKey[i] = big.NewInt(int64(0))
			Bob.PrivateKey[i].SetString(os.Args[i+1], 10)
		}
	}
	Bob.Q = big.NewInt(int64(0))
	for i := 0; i < Bob.KeyLen; i++ {
		Bob.Q.Add(Bob.Q, Bob.PrivateKey[i])
	}
	for !MillerRabinTest(Bob.Q, 100000) {
		Bob.Q.Add(Bob.Q, big.NewInt(int64(1)))
	}
	Bob.R = big.NewInt(int64(0))
	Bob.R.Rand(rand.New(rand.NewSource(time.Now().UnixNano())), Bob.Q)
	Bob.PublicKey = make([]*big.Int, Bob.KeyLen)
	for i := 0; i < Bob.KeyLen; i++ {
		Bob.PublicKey[i] = big.NewInt(int64(0))
		mul := big.NewInt(int64(0))
		mod := big.NewInt(int64(0))
		Bob.PublicKey[i].Set(mod.Mod(mul.Mul(Bob.PrivateKey[i], Bob.R), Bob.Q))
	}
	return Bob
}

func DecryptToBinary(Bob *T_Bob) []string {
	modInv := new(big.Int)
	modInv.ModInverse(Bob.R, Bob.Q)
	var newLen = len(Bob.CryptedMsg)
	var tmp []*big.Int = make([]*big.Int, newLen)
	for i := 0; i < newLen; i++ {
		tmp[i] = new(big.Int)
		mul := new(big.Int)
		mod := new(big.Int)
		tmp[i].Set(mod.Mod((mul.Mul(Bob.CryptedMsg[i], modInv)), Bob.Q))
	}
	var tmpBinary [][]byte = make([][]byte, newLen)
	var binary []string = make([]string, newLen)
	for i := 0; i < newLen; i++ {
		tmpBinary[i] = make([]byte, Bob.KeyLen)
		for j := Bob.KeyLen - 1; j >= 0; j-- {
			if Bob.PrivateKey[j].Cmp(tmp[i]) == 1 {
				tmpBinary[i][j] = '0'
			} else {
				tmpBinary[i][j] = '1'
				neg := big.NewInt(int64(0))
				tmp[i].Add(tmp[i], neg.Neg(Bob.PrivateKey[j]))
			}
		}
		binary[i] = string(tmpBinary[i])
	}
	return binary
}

func OldBinary(Bob *T_Bob) []string {
	var newBinary []string = DecryptToBinary(Bob)
	var tmp string
	var b strings.Builder
	newLen := len(newBinary)
	b.Grow(newLen * Bob.KeyLen)
	for i := 0; i < newLen; i++ {
		b.WriteString(newBinary[i])
	}
	tmp = b.String()
	zeros := len(tmp) % (Bob.MsgLen * 7)
	n := len(tmp) / (Bob.MsgLen * 7)
	if n > 1 {
		zeros += (n - 1) * (Bob.MsgLen * 7)
	}
	tmp = tmp[zeros :]
	var oldBinary []string = make([]string, Bob.MsgLen)
	for i := 0 ; i < Bob.MsgLen; i++ {
		oldBinary[i] = tmp[i * 7 : (i + 1) * 7]
	}
	return oldBinary
}

func (Bob *T_Bob) Decrypting() string {
	var binary []string = OldBinary(Bob)
	var tmp []byte = make([]byte, Bob.MsgLen)
	for i := 0; i < Bob.MsgLen; i++ {
		a := 0
		for j := 6; j >= 0; j-- {
			if binary[i][j] == '1' {
				a += int(math.Pow(2,float64(6 - j)))
			}
		}
		tmp[i] = byte(a)
	}
	encryptedMsg := string(tmp)
	return encryptedMsg

}