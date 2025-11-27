const { ethers } = require("hardhat");

// 已部署的合约地址 - 替换为你实际部署的地址
const CONTRACT_ADDRESS = "0xc6e7DF5E7b4f2A278906862b61205850344D4e7d";

async function main() {
    console.log("🚀 开始基于已部署合约的交互测试...\n");
    console.log("📄 使用合约地址:", CONTRACT_ADDRESS);

    // 获取测试账户  这里获取的签名器数组，第一个就是默认操作账户
    const [signer1, signer2, signer3] = await ethers.getSigners();
    console.log("📝 测试账户:");
    console.log("  账户1 (操作账户):", signer1.address);
    console.log("  账户2 (接收账户):", signer2.address);
    console.log("  账户3 (观察账户):", signer3.address);

    // 检查初始余额
    console.log("\n💰 初始余额:");
    const balance1 = await ethers.provider.getBalance(signer1.address);
    const balance2 = await ethers.provider.getBalance(signer2.address);
    const contractBalance = await ethers.provider.getBalance(CONTRACT_ADDRESS);

    console.log("  账户1余额:", ethers.formatEther(balance1), "ETH");
    console.log("  账户2余额:", ethers.formatEther(balance2), "ETH");
    console.log("  合约余额:", ethers.formatEther(contractBalance), "ETH");

    // 连接到已部署的合约
    console.log("\n🔗 连接到已部署合约...");
    const TransferDemo = await ethers.getContractFactory("TransferDemo");
    //连接到合约时，使用的是 signer1 作为默认签名器
    const transferDemo = TransferDemo.attach(CONTRACT_ADDRESS);
    // 等价于：const transferDemo = TransferDemo.connect(signer1).attach(CONTRACT_ADDRESS);


    // 验证合约连接
    try {
        const contractOwner = await transferDemo.owner();
        console.log("  ✅ 合约连接成功");
        console.log("  合约所有者:", contractOwner);
    } catch (error) {
        console.log("  ❌ 合约连接失败:", error.message);
        return;
    }

    // 测试1: 存款到合约
    console.log("\n1️⃣ 测试存款功能...");
    const depositAmount = ethers.parseEther("0.5");

    console.log("  准备存款:", ethers.formatEther(depositAmount), "ETH");

    const depositTx = await transferDemo.deposit({ value: depositAmount });
    console.log("  交易已发送，等待确认...");
    await depositTx.wait();
    console.log("  ✅ 存款成功");

    // 验证存款结果
    const userBalanceAfterDeposit = await transferDemo.getBalance(signer1.address);
    const contractBalanceAfterDeposit = await transferDemo.getContractBalance();
    console.log("  用户合约余额:", ethers.formatEther(userBalanceAfterDeposit), "ETH");
    console.log("  合约总余额:", ethers.formatEther(contractBalanceAfterDeposit), "ETH");

    // 测试2: 账户间转账（在合约内）
    console.log("\n2️⃣ 测试合约内转账...");
    const transferAmount = ethers.parseEther("0.1");

    console.log("  准备转账:", ethers.formatEther(transferAmount), "ETH");
    console.log("  从:", signer1.address);
    console.log("  到:", signer2.address);

    const transferTx = await transferDemo.transferTo(signer2.address, transferAmount);
    console.log("  交易已发送，等待确认...");
    await transferTx.wait();
    console.log("  ✅ 转账成功");

    // 验证转账结果
    const balanceAfterTransfer1 = await transferDemo.getBalance(signer1.address);
    const balanceAfterTransfer2 = await transferDemo.getBalance(signer2.address);
    console.log("  转账后账户1合约余额:", ethers.formatEther(balanceAfterTransfer1), "ETH");
    console.log("  转账后账户2合约余额:", ethers.formatEther(balanceAfterTransfer2), "ETH");

    // 测试3: 直接ETH转账
    console.log("\n3️⃣ 测试直接ETH转账...");
    const directAmount = ethers.parseEther("0.05");

    console.log("  准备直接转账:", ethers.formatEther(directAmount), "ETH");

    const directTx = await transferDemo.directTransfer(signer2.address, { value: directAmount });
    console.log("  交易已发送，等待确认...");
    await directTx.wait();
    console.log("  ✅ 直接转账成功");

    // 测试4: 取款
    console.log("\n4️⃣ 测试取款功能...");
    const withdrawAmount = ethers.parseEther("0.2");

    console.log("  准备取款:", ethers.formatEther(withdrawAmount), "ETH");

    const withdrawTx = await transferDemo.withdraw(withdrawAmount);
    console.log("  交易已发送，等待确认...");
    await withdrawTx.wait();
    console.log("  ✅ 取款成功");

    // 最终余额检查
    console.log("\n📊 最终余额统计:");

    // 区块链余额
    const finalBalance1 = await ethers.provider.getBalance(signer1.address);
    const finalBalance2 = await ethers.provider.getBalance(signer2.address);
    const finalContractBalance = await ethers.provider.getBalance(CONTRACT_ADDRESS);

    console.log("💰 区块链余额:");
    console.log("  账户1余额:", ethers.formatEther(finalBalance1), "ETH");
    console.log("  账户2余额:", ethers.formatEther(finalBalance2), "ETH");
    console.log("  合约余额:", ethers.formatEther(finalContractBalance), "ETH");

    // 合约内余额
    const contractStateBalance = await transferDemo.getContractBalance();
    const user1StateBalance = await transferDemo.getBalance(signer1.address);
    const user2StateBalance = await transferDemo.getBalance(signer2.address);

    console.log("\n📋 合约内余额记录:");
    console.log("  合约记录的总余额:", ethers.formatEther(contractStateBalance), "ETH");
    console.log("  合约记录的账户1余额:", ethers.formatEther(user1StateBalance), "ETH");
    console.log("  合约记录的账户2余额:", ethers.formatEther(user2StateBalance), "ETH");

    // 余额变化统计
    console.log("\n📈 余额变化统计:");
    const initialTotal = balance1 + balance2 + contractBalance;
    const finalTotal = finalBalance1 + finalBalance2 + finalContractBalance;
    const gasCost = initialTotal - finalTotal;

    console.log("  初始总余额:", ethers.formatEther(initialTotal), "ETH");
    console.log("  最终总余额:", ethers.formatEther(finalTotal), "ETH");
    console.log("  Gas 总消耗:", ethers.formatEther(gasCost), "ETH");

    console.log("\n🎉 基于已部署合约的测试完成!");
}

main().catch((error) => {
    console.error("❌ 错误:", error);
    process.exitCode = 1;
});