# existing-user

In this version, a system user has to exist with the same username as your IAM
user.

## Example sshd_config

```
AuthorizedKeysFile none
AuthorizedKeysCommand /path/to/compiled/binary
AuthorizedKeysCommandUser nobody
```

If you still want to be able to use the authorized_keys file for some users,
e.g. in case IAM is experiencing downtime, you can add something like the
following:

```
Match User ubuntu
  AuthorizedKeysFile %h/.ssh/authorized_keys
```

Don't forget to restart the ssh service:

```shell
service ssh restart
```

## IAM Role Permissions

This script needs the following policy to execute properly, so make sure you
apply it to your EC2 Role:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "iam:ListSSHPublicKeys",
                "iam:GetSSHPublicKey"
            ],
            "Resource": "*"
        }
    ]
}
```
