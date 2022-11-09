def ReverseList(self, head: ListNode, val: int) -> ListNode:
    # Time Complexity: O(N) 
    # Space Complexity: O(1) 
    dummy = ListNote(0)
    dummy.next = head
    while head is not None and head.next is not None:
        dnext = dummy.next
        hnext = head.next
        dummy.next = hnext
        head.next = hnext.next
        hnext.next = dnext
    return dummy.next