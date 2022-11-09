class Solution:
    def numRescueBoats(self, people: List[int], limit: int) -> int:
        
        # Time Complexity: O(NlogN) 
        # Space Complexity: O(1)
        if people is None or len(people) == 0:
            return 0
        people.sort()
        i = 0
        j = len(people) - 1
        res = 0
        while(i <= j):
            if people[i] + people[j] <= limit:
                i = i + 1
            j = j - 1
            res = res + 1
        return res
