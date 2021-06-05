#include "ast.hpp"

#include <sstream>

eresults compsetstmt::evaluate(std::unordered_map<std::string, eresults>* ctxt) {
  eresults results = _expression->evaluate(ctxt);
  (*ctxt)[_id] = results;
  return results;
}

eresults outputstmt::evaluate(std::unordered_map<std::string, eresults>* ctxt) {
  eresults results = _expression->evaluate(ctxt);
  results.type += 4;
  return results;
}

eresults funcdec::evaluate(std::unordered_map<std::string, eresults>* ctxt) {
  return {3, false, "", 0, this};
}

bool boolean(eresults r) {
  return !((r.type == 0 && r.b_value == false) ||
    (r.type == 1 && r.s_value == "") ||
    (r.type == 2 && r.n_value == 0) ||
    (r.type == 3 && r.e_value == nullptr));
}

std::string toString(eresults r) {
  if (r.type == 0) {
    return r.b_value ? "true" : "false";
  } else if (r.type == 1) {
    return r.s_value;
  } else if (r.type == 2) {
    std::stringstream ss = std::stringstream();
    ss << r.n_value;
    return ss.str();
  } else {
    return "";
  }
}

eresults binop::evaluate(std::unordered_map<std::string, eresults>* ctxt) {
  eresults results, lhs, rhs;
  std::string tempA, tempB;

  switch(_op) {
    case ops::AND:
      lhs = _lhs->evaluate(ctxt);
      if (boolean(lhs)) {
        rhs = _rhs->evaluate(ctxt);
        if (boolean(rhs)) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      } else {
        results = {0, false, "", 0, nullptr};
      }
      break;
    case ops::OR:
      lhs = _lhs->evaluate(ctxt);
      if (boolean(lhs)) {
        results = {0, true, "", 0, nullptr};
      } else {
        rhs = _rhs->evaluate(ctxt);
        if (boolean(rhs)) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      }
      break;
    case ops::EQ:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      tempA = toString(lhs);
      tempB = toString(rhs);
      if (tempA == tempB) {
        results = {0, true, "", 0, nullptr};
      } else {
        results = {0, false, "", 0, nullptr};
      }
      break;
    case ops::NEQ:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      tempA = toString(lhs);
      tempB = toString(rhs);
      if (tempA != tempB) {
        results = {0, true, "", 0, nullptr};
      } else {
        results = {0, false, "", 0, nullptr};
      }
      break;
    case ops::LT:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      if (lhs.type == 2 && rhs.type == 2) {
        if(lhs.n_value < rhs.n_value) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      } else if (lhs.type == 0 && rhs.type == 0) {
        if(!lhs.b_value && rhs.b_value) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      } else {
        tempA = toString(lhs);
        tempB = toString(rhs);
        if(tempA < tempB) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      }
      break;
    case ops::GT:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      if (lhs.type == 2 && rhs.type == 2) {
        if(lhs.n_value > rhs.n_value) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      } else if (lhs.type == 0 && rhs.type == 0) {
        if(lhs.b_value && !rhs.b_value) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      } else {
        tempA = toString(lhs);
        tempB = toString(rhs);
        if(tempA > tempB) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      }
      break;
    case ops::LTE:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      if (lhs.type == 2 && rhs.type == 2) {
        if(lhs.n_value <= rhs.n_value) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      } else if (lhs.type == 0 && rhs.type == 0) {
        if((!lhs.b_value && rhs.b_value) || (lhs.b_value == rhs.b_value)) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      } else {
        tempA = toString(lhs);
        tempB = toString(rhs);
        if(tempA <= tempB) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      }
      break;
    case ops::GTE:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      if (lhs.type == 2 && rhs.type == 2) {
        if(lhs.n_value >= rhs.n_value) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      } else if (lhs.type == 0 && rhs.type == 0) {
        if((lhs.b_value && !rhs.b_value) || (lhs.b_value == rhs.b_value)) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      } else {
        tempA = toString(lhs);
        tempB = toString(rhs);
        if(tempA >= tempB) {
          results = {0, true, "", 0, nullptr};
        } else {
          results = {0, false, "", 0, nullptr};
        }
      }
      break;
    case ops::ADD:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      if(lhs.type == 2 && rhs.type == 2) {
        results = {2, false, "", lhs.n_value + rhs.n_value, nullptr};
      } else if(lhs.type == 2 && rhs.type == 0) {
        results = {2, false, "", lhs.n_value + (rhs.b_value ? 1 : 0), nullptr};
      } else if(rhs.type == 2 && lhs.type == 0) {
        results = {2, false, "", rhs.n_value + (lhs.b_value ? 1 : 0), nullptr};
      } else if(lhs.type == 0 && rhs.type == 0) {
        results = {2, false, "", (u_int64_t)((lhs.b_value ? 1 : 0) + (rhs.b_value ? 1 : 0)), nullptr};
      } else {
        tempA = toString(lhs);
        tempB = toString(rhs);
        results = {1, false, tempA + tempB, 0, nullptr};
      }
      break;
    case ops::SUB:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      if(lhs.type == 2 && rhs.type == 2) {
        results = {2, false, "", lhs.n_value - rhs.n_value, nullptr};
      } else if(lhs.type == 2 && rhs.type == 0) {
        results = {2, false, "", lhs.n_value - (rhs.b_value ? 1 : 0), nullptr};
      } else if(rhs.type == 2 && lhs.type == 0) {
        results = {2, false, "", (lhs.b_value ? 1 : 0) - rhs.n_value, nullptr};
      } else if(lhs.type == 0 && rhs.type == 0) {
        results = {2, false, "", (u_int64_t)(lhs.b_value ? 1 : 0) - (rhs.b_value ? 1 : 0), nullptr};
      } else {
        results = {0, false, "", 0, nullptr};
      }
      break;
    case ops::MULT:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      if(lhs.type == 2 && rhs.type == 2) {
        results = {2, false, "", lhs.n_value * rhs.n_value, nullptr};
      } else if(lhs.type == 2 && rhs.type == 0) {
        results = {2, false, "", lhs.n_value * (rhs.b_value ? 1 : 0), nullptr};
      } else if(rhs.type == 2 && lhs.type == 0) {
        results = {2, false, "", rhs.n_value * (lhs.b_value ? 1 : 0), nullptr};
      } else if(lhs.type == 0 && rhs.type == 0) {
        results = {2, false, "", (u_int64_t)(lhs.b_value ? 1 : 0) * (rhs.b_value ? 1 : 0), nullptr};
      } else {
        results = {0, false, "", 0, nullptr};
      }
      break;
    case ops::DIV:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      if(lhs.type == 2 && rhs.type == 2) {
        results = {2, false, "", lhs.n_value / rhs.n_value, nullptr};
      } else if(lhs.type == 2 && rhs.type == 0) {
        results = {rhs.b_value ? 2 : 0, false, "", rhs.b_value ? lhs.n_value : 0, nullptr};
      } else if(rhs.type == 2 && lhs.type == 0) {
        results = {2, false, "", 0, nullptr};
      } else if(lhs.type == 0 && rhs.type == 0) {
        results = {(lhs.b_value && rhs.b_value) ? 2 : 0, false, "", (u_int64_t)((lhs.b_value && rhs.b_value) ? 1 : 0), nullptr};
      } else {
        results = {0, false, "", 0, nullptr};
      }
      break;
    case ops::MOD:
      lhs = _lhs->evaluate(ctxt);
      rhs = _rhs->evaluate(ctxt);
      if(lhs.type == 2 && rhs.type == 2) {
        results = {2, false, "", lhs.n_value % rhs.n_value, nullptr};
      } else if(lhs.type == 2 && rhs.type == 0) {
        results = {rhs.b_value ? 2 : 0, false, "", (u_int64_t)(rhs.b_value ? 1 : 0), nullptr};
      } else if(rhs.type == 2 && lhs.type == 0) {
        results = {2, false, "", 0, nullptr};
      } else if(lhs.type == 0 && rhs.type == 0) {
        results = {(lhs.b_value && rhs.b_value) ? 2 : 0, false, "", (u_int64_t)((lhs.b_value && rhs.b_value) ? 1 : 0), nullptr};
      } else {
        results = {0, false, "", 0, nullptr};
      }
      break;
    default:
      results = {0, false, "", 0, nullptr};
  }
  return results;
}

eresults fixOutput(eresults r) {
  r.type -= 4;
  return r;
}

eresults call::evaluate(std::unordered_map<std::string, eresults>* ctxt) {
  eresults found = (*ctxt)[_id];
  if(found.type != 3) return {1, false, "", 0, nullptr};

  auto _ps = *(found.e_value->_params);
  for (int i = 0; i < _ps.size(); i++) {
    auto p = _ps[i];
    if(i < _params->size()) {
      (*ctxt)[p] = (*_params)[i]->evaluate(ctxt);
    } else {
      (*ctxt)[p] = {1, false, "", 0, nullptr};
    }
  }

  auto _stmts = *(found.e_value->_stmts);
  eresults r;
  for (auto e : _stmts)
  {
    r = e->evaluate(ctxt);
    if (r.type >= 4) {
      return fixOutput(r);
    }
  }
  return r;  
}

eresults when::evaluate(std::unordered_map<std::string, eresults>* ctxt) {
  auto c = _condition->evaluate(ctxt);
  if (boolean(c)) {
    auto r = _then->evaluate(ctxt);
    r.type += 4;
    return r;
  }
  return c;
}

eresults caseexpr::evaluate(std::unordered_map<std::string, eresults>* ctxt) {
  for (auto w : *_when) {
    auto wr = w->evaluate(ctxt);
    if (wr.type >= 4) {
      return fixOutput(wr);
    }
  }
  return _expr->evaluate(ctxt);
}

eresults compnum::evaluate(std::unordered_map<std::string, eresults>* ctxt) {
  return {2, false, "", _value, nullptr};
}

eresults compstr::evaluate(std::unordered_map<std::string, eresults>* ctxt) {
  return {1, false, _value, 0, nullptr};
}

eresults compid::evaluate(std::unordered_map<std::string, eresults>* ctxt) {
  return (*ctxt)[_value];
}
