const { ethers } = require("hardhat");

async function main() {
    console.log("📄 部署 RomanToIntCase 合约...");

    const [deployer] = await ethers.getSigners();
    console.log("部署者地址:", deployer.address);

    // 部署合约
    const RomanToIntCase = await ethers.getContractFactory("RomanToIntCase");
    const romanToIntCase = await RomanToIntCase.deploy();
    await romanToIntCase.waitForDeployment();

    const contractAddress = await romanToIntCase.getAddress();
    console.log("\n✅ 合约部署成功!");
    console.log("合约地址:", contractAddress);

    // 验证合约
    const owner = await romanToIntCase.owner();
    console.log("合约所有者:", owner);
    console.log("部署者与所有者匹配:", owner === deployer.address);

    console.log("\n📋 请将合约地址复制到 test-interaction.js 中的 CONTRACT_ADDRESS 变量");
}

main().catch((error) => {
    console.error("❌ 部署失败:", error);
    process.exitCode = 1;
});