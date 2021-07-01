#include "context.hpp"

#include <iostream>
#include "ast.hpp"

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
  if(getChars(1) == "") {
    return true;
  } else {
    seekBack(1);
    return false;
  }
}
