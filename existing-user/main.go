package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	exitCodeOk    int = 0
	exitCodeError int = 1
)

func main() {
	// Exit cleanly on SIGPIPE
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGPIPE)
	go func() {
		_ = <-c
		os.Exit(exitCodeOk)
	}()

	// Check for user name
	if len(os.Args) <= 1 {
		os.Exit(exitCodeOk)
	}
	sshUserName := os.Args[1]

	// Get the user's SSH keys
	params := &iam.ListSSHPublicKeysInput{
		UserName: &sshUserName,
	}
	sess, _ := session.NewSession()
	svc := iam.New(sess)
	if resp, err := svc.ListSSHPublicKeys(params); err == nil {
		for _, key := range resp.SSHPublicKeys {
			if *key.Status != "Active" {
				continue
			}
			params := &iam.GetSSHPublicKeyInput{
				Encoding:       aws.String("SSH"),
				SSHPublicKeyId: key.SSHPublicKeyId,
				UserName:       &sshUserName,
			}
			resp, _ := svc.GetSSHPublicKey(params)
			fmt.Printf("# %s\n", sshUserName)
			fmt.Println(*resp.SSHPublicKey.SSHPublicKeyBody)
		}
	}
	os.Exit(exitCodeOk)
}
