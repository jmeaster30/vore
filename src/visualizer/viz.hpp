#pragma once

#ifdef WITH_VIZ
#define VIZ_EXTEND : public Viz::Viz
#define VIZ_FUNC void visualize(Agraph_t* subgraph);
#define VIZ_VFUNC virtual void visualize(Agraph_t* subgraph) {};
#else
#define VIZ_EXTEND
#define VIZ_FUNC
#define VIZ_VFUNC
#endif

#ifdef WITH_VIZ
#include <graphviz/cgraph.h>

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
    Agnode_t* node = nullptr;

    virtual void visualize(Agraph_t* subgraph) = 0;
  };

  void render(std::string filename, std::vector<Compiler::Statement*> statements); 
}
#endif
