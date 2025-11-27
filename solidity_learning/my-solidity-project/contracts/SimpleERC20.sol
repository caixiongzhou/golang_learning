// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleERC20 {
    // 代币基本信息
    string public name;
    string public symbol;
    uint8 public decimals;
    uint256 public totalSupply;

    // 合约所有者
    address public owner;
    // 余额映射
    mapping(address => uint256) public balanceOf;

    // 授权映射（owner =>(spender => amount)）
    mapping(address => mapping(address => uint256)) public allowance;

    // 事件定义
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
    event Mint(address indexed to, uint256 value);

    // 修饰器:只有所有者可以调用
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    /**
     * @dev 构造函数
     * @param _name 代币名称
     * @param _symbol 代币符号
     * @param _decimals 小数位数
     * @param _initialSupply 初始供应量
     */
    constructor(
        string memory _name,
        string memory _symbol,
        uint8 _decimals,
        uint256 _initialSupply
    ) {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        owner = msg.sender;

        // 初始发行给合约部署者
        _mint(msg.sender, _initialSupply * 10 ** decimals);
    }

    /**
     * @dev 查询账户余额
     * @param account 要查询的账户地址
     * @return 账户余额
     */
    function getBalance(address account) public view returns (uint256) {
        return balanceOf[account];
    }

    /**
     * @dev 转账功能
     * @param to 收款人地址
     * @param value 转账金额
     * @return 是否成功
     *
     * 功能：调用者将自己的代币转给别人,msg.sender → to,不需要授权，直接操作自己的资金
     */
    function transfer(address to, uint256 value) public returns (bool) {
        require(to != address(0), "Transfer to zero address");
        require(balanceOf[msg.sender] >= value, "Insufficient balance");

        balanceOf[msg.sender] -= value;
        balanceOf[to] += value;

        emit Transfer(msg.sender, to, value);
        return true;
    }

    /**
     * @dev 授权功能
     * @param spender 被授权人地址
     * @param value 授权金额
     * @return 是否成功
     */
    function approve(address spender, uint256 value) public returns (bool) {
        require(spender != address(0), "Approve to zero address");

        allowance[msg.sender][spender] = value;

        emit Approval(msg.sender, spender, value);
        return true;
    }

    /**
     * @dev 代扣转账功能
     * @param from 付款人地址
     * @param to 收款人地址
     * @param value 转账金额
     * @return 是否成功
     *
     * 功能：调用者代替别人转账（需要授权）from → to ,（但由 msg.sender 操作）,需要 from 提前授权给 msg.sender
     */
    function transferFrom(address from, address to, uint256 value) public returns (bool) {
        require(from != address(0), "Transfer from zero address");
        require(to != address(0), "Transfer to zero address");
        require(balanceOf[from] >= value, "Insufficient balance");
        require(allowance[from][msg.sender] >= value, "Insufficient allowance");

        balanceOf[from] -= value;
        balanceOf[to] += value;
        allowance[from][msg.sender] -= value;

        emit Transfer(from, to, value);
        return true;
    }

    /**
     * @dev 内部增发函数
     * @param to 接收地址
     * @param value 增发金额
     */
    function _mint(address to, uint256 value) internal {
        require(to != address(0), "Mint to zero address");

        balanceOf[to] += value;
        totalSupply += value;

        emit Transfer(address(0), to, value);
        emit Mint(to, value);
    }

    /**
     * @dev 增发代币（仅所有者）
     * @param to 接收地址
     * @param value 增发金额
     * @return 是否成功
     */
    function mint(address to, uint256 value) public onlyOwner returns (bool) {
        _mint(to, value);
        return true;
    }

    /**
     * @dev 批量转账
     * @param recipients 收款人地址数组
     * @param values 转账金额数组
     * @return 是否成功
     */
    function batchTransfer(address[] memory recipients, uint256[] memory values) public returns (bool) {
        require(recipients.length == values.length, "Arrays length mismatch");

        uint256 total = 0;
        for (uint256 i = 0; i < values.length; i++) {
            total += values[i];
        }

        require(balanceOf[msg.sender] >= total, "Insufficient balance");

        for (uint256 i = 0; i < recipients.length; i++) {
            require(recipients[i] != address(0), "Transfer to zero address");
            balanceOf[msg.sender] -= values[i];
            balanceOf[recipients[i]] += values[i];
            emit Transfer(msg.sender, recipients[i], values[i]);
        }

        return true;
    }

    /**
     * @dev 查询授权额度
     * @param _owner 授权人
     * @param spender 被授权人
     * @return 授权额度
     */
    function getAllowance(address _owner, address spender) public view returns (uint256) {
        return allowance[_owner][spender];
    }

    /**
     * @dev 转移合约所有权
     * @param newOwner 新所有者地址
     */
    function transferOwnership(address newOwner) public onlyOwner {
        require(newOwner != address(0), "New owner is zero address");
        owner = newOwner;
    }
}