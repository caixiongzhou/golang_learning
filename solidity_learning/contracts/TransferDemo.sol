// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract TransferDemo {
    address public owner;
    mapping(address => uint256) public balances;
    uint256 public contractBalance;

    event Deposit(address indexed from, uint256 amount);
    event Withdraw(address indexed to, uint256 amount);
    event Transfer(address indexed from, address indexed to, uint256 amount);

    constructor() {
        owner = msg.sender;
    }

    // 存款到合约
    function deposit() public payable {
        require(msg.value > 0, "Amount must be greater than 0");
        balances[msg.sender] += msg.value;
        contractBalance += msg.value;
        emit Deposit(msg.sender, msg.value);
    }

    // 从合约取款
    function withdraw(uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");
        require(address(this).balance >= amount, "Contract insufficient balance");

        balances[msg.sender] -= amount;
        contractBalance -= amount;

        payable(msg.sender).transfer(amount);
        emit Withdraw(msg.sender, amount);
    }

    // 转账给其他用户
    function transferTo(address payable to, uint256 amount) public {
        require(balances[msg.sender] >= amount, "Insufficient balance");

        balances[msg.sender] -= amount;
        balances[to] += amount;

        emit Transfer(msg.sender, to, amount);
    }

    // 直接ETH转账
    function directTransfer(address payable to) public payable {
        require(msg.value > 0, "Amount must be greater than 0");
        to.transfer(msg.value);
        emit Transfer(msg.sender, to, msg.value);
    }

    // 获取用户余额
    function getBalance(address user) public view returns (uint256) {
        return balances[user];
    }

    // 获取合约ETH余额
    function getContractBalance() public view returns (uint256) {
        return address(this).balance;
    }

    // 接收ETH的回退函数
    receive() external payable {
        deposit();
    }
}