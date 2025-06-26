package test

/*
给你两个字符串：ransomNote 和 magazine ，判断 ransomNote 能不能由 magazine 里面的字符构成。

如果可以，返回 true ；否则返回 false 。

magazine 中的每个字符只能在 ransomNote 中使用一次。

示例 1：

输入：ransomNote = "a", magazine = "b"
输出：false
示例 2：

输入：ransomNote = "aa", magazine = "ab"
输出：false
示例 3：

输入：ransomNote = "aa", magazine = "aab"
输出：true

提示：

1 <= ransomNote.length, magazine.length <= 105
ransomNote 和 magazine 由小写英文字母组成
*/
func CanConstruct(ransomNote string, magazine string) bool {
	mapM := make(map[byte]int)
	for i := 0; i < len(magazine); i++ {
		mapM[magazine[i]]++
	}
	for i := 0; i < len(ransomNote); i++ {
		if mapM[ransomNote[i]] == 0 {
			return false
		}
		mapM[ransomNote[i]]--
	}
	return true
	/*
	 if len(ransomNote) > len(magazine) {
	        return false
	    }
	    cnt := [26]int{}
	    for _, ch := range magazine {
	        cnt[ch-'a']++
	    }
	    for _, ch := range ransomNote {
	        cnt[ch-'a']--
	        if cnt[ch-'a'] < 0 {
	            return false
	        }
	    }
	    return true
	*/
}
