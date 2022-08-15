#!/bin/bash
#拉取镜像
docker pull hyperledger/fabric-ca:1.5.0
docker tag hyperledger/fabric-ca:1.5.0 hyperledger/fabric-ca:latest

docker pull hyperledger/fabric-tools:2.2
docker tag hyperledger/fabric-tools:2.2 hyperledger/fabric-tools:latest

docker pull hyperledger/fabric-ccenv:2.2
docker tag hyperledger/fabric-ccenv:2.2 hyperledger/fabric-ccenv:latest

docker pull hyperledger/fabric-orderer:2.2
docker tag hyperledger/fabric-orderer:2.2 hyperledger/fabric-orderer:latest

docker pull hyperledger/fabric-peer:2.2
docker tag hyperledger/fabric-peer:2.2 hyperledger/fabric-peer:latest

docker pull hyperledger/fabric-couchdb
docker tag hyperledger/fabric-couchdb hyperledger/fabric-couchdb:latest
docker pull hyperledger/fabric-couchdb:2.2
docker tag hyperledger/fabric-couchdb:2.2 hyperledger/fabric-couchdb:latest

docker pull hyperledger/fabric-baseos:2.2
docker tag hyperledger/fabric-baseos:2.2 hyperledger/fabric-baseos:latest

#生成创世区块和通道文件
cryptogen generate --config ./crypto-config.yaml
configtxgen -profile configtx.yaml -profile SampleMultiNodeEtcdRaft -channelID systemchannel -outputBlock ./channel-artifacts/genesis.block
configtxgen -profile designChannel -outputCreateChannelTx ./channel-artifacts/designchannel.tx -channelID designchannel
configtxgen -profile processChannel -outputCreateChannelTx ./channel-artifacts/processchannel.tx -channelID processchannel
configtxgen -profile assembleChannel -outputCreateChannelTx ./channel-artifacts/assemblechannel.tx -channelID assemblechannel
configtxgen -profile qualityChannel -outputCreateChannelTx ./channel-artifacts/qualitychannel.tx -channelID qualitychannel
configtxgen -profile contractChannel -outputCreateChannelTx ./channel-artifacts/contractchannel.tx -channelID contractchannel
configtxgen -profile manufactureChannel -outputCreateChannelTx ./channel-artifacts/manufacturechannel.tx -channelID manufacturechannel
configtxgen -profile superviseChannel -outputCreateChannelTx ./channel-artifacts/supervisechannel.tx -channelID supervisechannel

#设置各通道锚节点
configtxgen -profile designChannel -outputAnchorPeersUpdate ./channel-artifacts/DesignMSPanchors_designchannel.tx -channelID designchannel -asOrg DesignMSP
configtxgen -profile designChannel -outputAnchorPeersUpdate ./channel-artifacts/DesignMSPanchors_designchannel.tx -channelID designchannel -asOrg DesignMSP

configtxgen -profile processChannel -outputAnchorPeersUpdate ./channel-artifacts/ProcessMSPanchors.tx -channelID processchannel -asOrg ProcessMSP
configtxgen -profile processChannel -outputAnchorPeersUpdate ./channel-artifacts/ProcessMSPanchors.tx -channelID processchannel -asOrg ProcessMSP

configtxgen -profile assembleChannel -outputAnchorPeersUpdate ./channel-artifacts/AssembleMSPanchors.tx -channelID assemblechannel -asOrg AssembleMSP
configtxgen -profile assembleChannel -outputAnchorPeersUpdate ./channel-artifacts/AssembleMSPanchors.tx -channelID assemblechannel -asOrg AssembleMSP

configtxgen -profile qualityChannel -outputAnchorPeersUpdate ./channel-artifacts/QualityMSPanchors.tx -channelID qualitychannel -asOrg QualityMSP
configtxgen -profile qualityChannel -outputAnchorPeersUpdate ./channel-artifacts/QualityMSPanchors.tx -channelID qualitychannel -asOrg QualityMSP

configtxgen -profile contractChannel -outputAnchorPeersUpdate ./channel-artifacts/ContractMSPanchors_contractchannel.tx -channelID contractchannel -asOrg ContractMSP
configtxgen -profile contractChannel -outputAnchorPeersUpdate ./channel-artifacts/ContractMSPanchors_contractchannel.tx -channelID contractchannel -asOrg ContractMSP

configtxgen -profile manufactureChannel -outputAnchorPeersUpdate ./channel-artifacts/ManufactureMSPanchors_manufacturechannel.tx -channelID manufacturechannel -asOrg ManufactureMSP
configtxgen -profile manufactureChannel -outputAnchorPeersUpdate ./channel-artifacts/ManufactureMSPanchors_manufacturechannel.tx -channelID manufacturechannel -asOrg ManufactureMSP

configtxgen -profile superviseChannel -outputAnchorPeersUpdate ./channel-artifacts/SuperviseMSPanchors_supervisechannel.tx -channelID supervisechannel -asOrg SuperviseMSP

#启动cli
docker-compose -f ./docker-compose-orderer.yaml up -d
docker-compose -f ./docker-compose-couchdb.yaml up -d
docker-compose -f ./docker-compose-cli.yaml up -d

#创建通道文件
docker exec design_peer0 peer channel create -o orderer.huniuping.com:7050 -c designchannel -f ./channel-artifacts/designchannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec process_peer0 peer channel create -o orderer.huniuping.com:7050 -c processchannel -f ./channel-artifacts/processchannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec assemble_peer0 peer channel create -o orderer.huniuping.com:7050 -c assemblechannel -f ./channel-artifacts/assemblechannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec quality_peer0 peer channel create -o orderer.huniuping.com:7050 -c qualitychannel -f ./channel-artifacts/qualitychannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec contract_peer0 peer channel create -o orderer.huniuping.com:7050 -c contractchannel -f ./channel-artifacts/contractchannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec manufacture_peer0 peer channel create -o orderer.huniuping.com:7050 -c manufacturechannel -f ./channel-artifacts/manufacturechannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec supervise_peer0 peer channel create -o orderer.huniuping.com:7050 -c supervisechannel -f ./channel-artifacts/supervisechannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
#复制通道文件到宿主机
docker cp design_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/designchannel.block /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block
docker cp process_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/processchannel.block /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block
docker cp assemble_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/assemblechannel.block /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block
docker cp quality_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/qualitychannel.block /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block
docker cp contract_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/contractchannel.block /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block
docker cp manufacture_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/manufacturechannel.block /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block
#复制通道文件到目标主机
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block/designchannel.block supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block/processchannel.block supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block/assemblechannel.block supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block/qualitychannel.block supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block/contractchannel.block supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/block/manufacturechannel.block supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer


#更新组织锚节点
docker exec design_peer0 peer channel update -o orderer.huniuping.com:7050 -c designchannel -f ./channel-artifacts/DesignMSPanchors_designchannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec process_peer0 peer channel update -o orderer.huniuping.com:7050 -c processchannel -f ./channel-artifacts/ProcessMSPanchors.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec assemble_peer0 peer channel update -o orderer.huniuping.com:7050 -c assemblechannel -f ./channel-artifacts/AssembleMSPanchors.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec quality_peer0 peer channel update -o orderer.huniuping.com:7050 -c qualitychannel -f ./channel-artifacts/QualityMSPanchors.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec contract_peer0 peer channel update -o orderer.huniuping.com:7050 -c contractchannel -f ./channel-artifacts/ContractMSPanchors_contractchannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec manufacture_peer0 peer channel update -o orderer.huniuping.com:7050 -c manufacturechannel -f ./channel-artifacts/ManufactureMSPanchors_manufacturechannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem
docker exec supervise_peer0 peer channel update -o orderer.huniuping.com:7050 -c supervisechannel -f ./channel-artifacts/SuperviseMSPanchors_supervisechannel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem

#加入通道
docker exec supervise_peer0 peer channel join -b supervisechannel.block

docker exec design_peer0 peer channel join -b designchannel.block
docker exec supervise_peer0 peer channel join -b designchannel.block

docker exec process_peer0 peer channel join -b processchannel.block
docker exec supervise_peer0 peer channel join -b processchannel.block

docker exec assemble_peer0 peer channel join -b assemblechannel.block
docker exec supervise_peer0 peer channel join -b assemblechannel.block

docker exec quality_peer0 peer channel join -b qualitychannel.block
docker exec supervise_peer0 peer channel join -b qualitychannel.block

docker exec contract_peer0 peer channel join -b contractchannel.block
docker exec supervise_peer0 peer channel join -b contractchannel.block

docker exec manufacture_peer0 peer channel join -b manufacturechannel.block
docker exec supervise_peer0 peer channel join -b manufacturechannel.block

#查看通道
docker exec supervise_peer0 peer channel list
docker exec design_peer0 peer channel list
docker exec process_peer0 peer channel list
docker exec assemble_peer0 peer channel list
docker exec quality_peer0 peer channel list
docker exec contract_peer0 peer channel list
docker exec manufacture_peer0 peer channel list

#打包链码
docker exec supervise_peer0 peer lifecycle chaincode package trancecc.tar.gz --path /opt/gopath/src/github.com/hyperledger/fabric-samples/chaincode/trance --lang golang --label trance

docker exec design_peer0 peer lifecycle chaincode package drawingcc.tar.gz --path /opt/gopath/src/github.com/hyperledger/fabric-samples/chaincode/drawing --lang golang --label drawing

docker exec process_peer0 peer lifecycle chaincode package processcc.tar.gz --path /opt/gopath/src/github.com/hyperledger/fabric-samples/chaincode/process --lang golang --label process

docker exec assemble_peer0 peer lifecycle chaincode package assemblecc.tar.gz --path /opt/gopath/src/github.com/hyperledger/fabric-samples/chaincode/assemble --lang golang --label assemble

docker exec quality_peer0 peer lifecycle chaincode package qualitycc.tar.gz --path /opt/gopath/src/github.com/hyperledger/fabric-samples/chaincode/quality --lang golang --label quality

docker exec contract_peer0 peer lifecycle chaincode package contractcc.tar.gz --path /opt/gopath/src/github.com/hyperledger/fabric-samples/chaincode/contract --lang golang --label contract

docker exec manufacture_peer0 peer lifecycle chaincode package gongdancc.tar.gz --path /opt/gopath/src/github.com/hyperledger/fabric-samples/chaincode/gongdan --lang golang --label gongdan

#复制打包的链码到其他节点
docker cp supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/trancecc.tar.gz /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/trancecc.tar.gz design_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/trancecc.tar.gz process_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/trancecc.tar.gz assemble_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/trancecc.tar.gz quality_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/trancecc.tar.gz contract_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/trancecc.tar.gz manufacture_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer

docker cp design_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/drawingcc.tar.gz /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/drawingcc.tar.gz supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer

docker cp process_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/processcc.tar.gz /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/processcc.tar.gz supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer

docker cp assemble_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/assemblecc.tar.gz /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/assemblecc.tar.gz supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer

docker cp quality_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/qualitycc.tar.gz /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/qualitycc.tar.gz supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer

docker cp contract_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/contractcc.tar.gz /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/contractcc.tar.gz supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer

docker cp manufacture_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer/gongdancc.tar.gz /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp
docker cp /home/myfabric/code/go/src/raft-fabric-project/fabric-network/chaincode-temp/gongdancc.tar.gz supervise_peer0:/opt/gopath/src/github.com/hyperledger/fabric/peer


########

#安装链码
docker exec supervise_peer0 peer lifecycle chaincode install drawingcc.tar.gz
docker exec supervise_peer0 peer lifecycle chaincode install processcc.tar.gz
docker exec supervise_peer0 peer lifecycle chaincode install assemblecc.tar.gz
docker exec supervise_peer0 peer lifecycle chaincode install qualitycc.tar.gz
docker exec supervise_peer0 peer lifecycle chaincode install contractcc.tar.gz
docker exec supervise_peer0 peer lifecycle chaincode install gongdancc.tar.gz

docker exec design_peer0 peer lifecycle chaincode install drawingcc.tar.gz
docker exec process_peer0 peer lifecycle chaincode install processcc.tar.gz
docker exec assemble_peer0 peer lifecycle chaincode install assemblecc.tar.gz
docker exec quality_peer0 peer lifecycle chaincode install qualitycc.tar.gz
docker exec contract_peer0 peer lifecycle chaincode install contractcc.tar.gz
docker exec manufacture_peer0 peer lifecycle chaincode install gongdancc.tar.gz

#查询链码安装
docker exec supervise_peer0 peer lifecycle chaincode queryinstalled
docker exec design_peer0 peer lifecycle chaincode queryinstalled
docker exec process_peer0 peer lifecycle chaincode queryinstalled
docker exec assemble_peer0 peer lifecycle chaincode queryinstalled
docker exec quality_peer0 peer lifecycle chaincode queryinstalled
docker exec contract_peer0 peer lifecycle chaincode queryinstalled
docker exec manufacture_peer0 peer lifecycle chaincode queryinstalled


#审批认证链码
docker exec design_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID designchannel --name drawingcc --version 1.0 --init-required --package-id drawing:4a4677d65aa543bde1bb008a7d2957d65be6e5c57db50c8568e4854d846c4212 --sequence 1 --waitForEvent
docker exec supervise_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID designchannel --name drawingcc --version 1.0 --init-required --package-id drawing:4a4677d65aa543bde1bb008a7d2957d65be6e5c57db50c8568e4854d846c4212 --sequence 1 --waitForEvent

docker exec process_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID processchannel --name processcc --version 1.0 --init-required --package-id process:0f166ae0e39693f92b883b8a867ec71ec3f32a09081d84c1ce92cbaab4c4ff4f --sequence 1 --waitForEvent
docker exec supervise_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID processchannel --name processcc --version 1.0 --init-required --package-id process:0f166ae0e39693f92b883b8a867ec71ec3f32a09081d84c1ce92cbaab4c4ff4f --sequence 1 --waitForEvent

docker exec assemble_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID assemblechannel --name assemblecc --version 1.0 --init-required --package-id assemble:ec0fb9a38230498a077b22126791457e79247107057313700ac29244d279312a --sequence 1 --waitForEvent
docker exec supervise_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID assemblechannel --name assemblecc --version 1.0 --init-required --package-id assemble:ec0fb9a38230498a077b22126791457e79247107057313700ac29244d279312a --sequence 1 --waitForEvent

docker exec quality_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID qualitychannel --name qualitycc --version 1.0 --init-required --package-id quality:ae44152194cade7abbba1f10053d59d391fa03e77ced05ab65dcfba842a645a5 --sequence 1 --waitForEvent
docker exec supervise_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID qualitychannel --name qualitycc --version 1.0 --init-required --package-id quality:ae44152194cade7abbba1f10053d59d391fa03e77ced05ab65dcfba842a645a5 --sequence 1 --waitForEvent

docker exec contract_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID contractchannel --name contractcc --version 1.0 --init-required --package-id contract:53a5e530d1ce65cc325dfe7c5e66cb81605d9cb9f8af1ceeb8e1b0470e05df6f --sequence 1 --waitForEvent
docker exec supervise_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID contractchannel --name contractcc --version 1.0 --init-required --package-id contract:53a5e530d1ce65cc325dfe7c5e66cb81605d9cb9f8af1ceeb8e1b0470e05df6f --sequence 1 --waitForEvent

docker exec manufacture_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID manufacturechannel --name gongdancc --version 1.0 --init-required --package-id gongdan:8200d13887630396e536a72255cb2ff0483c0075a7b7c9d965f8d93b7dc9c639 --sequence 1 --waitForEvent
docker exec supervise_peer0 peer lifecycle chaincode approveformyorg --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID manufacturechannel --name gongdancc --version 1.0 --init-required --package-id gongdan:8200d13887630396e536a72255cb2ff0483c0075a7b7c9d965f8d93b7dc9c639 --sequence 1 --waitForEvent


#查看认证结果
docker exec design_peer0 peer lifecycle chaincode checkcommitreadiness --channelID designchannel --name drawingcc --version 1.0 --sequence 1 --output json --init-required

docker exec process_peer0 peer lifecycle chaincode checkcommitreadiness --channelID processchannel --name processcc --version 1.0 --sequence 1 --output json --init-required

docker exec assemble_peer0 peer lifecycle chaincode checkcommitreadiness --channelID assemblechannel --name assemblecc --version 1.0 --sequence 1 --output json --init-required

docker exec quality_peer0 peer lifecycle chaincode checkcommitreadiness --channelID qualitychannel --name qualitycc --version 1.0 --sequence 1 --output json --init-required

docker exec contract_peer0 peer lifecycle chaincode checkcommitreadiness --channelID contractchannel --name contractcc --version 1.0 --sequence 1 --output json --init-required

docker exec manufacture_peer0 peer lifecycle chaincode checkcommitreadiness --channelID manufacturechannel --name gongdancc --version 1.0 --sequence 1 --output json --init-required


#提交链码
docker exec design_peer0 peer lifecycle chaincode commit -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID designchannel --name drawingcc --peerAddresses peer0.design.huniuping.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/design.huniuping.com/peers/peer0.design.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --version 1.0 --sequence 1 --init-required

docker exec process_peer0 peer lifecycle chaincode commit -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID processchannel --name processcc --peerAddresses peer0.process.huniuping.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/process.huniuping.com/peers/peer0.process.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --version 1.0 --sequence 1 --init-required

docker exec assemble_peer0 peer lifecycle chaincode commit -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID assemblechannel --name assemblecc --peerAddresses peer0.assemble.huniuping.com:11051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/assemble.huniuping.com/peers/peer0.assemble.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --version 1.0 --sequence 1 --init-required

docker exec quality_peer0 peer lifecycle chaincode commit -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID qualitychannel --name qualitycc --peerAddresses peer0.quality.huniuping.com:13051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/quality.huniuping.com/peers/peer0.quality.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --version 1.0 --sequence 1 --init-required

docker exec contract_peer0 peer lifecycle chaincode commit -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID contractchannel --name contractcc --peerAddresses peer0.contract.huniuping.com:15051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/contract.huniuping.com/peers/peer0.contract.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --version 1.0 --sequence 1 --init-required

docker exec manufacture_peer0 peer lifecycle chaincode commit -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem --channelID manufacturechannel --name gongdancc --peerAddresses peer0.manufacture.huniuping.com:19051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/manufacture.huniuping.com/peers/peer0.manufacture.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --version 1.0 --sequence 1 --init-required

#链码提交查询
docker exec supervise_peer0 peer lifecycle chaincode querycommitted --channelID supervisechannel --name trancecc
docker exec design_peer0 peer lifecycle chaincode querycommitted --channelID designchannel --name trancecc
docker exec design_peer0 peer lifecycle chaincode querycommitted --channelID designchannel --name drawingcc
docker exec process_peer0 peer lifecycle chaincode querycommitted --channelID processchannel --name processcc
docker exec assemble_peer0 peer lifecycle chaincode querycommitted --channelID assemblechannel --name assemblecc
docker exec account_peer0 peer lifecycle chaincode querycommitted --channelID accountchannel --name accountcc
docker exec manufacture_peer0 peer lifecycle chaincode querycommitted --channelID manufacturechannel --name gongdancc
#初始化链码
docker exec design_peer0 peer chaincode invoke -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem -C designchannel -n drawingcc --peerAddresses peer0.design.huniuping.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/design.huniuping.com/peers/peer0.design.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --isInit -c '{"Args":["Init","drawing-001","drawing name","f3ahn213gjvhr5hn323hj124bh123jntytrghfgh","contract-001","technology-001"]}'

docker exec process_peer0 peer chaincode invoke -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem -C processchannel -n processcc --peerAddresses peer0.process.huniuping.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/process.huniuping.com/peers/peer0.process.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --isInit -c '{"Args":["Init","component-001","component name","work_order-001","product-001","qwdsadf32werfdfgertffbvcberedfdsfdsa","a-e-f","c82yoIFw0VKrr2YYhP1RL1l7PxqrQQwwqVuHqcsOD2b1rDjdrCMvCnEnQJPBzYFC9xpc+jxfnysihmxzGpNK55MLyW/24ljwkNn/Lz209kidwn2pthBJD9gK9rAk7Y2TcWo1aZ75tAPHvFiaOQKarIUdicgsKrV4PCE/w33XJNA"]}'

docker exec assemble_peer0 peer chaincode invoke -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem -C assemblechannel -n assemblecc --peerAddresses peer0.assemble.huniuping.com:11051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/assemble.huniuping.com/peers/peer0.assemble.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --isInit -c '{"Args":["Init","product-001","product name","work_order-002","assembleline-001","2022/01/01","hjfdh54fdsbgbu6dfddvgi86ygiungb233rf45gfbd","uyiujadnfbhsfng8uij23kref768yuhb 247gfihb3h9u238u2bj3bd74fbriuj32"]}'

docker exec quality_peer0 peer chaincode invoke -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem -C qualitychannel -n qualitycc --peerAddresses peer0.quality.huniuping.com:13051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/quality.huniuping.com/peers/peer0.quality.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --isInit -c '{"Args":["Init","product-001","2022/4/28","good","123456","alice"]}'

docker exec contract_peer0 peer chaincode invoke -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem -C contractchannel -n contractcc --peerAddresses peer0.contract.huniuping.com:15051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/contract.huniuping.com/peers/peer0.contract.huniuping.com/tls/ca.crt --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt --isInit -c '{"Args":["Init","contract-001","contract_name","buyer_company","seller_company","product name","50","fvdj576uy2j1ehd7fyiguh2n3efrhuobjk2092oiu3brjef"]}'

docker exec manufacture_peer0 peer chaincode invoke -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem -C manufacturechannel -n gongdancc --peerAddresses peer0.manufacture.huniuping.com:19051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/manufacture.huniuping.com/peers/peer0.manufacture.huniuping.com/tls/ca.crt --isInit --peerAddresses peer0.supervise.huniuping.com:17051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/supervise.huniuping.com/peers/peer0.supervise.huniuping.com/tls/ca.crt -c '{"Args":["Init","work_order-001","uiyfhusbdvhvushuiffnweouhisdfjnhufejfrew","2 week","drawing-001","contract-001"]}'

#查询链码
docker exec design_peer0 peer chaincode query -C designchannel -n drawingcc -c '{"Args":["query","drawing-001"]}'
docker exec supervise_peer0 peer chaincode query -C designchannel -n drawingcc -c '{"Args":["trance","query","product-001","assemblecc","assemblechannel"]}'

#新增链码调用
docker exec process_peer0 peer chaincode invoke -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem -C processchannel -n processcc --peerAddresses peer0.process.huniuping.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/process.huniuping.com/peers/peer0.process.huniuping.com/tls/ca.crt -c '{"Args":["set","component-001","component name","work_order-001","product-001","qwdsadf32werfdfgertffbvcberedfdsfdsa","a-e-f","c82yoIFw0VKrr2YYhP1RL1l7PxqrQQwwqVuHqcsOD2b1rDjdrCMvCnEnQJPBzYFC9xpc+jxfnysihmxzGpNK55MLyW/24ljwkNn/Lz209kidwn2pthBJD9gK9rAk7Y2TcWo1aZ75tAPHvFiaOQKarIUdicgsKrV4PCE/w33XJNA"]}'
docker exec process_peer0 peer chaincode invoke -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem -C processchannel -n processcc --peerAddresses peer0.process.huniuping.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/process.huniuping.com/peers/peer0.process.huniuping.com/tls/ca.crt -c '{"Args":["set","component-002","component name","work_order-001","product-001","qwdsadf32werfdfgertffbvcberedfdsfdsa","c-d-f-c-b","JJMaurW5o6PXTjrjRkIUEYQJRhflVkoqcA2+r6zj+ZCCk9JTjGC8Y/g+EyZd9JKP+JhVrEuF3XKYRY4tjOcfNzidpXKUdEJvhWnGbczaWdLpLmniGR/vqcTglBU1qd9ec+IsSQfQEUE2XEbPR7124POGQ9yt1Y51MVQL+1AKcH0"]}'
docker exec process_peer0 peer chaincode invoke -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem -C processchannel -n processcc --peerAddresses peer0.process.huniuping.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/process.huniuping.com/peers/peer0.process.huniuping.com/tls/ca.crt -c '{"Args":["set","component-003","component name","work_order-001","product-001","qwdsadf32werfdfgertffbvcberedfdsfdsa","f-a-f","f95LMo6Nw3psClFF2eYQ2UbbsJhxYomfd7RkCz2k/WTYplKOlwR6R8+LuLtMbelmFIwso0NJqKSNdslGjjT6SnE//gYQfXY4aReoxYfc7U1Q/cSoQ8l3tM0Z1FP0dBc4Nv9+LASzlmu6P/fLSRoW0mt5fP2RNLRJ+zSeD6Oeo7c"]}'

docker exec manufacture_peer0 peer chaincode invoke -o orderer.huniuping.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/huniuping.com/orderers/orderer.huniuping.com/msp/tlscacerts/tlsca.huniuping.com-cert.pem -C manufacturechannel -n gongdancc --peerAddresses peer0.manufacture.huniuping.com:23051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/manufacture.huniuping.com/peers/peer0.manufacture.huniuping.com/tls/ca.crt  -c '{"Args":["set","work_order-002","uifhjfjdsd34r23rfwegrgu12e1r3tfgb2gggh4t","1 week","drawing-001","contract-001"]}'
