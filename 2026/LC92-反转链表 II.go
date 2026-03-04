package _026

func reverseBetween(head *ListNode, left int, right int) *ListNode {
	nullNode := &ListNode{Val: -1}
	nullNode.Next = head
	// 找到leftNode和rightNode
	// pre是left的前一个节点 last是right的后一节点
	// pre.next right.next 变成nil
	// 翻转leftNode和rightNode之间的节点
	// pre和leftNode连接，rightNode和lastNode连接
	pre := nullNode
	for i := 0; i < left-1; i++ {
		pre = pre.Next
	}
	rightNode := pre
	for i := 0; i < right-left+1; i++ {
		rightNode = rightNode.Next
	}

	leftNode := pre.Next
	lastNode := rightNode.Next

	pre.Next = nil
	rightNode.Next = nil

	ReverseList(leftNode)

	pre.Next = rightNode
	leftNode.Next = lastNode

	return nullNode.Next
}
