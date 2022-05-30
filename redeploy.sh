#!/bin/bash

go run ./ deploy
sleep 5
go run ./ mint

go run ./ approval ewom && go run ./ approval nft