# ec2-cli
CLI to help lookup ec2 instances in aws acount

## Installation steps

Set architecture and OS type

### For MAC M1 (arm64 architecture)

```shell
EC2CLI_OS_ARCH=ec2-cli-darwin-arm64
```

### For MAC (amd64)

```shell
EC2CLI_OS_ARCH=ec2-cli-darwin-amd64
```

### For linux (amd64)

```shell
EC2CLI_OS_ARCH=ec2-cli-linux-amd64
```

### Execute after setting the Arch and OS type

```shell
LATEST_RELEASE=$(curl -L -s -H 'Accept: application/json' https://github.com/prashantkalkar/ec2-cli/releases/latest) && \
LATEST_VERSION=$(echo "$LATEST_RELEASE" | jq -r ".tag_name") && \
curl -O -L https://github.com/prashantkalkar/ec2-cli/releases/download/"$LATEST_VERSION"/"$EC2CLI_OS_ARCH" && \
mv "$EC2CLI_OS_ARCH" ec2-cli && \
chmod +x ec2-cli && \
mv ec2-cli /usr/local/bin/
```

# Usage

## Search instances by IP

```shell
$ ec2-cli --ip 172.16.110.76
INSTANCE_ID             NAME              IP_ADDRESS
i-06c6ea69e00f4787f     INSTANCE_NAME     172.16.110.76
```

## Only get the ID instead of table

```shell
$ ec2-cli --ip 172.16.110.76 --id
i-06c6ea69e00f4787f
```

## Lookup instances by tags

```shell
$ ec2-cli --tags tag1,tag2
INSTANCE_ID             NAME                                 IP_ADDRESS
i-06c6ea69e00f4787f     instance_name1                       172.16.110.76
i-020bbf5c4fe2cf1b4     instance_name2                       172.16.110.12
i-03fe16d5667affe46     instance_name3                       172.16.110.44
```

## SSH using SSM

Generate public and private key pair if not already exists. 
(Sample steps: https://docs.github.com/en/authentication/connecting-to-github-with-ssh/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent)

Configure aws cli and also install SSM plugin. https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html

Ensure following permissions are provided to the user.
`ec2-instance-connect:SendSSHPublicKey`
`ssm:StartSession`
`ssm:TerminateSession`

Add following snippet at the end of the file `~/.ssh/config`

```
# SSH over Session Manager
host i-*
 IdentityFile $PRIVATE_KEY
 User ec2-user
 ProxyCommand sh -c "aws ec2-instance-connect send-ssh-public-key --instance-id %h --instance-os-user %r --ssh-public-key 'file://$PUBLIC_KEY' --availability-zone '$(aws ec2 describe-instances --instance-ids %h --query 'Reservations[0].Instances[0].Placement.AvailabilityZone' --output text)' && aws ssm start-session --target %h --document-name AWS-StartSSHSession --parameters 'portNumber=%p'"
```

Do not forget to replace $PRIVATE_KEY and $PUBLIC_KEY with the path to your private and public key.

## SSH using IP lookup

Once SSM configuration is done. Use ec2-cli to ssh into instances using IP. 

```shell
ssh $(ec2-cli --ip <IP> --id)
```

## SSH using tag lookup

Once SSM configuration is done. Use ec2-cli to ssh into instances using tags.

```shell
# get instance Id using ec2-cli
ssh $(ec2-cli --tags <tags> | grep <IPOrPart> | awk '{print $1}')
```

`<tags>` are comma separated instance tags (case insensitive).  

Reference:
https://cloudonaut.io/connect-to-your-ec2-instance-using-ssh-the-modern-way/