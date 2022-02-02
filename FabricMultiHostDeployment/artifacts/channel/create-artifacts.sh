
chmod -R 0755 ./crypto-config
# Delete existing artifacts
rm -rf ./crypto-config
rm genesis.block channel_1.tx
rm -rf ../../channel-artifacts/*

#Generate Crypto artifactes for organizations
# cryptogen generate --config=./crypto-config.yaml --output=./crypto-config/



# System channel
SYS_CHANNEL="sys-channel"

# channel name defaults to "channel_1"
CHANNEL_NAME="channel_1"

echo $CHANNEL_NAME

# Generate System Genesis block
configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL  -outputBlock ./genesis.block


# Generate channel configuration block
configtxgen -profile BasicChannel -configPath . -outputCreateChannelTx ./channel_1.tx -channelID $CHANNEL_NAME

echo "#######    Generating anchor peer update for DiggipetMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./DiggipetMSPanchors.tx -channelID $CHANNEL_NAME -asOrg DiggipetMSP

echo "#######    Generating anchor peer update for Vet_orgMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./VetorgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg VetorgMSP

echo "#######    Generating anchor peer update for Breeders_orgMSP  ##########"
configtxgen -profile BasicChannel -configPath . -outputAnchorPeersUpdate ./BreedersorgMSPanchors.tx -channelID $CHANNEL_NAME -asOrg BreedersorgMSP