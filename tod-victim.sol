// SPDX-License-Identifier: MIT
pragma solidity >=0.7.0;

contract TODVictim {
    uint public stockQuantity;
    address public owner;
    uint public price;
    constructor() {
        owner = msg.sender;
        price = 1;
        stockQuantity = 100;
    }
    function updatePrice(uint _price) public {
        price = _price;
    }
    function buy(uint _quantity) public payable {
        require(msg.value >= _quantity * price && _quantity <= stockQuantity, "");
        stockQuantity -= _quantity;
        msg.sender.transfer(msg.value - _quantity * price);
    }
}
