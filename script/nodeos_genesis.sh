#!/bin/bash

# $1 = data-dir path
# $2 = producer_name
# $3 = producer_publickey
# $4 = producer_privatekey
# $5 = filepath

docker kill nodeos-sic || true
docker rm nodeos-sic || true

if [ ! -d "$1" ]; then
mkdir $1
fi

echo "Copying base config"
cp $5/script/base_config.ini $5/config.ini

echo "Copy genesis.json"
cp $5/script/genesis.json $1/genesis.json

echo "Import eosio to config.ini"
eosio_pub="EOS8Znrtgwt8TfpmbVpTKvA2oB8Nqey625CLN8bCN3TEbgx86Dsvr"
eosio_pri="5K463ynhZoCDDa4RDcr63cUwWLTnKqmdcoTKTHBjqoKfv4u5V7p"
echo "producer-name = eosio" >> config.ini
echo "signature-provider=${eosio_pub}=KEY:${eosio_pri}" >> config.ini
echo "enable-stale-production = true" >> config.ini

echo "Import Producer to config.ini"
echo "producer_name = $2" >> config.ini
echo "signature-provider=$3=KEY:$4" >> config.ini

mv config.ini $1/

echo "Running 'nodeos' through Docker."
docker run -ti --detach --name nodeos-sic \
       -v $1:/opt/eosio/bin/data-dir \
       -p 8888:8888 -p 9876:9876 \
       docker.banmadata.com/eosio:v1.0.1 \
       /opt/eosio/bin/nodeos --data-dir=/opt/eosio/bin/data-dir \
                             --genesis-json=/opt/eosio/bin/data-dir/genesis.json \
                             --config-dir=/opt/eosio/bin/data-dir

echo ""
echo "   View logs with: docker logs -f nodeos-sic"
echo ""

echo "Waiting 3 secs for nodeos to launch through Docker"
sleep 3

echo "Hit ENTER to continue"
read