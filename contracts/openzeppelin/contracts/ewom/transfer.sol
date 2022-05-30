// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./womfnt.sol";
import "./ewomtoken.sol";
import "../registry/AuthenticatedProxy.sol";

contract WomTransfer {
    EWOMToken public ewomToken;
    WomNFT public womNFT;

    constructor(EWOMToken _ewomAddr, WomNFT _nftAddr) {
        ewomToken = _ewomAddr;
        womNFT = _nftAddr;
    }
    
    // 发送交易
    function send(address seller,uint256 token,uint256 number,uint256 price) external {
        // 转账给卖家
        bool boo = ewomToken.transferFrom(msg.sender,seller,price);
        require(boo,"transfer fail");

        // 发NFT给买家(需要提前调用 womNFT的setApprovalForAll给当前合约)
        womNFT.safeTransferFrom(seller,msg.sender,token,number,"0x00");
    }

    // 赠送cion
    function coin(address _to,uint256 price) external{
        // 转账给卖家(需要提前调womToken合约中的approve)
        ewomToken.transferFrom(msg.sender,_to,price);
    }

    // 购买NFT
    function nft(address _to,uint256 token,uint256 number) external{
        //function send(address to,uint256 token,uint256 number) external{
        // 发NFT给买家(需要提前调用 womNFT的setApprovalForAll给当前合约)，扣
        womNFT.safeTransferFrom(msg.sender,_to,token,number,"0x00");
    }

    function registerProxy()
    public
    returns (OwnableDelegateProxy proxy)
    {
        require(proxies[msg.sender] == address(0));
        // 创建一个新的代理合约
        proxy = new OwnableDelegateProxy(msg.sender, delegateProxyImplementation, abi.encodeWithSignature("initialize(address,address)", msg.sender, address(this)));
        proxies[msg.sender] = proxy;
        return proxy;
    }
}

// dj 0x73A0BE25B795d588354f9663A38A138c0D54445F
// cs 0x7005EF493499EF1bD76584C889263373D329A1Fa
// cs 0xAbea2142ac4EeC7a4dc46B228258349CE733aC2b
