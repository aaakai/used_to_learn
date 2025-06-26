package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

type S struct {
	A string `json:"a"`
}
type User struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func main() {
	//jsonData := `{"name": null, "age": null}`
	jsonData := ""

	var user User
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(user)
}

func extractDescription(text string) string {
	// 定义需要保留的文本模式（正则表达式）
	pattern := `(?s)(.*?<name>.*?</name>.*?)(?:<.*?>)?$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(text)

	if len(matches) > 1 {
		// 移除<name>标签
		result := strings.ReplaceAll(matches[1], "<name>", "")
		result = strings.ReplaceAll(result, "</name>", "")
		return result
	}
	return ""
}

func byteTest() {
	a := 'a'
	b := '0'
	//当我们直接输出byte（字符）的时候输出的是这个字符对应的码值
	fmt.Println(a)
	fmt.Println(b)

	//如果我们要输出这个字符，需要格式化输出
	fmt.Printf("%c--%c", a, b) //%c	相应Unicode码点所表示的字符
}

func anyway() {
	var sp interface{} = S{A: "AAA"}
	fmt.Printf("==%#v==\n", sp)
	_ = json.Unmarshal([]byte("{\"a\": \"xxxx\"}"), &sp)
	fmt.Printf("==%#v==\n", sp)
	var sp2 S
	_ = json.Unmarshal([]byte(`{"a": "xxxx"}`), &sp2)
	fmt.Printf("==%#v==\n", sp2)
}

func merge() {
	a := []int{1, 3, 5}
	b := []int{2, 4, 6}
	n, m := 3, 3
	lenA, lenB := 0, 0
	res := make([]int, 0, n+m)
	for {
		if lenA == n {
			res = append(res, b[lenB:]...)
			break
		}
		if lenB == m {
			res = append(res, a[lenA:]...)
			break
		}
		if a[lenA] < b[lenB] {
			res = append(res, a[lenA])
			lenA++
		} else if a[lenA] > b[lenB] {
			res = append(res, b[lenB])
			lenB++
		}
	}
	fmt.Println(res)
}

func revert() {
	nums := []int{}
	lenA := len(nums)
	left, right := 0, lenA
	val := 3
	for left < right {
		if nums[left] == val {
			nums[left] = nums[right-1]
			right--
		} else {
			left++
		}
	}
	fmt.Println(nums, left)
}

func removeReq() {
	nums := []int{1, 1, 2, 2, 2, 3, 3, 3}
	i, j := 0, 1
	if len(nums) < 1 {
		fmt.Println(nums, 0)
	}
	for j < len(nums) {
		if nums[i] != nums[j] {
			i++
			nums[i] = nums[j]
		}
		j++
	}
	fmt.Println(nums, i+1)
}

func removeReq2() {
	nums := []int{1, 1, 1, 2, 2, 3, 3, 3}
	if len(nums) <= 2 {
		fmt.Println(nums, len(nums))
	}
	i, j := 2, 2
	for j < len(nums) {
		if nums[i-2] != nums[j] {
			nums[i] = nums[j]
			i++
		}
		j++
	}
	fmt.Println(nums, i)
}

func maxMoneyTwo(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	res := 0
	for i := 1; i < len(prices); i++ {
		res = maxInt(res, prices[i]-prices[i-1])
	}
	return res
}

func maxMoney(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	res, minTmp := 0, prices[0]
	for _, v := range prices {
		if v >= minTmp {
			res = maxInt(res, v-minTmp)
		} else {
			minTmp = v
		}
	}
	return res
}
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}
