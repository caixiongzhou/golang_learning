const { ethers } = require("hardhat");

async function main() {
    console.log("📄 部署 TransferDemo 合约...");

    const [deployer] = await ethers.getSigners();
    console.log("部署者地址:", deployer.address);

    // 部署合约
    const TransferDemo = await ethers.getContractFactory("TransferDemo");
    const transferDemo = await TransferDemo.deploy();
    await transferDemo.waitForDeployment();

    const contractAddress = await transferDemo.getAddress();
    console.log("\n✅ 合约部署成功!");
    console.log("合约地址:", contractAddress);

    // 验证合约
    const owner = await transferDemo.owner();
    console.log("合约所有者:", owner);
    console.log("部署者与所有者匹配:", owner === deployer.address);

    console.log("\n📋 请将合约地址复制到 test-interaction.js 中的 CONTRACT_ADDRESS 变量");
}

main().catch((error) => {
    console.error("❌ 部署失败:", error);
    process.exitCode = 1;
});