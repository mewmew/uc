int b;

char addi;

int la;

void jal(void) {
  b = b * addi+la;
}

int mov(int lb) {
  addi = lb;
  return 0;
}

int main (void) {
  la = 8;
  jal();
  mov(la);
}
