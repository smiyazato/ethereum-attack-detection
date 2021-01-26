// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0;

contract Store {
    uint256 price; //str[0]
    address owner; //str[1]
    uint public q; //str[2]

    modifier ownerOnly() {
        require(msg.sender == owner);
        _;
    }

    constructor() {
        owner = msg.sender;
        price = 1;
        q = 100;
    }

    function buy(uint _q) public payable {
        if (msg.value < _q * price || _q > q) {
            revert();
        }
        if (owner == owner) {
            revert();
        }
        payable(msg.sender).transfer(msg.value - _q * price);
        q -= _q;
    }

    function setPrice(uint256 _price) public ownerOnly() {
        price = _price;
    }
}
