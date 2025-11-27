console.log("=== 环境变量检查 ===");
console.log("PRIVATE_KEY:", process.env.PRIVATE_KEY ? "已设置" : "未设置");
console.log("INFURA_PROJECT_ID:", process.env.INFURA_PROJECT_ID ? "已设置" : "未设置");
console.log("ETHERSCAN_API_KEY:", process.env.ETHERSCAN_API_KEY ? "已设置" : "未设置");
console.log("SEPOLIA_URL:", process.env.SEPOLIA_URL ? "已设置" : "未设置");

if (process.env.PRIVATE_KEY) {
    console.log("PRIVATE_KEY 长度:", process.env.PRIVATE_KEY.length);
}
if (process.env.INFURA_PROJECT_ID) {
    console.log("INFURA_PROJECT_ID:", process.env.INFURA_PROJECT_ID);
}
if (process.env.SEPOLIA_URL) {
    console.log("SEPOLIA_URL:", process.env.SEPOLIA_URL);
}