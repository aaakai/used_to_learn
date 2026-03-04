package _026

func detectCycle(head *ListNode) *ListNode {
	tmp := map[*ListNode]bool{}
	for head != nil {
		if _, ok := tmp[head]; ok {
			return head
		}
		tmp[head] = true
		head = head.Next
	}
	return nil
}

func detectCycle2(head *ListNode) *ListNode {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			p := head
			for p != slow {
				p = p.Next
				slow = slow.Next
			}
			return p
		}
	}
	return nil
}
