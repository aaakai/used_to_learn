package test

import (
	"fmt"
	"sort"
)

type IntensitySegment struct {
	is map[int]int
}

func NewIntensitySegment() *IntensitySegment {
	return &IntensitySegment{
		is: make(map[int]int),
	}
}

func (is *IntensitySegment) add(from, to, amount int) {
	keys := is.orderedKeys()
	// 初始化
	if len(keys) == 0 {
		is.set(from, to, amount)
		return
	}

	// from to 超出已有范围
	if to < keys[0] || from > keys[len(keys)-1] {
		is.set(from, to, amount)
		return
	}

	//  to 在已有范围中 不存在 这等于第一个比to 小的
	if _, f := is.is[to]; !f {
		l := is.leftKey(to)
		is.is[to] = is.is[l]
	}

	// from 在区间外
	if from < keys[0] {
		is.is[from] = amount
	} else if from == keys[0] {
		// from是第一个
		is.is[from] += amount
	} else {
		if _, f := is.is[from]; !f {
			// from不存在 from=左边的第一个+amount
			is.is[from] = amount
			if l := is.leftKey(from); l != -1 {
				is.is[from] += is.is[l]
			}
		} else {
			// from存在 from=is[from]+amount
			is.is[from] += amount
		}
	}

	// from to 区间数据 加amount
	for _, k := range keys {
		if k > from && k < to {
			is.is[k] += amount
		}
	}

	is.merge()
}

func (is *IntensitySegment) set(from, to, amount int) {

	keys := is.orderedKeys()
	// 初始化场景
	if len(keys) == 0 {
		is.is[to] = 0
		is.is[from] = amount
		return
	}

	// 范围无交集场景
	if to < keys[0] || to > keys[len(keys)-1] {
		is.is[to] = 0
	} else if _, f := is.is[to]; !f {
		// to存在 找到to的左边的key，并复制
		l := is.leftKey(to)
		is.is[to] = is.is[l]
	}

	is.is[from] = amount

	// 重制from to 之间的数据
	for _, k := range is.orderedKeys() {
		if k > from && k < to {
			delete(is.is, k)
		}
	}

	is.merge()
}

// 重新处理数据
func (is *IntensitySegment) merge() {
	keys := is.orderedKeys()
	if len(keys) == 0 {
		return
	}

	// 清楚前序连续0值 删除前者【0，0】【5，0】【10，1】【20，2】 -> 【5，0】【10，1】【20，2】
	i := 0
	for ; i < len(keys) && is.is[keys[i]] == 0; i++ {
		delete(is.is, keys[i])
	}
	keys = keys[i:]

	// 清除后续连续0 删除后者【5，1】【10，0】【20，0】 -> 【5，1】【10，0】
	for i = len(keys) - 1; i >= 0; i-- {
		if is.is[keys[i]] == 0 {
			if i-1 >= 0 && is.is[keys[i-1]] == 0 {
				delete(is.is, keys[i])
				keys = append(keys[0:i], keys[i+1:]...)
			}
		}
	}

	// 清楚相同v的k,清楚后者 【10，1】【20，2】【30，2】【40，3】 -> 【10，1】【20，2】【40，3】
	lastIntensity := is.is[keys[0]]
	for i := 1; i < len(keys); i++ {
		if is.is[keys[i]] == lastIntensity && lastIntensity != 0 {
			delete(is.is, keys[i])
		} else {
			lastIntensity = is.is[keys[i]]
		}
	}

}

// key 升序排序
func (is *IntensitySegment) orderedKeys() []int {
	keys := make([]int, 0)
	for k := range is.is {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

// 差到第一个比index小的key
func (is *IntensitySegment) leftKey(v int) int {
	index := -1
	for _, k := range is.orderedKeys() {
		if k < v {
			index = k
		}
	}
	return index
}

func (is *IntensitySegment) toString() {
	fmt.Printf("%v", is.toPrintData())
}

func (is *IntensitySegment) toPrintData() string {
	keys := is.orderedKeys()
	res := make([][2]int, 0)
	for _, k := range keys {
		res = append(res, [2]int{k, is.is[k]})
	}
	return fmt.Sprintf("%v", res)
}
