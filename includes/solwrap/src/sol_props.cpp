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


const char** read_prop_params(const char* host, const char* vpn, const char* user, const char* pass, const char* clientName, const char* appDesc, const char* windowsize) {

    std::string SESSION_HOST_K = "SESSION_HOST";
    std::string SESSION_HOST_V = host;
    
    std::string SESSION_VPN_NAME_K = "SESSION_VPN_NAME";
    std::string SESSION_VPN_NAME_V = vpn;

    std::string SESSION_USERNAME_K = "SESSION_USERNAME";
    std::string SESSION_USERNAME_V = user;
    
    std::string SESSION_PASSWORD_K = "SESSION_PASSWORD";
    std::string SESSION_PASSWORD_V = pass;
    
    std::string SESSION_PUB_WINDOW_SIZE_K = "SESSION_PUB_WINDOW_SIZE";
    std::string SESSION_PUB_WINDOW_SIZE_V = windowsize;

    std::string SESSION_CLIENT_NAME_K = "SESSION_CLIENT_NAME";
    std::string SESSION_CLIENT_NAME_V = clientName;

    std::string SESSION_APP_DESC_K = "SESSION_APPLICATION_DESCRIPTION";
    std::string SESSION_APP_DESC_V = appDesc;

    std::map<std::string, std::string> props;
    props[SESSION_HOST_K]            = trim(SESSION_HOST_V);
    props[SESSION_VPN_NAME_K]        = trim(SESSION_VPN_NAME_V);
    props[SESSION_USERNAME_K]        = trim(SESSION_USERNAME_V);
    props[SESSION_PASSWORD_K]        = trim(SESSION_PASSWORD_V);
    props[SESSION_PUB_WINDOW_SIZE_K] = trim(SESSION_PUB_WINDOW_SIZE_V);
    props[SESSION_APP_DESC_K]        = trim(SESSION_APP_DESC_V);
    props[SESSION_CLIENT_NAME_K]     = trim(SESSION_CLIENT_NAME_V);

    // Retry Forever
    props["SESSION_CONNECT_RETRIES_PER_HOST"] = "-1";
    props["SESSION_CONNECT_RETRIES"] = "-1";
    props["SESSION_RECONNECT_RETRIES"] = "-1";
    props["SESSION_RECONNECT_RETRY_WAIT_MS"] = "10000";


    const char** sp = new const char*[ props.size()*2 +1 ];
    std::for_each( props.begin(), props.end(), inserter(sp) );
    sp[ 2*props.size() ] = 0;
    return sp;
}
