package encryptor

import (
	"bufio"
	"fmt"
	"math/big"
	"os"
	"strings"
	"utils/publisher"
	"utils/subscriber"

	"github.com/nats-io/nats.go"
)

type Encryptor struct {
	PublicKey    []*big.Int
	KeyLen       int
	MsgLen       int

	KeySub *subscriber.Subscriber
	MsgPub *publisher.Publisher
}

func NewEncryptor(url string, keyTopic string) (*Encryptor, error) {
	var (
		enc = &Encryptor{}

		err error
	)

	enc.KeySub, err = subscriber.NewSubscriber(nats.DefaultURL, keyTopic)
	if err != nil {
		return nil, err
	}

	enc.MsgPub, err = publisher.NewPublisher(url)
	if err != nil {
		return nil, err
	}

	return enc, nil
}

func readMessage() string {
	var b strings.Builder
	var scanner bufio.Scanner = *bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	fmt.Print("Enter the message:\n<<< ")

	for scanner.Scan() {
		line := scanner.Text()
		if line == "/exit" {
			break
		}
		b.WriteString(line)
		b.WriteString("\n")

		fmt.Print("<<< ")
	}

	return b.String()
}

func (enc *Encryptor) messageToBinary() []string {
	var (
		msg    = readMessage()
		msgLen = len(msg)

		tmp    = make([][]byte, msgLen)
		binary = make([]string, msgLen)
	)

	enc.MsgLen = msgLen
	for i := 0; i < msgLen; i++ {
		ch := byte(msg[i])
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

	return binary
}

func (enc *Encryptor) newBinary() []string {
	var (
		binary = enc.messageToBinary()
		newLen = enc.MsgLen*7 + (enc.KeyLen - (enc.MsgLen*7)%enc.KeyLen)

		b strings.Builder
	)

	b.Grow(newLen)
	for i := 0; i < newLen-enc.MsgLen*7; i++ {
		b.WriteString("0")
	}
	for i := 0; i < enc.MsgLen; i++ {
		b.WriteString(binary[i])
	}

	new := b.String()
	var newBinary []string = make([]string, newLen/enc.KeyLen)
	for i := 0; i < newLen/enc.KeyLen; i++ {
		newBinary[i] = new[i*enc.KeyLen : (i+1)*enc.KeyLen]
	}

	fmt.Println("\n\nBinary reprasantation of knapsack:")

	for i := 0; i < newLen/enc.KeyLen; i++ {
		fmt.Printf("%s\n", newBinary[i])
	}
	fmt.Println()

	return newBinary
}

func (enc *Encryptor) Encrypting() []*big.Int {
	var (
		newBinary = enc.newBinary()
		newLen    = len(newBinary)
	)

	encryptedMsg := make([]*big.Int, newLen)
	for i := 0; i < newLen; i++ {
		encryptedMsg[i] = big.NewInt(int64(0))
		for j := 0; j < enc.KeyLen; j++ {
			mul := big.NewInt(int64(0))
			encryptedMsg[i].Add(encryptedMsg[i],
				mul.Mul(big.NewInt(int64((int(newBinary[i][j])-'0'))), enc.PublicKey[j]))
		}
	}

	return encryptedMsg
}
