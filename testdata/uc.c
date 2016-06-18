#include <stdio.h>

void putint(int x) {
   printf("%d",x);
}

void putstring(char s[]) {
  fputs(s,stdout);
}

int getint(void){
  int i;
  scanf("%d", &i);
  return i;
}


int getstring(char s[]) {
  scanf("%s",s);
}

