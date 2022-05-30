// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "./womnft.sol";
import "./ewom.sol";

contract WomTransfer {
    EWOMToken public ewomToken;
    WomNFT public womNFT;

    constructor(EWOMToken _ewomAddr, WomNFT _nftAddr) {
        ewomToken = _ewomAddr;
        womNFT = _nftAddr;
    }

    // 发送交易
    function send(address seller,uint256 token,uint256 number,uint256 price, bytes memory signature) external{
        bytes32 digest = keccak256(abi.encode(seller, token, price));

        address recoveredSigner = ECDSA.recover(digest, signature);
        require(recoveredSigner == seller,"seller error");

        // 转账给卖家
        bool boo = ewomToken.transferFrom(msg.sender,seller,price);
        require(boo,"transfer fail");

        // 发NFT给买家(需要提前调用 womNFT的setApprovalForAll给当前合约)
        womNFT.safeTransferFrom(seller,msg.sender,token,number,"0x00");
    }
}
