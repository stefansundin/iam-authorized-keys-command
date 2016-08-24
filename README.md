iam-authorized-keys-command
===========================

Authenticate SSH connections with public keys stored in Amazon IAM.

Building
========

The `AuthorizedKeysCommand` is executed by sshd and there is no way to control
the input arguments. You'll need to build this project yourself to configure
it the way you want.

There are two symbols in the source that can be specified
with build flags: `iamGroup` and `sshUserName`. The `iamGroup` variable let's
you limit IAM user lookups to members of a particular group. The `sshUserName`
let's you limit key lookups to a specific system user.

To build a binary that will only authenticate session for the `ubuntu` user you
can run the following build command:

```shell
go build -ldflags "-X main.sshUserName=ubuntu"
```

To build a binary that will only load keys for users in a specific IAM group:

```shell
go build -ldflags "-X main.iamGroup=Engineers"
```

To build a binary for 64-bit linux:

```shell
GOOS=linux GOARCH=amd64 go build
```

IAM Role Permissions
====================

This script needs the following policy to execute properly, so make sure you
apply it to your EC2 Role:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "iam:GetGroup",
                "iam:ListUsers",
                "iam:GetSSHPublicKey",
                "iam:ListSSHPublicKeys"
            ],
            "Resource": "*"
        }
    ]
}
```

Configuring sshd
================

To configure sshd to use your compiled authorized keys command, add the
following to `/etc/ssh/sshd_config`:

```shell
AuthorizedKeysCommand /path/to/compiled/binary
AuthorizedKeysCommandRunAs nobody
```

Then restart sshd:

```shell
service ssh restart
```

Why Go?
=======

Because go is awesome! No seriously though, we're looking up a large number
of users in IAM and doing so synchronously is really slow. Doing these lookups
concurrently within a goroutine significantly decreases the amount of time it
takes for this script to return.
