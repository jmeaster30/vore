#include "vore_util.hpp"

int lexico_compare(std::string a, std::string b)
{
  if (a.length() < b.length()) return -1;
  if (a.length() > b.length()) return 1;

  // the length of a and b are the same here
  for (uint64_t i = 0; i < a.length(); i++)
  {
    if (a[i] != b[i])
      return a[i] < b[i] ? -1 : 1;
  }
  return 0;
}

void swap_if_less(uint64_t* a, uint64_t* b)
{
  if (*a < *b) {
    uint64_t c = *a;
    *a = *b;
    *b = c;
  }
}