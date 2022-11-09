class MyHashSet:
    # Array
    # Space Complexity：O(N)

    def __init__(self):
        """Initialize your data structure here.
        """
        self.hashset = [0]*1000001

    def add(self, key: int) -> None:
        self.hashset[key] = 1

    def remove(self, key: int) -> None:
        self.hashset[key] = 0

    def contains(self, key: int) -> bool:
        """Returns true if this set contains the specified element
        """        
        return self.hashset[key]