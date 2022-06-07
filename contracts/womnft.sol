// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./openzeppelin/contracts/access/Ownable.sol";
import "./openzeppelin/contracts/token/ERC1155/presets/ERC1155PresetMinterPauser.sol";
import "./openzeppelin/contracts/utils/Counters.sol";
import "./openzeppelin/contracts/token/ERC1155/extensions/ERC1155URIStorage.sol";

contract WomNFT is ERC1155URIStorage, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIds;

    constructor() ERC1155(""){}

    function sendNFT(address creator, string memory _url,uint256 amount) public returns (uint256) {
        require(creator != address(0), "address is null");
        require(amount > 0, "amount is 0");

        uint256 tokenId = _tokenIds.current();
        _mint(creator, tokenId, amount, "");
        _setURI(tokenId,_url);

        _tokenIds.increment();
        return tokenId;
    }
}

// ewom : 0x780f24e640659681cE5Aa9c02B16a62C55a01227
// nft : 0x2B9B0BC37eD25ED18883C88ede07F97BC6F90296
// tx: 0xa05e0881aEB7a4d11f5af30Ee61E62E7106596D8
