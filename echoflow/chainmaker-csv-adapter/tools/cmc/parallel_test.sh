#!/usr/bin/env bash
#调用合约
./cmc parallel invoke \
--loopNum=1000  \
--printTime=5 \
--threadNum=200 \
--timeout=1800 \
--sleepTime=10 \
--climbTime=1  \
--use-tls=true \
--showKey=false \
--contract-name=fact \
--method="save" \
--pairs="" \
--hosts="127.0.0.1:12301" \
--pairs-file="./testdata/save.json" \
--org-IDs="wx-org1.chainmaker.org" \
--user-keys="./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key" \
--user-crts="./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt" \
--output-result=false \
--ca-path="./testdata/crypto-config/wx-org1.chainmaker.org/ca" \
--admin-sign-keys="./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.key" \
--admin-sign-crts="./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.crt"
