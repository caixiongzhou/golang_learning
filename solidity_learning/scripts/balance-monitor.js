// 余额监控脚本
const { ethers } = require("hardhat");

async function monitorBalances() {
    console.log("📈 余额监控启动...\n");

    const [signer1, signer2] = await ethers.getSigners();
    const contractAddress = "0x5FbDB2315678afecb367f032d93F642f64180aa3"; // 替换为你的合约地址

    const TransferDemo = await ethers.getContractFactory("TransferDemo");
    const transferDemo = TransferDemo.attach(contractAddress);

    async function printBalances() {
        console.log(`🕒 ${new Date().toLocaleTimeString()}`);

        // 账户ETH余额
        const ethBalance1 = await ethers.provider.getBalance(signer1.address);
        const ethBalance2 = await ethers.provider.getBalance(signer2.address);
        const contractEthBalance = await ethers.provider.getBalance(contractAddress);

        // 合约内余额
        const contractBalance1 = await transferDemo.getBalance(signer1.address);
        const contractBalance2 = await transferDemo.getBalance(signer2.address);
        const totalContractBalance = await transferDemo.getContractBalance();

        console.log("💰 账户ETH余额:");
        console.log(`  账户1: ${ethers.formatEther(ethBalance1)} ETH`);
        console.log(`  账户2: ${ethers.formatEther(ethBalance2)} ETH`);
        console.log(`  合约: ${ethers.formatEther(contractEthBalance)} ETH`);

        console.log("📊 合约内余额:");
        console.log(`  账户1: ${ethers.formatEther(contractBalance1)} ETH`);
        console.log(`  账户2: ${ethers.formatEther(contractBalance2)} ETH`);
        console.log(`  合约总: ${ethers.formatEther(totalContractBalance)} ETH`);

        console.log("---");
    }

    // 初始余额
    await printBalances();

    // 每10秒更新一次
    setInterval(printBalances, 10000);

    console.log("✅ 余额监控运行中，每10秒更新一次...");
    console.log("按 Ctrl+C 停止监控\n");
}

monitorBalances().catch(console.error);