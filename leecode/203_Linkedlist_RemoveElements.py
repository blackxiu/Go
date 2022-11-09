def RemoveElements(self, head: ListNode, val: int) -> ListNode:
    # Time Complexity: O(N) 
    # Space Complexity: O(1) 
    dummy = ListNote(0)
    dummy.next = head
    prev = dummy

    while head is not None:
        if head.val == val:
            prev.next == head.next
        else:
            prev = prev.next
        head = head.next

    return dummy.next