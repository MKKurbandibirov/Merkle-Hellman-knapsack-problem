package alice

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
)

type T_Alissa struct {
	PublicKey		[]*big.Int
	EncryptedMsg	[]*big.Int
	Message			string
	KeyLen     		int
	MsgLen			int
}

func CreateAlissa(keyLen int) *T_Alissa {
	var Alissa *T_Alissa = new(T_Alissa)
	Alissa.KeyLen = keyLen
	return (Alissa)
}

func readMessage(Alissa *T_Alissa) {
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
	Alissa.Message = b.String()
}

func messageToBinary(Alissa *T_Alissa) []string {
	readMessage(Alissa)
	Alissa.MsgLen = len(Alissa.Message)
	var tmp [][]byte = make([][]byte, Alissa.MsgLen)
	var binary []string = make([]string, Alissa.MsgLen)
	for i := 0; i < Alissa.MsgLen; i++ {
		ch := byte(Alissa.Message[i])
		tmp[i] = make([]byte, 7)
		for j := 6; j >= 0; j-- {
			if ch%2 == 1 {
				tmp[i][j] = '1'
			} else {
				tmp[i][j] = '0'
			}
			ch /= 2
		}
		binary[i] = string(tmp[i])
	}
	fmt.Println("\nBinary reprasantation of message:")
	for i := 0; i < Alissa.MsgLen; i++ {
		fmt.Printf("%s\n", binary[i])
	}
	fmt.Println()
	return binary
}

func newBinary(Alissa *T_Alissa) []string {
	var binary []string = messageToBinary(Alissa)
	var newLen = Alissa.MsgLen*7 + (Alissa.KeyLen - (Alissa.MsgLen*7)%Alissa.KeyLen)
	var b strings.Builder
	b.Grow(newLen)
	for i := 0; i < newLen-Alissa.MsgLen*7; i++ {
		b.WriteString("0")
	}
	for i := 0; i < Alissa.MsgLen; i++ {
		b.WriteString(binary[i])
	}
	new := b.String()
	var newBinary []string = make([]string, newLen/Alissa.KeyLen)
	for i := 0; i < newLen/Alissa.KeyLen; i++ {
		newBinary[i] = new[i*Alissa.KeyLen : (i+1)*Alissa.KeyLen]
	}
	fmt.Println("Binary reprasantation after changes:")
	for i := 0; i < newLen/Alissa.KeyLen; i++ {
		fmt.Printf("%s\n", newBinary[i])
	}
	fmt.Println()
	return newBinary
}

func (Alissa *T_Alissa) Encrypting() {
	var newBinary []string = newBinary(Alissa)
	newLen := len(newBinary)
	Alissa.EncryptedMsg = make([]*big.Int, newLen)
	for i := 0; i < newLen; i++ {
		Alissa.EncryptedMsg[i] = big.NewInt(int64(0))
		for j := 0; j < Alissa.KeyLen; j++ {
			mul := big.NewInt(int64(0))
			Alissa.EncryptedMsg[i].Add(Alissa.EncryptedMsg[i],
				mul.Mul(big.NewInt(int64((int(newBinary[i][j])-'0'))), Alissa.PublicKey[j]))
		}
	}
}
