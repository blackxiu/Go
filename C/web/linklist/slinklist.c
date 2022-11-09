/*
 *  程序名：slinklist.c，此程序演示单向链表及常用操作。
 *  作者：C语言技术网(www.freecplus.net) 日期：20201019
*/
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

struct st_girl
{
  int  no;               // 超女的编号。
  char name[31];         // 超女的姓名。
  struct st_girl *next;  // 存放下一个节点结构体的地址。
};

struct st_girl *head,*tail;  // 链表的头和尾指针。
int  linksize;  // 链表的大小（节点数）。

void initlinklist();   // 链表初始化。

void addnode1(struct st_girl *pst);          // 把节点插入到链表的头部。
void addnode2(struct st_girl *pst);          // 把节点追加到链表的尾部。
void addnode3(struct st_girl *pst,int iseq); // 在链表iseq位置插入节点，iseq从0开始，如果iseq的值大于等于linksize，则在链表的尾部追加。
void addnode4(struct st_girl *pst);          // 把节点插入到链表的中间，整个链表按结构体的no字段排序。

void delnode1();         // 从链表的头部删除一个节点。
void delnode2();         // 从链表的尾部删除一个节点。
void delnode3(int iseq); // 从链表中删除第iseq个节点，iseq的取值必须为：0<=iseq<linksize，否则不删除任何节点。
void delnode4(int ino);  // 从链表中删除结构体no字段等于ino的记录。

void printlinklist();    // 显示链表全部节点的数据。

void freelinklist();     // 释放链表占用的内存。

int main()
{
  initlinklist();   // 链表初始化。

  struct st_girl *pst=0;     // 超女结构体指针。

  // 循环5次，从链盘上输入超女记录。
  int ii=0;
  for (ii=0;ii<5;ii++)
  {
    // 创建一个新的节点。
    pst=(struct st_girl*)malloc(sizeof(struct st_girl));   // 申请一块新内存。
    memset(pst,0,sizeof(struct st_girl));

    printf("请输入超女的编号和姓名，中间用空格分隔：");
    scanf("%d %s",&pst->no,pst->name);

    // addnode1(pst); // 把节点插入到链表的头部。
    // addnode2(pst); // 把节点追加到链表的尾部。
    // 在链表iseq位置插入节点，iseq从0开始，如果iseq的值大于等于linksize，则在链表的尾部追加。
    // addnode3(pst,1);  // 在链表第3个位置插入节点。
    addnode4(pst); // 把节点插入到链表的中间，按结构体的no字段排序。
  }

  // delnode1();  // 从链表的头部删除一个节点。
  // delnode2();  // 从链表的尾部删除一个节点。
  // delnode3(3);  // 从链表中删除指定位置的节点。
  // delnode4(3); // 从链表中删除结构体no字段等于ino的记录。

  // 显示链表全部节点的数据。
  printlinklist();

  // 释放链表占用的内存。
  freelinklist();

  // system("pause");  // windows平台启用这行代码。
}

// 在实际开发中，链表的每个节点，使用的时候就动态分配，不用了就删除它，预先分配节点会增加编程的难度，没有意义。
// 链表初始化。
void initlinklist()
{
  linksize=0;
  head=tail=0;
}

// 把节点插入到链表的头部。
void addnode1(struct st_girl *pst)
{
  if (head==0) // 如果链表为空。
  { 
    head=tail=pst; pst->next=0;
  }
  else // 如果链表不为空。
  {
    pst->next=head;
    head=pst;
  }

  linksize++;
}

// 把节点追加到链表的尾部。
void addnode2(struct st_girl *pst)
{
  if (head==0) // 如果链表为空。
  { 
    head=tail=pst;  pst->next=0;
  }
  else // 如果链表不为空。
  {
    tail->next=pst;
    tail=pst;
    pst->next=0;
  }

  linksize++;
}

// 我的代码不是最精简的，但是高效易懂。

// 在链表iseq位置插入节点，iseq从0开始，如果iseq的值大于linksize，则在链表的尾部追加。
void addnode3(struct st_girl *pst,int iseq)
{
  // 如果链表为空。
  if (head==0) 
  { 
    head=tail=pst; pst->next=0; linksize=1; return;
  }

  // 如果iseq==0，在链表头部插入。
  if (iseq==0) 
  {
    addnode1(pst); return;
  }

  // 如果iseq>=linksize，把节点追加到链表的尾部。
  if (iseq>=linksize) 
  {
    addnode2(pst); return;
  }

  // 以下代码用于处理当前在头和尾之间插入的情况。
  struct st_girl *ptmp=head; // 从头节点开始。
  int itmp=1;

  while (1)
  {
    if ( iseq== itmp )
    {
      pst->next=ptmp->next;
      ptmp->next=pst;
      linksize++;
      return;
    }

    itmp++;
    ptmp=ptmp->next;  // 往后移动一个节点。
  } 
}

// 显示链表全部节点的数据。

// 把节点插入到链表的中间，按结构体的no字段排序。
void addnode4(struct st_girl *pst)
{
  // 如果链表为空。
  if (head==0) 
  { 
    head=tail=pst; pst->next=0; linksize=1; return;
  }

  // 如果当前节点的no小于等于头节点的no，在链表头部插入。
  if (pst->no<=head->no) 
  {
    addnode1(pst); return;
  }

  // 如果当前节点的no大于等于尾节点的no，把节点追加到链表的尾部。
  if (pst->no>=tail->no) 
  {
    addnode2(pst); return;
  }

  // 以下代码用于处理当前节点的no在头和尾之间的情况。
  struct st_girl *ptmp=head; // 从头节点开始。

  while (ptmp!=0)
  {
    if ( (pst->no>=ptmp->no) && (pst->no<=ptmp->next->no) )
    {
      pst->next=ptmp->next;
      ptmp->next=pst;
      linksize++;
      return;
    }

    ptmp=ptmp->next;  // 往后移动一个节点。
  } 
}

// 显示链表全部节点的数据。
void printlinklist()
{
  struct st_girl *ptmp=head; // 从头节点开始。

  while (ptmp!=0)
  {
    printf("no=%d,name=%s\n",ptmp->no,ptmp->name);
    ptmp=ptmp->next; // 指针移动到下一节点。
  } 

  printf("linksize=%d\n",linksize);
}

// 释放链表占用的内存。
void freelinklist()
{
  struct st_girl *ptmp; // 一个临时的指针。

  while (head!=0)
  {
    ptmp=head;
    head=head->next;  // 移动到下一节点。
    free(ptmp);       // 释放当前节点。
  } 

  initlinklist();
}


// 从链表的头部删除一个节点。
void delnode1()
{
  // 如果是一个空链表，直接返回。
  if (head==0) return;

  // 如果链表只有一个节点。
  if (head==tail)
  {
    free(head); initlinklist(); return;
  }

  // 如果链表有多个节点。
  struct st_girl *ptmp=head;
  head=head->next;
  free(ptmp);
  linksize--;
}

// 从链表的尾部删除一个节点。
void delnode2()
{
  // 如果是一个空链表，直接返回。
  if (head==0) return;

  // 如果链表只有一个节点。
  if (head==tail)
  {
    free(head); initlinklist(); return;
  }

  struct st_girl *ptmp=head; // 从头节点开始。

  while (ptmp!=0)
  {
    if (ptmp->next==tail)
    {
      free(tail); tail=ptmp; tail->next=0; linksize--; return;
    }

    ptmp=ptmp->next; 
  } 
}

// 从链表中删除第iseq个节点，iseq的取值必须为：0<=iseq<linksize，否则不删除任何节点。
void delnode3(int iseq) 
{
  // 判断iseq的取值是否合法，如果不合法，直接返回。
  if ( (iseq<0) || (iseq>=linksize) ) return;

  // 如果iseq==0，删除头节点。
  if (iseq==0)
  {
    delnode1(); return;
  }

  // 如果iseq=linksize-1，删除尾节点。
  if (iseq==linksize-1)
  {
    delnode2(); return;
  }

  struct st_girl *ptmp=head; // 从头节点开始。
  int itmp=0;

  while (ptmp!=0)
  {
    if (itmp==iseq-1)
    {
      struct st_girl *ptmp1=ptmp->next;
      ptmp->next=ptmp->next->next;
      free(ptmp1);
      linksize--;
      return;
    }

    itmp++;
    ptmp=ptmp->next;
  } 
}

// 从链表中删除结构体no字段等于ino的记录。
void delnode4(int ino) 
{
  // 如果头节点的no字段等于ino，删除头节点。
  if (head->no==ino)
  {
    delnode1(); return;
  }

  // 如果尾节点的no字段等于ino，删除尾节点。
  if (tail->no==ino)
  {
    delnode2(); return;
  }

  struct st_girl *ptmp=head; // 从头节点开始。

  while (ptmp!=0)
  {
    if (ptmp->next->no==ino)
    {
      struct st_girl *ptmp1=ptmp->next;
      ptmp->next=ptmp->next->next;
      free(ptmp1);
      linksize--;
      return;
    }

    ptmp=ptmp->next;
  } 
}
