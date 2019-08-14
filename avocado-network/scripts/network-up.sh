cd ..

../bin/configtxgen -profile ThreeOrgsOrdererGenesis -channelID avocado-sys-channel -outputBlock ./channel-artifacts/genesis.block

export CHANNEL_NAME=avocadochannel  && ../bin/configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME

../bin/configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/FAFMSPanchors.tx -channelID $CHANNEL_NAME -asOrg FAFMSP

../bin/configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/SupplierMSPanchors.tx -channelID $CHANNEL_NAME -asOrg SupplierMSP

../bin/configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/TransportMSPanchors.tx -channelID $CHANNEL_NAME -asOrg TransportMSP

docker-compose -f docker-compose-cli.yaml up -d


sleep 7

echo "Network is up and running..."

docker exec -it cli bash -c "./scripts/chaincode-install.sh"
