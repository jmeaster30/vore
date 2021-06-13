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

inline eresults ebool(bool value) { return {0, value, "", 0, nullptr}; }
inline eresults estring(std::string value) { return {1, false, value, 0, nullptr}; }
inline eresults enumber(u_int64_t value) { return {2, false, "", value, nullptr}; }
inline eresults efunc(funcdec* value) { return {3, false, "", 0, value}; }

class match {
public:
  std::string value = "";
  std::string replacement = "";
  std::string lastMatch = "";
  u_int64_t file_offset = 0;
  u_int64_t lineNumber = 0;
  u_int64_t matchNumber = 0;
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
  bool appendFile = false;
  bool changeFile = false;
  bool dontStore = false;
  FILE* file = nullptr;
  std::string filename = "";
  std::string input = "";
  std::vector<match*> matches = std::vector<match*>();
  std::unordered_map<std::string, eresults> global  = std::unordered_map<std::string, eresults>();

  context() {
    context(nullptr);
  }

  context(std::string _filename, FILE* _file){
    file = _file;
    filename = _filename;
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
