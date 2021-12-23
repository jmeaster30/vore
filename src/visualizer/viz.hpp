#pragma once

#ifdef WITH_VIZ
#include <graphviz/cgraph.h>
#endif

#include <vector>
#include <string>

//forward declare
namespace Compiler {
  class Statement;
}

namespace Viz
{
  class Viz
  {
  public:
#ifdef WITH_VIZ
    Agnode_t* node = nullptr;
#endif
    virtual void visualize() = 0;
  };

  void render(std::string filename, std::vector<Compiler::Statement*> statements); 
}
