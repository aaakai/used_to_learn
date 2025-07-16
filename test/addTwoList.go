package test

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	var res *ListNode
	var head *ListNode
	carry := 0
	for l1 != nil || l2 != nil {
		value1, value2 := 0, 0
		if l1 != nil {
			value1 = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			value2 = l2.Val
			l2 = l2.Next
		}
		sum := value1 + value2 + carry
		sum, carry = sum%10, sum/10
		if head == nil {
			head = &ListNode{Val: sum}
			res = head
		} else {
			res.Next = &ListNode{Val: sum}
			res = res.Next
		}
	}
	if carry > 0 {
		res.Next = &ListNode{Val: carry}
	}
	return head
}
