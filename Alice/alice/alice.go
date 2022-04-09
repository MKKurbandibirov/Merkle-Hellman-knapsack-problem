package alice

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

type T_Alice struct {
	PublicKey	[]*big.Int
	CryptedMsg	[]*big.Int
	Message		string
	KeyLen		int
	MsgLen		int
}

func CreateAlice(keyLen int) *T_Alice {
	var Alice *T_Alice = new(T_Alice)
	Alice.KeyLen = keyLen
	return (Alice)
}

func readMessage(Alice *T_Alice) {
	var b strings.Builder
	var scanner bufio.Scanner = *bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	fmt.Print("Enter the message:\n>>> ")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "/exit" {
			break
		}
		b.WriteString(line)
		b.WriteString("\n")
		fmt.Print(">>> ")
	}
	Alice.Message = b.String()
}

func messageToBinary(Alice *T_Alice) []string {
	readMessage(Alice)
	Alice.MsgLen = len(Alice.Message)
	var tmp [][]byte = make([][]byte, Alice.MsgLen)
	var binary []string = make([]string, Alice.MsgLen)
	for i := 0; i < Alice.MsgLen; i++ {
		ch := byte(Alice.Message[i])
		tmp[i] = make([]byte, 7)
		for j := 6; j >= 0; j-- {
			if ch % 2 == 1 {
				tmp[i][j] = '1'
			} else {
				tmp[i][j] = '0'
			}
			ch /= 2
		}
		binary[i] = string(tmp[i])
	}
	return binary
}

func newBinary(Alice *T_Alice) []string {
	var binary []string = messageToBinary(Alice)
	var newLen = Alice.MsgLen * 7 + (Alice.KeyLen - (Alice.MsgLen * 7) % Alice.KeyLen)
	var b strings.Builder
	b.Grow(newLen)
	for i := 0; i < newLen - Alice.MsgLen * 7; i++ {
		b.WriteString("0")
	}
	for i := 0; i < Alice.MsgLen; i++ {
		b.WriteString(binary[i])
	}
	new := b.String()
	var newBinary []string = make([]string, newLen / Alice.KeyLen)
	for i := 0; i < newLen / Alice.KeyLen; i++ {
		newBinary[i] = new[i * Alice.KeyLen : (i + 1) * Alice.KeyLen]
	}
	return newBinary
}

func (Alice *T_Alice) Crypting() {
	var newBinary []string = newBinary(Alice)
	newLen := len(newBinary)
	Alice.CryptedMsg = make([]*big.Int, newLen)
	for i := 0; i < newLen; i++ {
		Alice.CryptedMsg[i] = big.NewInt(int64(0))
		for j := 0; j < Alice.KeyLen; j++ {
			mul := big.NewInt(int64(0))
			Alice.CryptedMsg[i].Add(Alice.CryptedMsg[i], 
				mul.Mul(big.NewInt(int64((int(newBinary[i][j]) - '0'))), Alice.PublicKey[j]))
		}
	}
}