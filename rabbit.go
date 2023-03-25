package rabbit

import (
	"fmt"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

func NewRabbit() string {
	return "new rabbit created"
}

func readAndParsePrivateKey(keyPath string) (ssh.Signer, error) {
	pemBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read PEM file: %w", err)
	}

	pemSigner, err := ssh.ParsePrivateKey(pemBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PEM file: %w", err)
	}

	return pemSigner, nil
}

// How to use:
// session, err := rabbit.Connect("ubuntu", "path_to_pem_file", "ec2-3-12-74-6.us-east-2.compute.amazonaws.com")
// if err != nil {
// 	fmt.Println("Failed to connect to remote system:", err)
// 	return
// }
// defer session.Close()
//
// // Now you can use the session to execute commands on the remote system
// output, err := session.Output("ls -la")
// if err != nil {
// 	fmt.Println("Failed to execute command:", err)
// 	return
// }
// fmt.Println(string(output))

func Connect(user, password, pemPath, host string) (*ssh.Session, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		// HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			// TODO: Implement custom host key validation here like checking trusted keys or fingerprints
			return nil
		},
	}

	// If a password is provided, add it to the authentication methods
	if password != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(password))
	}

	// If a PEM file is provided, parse it and add it to the authentication methods
	if pemPath != "" {
		pemSigner, _ := readAndParsePrivateKey(pemPath)
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeys(pemSigner))
	}

	// Dial the remote host using the SSH client configuration
	sshClient, err := ssh.Dial("tcp", host+":22", sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to host: %w", err)
	}

	// Create a new session over the SSH connection
	session, err := sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}
