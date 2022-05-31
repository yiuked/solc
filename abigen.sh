#!/bin/bash

# shellcheck disable=SC2164
cd contracts
ls
abigen -sol womtx.sol -pkg womtx -out ../womtx/womtx.go
abigen -sol womtx.sol -pkg womnft -out ../womnft/womnft.go
abigen -sol womtx.sol -pkg ewom -out ../ewom/ewom.go