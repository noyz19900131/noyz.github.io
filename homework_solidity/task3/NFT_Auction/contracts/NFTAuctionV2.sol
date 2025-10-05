// SPDX-License-Identifier: SEE LICENSE IN LICENSE
pragma solidity ^0.8.20;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract NFTAuctionV2 is Initializable, UUPSUpgradeable, OwnableUpgradeable {
    struct Auction {
        address seller; // 卖家
        uint256 duration; // 拍卖持续时间
        uint256 startPrice; // 起始价格
        uint256 startTime; // 开始时间
        bool ended; // 是否结束
        address highestBidder; // 最高出价者
        uint256 highestBid; // 最高出价
        address nftContract; // NFT合约地址
        uint256 tokenId; // NFT ID
        address bidTokenAddress; // 参与竞价的资产类型
    }

    mapping(uint256 => Auction) public auctions; // 状态变量
    uint256 public nextAuctionId; // 下一个拍卖ID
    mapping(address => AggregatorV3Interface) public priceFeeds; // 价格预言机
    mapping(address => uint256) public priceFeedDecimals; // 价格预言机小数位数

    // 事件
    event AuctionCreated(
        uint256 indexed auctionId,
        address indexed nftContract,
        uint256 indexed tokenId,
        uint256 startPrice,
        uint256 duration,
        address bidTokenAddress
    );
    event BidPlaced(
        uint256 indexed auctionId,
        address indexed bidder,
        uint256 bidAmount,
        address bidTokenAddress
    );
    event AuctionEnded(
        uint256 indexed auctionId,
        address indexed highestBidder,
        uint256 highestBid,
        address bidTokenAddress
    );

    function initialize() public initializer {
        __Ownable_init();
        nextAuctionId = 0;
    }

    // 设置 ETH 价格预言机
    function setPriceFeed(
        address tokenAddress,
        address _priceETHFeed,
        uint256 _decimals
    ) public onlyOwner {
        require(_priceETHFeed != address(0), "Invalid price feed address");
        require(_decimals > 0, "Decimals must be greater than 0");
        priceFeeds[tokenAddress] = AggregatorV3Interface(_priceETHFeed);
        priceFeedDecimals[tokenAddress] = _decimals;
    }

    // 获取最新的 ETH 价格
    function getLatestPrice(
        address tokenAddress
    ) public view returns (uint256 price, uint256 decimals) {
        AggregatorV3Interface priceFeed = priceFeeds[tokenAddress];
        require(address(priceFeed) != address(0), "Price feed not set");
        decimals = priceFeedDecimals[tokenAddress];
        require(decimals > 0, "Invalid decimals");

        (, int256 answer, , , ) = priceFeed.latestRoundData();
        require(answer > 0, "Invalid price");
        price = uint256(answer);
        return (price, decimals);
    }

    // 创建拍卖
    function createAuction(
        uint256 _duration,
        uint256 _startPrice,
        address _nftContract,
        uint256 _tokenId,
        address _bidTokenAddress
    ) external onlyOwner {
        // 参数校验
        require(_duration >= 10, "Duration must be greater than 10s");
        require(_startPrice > 0, "Start price must be greater than 0");
        require(_nftContract != address(0), "Invalid NFT contract address");
        require(_tokenId > 0, "Invalid token ID");

        // NFT 授权校验
        IERC721 nft = IERC721(_nftContract);
        require(
            nft.isApprovedForAll(msg.sender, address(this)) ||
                nft.getApproved(_tokenId) == address(this),
            "Contract not approved to transfer NFT"
        );

        // 转移NFT到合约
        IERC721(_nftContract).safeTransferFrom(
            msg.sender,
            address(this),
            _tokenId
        );

        auctions[nextAuctionId] = Auction(
            msg.sender,
            _duration,
            _startPrice,
            block.timestamp,
            false,
            address(0),
            0,
            _nftContract,
            _tokenId,
            _bidTokenAddress
        );

        emit AuctionCreated(
            nextAuctionId,
            _nftContract,
            _tokenId,
            _startPrice,
            _duration,
            _bidTokenAddress
        );

        nextAuctionId++;
    }

    // 买家参与买单
    function placeBid(uint256 _auctionId, uint256 _bidAmount) external payable {
        Auction storage auction = auctions[_auctionId];
        // 判断当前拍卖是否结束
        require(!auction.ended, "Auction has ended");
        require(
            block.timestamp < auction.startTime + auction.duration,
            "Auction timed out"
        );
        // 判断出价是否大于当前最高出价
        require(
            _bidAmount > auction.highestBid,
            "Bid must be higher than current highest bid"
        );
        // 判断出价是否大于起拍价
        require(
            _bidAmount >= auction.startPrice,
            "Bid must be higher than start price"
        );

        // 1. 退还前最高出价者的资金（区分ETH/ERC20）
        if (auction.highestBidder != address(0)) {
            if (auction.bidTokenAddress == address(0)) {
                // ETH退款
                (bool success, ) = payable(auction.highestBidder).call{
                    value: auction.highestBid
                }("");
                require(success, "Failed to refund ETH");
            } else {
                // ERC20退款
                IERC20 erc20 = IERC20(auction.bidTokenAddress);
                require(
                    erc20.transfer(auction.highestBidder, auction.highestBid),
                    "Failed to refund ERC20"
                );
            }
        }

        // 2. 接收当前出价者的资金（区分ETH/ERC20）
        if (auction.bidTokenAddress == address(0)) {
            // ETH出价
            require(msg.value == _bidAmount, "ETH amount mismatch");
        } else {
            // ERC20出价
            IERC20 erc20 = IERC20(auction.bidTokenAddress);
            // 校验出价者是否授权了足够的ERC20代币
            require(
                erc20.allowance(msg.sender, address(this)) >= _bidAmount,
                "Insufficient ERC20 allowance"
            );
            require(
                erc20.transferFrom(msg.sender, address(this), _bidAmount),
                "Failed to transfer ERC20"
            );
        }

        // 更新拍卖信息
        auction.highestBid = _bidAmount;
        auction.highestBidder = msg.sender;

        emit BidPlaced(
            _auctionId,
            msg.sender,
            _bidAmount,
            auction.bidTokenAddress
        );
    }

    // 结束拍卖
    function endAuction(uint256 _auctionId) external {
        Auction storage auction = auctions[_auctionId];
        // 判断拍卖是否结束
        require(!auction.ended, "Auction already ended");
        require(
            block.timestamp >= auction.startTime + auction.duration,
            "Auction not timed out"
        );
        require(auction.nftContract != address(0), "Auction does not exist");

        // 标记拍卖已结束
        auction.ended = true;

        // 处理NFT和资金：分两种情况（有出价/无出价）
        if (auction.highestBidder != address(0)) {
            // 情况1：有出价 → NFT转给最高出价者，资金转给卖家
            IERC721(auction.nftContract).safeTransferFrom(
                address(this),
                auction.highestBidder,
                auction.tokenId
            );

            // 转移资金（区分ETH/ERC20）
            if (auction.bidTokenAddress == address(0)) {
                // ETH资金转移
                (bool success, ) = payable(auction.seller).call{
                    value: address(this).balance
                }("");
                require(success, "Failed to transfer ETH to seller");
            } else {
                // ERC20资金转移
                IERC20 erc20 = IERC20(auction.bidTokenAddress);
                uint256 erc20Balance = erc20.balanceOf(address(this));
                require(
                    erc20.transfer(auction.seller, erc20Balance),
                    "Failed to transfer ERC20 to seller"
                );
            }
        } else {
            // 情况2：无出价 → NFT还给卖家（避免NFT卡住）
            IERC721(auction.nftContract).safeTransferFrom(
                address(this),
                auction.seller,
                auction.tokenId
            );
        }

        emit AuctionEnded(
            _auctionId,
            auction.highestBidder,
            auction.highestBid,
            auction.bidTokenAddress
        );
    }

    // UUPS升级授权（仅owner可升级）
    function _authorizeUpgrade(
        address newImplementation
    ) internal override onlyOwner {}

    // 接收ETH
    receive() external payable {}

    // 为区别合约V2和V1而设置该函数
    function testHello() public pure returns (string memory) {
        return "Hello";
    }
}
