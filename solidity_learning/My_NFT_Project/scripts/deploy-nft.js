const { ethers } = require("hardhat");

async function main() {
  console.log("🎨 开始部署 MyNFT 合约...");
  
  const [deployer] = await ethers.getSigners();
  console.log("部署者地址:", deployer.address);

  // 替换为您上传的元数据基础URI
  const baseTokenURI = "ipfs://您的元数据文件夹CID/";

  console.log("📝 合约参数:");
  console.log("  - 名称: MyNFT Collection");
  console.log("  - 符号: MNFT");
  console.log("  - 基础URI:", baseTokenURI);

  const MyNFT = await ethers.getContractFactory("MyNFT");
  console.log("⏳ 正在部署合约...");
  
  const nft = await MyNFT.deploy(
    "MyNFT Collection",    // name
    "MNFT",               // symbol
    baseTokenURI          // baseTokenURI
  );

  // 等待合约部署完成（兼容旧版本）
  console.log("⏳ 等待合约部署确认...");
  await nft.deployed(); // 使用 deployed() 而不是 waitForDeployment()
  
  const contractAddress = nft.address;
  console.log("✅ MyNFT 合约部署成功!");
  console.log("📄 合约地址:", contractAddress);
  console.log("👤 合约所有者:", deployer.address);
  console.log("🔗 在 Etherscan 查看: https://sepolia.etherscan.io/address/" + contractAddress);

  // 等待几个区块确认
  console.log("⏳ 等待区块确认...");
  await new Promise(resolve => setTimeout(resolve, 30000));

  return contractAddress;
}

main().catch((error) => {
  console.error("❌ 部署失败:", error);
  process.exitCode = 1;
});