// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/IERC20Metadata.sol";
import "./INFTAuctionFactory.sol";
import "./INFTAuction.sol";

contract NFTAuction is
    Initializable,
    UUPSUpgradeable,
    OwnableUpgradeable,
    INFTAuction
{
    struct Auction {
        address seller;
        uint256 duration;
        uint256 startPrice;
        uint256 startPriceUSD;
        uint256 startTime;
        bool ended;
        address highestBidder;
        uint256 highestBid;
        uint256 highestBidUSD;
        address nftContract;
        uint256 tokenId;
        address bidTokenAddress; // 参与竞价的资产地址（ETH=address(0)）
    }

    Auction public auction;
    address public factory;
    mapping(address => AggregatorV3Interface) public priceFeeds;
    mapping(address => uint256) public priceFeedDecimals;

    event AuctionCreated(
        address indexed auctionContract,
        address indexed nftContract,
        uint256 indexed tokenId,
        address seller,
        uint256 startPrice,
        uint256 startPriceUSD,
        uint256 duration,
        address bidTokenAddress
    );
    event BidPlaced(
        address indexed auctionContract,
        address indexed bidder,
        uint256 bidAmount,
        uint256 bidAmountUSD,
        address bidTokenAddress
    );
    event AuctionEnded(
        address indexed auctionContract,
        address indexed highestBidder,
        address indexed seller,
        uint256 highestBid,
        address bidTokenAddress
    );

    // 初始化
    function initialize(AuctionInitParams calldata params) public initializer {
        __Ownable_init();
        _transferOwnership(params.seller);
        __UUPSUpgradeable_init();

        require(params.duration > 10, "Duration>10");
        require(params.startPrice > 0, "StartPrice>0");
        require(params.nftContract != address(0), "Invalid NFT");
        require(params.tokenId > 0, "TokenId>0");
        require(params.factory != address(0), "Invalid Factory");
        require(params.priceFeed != address(0), "Invalid Feed");
        require(params.feedDecimals > 0, "FeedDec>0");
        require(
            params.seller ==
                IERC721(params.nftContract).ownerOf(params.tokenId),
            "Seller is not NFT Owner"
        );

        priceFeeds[params.bidTokenAddress] = AggregatorV3Interface(
            params.priceFeed
        );
        priceFeedDecimals[params.bidTokenAddress] = params.feedDecimals;

        uint256 startPriceUSD = _calculateBidUSDValue(
            params.bidTokenAddress,
            params.startPrice
        );
        require(startPriceUSD > 0, "startPriceUSD > 0");

        IERC721 nft = IERC721(params.nftContract);
        require(
            nft.isApprovedForAll(params.seller, address(this)) ||
                nft.getApproved(params.tokenId) == address(this),
            "No NFT Approval"
        );

        IERC721(params.nftContract).safeTransferFrom(
            params.seller,
            address(this),
            params.tokenId
        );

        auction = Auction(
            params.seller,
            params.duration,
            params.startPrice,
            startPriceUSD,
            block.timestamp,
            false,
            address(0),
            0,
            0,
            params.nftContract,
            params.tokenId,
            params.bidTokenAddress
        );

        factory = params.factory;

        emit AuctionCreated(
            address(this),
            params.nftContract,
            params.tokenId,
            params.seller,
            params.startPrice,
            startPriceUSD,
            params.duration,
            params.bidTokenAddress
        );
    }

    // 买家竞拍
    function placeBid(uint256 _bidAmount) external payable {
        require(!auction.ended, "Auction ended");
        require(
            block.timestamp < auction.startTime + auction.duration,
            "Auction timeout"
        );

        uint256 bidUSD = _calculateBidUSDValue(
            auction.bidTokenAddress,
            _bidAmount
        );
        require(bidUSD > auction.highestBidUSD, "BidUSD not highest");
        require(bidUSD >= auction.startPriceUSD, "BidUSD < StartUSD");

        // 1. 退还前最高出价者的资金
        if (auction.highestBidder != address(0)) {
            if (auction.bidTokenAddress == address(0)) {
                // ETH
                (bool success, ) = payable(auction.highestBidder).call{
                    value: auction.highestBid
                }("");
                require(success, "Failed to refund ETH");
            } else {
                // ERC20
                IERC20 erc20 = IERC20(auction.bidTokenAddress);
                require(
                    erc20.transfer(auction.highestBidder, auction.highestBid),
                    "Failed to refund ERC20"
                );
            }
        }

        // 2. 接收当前出价者的资金
        if (auction.bidTokenAddress == address(0)) {
            // ETH
            require(msg.value == _bidAmount, "ETH amount mismatch");
        } else {
            // ERC20
            IERC20 erc20 = IERC20(auction.bidTokenAddress);
            require(
                erc20.allowance(msg.sender, address(this)) >= _bidAmount,
                "Insufficient ERC20 allowance"
            );
            require(
                erc20.transferFrom(msg.sender, address(this), _bidAmount),
                "Failed to transfer ERC20"
            );
        }

        auction.highestBid = _bidAmount;
        auction.highestBidUSD = bidUSD;
        auction.highestBidder = msg.sender;

        emit BidPlaced(
            address(this),
            msg.sender,
            _bidAmount,
            bidUSD,
            auction.bidTokenAddress
        );
    }

    // 结束拍卖
    function endAuction() external {
        address caller = msg.sender;
        require(
            caller == auction.seller ||
                caller == auction.highestBidder ||
                caller == factory,
            "No end perm"
        );

        require(!auction.ended, "Auction ended");
        require(
            block.timestamp >= auction.startTime + auction.duration,
            "Auction not timeout"
        );
        require(auction.nftContract != address(0), "Auction not exist");

        auction.ended = true;

        if (auction.highestBidder != address(0)) {
            IERC721(auction.nftContract).safeTransferFrom(
                address(this),
                auction.highestBidder,
                auction.tokenId
            );

            // 转移资金
            if (auction.bidTokenAddress == address(0)) {
                // ETH
                (bool success, ) = payable(auction.seller).call{
                    value: auction.highestBid
                }("");
                require(success, "Failed to transfer ETH to seller");
            } else {
                // ERC20
                IERC20 erc20 = IERC20(auction.bidTokenAddress);
                uint256 erc20Balance = erc20.balanceOf(address(this));
                require(
                    erc20.transfer(auction.seller, erc20Balance),
                    "Failed to transfer ERC20 to seller"
                );
            }
        } else {
            IERC721(auction.nftContract).safeTransferFrom(
                address(this),
                auction.seller,
                auction.tokenId
            );
        }

        INFTAuctionFactory(factory).clearAuctionMap(
            auction.nftContract,
            auction.tokenId
        );

        emit AuctionEnded(
            address(this),
            auction.highestBidder,
            auction.seller,
            auction.highestBid,
            auction.bidTokenAddress
        );
    }

    // 辅助函数：计算出价的 USD 价值
    function _calculateBidUSDValue(
        address _tokenAddress,
        uint256 _amount
    ) internal view returns (uint256) {
        AggregatorV3Interface feed = priceFeeds[_tokenAddress];
        require(
            address(feed) != address(0),
            "Price feed not set for bid token"
        );

        (, int256 priceRaw, , , ) = feed.latestRoundData();
        require(priceRaw > 0, "Invalid price from feed");
        uint256 price = uint256(priceRaw);
        uint256 feedDecimal = priceFeedDecimals[_tokenAddress];

        uint256 tokenDecimal = (_tokenAddress == address(0))
            ? 18
            : IERC20Metadata(_tokenAddress).decimals();

        return ((_amount * price) / (10 ** (tokenDecimal + feedDecimal))) * 1e8;
    }

    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyOwner {}

    receive() external payable {}
}
