// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MergeSortedArray{
    /**
     * @dev 合并两个有序数组
     * @param nums1 第一个有序数组
     * @param nums2 第二个有序数组
     * @return 合并后的有序数组
     */

    function mergeSortedArray(uint256[] memory nums1,uint256[] memory nums2)public pure returns (uint256[] memory){
        uint256 m = nums1.length;
        uint256 n = nums2.length;
        // 如果其中一个数组为空，直接返回另一个
        if (m == 0)  return nums2;
        if (n == 0)  return nums1;

        // 创建结果数组
        uint256[] memory result = new uint256[](m+n);

        uint256 i = 0; // nums1 的指针
        uint256 j = 0; // nums2 的指针
        uint256 k = 0; // result 的指针

        // 双指针法合并两个有序数组
        while(i<m && j<n){
            if (nums1[i] < nums2[j]){
                result[k] = nums1[i];
                i++;
            }else{
                result[k] = nums2[j];
                j++;
            }
            k++;
        }
        // 将剩余元素添加到结果中
        while(i < m){
            result[k] = nums1[i];
            i++;
            k++;
        }
        while(j < n){
            result[k] = nums2[j];
            j++;
            k++;
        }
        return result;
    }
}