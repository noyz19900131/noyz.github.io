// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

interface INFTAuction {
    struct AuctionInitParams {
        address seller;
        uint256 duration;
        uint256 startPrice;
        address nftContract;
        uint256 tokenId;
        address bidTokenAddress;
        address factory;
        address priceFeed;
        uint256 feedDecimals;
    }

    function initialize(AuctionInitParams calldata params) external;

    function placeBid(uint256 _bidAmount) external payable;

    function endAuction() external;
}
