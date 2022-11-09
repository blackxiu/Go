def MoveZeroes(self, nums: List[int]) -> None:
    # Time Complexity: O(N)
    # Space Complexity: O(1)

    index = 0
    for num in nums:
        if num != 0:
            nums[index] = num
            index += 1
    for i in range(index,len(nums)):
        nums[i] = 0