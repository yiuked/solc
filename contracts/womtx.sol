// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./openzeppelin/contracts/utils/cryptography/ECDSA.sol";
import "./womnft.sol";
import "./ewom.sol";

contract WomTransfer is Ownable{
    EWOMToken public ewomToken;
    WomNFT public womNFT;
    // token limited
    mapping(uint256 => uint256) private _tokenLimited;


    constructor(EWOMToken _ewomAddr, WomNFT _nftAddr) {
        ewomToken = _ewomAddr;
        womNFT = _nftAddr;
    }

    //
    function MintTransfer(address seller, uint256 token, string memory _url, uint256 number, uint256 price, bytes memory signature) external {
        bytes32 digest = keccak256(abi.encode(seller, token, price));

        // 签名验证
        address recoveredSigner = ECDSA.recover(digest, signature);
        require(recoveredSigner == seller, "seller error");

        // 授权
        bool success = ewomToken.approve(address(this),price);
        require(success, "ewom approval error");

        // 铸造NFT
        uint256 tokenID = womNFT.sendNFT(seller, _url,number);
        require(tokenID >= 0, "mint token error");

        // 转账给卖家
        bool boo = ewomToken.transferFrom(msg.sender, seller, price);
        require(boo, "transfer fail");

        // 扣除手续费

        // 发NFT给买家(需要提前调用 womNFT的setApprovalForAll给当前合约)
        womNFT.safeTransferFrom(seller, msg.sender, token, number, "0x00");
    }
}
