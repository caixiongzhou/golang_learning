// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanToIntCase {
    address public owner;

    constructor() {
        owner = msg.sender;
    }

    /**
     * @dev 将罗马数字转换为整数 - 修复版本
     */
    function RomanToInt(string memory s) public pure returns (uint256) {
        bytes memory romanBytes = bytes(s);
        uint256 length = romanBytes.length;
        require(length > 0, "Empty string");

        uint256 result = 0;
        uint256 i = 0;

        while (i < length) {
            uint256 current = _charToValue(romanBytes[i]);

            // 检查是否有下一个字符且需要减法
            if (i < length - 1) {
                uint256 next = _charToValue(romanBytes[i + 1]);

                // 如果当前字符小于下一个字符，使用减法规则
                if (current < next) {
                    // 直接计算减法组合的值，然后加到结果中
                    result += (next - current);
                    i += 2; // 跳过两个字符
                    continue;
                }
            }

            // 否则直接加当前字符的值
            result += current;
            i += 1;
        }

        require(result > 0 && result < 4000, "Invalid Roman numeral");
        return result;
    }

    /**
     * @dev 将单个罗马数字字符转换为对应的数值
     */
    function _charToValue(bytes1 c) private pure returns (uint256) {
        if (c == 'I') return 1;
        if (c == 'V') return 5;
        if (c == 'X') return 10;
        if (c == 'L') return 50;
        if (c == 'C') return 100;
        if (c == 'D') return 500;
        if (c == 'M') return 1000;
        revert("Invalid Roman character");
    }

    /**
     * @dev 批量转换多个罗马数字
     */
    function batchRomanToInt(string[] memory romans) public pure returns (uint256[] memory) {
        uint256[] memory results = new uint256[](romans.length);

        for (uint256 i = 0; i < romans.length; i++) {
            results[i] = RomanToInt(romans[i]);
        }

        return results;
    }
}