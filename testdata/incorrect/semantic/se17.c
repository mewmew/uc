/* Test file for semantic errors. Contains exactly one error. */


int main(void) {
  char hello[5];
  hello+1; //  Attempt to use char array in arithmetic. (legal in C)
  
}

