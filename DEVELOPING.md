
# Installing a dev environment

The steps below can be used to set up a develop environment using an AWS EC2 instance running Amazon Linux 2.

*Assuming running as root* 

## Install Go and Python 3:

```
amazon-linux-extras install golang1.9
amazon-linux-extras install  python3
```

## Install git
```
yum install -y git
```

## Download the proto compiler
```
wget https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-linux-x86_64.zip
unzip protoc-3.5.1-linux-x86_64.zip
cp bin/protoc /usr/bin/
```

## Download Go dependencies
```
go get -u github.com/gobuffalo/packr/...
go get github.com/square/certstrap
go get google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
```

## Add Go bin dir to path
```
PATH=$PATH:/root/go/bin/
```

## Clone repo
```
cd /root/go/src/github.com
mkdir ebastos
cd ebastos
git clone 
https://github.com/ebastos/shell-history.git
cd shell-history
```

## Install Python requirements
```
pip3 install -r requirements.txt
```

## Generate certificates
```
certstrap --depot-path certs init --common-name localhost
```