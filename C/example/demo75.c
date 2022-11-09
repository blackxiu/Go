/*
75.自定义整数反转后输出

题目：输入一个整数，并将其反转后输出。

程序分析：无。
*/

#include <stdio.h>

int main()
{
  int n, reversedNumber = 0, remainder;
 
  printf("输入一个整数: ");
  scanf("%d", &n);
 
  while(n != 0)
  {
    remainder = n%10;
    reversedNumber = reversedNumber*10 + remainder;
    n /= 10;
  }
 
  printf("反转后的整数: %d\n", reversedNumber);
 
  return 0;
}
