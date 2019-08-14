## Build Your First Network (BYFN)

### Create Genesis Block

`../bin/configtxgen -profile ThreeOrgsOrdererGenesis -channelID avocado-sys-channel -outputBlock ./channel-artifacts/genesis.block`

### Create Channel Transaction block

`export CHANNEL_NAME=avocadochannel  && ../bin/configtxgen -profile ThreeOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME`

### Create Anchor Peer Blocks

`../bin/configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/FAFMSPanchors.tx -channelID $CHANNEL_NAME -asOrg FAFMSP`

`../bin/configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/SupplierMSPanchors.tx -channelID $CHANNEL_NAME -asOrg SupplierMSP`

`../bin/configtxgen -profile ThreeOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/TransportMSPanchors.tx -channelID $CHANNEL_NAME -asOrg TransportMSP`

### Start up the network
Run the following command:
`docker-compose -f docker-compose-cli.yaml up -d`

To stop the network run the following command:
`docker-compose -f docker-compose-cli.yaml down --volumes --remove-orphan`

### The CLI container
The CLI is used for instantiating chaincode and joining the channel. Keep this container window open as all the commands will be exectued from here when in development.

1. Go into the cli container
`docker exec -it cli bash`

2. Create the channel that the peers can join and contains the corresponding orderer which handles transactions.

`export CHANNEL_NAME=avocadochannel`

`peer channel create -o orderer.avocadofederation.com:7050 -c avocadochannel -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/avocadofederation.com/orderers/orderer.avocadofederation.com/msp/tlscacerts/tlsca.avocadofederation.com-cert.pem`


Running this command returns avocadochannel.block which will be used by peers to join the channel.

### Choose a peer to join the channel:

1. Joining peer0.enerdynamic.avocadofederation.com to the avocadochannel
`peer channel join -b avocadochannel.block`

2. Add the other peers from the other organizations to the channel
    (The current order of executions to add these organizations to the channel is as follows)

    1. Logistics
    2. SolarSupplier

    You can run all these commands at once
    `
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/faf.avocadofederation.com/users/Admin@faf.avocadofederation.com/msp CORE_PEER_ADDRESS=peer0.faf.avocadofederation.com:7051 CORE_PEER_LOCALMSPID="FAFMSP"
    CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/faf.avocadofederation.com/peers/peer0.faf.avocadofederation.com/tls/ca.crt peer channel join -b avocadochannel.block
    `
3. Install chaincode onto peers

    export CCVERSION=1.0;

    `peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/cc`

    Second Peer

    `CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/faf.avocadofederation.com/users/Admin@faf.avocadofederation.com/msp CORE_PEER_ADDRESS=peer0.faf.avocadofederation.com:7051 CORE_PEER_LOCALMSPID="FAFMSP"
    CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/faf.avocadofederation.com/peers/peer0.faf.avocadofederation.com/tls/ca.crt peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/cc`
    `

5. Instantiate Chaincode onto channel (right now only Enerdynamic and Logistics will instantiate chaincode)

    `
    peer chaincode instantiate -o orderer.avocadofederation.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/avocadofederation.com/orderers/orderer.avocadofederation.com/msp/tlscacerts/tlsca.avocadofederation.com-cert.pem -C $CHANNEL_NAME -n mycc -v 1.0 -c '{"Args":["init","a", "100"]}' -P "AND ('FAFMSP.peer')"
    `

6. Invoke Chaincode to create a box and store it to the blockchain network
`peer chaincode invoke -o orderer.avocadofederation.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/avocadofederation.com/orderers/orderer.avocadofederation.com/msp/tlscacerts/tlsca.avocadofederation.com-cert.pem -C $CHANNEL_NAME -n mycc  -c '{"Args":["sortToBox","1-01","Test Producer", "2019-04-03:23:21", "30"]}'`

7. Query the network for the box with ID 1-01. This should show all the fields of the box in a JSON query.
`peer chaincode query -o orderer.avocadofederation.com:7050 -C $CHANNEL_NAME -n mycc  -c '{"Args":["getBox","1-01"]}'`

8. Invoke Chaincode to start the precooling for a box with ID '1-01'
`peer chaincode invoke -o orderer.avocadofederation.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/avocadofederation.com/orderers/orderer.avocadofederation.com/msp/tlscacerts/tlsca.avocadofederation.com-cert.pem -C $CHANNEL_NAME -n mycc  -c '{"Args":["preCoolBox","1-01","2019-04-03:03:50:21", "4C"]}'`

9. Invoke Chaincode to start the load for refigerated transport for a box with ID '1-01'
`peer chaincode invoke -o orderer.avocadofederation.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/avocadofederation.com/orderers/orderer.avocadofederation.com/msp/tlscacerts/tlsca.avocadofederation.com-cert.pem -C $CHANNEL_NAME -n mycc  -c '{"Args":["LoadForRefigeratedTransport","1-01", "Ford", "TRANSIT CARGO 250", "Catering Van", "Harry", "-4C"]}'`
