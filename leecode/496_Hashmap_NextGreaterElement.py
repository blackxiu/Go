class Solution:
    # Time Complexity: O(MN)
    # Space Complexity: O(N)
    def nextGreaterElement(self, nums1: List[int], nums2: List[int]) -> List[int]:
        res = []
        stack = []
        mapping = {}
        for num in nums2:
            while(len(stack) != 0 and num > stack[-1]):
                temp = stack.pop()
                mapping[temp] = num
        stack.append(num)

        while len(temp) != 0:
            mapping[stack.pop()] = -1
        
        for num in nums1:
            res.append(mapping[num])

        return res