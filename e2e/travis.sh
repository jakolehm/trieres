#!/bin/bash
# shellcheck disable=SC2039 disable=SC1091

set -ue

go build -o trieres main.go

ssh-keygen -t rsa -f ~/.ssh/id_rsa_travis -N ""
cat ~/.ssh/id_rsa_travis.pub > ~/.ssh/authorized_keys
chmod 0600 ~/.ssh/authorized_keys

envsubst < e2e/cluster.yml > cluster.yml
envsubst < e2e/footloose.yaml > footloose.yaml

curl -L https://github.com/weaveworks/footloose/releases/download/0.6.3/footloose-0.6.3-linux-x86_64 > ./footloose
chmod +x ./footloose
./footloose create
./footloose ssh root@master0 -- 'apt-get install -y curl || yum install -y curl which openssh-clients'
./footloose ssh root@worker0 -- 'apt-get install -y curl || yum install -y curl which openssh-clients'
./footloose ssh root@worker1 -- 'apt-get install -y curl || yum install -y curl which openssh-clients'

./trieres
./trieres -v
./trieres --debug up
./trieres kubeconfig > kubeconfig.e2e

export KUBECONFIG=./kubeconfig.e2e

echo "==> Test with sonobuoy"
curl -L https://github.com/vmware-tanzu/sonobuoy/releases/download/v0.18.0/sonobuoy_0.18.0_linux_amd64.tar.gz | tar xzv
chmod +x ./sonobuoy
sleep 30
./sonobuoy run --mode quick --timeout 600 --wait

