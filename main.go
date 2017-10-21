package main // import "github.com/denisbrodbeck/memgrdpeek"

import (
	"fmt"
	"os"

	"github.com/awnumar/memguard"
)

type message struct {
	ID        []byte                 // 32 byte
	Recipient []byte                 // 32 byte
	Sender    []byte                 // 32 byte
	Message   *memguard.LockedBuffer // Gets 32 byte input
	Meta      []byte                 // 32 byte
}

func newMessage(id, recipient, sender, msg, meta []byte) *message {
	msgBuf, err := memguard.NewImmutableFromBytes(assert32Bytes(msg))
	if err != nil {
		panic(err)
	}
	return &message{
		ID:        assert32Bytes(id),
		Recipient: assert32Bytes(recipient),
		Sender:    assert32Bytes(sender),
		Message:   msgBuf,
		Meta:      assert32Bytes(meta),
	}
}

func main() {
	defer memguard.DestroyAll()
	msg := newMessage(
		[]byte("(((id    | 1234567890123456  )))"),
		[]byte("(((recip | bobby@spacer.com  )))"),
		[]byte("(((sender| eva@secretmoon.m  )))"),
		[]byte("(((msg   | VERY SECRET VERY  )))"),
		[]byte("(((meta  | 2017-10-21 10:26  )))"),
	)
	pid := os.Getpid()
	fmt.Println("                                                              ") // easier to spot message
	fmt.Println("Running under PID:", pid)
	fmt.Println("Press enter after running 'memory_reader.py' to exit...")
	line := ""
	fmt.Scanln(&line)
	fmt.Printf("%T: %s\n", msg, msg)
}

// helper
func assert32Bytes(in []byte) []byte {
	if len(in) != 32 {
		panic("invalid length")
	}
	return in
}
