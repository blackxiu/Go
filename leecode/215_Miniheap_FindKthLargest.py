import heapq

class Solution:
    # Heap
    # Time Complexity: O(NlogN)
    # Space Complexity: O(N)

    def findKthLargest(self, nums: List[int], k: int) -> int:
        # return heapq.nlargest(k, nums)[-1]
        heap = []
        heapq.heapify(heap)
        for num in nums:
            heapq.heappush(heap, num*-1)

        while k > 1:
            heapq.heappop(heap)
            k = k - 1

        return heapq.heappop(heap)*-1