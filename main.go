package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

const (
	exitCodeOk    int = 0
	exitCodeError int = 1
)

var (
	wg sync.WaitGroup

	iamGroup    = ""
	sshUserName = ""
)

func main() {
	sess, _ := session.NewSession()
	svc := iam.New(sess)

	// check for valid user name
	if sshUserName != "" && (len(os.Args) < 2 || os.Args[1] != sshUserName) {
		os.Exit(exitCodeOk)
	}

	users, err := users(svc, iamGroup)
	if err != nil {
		os.Exit(exitCodeError)
	}

	for _, u := range users {
		go func(userName *string) {
			wg.Add(1)
			params := &iam.ListSSHPublicKeysInput{
				UserName: userName,
			}
			if resp, err := svc.ListSSHPublicKeys(params); err == nil {
				for _, k := range resp.SSHPublicKeys {
					params := &iam.GetSSHPublicKeyInput{
						Encoding:       aws.String("SSH"),
						SSHPublicKeyId: k.SSHPublicKeyId,
						UserName:       userName,
					}
					resp, _ := svc.GetSSHPublicKey(params)
					if *resp.SSHPublicKey.Status == "Active" {
						fmt.Printf("# %s\n", *userName)
						fmt.Println(*resp.SSHPublicKey.SSHPublicKeyBody)
					}
				}
			}
			wg.Done()
		}(u.UserName)
	}
	wg.Wait()
}

// get all IAM users, or just those that are part of the defined group
func users(svc *iam.IAM, iamGroup string) ([]*iam.User, error) {
	if iamGroup != "" {
		params := &iam.GetGroupInput{
			GroupName: aws.String(iamGroup),
		}
		resp, err := svc.GetGroup(params)
		return resp.Users, err
	}
	params := &iam.ListUsersInput{
		MaxItems: aws.Int64(100),
	}
	resp, err := svc.ListUsers(params)
	return resp.Users, err
}
