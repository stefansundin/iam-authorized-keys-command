package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func main() {
	// Check for user name
	if len(os.Args) <= 1 {
		os.Exit(0)
	}
	sshUserName := os.Args[1]

	// Get the user's SSH keys
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		},
	}))
	svc := iam.New(sess)
	if resp, err := svc.ListSSHPublicKeys(&iam.ListSSHPublicKeysInput{
		UserName: &sshUserName,
	}); err == nil {
		for _, key := range resp.SSHPublicKeys {
			if *key.Status != "Active" {
				continue
			}
			resp, _ := svc.GetSSHPublicKey(&iam.GetSSHPublicKeyInput{
				Encoding:       aws.String("SSH"),
				SSHPublicKeyId: key.SSHPublicKeyId,
				UserName:       &sshUserName,
			})
			fmt.Printf("# %s\n", sshUserName)
			fmt.Println(*resp.SSHPublicKey.SSHPublicKeyBody)
		}
	} else {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(0)
}
