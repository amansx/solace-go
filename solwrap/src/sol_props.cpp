#include "sol_props.h"

#include <map>
#include <string>
#include <iostream>
#include <sstream>
#include <fstream>
#include <algorithm>
#include <functional> 
#include <cctype>
#include <locale>
#include <memory>
#include <stdexcept>

namespace
{
    // trim from start
    static inline std::string& 
    ltrim(std::string &s) {
        s.erase(s.begin(), 
                std::find_if(s.begin(), 
                             s.end(), 
                             std::not1(std::ptr_fun<int, int>(std::isspace))
                )
        );
        return s;
    }
    // trim from end
    static inline std::string& 
    rtrim(std::string &s) {
        s.erase(std::find_if(s.rbegin(), 
                             s.rend(), 
                             std::not1(std::ptr_fun<int, int>(std::isspace))
                ).base(), 
                s.end());
        return s;
    }
    // trim from both ends
    static inline std::string& 
    trim(std::string &s) {
        return ltrim(rtrim(s));
    }


    struct inserter
    {
        inserter(const char** props) : p_(props), pos_(0) {}
        void operator() (const std::pair<std::string,std::string>& pr) {
            p_[pos_++] = alloc_init(pr.first);
            p_[pos_++] = alloc_init(pr.second);
        }
        const char* alloc_init(const std::string& rhs)
        {
            size_t len = rhs.length();
            char* s = new char[ len+1 ];
            rhs.copy(s, len);
            s[len] = 0;
            return s;
        }
    private:
        const char** p_;
        int pos_;
    };
}


const char** read_props(const char* propsfile)
{
    std::map<std::string, std::string> props;

    std::ifstream inf(propsfile);
    if (inf.is_open())
    {
        std::string k, v;
        while (inf.good())
        {
            std::getline(inf, k, '=');
            std::getline(inf, v);
            if (k.length()) {
                props[trim(k)] = trim(v);
            }
        }
        inf.close();
    }

    const char** sp = new const char*[ props.size()*2 +1 ];
    std::for_each( props.begin(), props.end(), inserter(sp) );
    sp[ 2*props.size() ] = 0;
    return sp;
}
