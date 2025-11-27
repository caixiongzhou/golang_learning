const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("SimpleERC20 合约测试", function () {
    let SimpleERC20;
    let simpleERC20;
    let owner;
    let addr1;
    let addr2;
    let addrs;

    const tokenName = "MyToken";
    const tokenSymbol = "MTK";
    const decimals = 18;
    const initialSupply = "1000000";
    const initialSupplyWei = ethers.utils.parseUnits(initialSupply, decimals);

    beforeEach(async function () {
        // 获取签名者
        [owner, addr1, addr2, ...addrs] = await ethers.getSigners();

        // 部署合约
        SimpleERC20 = await ethers.getContractFactory("SimpleERC20");
        simpleERC20 = await SimpleERC20.deploy(
            tokenName,
            tokenSymbol,
            decimals,
            initialSupply  // 这里传递的是不带小数位的数量（1000000）
        );
    });

    describe("网络检查", function () {
        it("应该显示当前网络", async function () {
            const network = await ethers.provider.getNetwork();
            console.log("当前网络:", network.name, "链ID:", network.chainId);
        });
    });

    describe("部署", function () {
        it("应该正确设置代币信息", async function () {
            expect(await simpleERC20.name()).to.equal(tokenName);
            expect(await simpleERC20.symbol()).to.equal(tokenSymbol);
            expect(await simpleERC20.decimals()).to.equal(decimals);
        });

        it("应该将初始供应量分配给部署者", async function () {
            const ownerBalance = await simpleERC20.getBalance(owner.address);
            expect(ownerBalance).to.equal(initialSupplyWei);
        });

        it("应该正确设置总供应量", async function () {
            const totalSupply = await simpleERC20.totalSupply();
            expect(totalSupply).to.equal(initialSupplyWei);
        });

        it("应该设置正确的合约所有者", async function () {
            expect(await simpleERC20.owner()).to.equal(owner.address);
        });
    });

    describe("转账", function () {
        it("应该允许转账", async function () {
            const transferAmount = ethers.utils.parseUnits("100", decimals);

            // 转账前余额检查
            const initialOwnerBalance = await simpleERC20.getBalance(owner.address);
            const initialAddr1Balance = await simpleERC20.getBalance(addr1.address);

            // 执行转账
            await simpleERC20.transfer(addr1.address, transferAmount);

            // 转账后余额检查
            const finalOwnerBalance = await simpleERC20.getBalance(owner.address);
            const finalAddr1Balance = await simpleERC20.getBalance(addr1.address);

            expect(finalOwnerBalance).to.equal(initialOwnerBalance.sub(transferAmount));
            expect(finalAddr1Balance).to.equal(initialAddr1Balance.add(transferAmount));
        });

        it("应该触发 Transfer 事件", async function () {
            const transferAmount = ethers.utils.parseUnits("100", decimals);

            await expect(simpleERC20.transfer(addr1.address, transferAmount))
                .to.emit(simpleERC20, "Transfer")
                .withArgs(owner.address, addr1.address, transferAmount);
        });

        it("应该拒绝零地址转账", async function () {
            const transferAmount = ethers.utils.parseUnits("100", decimals);

            await expect(
                simpleERC20.transfer(ethers.constants.AddressZero, transferAmount)
            ).to.be.revertedWith("Transfer to zero address");
        });

        it("应该拒绝余额不足的转账", async function () {
            const excessiveAmount = ethers.utils.parseUnits("10000000", decimals);

            await expect(
                simpleERC20.connect(addr1).transfer(addr2.address, excessiveAmount)
            ).to.be.revertedWith("Insufficient balance");
        });
    });

    describe("授权和代扣转账", function () {
        const approveAmount = ethers.utils.parseUnits("500", decimals);

        beforeEach(async function () {
            // 所有者先给 addr1 转账一些代币用于测试
            await simpleERC20.transfer(addr1.address, ethers.utils.parseUnits("1000", decimals));
        });

        it("应该允许授权", async function () {
            await simpleERC20.connect(addr1).approve(addr2.address, approveAmount);

            const allowance = await simpleERC20.allowance(addr1.address, addr2.address);
            expect(allowance).to.equal(approveAmount);
        });

        it("应该触发 Approval 事件", async function () {
            await expect(simpleERC20.connect(addr1).approve(addr2.address, approveAmount))
                .to.emit(simpleERC20, "Approval")
                .withArgs(addr1.address, addr2.address, approveAmount);
        });

        it("应该允许代扣转账", async function () {
            const transferAmount = ethers.utils.parseUnits("300", decimals);

            // 先授权
            await simpleERC20.connect(addr1).approve(addr2.address, approveAmount);

            // 执行代扣转账
            await simpleERC20.connect(addr2).transferFrom(addr1.address, owner.address, transferAmount);

            // 检查余额变化
            const addr1Balance = await simpleERC20.getBalance(addr1.address);
            const ownerBalance = await simpleERC20.getBalance(owner.address);

            expect(addr1Balance).to.equal(ethers.utils.parseUnits("700", decimals));
            expect(ownerBalance).to.be.above(initialSupplyWei.sub(ethers.utils.parseUnits("1000", decimals)));

            // 检查授权额度减少
            const remainingAllowance = await simpleERC20.allowance(addr1.address, addr2.address);
            expect(remainingAllowance).to.equal(approveAmount.sub(transferAmount));
        });

        it("应该拒绝授权不足的代扣转账", async function () {
            const transferAmount = ethers.utils.parseUnits("600", decimals); // 超过授权额度

            // 先授权
            await simpleERC20.connect(addr1).approve(addr2.address, approveAmount);

            await expect(
                simpleERC20.connect(addr2).transferFrom(addr1.address, owner.address, transferAmount)
            ).to.be.revertedWith("Insufficient allowance");
        });
    });

    describe("增发代币", function () {
        it("只有所有者可以增发代币", async function () {
            const mintAmount = ethers.utils.parseUnits("50000", decimals);

            await expect(
                simpleERC20.connect(addr1).mint(addr1.address, mintAmount)
            ).to.be.revertedWith("Only owner can call this function");
        });

        it("所有者可以增发代币", async function () {
            const mintAmount = ethers.utils.parseUnits("50000", decimals);

            const initialTotalSupply = await simpleERC20.totalSupply();
            const initialAddr1Balance = await simpleERC20.getBalance(addr1.address);

            await simpleERC20.mint(addr1.address, mintAmount);

            const finalTotalSupply = await simpleERC20.totalSupply();
            const finalAddr1Balance = await simpleERC20.getBalance(addr1.address);

            expect(finalTotalSupply).to.equal(initialTotalSupply.add(mintAmount));
            expect(finalAddr1Balance).to.equal(initialAddr1Balance.add(mintAmount));
        });

        it("增发应该触发 Mint 和 Transfer 事件", async function () {
            const mintAmount = ethers.utils.parseUnits("50000", decimals);

            await expect(simpleERC20.mint(addr1.address, mintAmount))
                .to.emit(simpleERC20, "Mint")
                .withArgs(addr1.address, mintAmount)
                .and.to.emit(simpleERC20, "Transfer")
                .withArgs(ethers.constants.AddressZero, addr1.address, mintAmount);
        });

        it("应该拒绝向零地址增发", async function () {
            const mintAmount = ethers.utils.parseUnits("50000", decimals);

            await expect(
                simpleERC20.mint(ethers.constants.AddressZero, mintAmount)
            ).to.be.revertedWith("Mint to zero address");
        });
    });

    describe("批量转账", function () {
        it("应该执行批量转账", async function () {
            const recipients = [addr1.address, addr2.address];
            const values = [
                ethers.utils.parseUnits("100", decimals),
                ethers.utils.parseUnits("200", decimals)
            ];

            const initialOwnerBalance = await simpleERC20.getBalance(owner.address);

            await simpleERC20.batchTransfer(recipients, values);

            const finalOwnerBalance = await simpleERC20.getBalance(owner.address);
            expect(finalOwnerBalance).to.equal(initialOwnerBalance.sub(ethers.utils.parseUnits("300", decimals)));

            expect(await simpleERC20.getBalance(addr1.address)).to.equal(ethers.utils.parseUnits("100", decimals));
            expect(await simpleERC20.getBalance(addr2.address)).to.equal(ethers.utils.parseUnits("200", decimals));
        });

        it("应该拒绝长度不匹配的数组", async function () {
            const recipients = [addr1.address, addr2.address];
            const values = [ethers.utils.parseUnits("100", decimals)]; // 长度不匹配

            await expect(
                simpleERC20.batchTransfer(recipients, values)
            ).to.be.revertedWith("Arrays length mismatch");
        });

        it("应该拒绝余额不足的批量转账", async function () {
            const recipients = [addr1.address, addr2.address];
            const values = [
                ethers.utils.parseUnits("5000000", decimals),
                ethers.utils.parseUnits("5000000", decimals)
            ];

            await expect(
                simpleERC20.batchTransfer(recipients, values)
            ).to.be.revertedWith("Insufficient balance");
        });

        it("应该拒绝向零地址的批量转账", async function () {
            const recipients = [addr1.address, ethers.constants.AddressZero];
            const values = [
                ethers.utils.parseUnits("100", decimals),
                ethers.utils.parseUnits("200", decimals)
            ];

            await expect(
                simpleERC20.batchTransfer(recipients, values)
            ).to.be.revertedWith("Transfer to zero address");
        });
    });

    describe("查询功能", function () {
        it("应该正确查询授权额度", async function () {
            const approveAmount = ethers.utils.parseUnits("1000", decimals);

            await simpleERC20.approve(addr1.address, approveAmount);

            const allowance = await simpleERC20.getAllowance(owner.address, addr1.address);
            expect(allowance).to.equal(approveAmount);
        });

        it("应该正确查询余额", async function () {
            const balance = await simpleERC20.getBalance(owner.address);
            expect(balance).to.equal(initialSupplyWei);
        });
    });

    describe("所有权转移", function () {
        it("只有所有者可以转移所有权", async function () {
            await expect(
                simpleERC20.connect(addr1).transferOwnership(addr2.address)
            ).to.be.revertedWith("Only owner can call this function");
        });

        it("应该转移所有权", async function () {
            await simpleERC20.transferOwnership(addr1.address);
            expect(await simpleERC20.owner()).to.equal(addr1.address);
        });

        it("应该拒绝零地址作为新所有者", async function () {
            await expect(
                simpleERC20.transferOwnership(ethers.constants.AddressZero)
            ).to.be.revertedWith("New owner is zero address");
        });

        it("新所有者可以执行所有者功能", async function () {
            // 转移所有权
            await simpleERC20.transferOwnership(addr1.address);

            // 新所有者可以增发代币
            const mintAmount = ethers.utils.parseUnits("10000", decimals);
            await simpleERC20.connect(addr1).mint(addr2.address, mintAmount);

            const newBalance = await simpleERC20.getBalance(addr2.address);
            expect(newBalance).to.equal(mintAmount);
        });
    });
});