/*
82.自定义一个八进制转换为十进制

题目：八进制转换为十进制

程序分析：无。
*/

#include<stdio.h>
#include<stdlib.h>

int main()
{
  int n=0,i=0;
  char s[20];

  printf("请输入一个8进制数:\n");
  gets(s);

  while(s[i]!='\0')
  {
    n=n*8+s[i]-'0';
    i++;
  }

  printf("刚输入的8进制数转化为十进制为\n%d\n",n);
  
  return 0;
}
