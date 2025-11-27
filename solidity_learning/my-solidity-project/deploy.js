const { ethers } = require("hardhat");

async function main() {
    console.log("开始部署 SimpleERC20 合约...");

    const tokenName = "MyToken";
    const tokenSymbol = "MTK";
    const decimals = 18;
    const initialSupply = "1000000"; // 100万代币（不带小数位）

    const [deployer] = await ethers.getSigners();
    console.log("部署者地址:", deployer.address);

    const SimpleERC20 = await ethers.getContractFactory("SimpleERC20");
    const simpleERC20 = await SimpleERC20.deploy(
        tokenName,
        tokenSymbol,
        decimals,
        initialSupply
    );

    await simpleERC20.deployed();
    const contractAddress = simpleERC20.address;

    console.log("SimpleERC20 合约部署成功!");
    console.log("合约地址:", contractAddress);
    console.log("代币名称:", tokenName);
    console.log("代币符号:", tokenSymbol);
    console.log("小数位数:", decimals);
    console.log("初始供应量:", initialSupply, "MTK");
    console.log("总供应量（含小数位）:", ethers.utils.formatUnits(await simpleERC20.totalSupply(), decimals), "MTK");

    // 等待更多区块确认
    console.log("等待区块确认...");
    await simpleERC20.deployTransaction.wait(10); // 等待10个区块确认

    console.log("开始验证合约...");

    // 使用 hardhat 的 run 方法来验证
    try {
        await hre.run("verify:verify", {
            address: contractAddress,
            constructorArguments: [tokenName, tokenSymbol, decimals, initialSupply],
        });
        console.log("合约验证成功!");
    } catch (error) {
        if (error.message.includes("Already Verified")) {
            console.log("合约已经验证过了");
        } else {
            console.log("合约验证失败:", error.message);
            console.log("你可以手动验证，使用以下命令:");
            console.log(`npx hardhat verify --constructor-args arguments.js ${contractAddress} --network sepolia`);
        }
    }
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });