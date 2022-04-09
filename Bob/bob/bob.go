package bob

import (
	"flag"
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

func IsPrime(num *big.Int) bool {
	var sq, mod *big.Int = new(big.Int), new(big.Int)
	sq.Sqrt(num)
	var i int64
	if mod.Mod(num, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}
	for i = 3; sq.Cmp(big.NewInt(i)) >= 0; i += 2 {
		mod.Mod(num, big.NewInt(i)) 
		if mod.Cmp(big.NewInt(int64(0))) == 0 {
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
			Bob.PrivateKey[i].SetString(os.Args[i + 1], 10)
		}
	}
	Bob.Q = big.NewInt(int64(0))
	for i := 0; i < Bob.KeyLen; i++ {
		Bob.Q.Add(Bob.Q, Bob.PrivateKey[i])
	}
	for !IsPrime(Bob.Q) {
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

func CryptToBinary(Bob *T_Bob) []string {
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
	var newBinary []string = CryptToBinary(Bob)
	var tmp string
	var b strings.Builder
	newLen := len(newBinary)
	b.Grow(newLen * Bob.KeyLen)
	for i := 0; i < newLen; i++ {
		b.WriteString(newBinary[i])
	}
	tmp = b.String()
	zeros := len(tmp) % (Bob.MsgLen * 7)
	tmp = tmp[zeros :]
	var oldBinary []string = make([]string, Bob.MsgLen)
	for i := 0 ; i < Bob.MsgLen; i++ {
		oldBinary[i] = tmp[i * 7 : (i + 1) * 7]
	}
	return oldBinary
}

func (Bob *T_Bob) Encrypting() string {
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