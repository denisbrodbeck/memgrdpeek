package main

import (
	"fmt"

	"github.com/awnumar/memguard"
)

type publicMessage struct {
	Pre  []byte // 32 byte
	Mid  []byte // 32 byte
	Post []byte // 32 byte
}

func newPublicMessage(pre, mid, post []byte) *publicMessage {
	return &publicMessage{
		Pre:  assert32Bytes(pre),
		Mid:  assert32Bytes(mid),
		Post: assert32Bytes(post),
	}
}

type privateMessage struct {
	Pre  []byte                 // 32 byte
	Mid  *memguard.LockedBuffer // Gets 32 byte input
	Post []byte                 // 32 byte
}

func newPrivateMessage(pre, mid, post []byte) *privateMessage {
	middle, err := memguard.NewImmutableFromBytes(assert32Bytes(mid))
	if err != nil {
		panic(err)
	}
	return &privateMessage{
		Pre:  assert32Bytes(pre),
		Mid:  middle,
		Post: assert32Bytes(post),
	}
}

func main() {
	defer memguard.DestroyAll()
	publicMsg := newPublicMessage(
		[]byte("(((pre| hi bob bob bobby |pre)))"),
		[]byte("(((mid| public  message! |mid)))"),
		[]byte("(((pos| after the data.. |pos)))"),
	)
	empty := []byte("                                ")
	// create an empty msg as a memory separator (makes py script output easier to read)
	emptyMsg := newPublicMessage(empty, empty, empty)
	privateMsg := newPrivateMessage(
		[]byte("(((pre| hi eva eva eva!! |pre)))"),
		[]byte("(((mid| private message! |mid)))"),
		[]byte("(((pos| end of pri msg.. |pos)))"),
	)
	fmt.Println("Press enter after running 'memory_reader.py' to exit...")
	line := ""
	fmt.Scanln(&line)
	fmt.Printf("Pub\nRaw byte output: %v\nString output: %s\n\n", publicMsg, publicMsg)
	fmt.Printf("Pri\nRaw byte output: %v\nString output: %s\n\n", privateMsg, privateMsg)
	fmt.Println(emptyMsg)
}

func assert32Bytes(in []byte) []byte {
	if len(in) != 32 {
		panic("invalid length")
	}
	return in
}
