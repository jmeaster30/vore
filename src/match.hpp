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

class match {
public:
  std::string value;
  std::string lastMatch;
  u_int64_t file_offset;
  u_int64_t match_length;
  std::unordered_map<std::string, std::string> variables;
  std::unordered_map<std::string, primary*> subroutines;

  match(u_int64_t startOffset){
    value = "";
    lastMatch = "";
    file_offset = startOffset;
    match_length = 0;
    variables = std::unordered_map<std::string, std::string>();
    subroutines = std::unordered_map<std::string, primary*>();
  };

  match* copy();
  void print();
};

class context {
public:
  bool changeFile;
  bool dontStore;
  FILE* file;
  std::string input;
  std::vector<match*> matches;

  context(FILE* _file){
    file = _file;
    input = "";
    inputPointer = 0;
    matches = std::vector<match*>();
    changeFile = false;
    dontStore = false;
  };

  context(std::string _input) {
    file = nullptr;
    input = _input;
    inputPointer = 0;
    matches = std::vector<match*>();
    changeFile = false;
    dontStore = false;
  }

  std::string getChars(u_int64_t amount);
  void seekForward(u_int64_t value);
  void seekBack(u_int64_t value);
  void setPos(u_int64_t value);
  u_int64_t getPos();
  u_int64_t getSize();
  bool endOfFile();

  void print();

private:
  u_int64_t inputPointer;
};

#endif