package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func envVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func main() {
	host := envVariable("host")
	username := envVariable("username")
	password := envVariable("password")
	port := envVariable("port")

	addr := fmt.Sprintf("%s:%s", host, port)

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Printf("Failed to dial: %v", err)
		os.Exit(1)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		fmt.Printf("Failed to create SFTP client: %v", err)
		os.Exit(1)
	}

	statVF, err := client.StatVFS("/")
	if err != nil {
		fmt.Printf("Failed to get filesytem info: %v", err)
		os.Exit(1)
	}
	fmt.Printf("Space available: %d bytes\n", statVF.FreeSpace())
	defer client.Close()
}
