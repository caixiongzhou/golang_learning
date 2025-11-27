// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;


contract IntToRoman {
    // 定义罗马数字符号和对应的数值
    struct RomanSymbol{
        uint256 value;
        string symbol;
    }

    RomanSymbol[] private romanSymbols;


    constructor(){
        // 初始化罗马数字符号表（按从大到小排序）
        romanSymbols.push(RomanSymbol(1000,"M"));
        romanSymbols.push(RomanSymbol(900,"CM"));
        romanSymbols.push(RomanSymbol(500, "D"));
        romanSymbols.push(RomanSymbol(400, "CD"));
        romanSymbols.push(RomanSymbol(100, "C"));
        romanSymbols.push(RomanSymbol(90, "XC"));
        romanSymbols.push(RomanSymbol(50, "L"));
        romanSymbols.push(RomanSymbol(40, "XL"));
        romanSymbols.push(RomanSymbol(10, "X"));
        romanSymbols.push(RomanSymbol(9, "IX"));
        romanSymbols.push(RomanSymbol(5, "V"));
        romanSymbols.push(RomanSymbol(4, "IV"));
        romanSymbols.push(RomanSymbol(1, "I"));

    }
    /**
     * @dev 将整数转换为罗马数字
     * @param num 要转换的整数 (1-3999)
     * @return 罗马数字字符串
     */

    function toRoman(uint256 num) public view returns(string memory) {
        require(num >0 && num <4000, "Number must be between 1 and 3999");

        bytes memory result;
        for(uint256 i = 0;i<romanSymbols.length;i++){
            RomanSymbol memory symbol = romanSymbols[i];
            while (num > symbol.value){
                // 将字符串转换为bytes以便拼接
                bytes memory symbolBytes = bytes(symbol.symbol);

                // 扩展结果数组
               bytes memory newResult = new bytes(result.length + symbolBytes.length);

                // 复制原有结果
                for(uint256 j =0;j < result.length; j++){
                    newResult[j] = result[j];
                }

                // 添加新符号
                for (uint256 k = 0;k<symbolBytes.length;k++){
                    newResult[result.length + k] = symbolBytes[k];
                }

                result = newResult;
                num -= symbol.value;
            }

        }
        return string(result);
    }

    /**
     * @dev 优化的版本，使用字符串拼接
     * @param num 要转换的整数 (1-3999)
     * @return 罗马数字字符串
     */
    function toRomanOptimized(uint256 num) public pure returns (string memory) {
        require(num > 0 && num < 4000, "Number must be between 1 and 3999");

        // 使用固定大小的存储数组
        string[13] memory symbols = ["M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I"];
        uint256[13] memory values = _getValues();

        string memory result;

        for (uint256 i = 0; i < symbols.length; i++) {
            while (num >= values[i]) {
                result = string(abi.encodePacked(result, symbols[i]));
                num -= values[i];
            }
        }

        return result;
    }

    /**
     * @dev 返回数值数组
     */
    function _getValues() private pure returns (uint256[13] memory) {
        uint256[13] memory values;
        values[0] = 1000;
        values[1] = 900;
        values[2] = 500;
        values[3] = 400;
        values[4] = 100;
        values[5] = 90;
        values[6] = 50;
        values[7] = 40;
        values[8] = 10;
        values[9] = 9;
        values[10] = 5;
        values[11] = 4;
        values[12] = 1;
        return values;
    }

    /**
     * @dev 批量转换多个数字
     * @param nums 要转换的数字数组
     * @return 罗马数字字符串数组
     */
    function batchToRoman(uint256[] memory nums) public pure returns (string[] memory) {
        string[] memory results = new string[](nums.length);

        for (uint256 i = 0; i < nums.length; i++) {
            results[i] = toRomanOptimized(nums[i]);
        }

        return results;
    }




}