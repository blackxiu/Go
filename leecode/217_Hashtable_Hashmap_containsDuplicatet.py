class Solution:
    # Hash map
    # Time Complexity: O(N)
    # Space Complexity: O(N)
    def containsDuplicate(self, nums: List[int]) -> bool:
        if len(nums) == 0:
            return False
        mapping = {}
        for num in nums:
            if num not in mapping:
                mapping[num] = 1
            else:
                mapping[num] = mapping.get(num) + 1
        for v in mapping.values():
            if v >1:
                return True
        return False 