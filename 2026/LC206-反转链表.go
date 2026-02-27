package _026

type ListNode struct {
	Val  int
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	var tmp *ListNode
	cur := head
	for cur != nil {
		next := cur.Next
		cur.Next = tmp
		tmp = cur
		cur = next
	}
	return tmp
}
