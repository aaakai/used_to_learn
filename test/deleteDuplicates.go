package test

/*
给定一个已排序的链表的头 head ， 删除原始链表中所有重复数字的节点，只留下不同的数字 。返回 已排序的链表 。

提示：

链表中节点数目在范围 [0, 300] 内
-100 <= Node.val <= 100
题目数据保证链表已经按升序 排列
*/
/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */

type ListNode struct {
	Val  int
	Next *ListNode
}

func deleteDuplicates(head *ListNode) *ListNode {
	if head == nil {
		return head
	}
	res := &ListNode{0, head}
	tmp := res

	for tmp.Next != nil && tmp.Next.Next != nil {
		if tmp.Next.Val == tmp.Next.Next.Val {
			val := tmp.Next.Val
			for tmp.Next != nil && tmp.Next.Val == val {
				tmp.Next = tmp.Next.Next
			}
		} else {
			tmp = tmp.Next
		}
	}

	return res.Next
}
