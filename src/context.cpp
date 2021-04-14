#include "context.hpp"

std::string context::peek(size_t length) {
  char* buf = (char*)malloc((length + 1) * sizeof(char));
  memset(buf, 0, length + 1);

  peek_size = fread(buf, sizeof(char), length, file);
  if(fseek(file, -peek_size, SEEK_CUR))
  {
    printf("UH OH SOMETHING BAD HAPPENED WHEN DOING THE SEEK :(");
    exit(1);
  }

  peek_buffer = buf;
  return peek_buffer;
}

std::string context::consume(size_t length) {
  //update the latest match
  peek(length);

  auto latest_match = matches->back();
  if(latest_match->file_offset == -1)
  {
    latest_match->file_offset = ftell(file);
  }

  if(fseek(file, peek_size, SEEK_CUR))
  {
    printf("UH OH SOMETHING BAD HAPPENED WHEN DOING THE CONSUME SEEK :(");
    exit(1);
  }

  latest_match->value += peek_buffer;
  latest_match->match_length += peek_size;

  startOfLine = peek_buffer[peek_size - 1] == '\n';

  return peek_buffer;
}

u_int64_t context::filepos() {
  return ftell(file);
}

bool context::isStartOfLine() {
  return startOfLine;
}

bool context::isEndOfFile() {
  return feof(file) != 0;
}

void context::addvar(std::string name, std::string value)
{
  auto latest_match = matches->back();
  (*(latest_match->variables))[name] = value;
}

std::string context::getvar(std::string name)
{
  auto latest_match = matches->back();
  auto found = latest_match->variables->find(name);
  if(found == latest_match->variables->end())
    return "";
  else
    return found->second;
}