const { ethers } = require("hardhat");

async function setupEventListener() {
    console.log("👂 设置事件监听器...\n");

    // 获取合约实例（需要先部署）
    const contractAddress = "0xc6e7DF5E7b4f2A278906862b61205850344D4e7d"; // 替换为你的合约地址
    const TransferDemo = await ethers.getContractFactory("TransferDemo");
    const transferDemo = TransferDemo.attach(contractAddress);

    // 监听存款事件
    transferDemo.on("Deposit", (from, amount, event) => {
        console.log("💰 存款事件:");
        console.log("  来自:", from);
        console.log("  金额:", ethers.formatEther(amount), "ETH");
        console.log("  交易哈希:", event.transactionHash);
        console.log("  区块号:", event.blockNumber);
        console.log("---");
    });

    // 监听取款事件
    transferDemo.on("Withdraw", (to, amount, event) => {
        console.log("🏧 取款事件:");
        console.log("  给:", to);
        console.log("  金额:", ethers.formatEther(amount), "ETH");
        console.log("  交易哈希:", event.transactionHash);
        console.log("---");
    });

    // 监听转账事件
    transferDemo.on("Transfer", (from, to, amount, event) => {
        console.log("🔄 转账事件:");
        console.log("  从:", from);
        console.log("  到:", to);
        console.log("  金额:", ethers.formatEther(amount), "ETH");
        console.log("  交易哈希:", event.transactionHash);
        console.log("---");
    });

    console.log("✅ 事件监听器已启动，等待事件...");
    console.log("按 Ctrl+C 停止监听\n");
}

// 运行监听器
setupEventListener().catch(console.error);