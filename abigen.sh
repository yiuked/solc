#!/bin/bash

# shellcheck disable=SC2164
cd contracts
ls
abigen -sol womtx.sol -pkg womtx -out ../womtx/womtx.go