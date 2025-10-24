package test

import "strconv"

/*
https://leetcode.cn/problems/restore-ip-addresses/description/

有效 IP 地址 正好由四个整数（每个整数位于 0 到 255 之间组成，且不能含有前导 0），整数之间用 '.' 分隔。

例如："0.1.2.201" 和 "192.168.1.1" 是 有效 IP 地址，但是 "0.011.255.245"、"192.168.1.312" 和 "192.168@1.1" 是 无效 IP 地址。
给定一个只包含数字的字符串 s ，用以表示一个 IP 地址，返回所有可能的有效 IP 地址，这些地址可以通过在 s 中插入 '.' 来形成。你 不能 重新排序或删除 s 中的任何数字。你可以按 任何 顺序返回答案。



示例 1：

输入：s = "25525511135"
输出：["255.255.11.135","255.255.111.35"]
示例 2：

输入：s = "0000"
输出：["0.0.0.0"]
示例 3：

输入：s = "101023"
输出：["1.0.10.23","1.0.102.3","10.1.0.23","10.10.2.3","101.0.2.3"]
*/

const (
	MAX_SEG = 4
)

var (
	res []string
	num []int
)

func RestoreIpAddresses(s string) []string {
	num = make([]int, MAX_SEG)
	res = []string{}
	restoreIpAddressesDfs(s, 0, 0)
	return res
}

func restoreIpAddressesDfs(s string, numIndex, sIndex int) {
	if numIndex == MAX_SEG {
		if sIndex == len(s) {
			tmp := ""
			for i := 0; i < MAX_SEG; i++ {
				tmp += strconv.Itoa(num[i])
				if i < MAX_SEG-1 {
					tmp += "."
				}
			}
			res = append(res, tmp)
		}
		return
	}
	if sIndex == len(s) {
		return
	}

	if s[sIndex] == '0' {
		num[numIndex] = 0
		restoreIpAddressesDfs(s, numIndex+1, sIndex+1)
	}

	addr := 0
	for sIndex < len(s) {
		addr = addr*10 + int(s[sIndex]-'0')
		if addr > 0 && addr <= 0xFF {
			num[numIndex] = addr
			restoreIpAddressesDfs(s, numIndex+1, sIndex+1)
		} else {
			break
		}
		sIndex++
	}
}
