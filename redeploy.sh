#!/bin/bash

go run ./ clear

go run ./ deploy
res=$?
if [ $res -ne 0 ]; then
        echo "deploy fail"
        echo
        exit 1
fi
sleep 5

# go run ./ approval ewom && go run ./ approval nft
go run ./ approval nft
#res=$?
#if [ $res -ne 0 ]; then
#        echo "approval fail"
#        echo
#        exit 1
#fi