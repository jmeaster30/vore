#include "ast.hpp"

std::string FixEscapeCharacters(std::string val) {
  std::string fixed = "";
  for(u_int64_t i = 0; i < val.length(); i++)
  {
    char current = val[i];
    if(current == '\\' && i < val.length() - 1) {
      char next = val[i + 1];
      switch(next)
      {
        case 'a':
          fixed += (char)7;
          break;
        case 'b':
          fixed += (char)8;
          break;
        case 'e':
          fixed += (char)27;
          break;
        case 'f':
          fixed += (char)12;
          break;
        case 'n':
          fixed += (char)10;
          break;
        case 'r':
          fixed += (char)13;
          break;
        case 't':
          fixed += (char)9;
          break;
        case 'v':
          fixed += (char)11;
          break;
        case '\\':
          fixed += '\\';
          break;
        case '\'':
          fixed += '\'';
          break;
        case '\"':
          fixed += '\"';
          break;
        case '?':
          fixed += '?';
          break;
        case 'x':
          if(i < val.length() - 3) {
            char a = val[i + 2];
            char b = val[i + 3];
            if(!((a >= '0' && a <= '9') || (a >= 'A' && a <= 'F') || (a >= 'a' && a <= 'f'))) {
              fixed += 'x';
              break;
            }
            if(!((b >= '0' && b <= '9') || (b >= 'A' && b <= 'F') || (b >= 'a' && b <= 'f'))) {
              fixed += 'x';
              break;
            }
            int anum = (a >= '0' && a <= '9') ? a - 48 : (a >= 'A' && a <= 'F') ? a - 55 : a - 87;
            int bnum = (b >= '0' && b <= '9') ? b - 48 : (b >= 'A' && b <= 'F') ? b - 55 : b - 87;
            fixed += (char)(anum * 16 + bnum);
            i += 2;
            break;
          }
          fixed += 'x';
          break;
        default:
          fixed += next;
          break;
      }
      i += 1; //consume the next char that we just read
    } else {
      fixed += current;
    }
  }
  return fixed;
}

string::string(std::string val, bool n) : atom(n) {
  std::string fixed = FixEscapeCharacters(val);
  _value = fixed;
  _value_len = fixed.length();
}