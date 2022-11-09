class Solution:
    # 加法运算
    # Time Complexity: O(N)
    # Space Complexity: O(1)
    
    def FindDifference(self, s: str, t: str) -> str:
        if len(s) == 0:
            return t
        total = 0
        for i in range(len(t)):
            if i < len(s):
                total -= ord(s[i])
            total += ord(t[i])
        return chr(total)