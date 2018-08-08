void sum(int *a, int b) {
  *a = *a + b;
}

int main() {
  int a = 1;
  int *b = &a;
  sum(b, 2);
  return a;
}
