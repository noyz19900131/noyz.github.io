// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.20;

import "./INFTAuction.sol";
import "./NFTAuction.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

contract NFTAuctionFactory is
    Initializable,
    UUPSUpgradeable,
    OwnableUpgradeable
{
    address[] private _auctions;
    mapping(address nftContract => mapping(uint256 tokenId => INFTAuction))
        private _auctionMap;
    mapping(address nftContract => mapping(uint256 tokenId => bool))
        private _isAuctioned;
    mapping(address bidToken => address priceFeed) private _bidTokenPriceFeeds;
    mapping(address bidToken => uint256 feedDecimals)
        private _bidTokenFeedDecimals;

    event AuctionCreated(
        address indexed auctionAddr,
        address indexed nftContract,
        uint256 indexed tokenId,
        address seller,
        address bidToken
    );
    event BidTokenFeedSet(address indexed bidToken, address feed, uint256 dec);

    function initialize() public initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();
    }

    function createAuction(
        uint256 duration,
        uint256 startPrice,
        address nftContract,
        uint256 tokenId,
        address bidToken
    ) external returns (address auctionAddr) {
        _validateNFT(nftContract, tokenId);
        (address feed, uint256 dec) = _getBidTokenFeed(bidToken);
        require(duration > 10, "D>10");
        require(startPrice > 0, "S>0");
        require(!_isAuctioned[nftContract][tokenId], "NFT:A");

        address seller = msg.sender;
        require(IERC721(nftContract).ownerOf(tokenId) == seller, "S!=O");

        INFTAuction auction = INFTAuction(address(new NFTAuction()));
        auctionAddr = address(auction);

        require(
            IERC721(nftContract).isApprovedForAll(seller, auctionAddr) ||
                IERC721(nftContract).getApproved(tokenId) == auctionAddr,
            "NFT:App"
        );

        auction.initialize(
            INFTAuction.AuctionInitParams({
                seller: seller,
                duration: duration,
                startPrice: startPrice,
                nftContract: nftContract,
                tokenId: tokenId,
                bidTokenAddress: bidToken,
                factory: address(this),
                priceFeed: feed,
                feedDecimals: dec
            })
        );

        _auctions.push(auctionAddr);
        _auctionMap[nftContract][tokenId] = auction;
        _isAuctioned[nftContract][tokenId] = true;

        emit AuctionCreated(
            auctionAddr,
            nftContract,
            tokenId,
            seller,
            bidToken
        );
    }

    // 1. 获取竞价资产的预言机
    function getBidTokenFeed(
        address bidToken
    ) external view returns (address feed, uint256 decimals) {
        return _getBidTokenFeed(bidToken);
    }

    // 2. 获取指定NFT的拍卖信息
    function getNFTAuctionInfo(
        address nftContract,
        uint256 tokenId
    ) external view returns (address auctionAddr, bool isAuctioning) {
        _validateNFT(nftContract, tokenId);
        auctionAddr = address(_auctionMap[nftContract][tokenId]);
        isAuctioning = _isAuctioned[nftContract][tokenId];
    }

    // 配置预言机
    function setBidTokenFeed(
        address bidToken,
        address feed,
        uint256 dec
    ) external onlyOwner {
        require(bidToken != address(0), "BT!0");
        require(feed != address(0), "Feed!0");
        require(dec > 0, "D>0");
        _bidTokenPriceFeeds[bidToken] = feed;
        _bidTokenFeedDecimals[bidToken] = dec;
        emit BidTokenFeedSet(bidToken, feed, dec);
    }

    // 分页查询
    function getAuctions(
        uint256 start,
        uint256 limit
    ) external view returns (address[] memory) {
        uint256 total = _auctions.length;
        if (start >= total) return new address[](0);
        limit = (start + limit > total) ? (total - start) : limit;

        address[] memory res = new address[](limit);
        for (uint256 i = 0; i < limit; i++) {
            res[i] = _auctions[start + i];
        }
        return res;
    }

    // 拍卖总数
    function getAuctionCount() external view returns (uint256) {
        return _auctions.length;
    }

    function getAuction(
        address nftContract,
        uint256 tokenId
    ) external view returns (address) {
        _validateNFT(nftContract, tokenId);
        require(_isAuctioned[nftContract][tokenId], "NFT:!A");
        address auctionAddr = address(_auctionMap[nftContract][tokenId]);
        require(auctionAddr != address(0), "A!0");
        return auctionAddr;
    }

    // 清理映射
    function clearAuctionMap(address nftContract, uint256 tokenId) external {
        _validateNFT(nftContract, tokenId);
        require(
            address(_auctionMap[nftContract][tokenId]) == msg.sender,
            "Only:A"
        );
        _isAuctioned[nftContract][tokenId] = false;
    }

    function _validateNFT(address nftContract, uint256 tokenId) internal pure {
        require(nftContract != address(0), "NFT!0");
        require(tokenId > 0, "T>0");
    }

    function _getBidTokenFeed(
        address bidToken
    ) internal view returns (address feed, uint256 dec) {
        feed = _bidTokenPriceFeeds[bidToken];
        dec = _bidTokenFeedDecimals[bidToken];
        require(feed != address(0), "BF!0");
        require(dec > 0, "BD>0");
    }

    // UUPS升级
    function _authorizeUpgrade(address newImpl) internal override onlyOwner {}
}
