const { ethers } = require("hardhat");

async function batchTest() {
    console.log("🧪 开始批量测试...\n");

    const [signer1, signer2] = await ethers.getSigners();

    // 部署合约
    const TransferDemo = await ethers.getContractFactory("TransferDemo");
    const transferDemo = await TransferDemo.deploy();
    await transferDemo.waitForDeployment();
    const contractAddress = await transferDemo.getAddress();
    console.log("合约地址:", contractAddress);

    const testCases = [
        {
            name: "小额存款",
            func: () => transferDemo.deposit({ value: ethers.parseEther("0.1") })
        },
        {
            name: "大额存款",
            func: () => transferDemo.deposit({ value: ethers.parseEther("1.0") })
        },
        {
            name: "合约内转账",
            func: () => transferDemo.transferTo(signer2.address, ethers.parseEther("0.3"))
        },
        {
            name: "直接转账",
            func: () => transferDemo.directTransfer(signer2.address, { value: ethers.parseEther("0.05") })
        },
        {
            name: "部分取款",
            func: () => transferDemo.withdraw(ethers.parseEther("0.2"))
        }
    ];

    for (let i = 0; i < testCases.length; i++) {
        const testCase = testCases[i];
        console.log(`\n${i + 1}. 测试: ${testCase.name}`);

        try {
            const tx = await testCase.func();
            await tx.wait();
            console.log("  ✅ 成功");
        } catch (error) {
            console.log("  ❌ 失败:", error.reason || error.message);
        }

        // 短暂延迟
        await new Promise(resolve => setTimeout(resolve, 1000));
    }

    console.log("\n🎉 批量测试完成!");
}

batchTest().catch(console.error);