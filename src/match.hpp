#ifndef __match_h__
#define __match_h__

#include <stdio.h>
#include <string>
#include <string.h> // for memset
#include <vector>
#include <unordered_map>
#include <stack>

//forward declarations
class primary;
class element;
class funcdec;

struct eresults {
  int type;
  bool b_value;
  std::string s_value;
  u_int64_t n_value;
  funcdec* e_value;
};

class match {
public:
  std::string value = "";
  std::string replacement = "";
  std::string lastMatch = "";
  u_int64_t file_offset = 0;
  u_int64_t lineNumber = 0;
  u_int64_t match_length = 0;
  std::unordered_map<std::string, std::string> variables = std::unordered_map<std::string, std::string>();
  std::unordered_map<std::string, primary*> subroutines = std::unordered_map<std::string, primary*>();

  match(u_int64_t startOffset){
    file_offset = startOffset;
  };

  match* copy();
  void print();
};

class context {
public:
  bool changeFile = false;
  bool dontStore = false;
  FILE* file = nullptr;
  std::string input = "";
  std::vector<match*> matches = std::vector<match*>();
  std::unordered_map<std::string, eresults> global  = std::unordered_map<std::string, eresults>();

  context() {
    context(nullptr);
  }

  context(FILE* _file){
    file = _file;
  };

  context(std::string _input) {
    input = _input;
  }

  std::string getChars(u_int64_t amount);
  void seekForward(u_int64_t value);
  void seekBack(u_int64_t value);
  void setPos(u_int64_t value);
  u_int64_t getPos();
  u_int64_t getSize();
  bool endOfFile();

  void print();
  context* copy();

private:
  u_int64_t inputPointer = 0;
};

#endif
