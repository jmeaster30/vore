#include "match.hpp"

#include <iostream>
#include "ast.hpp"

match* match::copy()
{
  match* newMatch = new match(file_offset);
  newMatch->value = value;
  newMatch->lastMatch = lastMatch;
  newMatch->match_length = match_length;
  newMatch->variables = variables;
  newMatch->subroutines = subroutines;
  return newMatch;
}

void match::print()
{
  std::cout << "===== START MATCH =====" << std::endl;

  std::cout << "value = '" << value << "'" << std::endl;
  std::cout << "fileOffset = " << file_offset << std::endl;
  std::cout << "matchLength = " << match_length << std::endl;

  std::cout << "## variables " << std::endl;
  for(const auto& [key, value] : variables) {
    std::cout << "'" << key << "' = '" << value << "'" << std::endl;
  }

  std::cout << "## subroutines " << std::endl;
  for(auto [key, value] : subroutines) {
    std::cout << "'" << key << "' = (";
    value->print();
    std::cout << ")" << std::endl;
  }
 
  std::cout << "===== END MATCH   =====" << std::endl;
}

void context::print()
{
  std::cout << "---------- START CONTEXT -------------" << std::endl;
  std::cout << "change? " << changeFile << " no store? " << dontStore << std::endl;
  std::cout << "file pointer " << file << std::endl;

  std::cout << "START CONTEXT MATCHES" << std::endl;
  for (auto match : matches) {
    match->print();
  }
  std::cout << "END CONTEXT MATCHES" << std::endl;

  std::cout << "---------- END CONTEXT   -----------" << std::endl;
}

std::string context::getChars(u_int64_t amount) {
  std::string result;
  if(file != nullptr) {
    char* buffer = (char*)malloc(amount * sizeof(char));
    memset(buffer, 0, amount * sizeof(char));
    u_int64_t read_bytes = fread(buffer, sizeof(char), amount, file);
    result = std::string(buffer, read_bytes);
    free(buffer);
  } else {
    u_int64_t fixedAmount = amount;
    if (input.length() < inputPointer + amount) {
      fixedAmount = input.length() - inputPointer; //just get the remaining characters
    }
    result = input.substr(inputPointer, fixedAmount);
    inputPointer += fixedAmount;
  }
  return result;
}

void context::seekForward(u_int64_t value) {
  //make it so this doesn't go past the end of the file/input
  if(file != nullptr) {
    fseek(file, value, SEEK_CUR);
  }else{

    inputPointer += value;
  }
}

void context::seekBack(u_int64_t value) {
  if (file != nullptr) {
    u_int64_t current = ftell(file);
    if(current <= value) {
      fseek(file, 0, SEEK_SET);
    } else {
      fseek(file, current - value, SEEK_SET);
    }
  } else {
    if (inputPointer < value) {
      inputPointer = 0;
    } else {
      inputPointer -= value;
    }
  }
}

void context::setPos(u_int64_t value) {
  if (file != nullptr)
    fseek(file, value, SEEK_SET);
  else
    inputPointer = value;
}

u_int64_t context::getPos() {
  if(file != nullptr)
    return ftell(file);
  else
    return inputPointer;
}

u_int64_t context::getSize() {
  if(file != nullptr) {
    fseek(file, 0, SEEK_END);
    u_int64_t size = ftell(file);
    fseek(file, 0, SEEK_SET);
    return size;
  } else {
    return input.length();
  }
}

bool context::endOfFile() {
  if(file != nullptr)
    return feof(file); //I have had issues with this 
  else
    return inputPointer == input.length();
}
