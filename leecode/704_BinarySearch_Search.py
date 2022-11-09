class Solution:
    def numRescueBoats(self, people: List[int], limit: int) -> int:
        
        # Time Complexity: O(N) 
        # Space Complexity: O(1)
        def search(self, nums, target):
            """
            :type nums: List[int]
            :type target: int
            :rtype: int
            """
            for i in range(len(nums)):
                if nums[i] == target:
                    return i
                return -1