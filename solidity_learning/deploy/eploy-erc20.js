const { ethers } = require("hardhat");

async function main() {
    console.log("🚀 开始部署 ERC20 代币合约...");

    const [deployer] = await ethers.getSigners();
    console.log("部署者地址:", deployer.address);
    console.log("部署者余额:", (await deployer.provider.getBalance(deployer.address)).toString());

    // 代币参数
    const tokenName = "MyTestToken";
    const tokenSymbol = "MTT";
    const decimals = 18;
    const initialSupply = 1000000; // 100万代币

    console.log(`\n📝 代币信息:`);
    console.log(`名称: ${tokenName}`);
    console.log(`符号: ${tokenSymbol}`);
    console.log(`小数位: ${decimals}`);
    console.log(`初始供应: ${initialSupply} ${tokenSymbol}`);

    // 部署合约
    const SimpleERC20 = await ethers.getContractFactory("SimpleERC20");
    const token = await SimpleERC20.deploy(
        tokenName,
        tokenSymbol,
        decimals,
        initialSupply
    );

    await token.waitForDeployment();
    const tokenAddress = await token.getAddress();

    console.log("\n✅ ERC20 代币合约部署成功!");
    console.log("合约地址:", tokenAddress);
    console.log("合约所有者:", await token.owner());
    console.log("总供应量:", (await token.totalSupply()).toString());
    console.log("部署者余额:", (await token.balanceOf(deployer.address)).toString());

    // 保存部署信息到文件（可选）
    const fs = require('fs');
    const deploymentInfo = {
        network: "sepolia",
        timestamp: new Date().toISOString(),
        contractAddress: tokenAddress,
        tokenName: tokenName,
        tokenSymbol: tokenSymbol,
        decimals: decimals,
        initialSupply: initialSupply,
        deployer: deployer.address
    };

    fs.writeFileSync('deployment-info.json', JSON.stringify(deploymentInfo, null, 2));
    console.log("\n📄 部署信息已保存到 deployment-info.json");

    console.log("\n🎉 部署完成! 您可以将以下地址导入到钱包:");
    console.log(`合约地址: ${tokenAddress}`);
}

main().catch((error) => {
    console.error("❌ 部署失败:", error);
    process.exitCode = 1;
});