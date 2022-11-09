class Solution:
    # Hash Table
    # Time Complexity: O(N)
    # Space Complexity: O(1)

    def FindDifference(self, s: str, t: str) -> str:
        if len(s) == 0:
            return t
        table = [0]*26
        for i in range(len(t)):
            if i < len(s):
                table[ord(s[i]) - ord('a')] -= 1
            table[ord(t[i]) - ord('a')] += 1
        for i in range(26):
            if table[i] != 0:
                return chr(i+97)
        return 'a'