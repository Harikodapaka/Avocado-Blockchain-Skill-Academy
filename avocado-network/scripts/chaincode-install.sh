echo "Inside the CLI container..."
export CHANNEL_NAME=avocadochannel

echo "Creating avocadochannel..."
peer channel create -o orderer.avocadofederation.com:7050 -c avocadochannel -f ./channel-artifacts/channel.tx --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/avocadofederation.com/orderers/orderer.avocadofederation.com/msp/tlscacerts/tlsca.avocadofederation.com-cert.pem

echo "Joining peers into avocadochannel..."
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/faf.avocadofederation.com/users/Admin@faf.avocadofederation.com/msp
CORE_PEER_ADDRESS=peer0.faf.avocadofederation.com:7051 CORE_PEER_LOCALMSPID="FAFMSP"
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/faf.avocadofederation.com/peers/peer0.faf.avocadofederation.com/tls/ca.crt peer channel join -b avocadochannel.block

echo "Installing Chaincode in peers..."

CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/faf.avocadofederation.com/users/Admin@faf.avocadofederation.com/msp
CORE_PEER_ADDRESS=peer0.faf.avocadofederation.com:7051 CORE_PEER_LOCALMSPID="FAFMSP"
CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/faf.avocadofederation.com/peers/peer0.faf.avocadofederation.com/tls/ca.crt peer chaincode install -n mycc -v 1.0 -p github.com/chaincode/cc

echo "Instantiating Chaincode..."
peer chaincode instantiate -o orderer.avocadofederation.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/avocadofederation.com/orderers/orderer.avocadofederation.com/msp/tlscacerts/tlsca.avocadofederation.com-cert.pem -C $CHANNEL_NAME -n mycc -v 1.0 -c '{"Args":["init","a", "100"]}' -P "AND ('FAFMSP.peer')"

echo "Inserting Box in the ClockChain..."
sleep 5
peer chaincode invoke -o orderer.avocadofederation.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/avocadofederation.com/orderers/orderer.avocadofederation.com/msp/tlscacerts/tlsca.avocadofederation.com-cert.pem -C $CHANNEL_NAME -n mycc  -c '{"Args":["sortToBox","1-01","Test Producer", "2019-04-03:23:21", "30"]}'

sleep 3
echo "Retreving the Box State..."
peer chaincode query -o orderer.avocadofederation.com:7050 -C $CHANNEL_NAME -n mycc  -c '{"Args":["getBox","1-01"]}'

sleep 3
echo "Pre Cooling The Box..."
peer chaincode invoke -o orderer.avocadofederation.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/avocadofederation.com/orderers/orderer.avocadofederation.com/msp/tlscacerts/tlsca.avocadofederation.com-cert.pem -C $CHANNEL_NAME -n mycc  -c '{"Args":["preCoolBox","1-01","2019-04-03:03:50:21", "4C"]}'

sleep 3
echo "Retreving the Box State..."
peer chaincode query -o orderer.avocadofederation.com:7050 -C $CHANNEL_NAME -n mycc  -c '{"Args":["getBox","1-01"]}'

sleep 3
echo "Transporting the Box As Lot..."
peer chaincode invoke -o orderer.avocadofederation.com:7050 --tls --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/avocadofederation.com/orderers/orderer.avocadofederation.com/msp/tlscacerts/tlsca.avocadofederation.com-cert.pem -C $CHANNEL_NAME -n mycc  -c '{"Args":["LoadForRefigeratedTransport","1-01", "Ford", "TRANSIT CARGO 250", "Catering Van", "Harry", "-4"]}'

sleep 3
echo "Retreving the Box State..."
peer chaincode query -o orderer.avocadofederation.com:7050 -C $CHANNEL_NAME -n mycc  -c '{"Args":["getBox","1-01"]}'
