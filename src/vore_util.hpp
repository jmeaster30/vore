#pragma once

#include <string>

/*
 *  Compares two strings lexicographically
 *  1 -> a > b
 *  0 -> a = b
 * -1 -> a < b
 */
int lexico_compare(std::string a, std::string b);


/*
 * Swaps the values of a and b if a is less than b
 */
void swap_if_less(uint64_t* a, uint64_t* b);