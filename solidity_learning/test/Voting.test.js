const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Voting 合约测试", function () {
    let Voting;
    let voting;
    let owner;
    let voter1;
    let voter2;

    beforeEach(async function () {
        // 获取签名者
        [owner, voter1, voter2] = await ethers.getSigners();

        // 部署合约
        Voting = await ethers.getContractFactory("Voting");
        voting = await Voting.deploy();
        await voting.waitForDeployment();
    });

    describe("投票功能", function () {
        it("应该允许用户投票给候选人", async function () {
            // voter1 投票给 Alice
            await voting.connect(voter1).vote("Alice");

            // 检查 Alice 的票数
            const votes = await voting.getVotes("Alice");
            expect(votes).to.equal(1);
        });

        it("应该允许多个用户投票给同一个候选人", async function () {
            // 多个用户投票给 Alice
            await voting.connect(voter1).vote("Alice");
            await voting.connect(voter2).vote("Alice");

            const votes = await voting.getVotes("Alice");
            expect(votes).to.equal(2);
        });

        it("应该允许投票给不同的候选人", async function () {
            await voting.connect(voter1).vote("Alice");
            await voting.connect(voter2).vote("Bob");

            const aliceVotes = await voting.getVotes("Alice");
            const bobVotes = await voting.getVotes("Bob");

            expect(aliceVotes).to.equal(1);
            expect(bobVotes).to.equal(1);
        });

        it("不应该允许空字符串的候选人名称", async function () {
            await expect(
                voting.connect(voter1).vote("")
            ).to.be.revertedWith("Candidate name cannot be empty");
        });
    });

    describe("查询功能", function () {
        it("应该正确返回候选人的票数", async function () {
            // 先投几票
            await voting.connect(voter1).vote("Alice");
            await voting.connect(voter2).vote("Alice");
            await voting.connect(owner).vote("Bob");

            // 查询票数
            const aliceVotes = await voting.getVotes("Alice");
            const bobVotes = await voting.getVotes("Bob");
            const charlieVotes = await voting.getVotes("Charlie"); // 不存在的候选人

            expect(aliceVotes).to.equal(2);
            expect(bobVotes).to.equal(1);
            expect(charlieVotes).to.equal(0); // 不存在的候选人应该返回0
        });

        it("应该可以通过公共mapping直接查询", async function () {
            await voting.connect(voter1).vote("Alice");

            // 使用自动生成的公共getter
            const votes = await voting.candidateVotes("Alice");
            expect(votes).to.equal(1);
        });
    });

    describe("重置功能", function () {
        it("应该重置所有候选人的票数为0", async function () {
            // 先投一些票
            await voting.connect(voter1).vote("Alice");
            await voting.connect(voter2).vote("Bob");
            await voting.connect(owner).vote("Alice");

            // 验证投票结果
            expect(await voting.getVotes("Alice")).to.equal(2);
            expect(await voting.getVotes("Bob")).to.equal(1);

            // 重置票数
            await voting.connect(owner).resetVotes();

            // 验证重置后票数为0
            expect(await voting.getVotes("Alice")).to.equal(0);
            expect(await voting.getVotes("Bob")).to.equal(0);
        });

        it("重置后应该可以重新投票", async function () {
            // 第一轮投票
            await voting.connect(voter1).vote("Alice");
            await voting.resetVotes();

            // 第二轮投票
            await voting.connect(voter2).vote("Alice");

            const votes = await voting.getVotes("Alice");
            expect(votes).to.equal(1);
        });
    });

    describe("事件", function () {
        it("投票时应该触发 Voted 事件", async function () {
            await expect(voting.connect(voter1).vote("Alice"))
                .to.emit(voting, "Voted")
                .withArgs(voter1.address, "Alice", 1);
        });

        it("重置时应该触发 VotesReset 事件", async function () {
            await voting.connect(voter1).vote("Alice");

            await expect(voting.resetVotes())
                .to.emit(voting, "VotesReset")
                .withArgs(owner.address);
        });
    });

    describe("额外功能", function () {
        it("应该能获取所有候选人列表", async function () {
            await voting.connect(voter1).vote("Alice");
            await voting.connect(voter2).vote("Bob");
            await voting.connect(owner).vote("Charlie");

            const candidates = await voting.getAllCandidates();
            expect(candidates).to.have.lengthOf(3);
            expect(candidates).to.include.members(["Alice", "Bob", "Charlie"]);
        });

        it("应该能获取候选人数量", async function () {
            expect(await voting.getCandidateCount()).to.equal(0);

            await voting.connect(voter1).vote("Alice");
            expect(await voting.getCandidateCount()).to.equal(1);

            await voting.connect(voter2).vote("Bob");
            expect(await voting.getCandidateCount()).to.equal(2);
        });
    });
});