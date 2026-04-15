package sftpclient

import (
	"fmt"
	"io"
	"log"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
	"piamonte-proxy/config"
)

func ForwardToSFTP(cfg *config.Config, fileReader io.Reader, filename string) error {
	var hostKeyCallback ssh.HostKeyCallback
	if cfg.SftpInsecureSkipVerify == "true" {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
		log.Println("WARNING: SSH HostKey verification is disabled. This is insecure.")
	} else {
		var err error
		hostKeyCallback, err = knownhosts.New(cfg.SftpKnownHostsFile)
		if err != nil {
			return fmt.Errorf("failed to load known_hosts file at %s: %v", cfg.SftpKnownHostsFile, err)
		}
	}

	sshConfig := &ssh.ClientConfig{
		User: cfg.SftpUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.SftpPassword),
		},
		HostKeyCallback: hostKeyCallback,
	}

	log.Printf("Attempting SSH connection to: %s", cfg.SftpHost)
	sshClient, err := ssh.Dial("tcp", cfg.SftpHost, sshConfig)
	if err != nil {
		return fmt.Errorf("SSH connection failed: %v", err)
	}
	defer sshClient.Close()

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return fmt.Errorf("SFTP session initialization failed: %v", err)
	}
	defer sftpClient.Close()

	destPath := filepath.Join(cfg.SftpDestDir, filename)
	destFile, err := sftpClient.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %v", err)
	}
	defer destFile.Close()

	// pasamos el churro de datos directo
	bytesCopied, err := io.Copy(destFile, fileReader)
	if err != nil {
		return fmt.Errorf("failed to write data to remote file: %v", err)
	}

	log.Printf("Transfer successful. %d bytes written to %s", bytesCopied, destPath)
	return nil
}
