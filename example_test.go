package smb2_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"

	"github.com/hermanntoast/go-smb2"
)

func Example() {
	conn, err := net.Dial("tcp", "localhost:445")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     "Guest",
			Password: "",
			Domain:   "MicrosoftAccount",
		},
	}

	c, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer c.Logoff()

	fs, err := c.Mount(`\\localhost\share`)
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	f, err := fs.Create("hello.txt")
	if err != nil {
		panic(err)
	}
	defer fs.Remove("hello.txt")
	defer f.Close()

	_, err = f.Write([]byte("Hello world!"))
	if err != nil {
		panic(err)
	}

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		panic(err)
	}

	bs, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bs))

	// Hello world!
}

func ExampleShare_WhoAmI() {
	conn, err := net.Dial("tcp", "localhost:445")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	d := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     "username",
			Password: "password",
			Domain:   "",
		},
	}

	c, err := d.Dial(conn)
	if err != nil {
		panic(err)
	}
	defer c.Logoff()

	fs, err := c.Mount(`\\localhost\share`)
	if err != nil {
		panic(err)
	}
	defer fs.Umount()

	// Retrieve current user's SID and group SIDs
	identity, err := fs.WhoAmI()
	if err != nil {
		panic(err)
	}

	fmt.Printf("User SID: %s\n", identity.UserSID)
	fmt.Printf("Group SIDs:\n")
	for _, groupSID := range identity.GroupSIDs {
		fmt.Printf("  - %s\n", groupSID)
	}
}
