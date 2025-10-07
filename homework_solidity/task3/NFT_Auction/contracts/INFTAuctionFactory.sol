// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.20;

interface INFTAuctionFactory {
    function clearAuctionMap(address nftContract, uint256 tokenId) external;
}
