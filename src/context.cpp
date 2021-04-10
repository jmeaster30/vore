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

void context::consume() {
  //update the latest match
  auto latest_match = matches->back();
  if(latest_match->file_offset == -1)
  {
    latest_match->file_offset = ftell(file);
    latest_match->start_line_number = line_number;
    latest_match->start_column_number = column_number;
    latest_match->end_line_number = line_number;
    latest_match->end_column_number = column_number;
  }

  latest_match->value += peek_buffer;
  latest_match->match_length += peek_size;

  //update the end line and end column values here

}