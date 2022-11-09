class Solution:
    # Set
    # Time Complexity: O(N)
    # Space Complexity: O(N)
    def containsDuplicate(self, nums: List[int]) -> bool:
        if len(nums) == 0:
            return False
        hashset = set(nums)
        return False if len(nums) == len(hashset) else True