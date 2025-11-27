// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract StringTest{
    function reverseString(string memory _str)public pure returns(string memory)  {
/*        // 将字符串转换为 bytes 类型
        bytes memory strBytes = bytes(_str);
        bytes memory reversed = new bytes(strBytes.length);  // new一个长度为strBytes.length的bytes数组准备用来存反转后的bytes


        for(uint i = 0;i < strBytes.length; i++){
            reversed[i] = strBytes[strBytes.length-1-i];
        }*/

        bytes memory strBytes = bytes(_str);
        uint length = strBytes.length;

        // 如果字符串为空或只有一个字符，直接返回
        if (length <= 1) {
            return _str;
        }

        bytes memory reversed = new bytes(length);


        // 使用两个指针同时移动，减少计算
        for(uint i = 0 ;i< length;i++){
            reversed[i] = strBytes[length - 1 - i];
        }
        return string(reversed);
    }
}