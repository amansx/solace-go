#include "sol_api.h"
#include <iostream>
#include <cstring>
#include <cstdlib>

void on_msg(SOLHANDLE handle, message_event* msg)
{
    std::cout << "MSG ON HANDLE "    << handle << std::endl
              << "Msg("<<msg->buflen<<"): " << (const char*)msg->buffer <<std::endl
              << "\tRedelivered? " << msg->redelivered_flag <<std::endl
              << "\tDiscards? "    << msg->discard_flag <<std::endl
              << "\tRequestID: "   << msg->req_id << std::endl
              << "\tDest: "        << (msg->destination != 0 ?  msg->destination : "") << std::endl
              << "\tDest-Type: "   << (msg->desttype == QUEUE ? "QUEUE" : "TOPIC") << std::endl
              << "\tFLOW: "        << msg->flow << std::endl
              << "\tID: "          << msg->id << std::endl ;
}

void on_error(SOLHANDLE handle, error_event* err)
{
    std::cerr << err->fn_name  
              << ": ReturnCode("  << err->return_code <<")=" << err->rc_str << std::endl
              << "\tSubCode("     << err->sub_code    <<")=" << err->sc_str << std::endl
              << "\tresponseCode="<< err->resp_code                         << std::endl
              << "\tInfo="        << err->err_str                           << std::endl;
}

int main(int c, char** a) 
{
    if (c < 2) {
        std::cout << "USAGE: " << a[0] << " <PROPERTIES-FILE>" << std::endl;
        ::exit(1);
    }

    SOLHANDLE session = sol_init( on_msg, on_error, NULL, NULL, NULL );
    std::cout << "\tSESSION: " << session << std::endl;

    int rc = sol_connect( session, a[1] );

    void* buffer = (void*)"howdy doody";
    int buflen = strlen((const char*)buffer) + 1; // include the null-terminator

    sol_send_direct( session, "cache/topic/1", buffer, buflen );
    sol_send_direct( session, "cache/topic/2", buffer, buflen );
    sol_send_direct( session, "cache/topic/3", buffer, buflen );
    sol_send_direct( session, "cache/topic/4", buffer, buflen );
    sol_send_direct( session, "cache/topic/5", buffer, buflen );

    sol_cache_req( session, "pysolcache", "cache/topic/>", 4321 );

    std::cout << "DONE" << std::endl;
    return 0;
}
