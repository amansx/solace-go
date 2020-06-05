#include "sol_api.h"
#include <iostream>
#include <cstring>
#include <cstdlib>

void mysleep(int secs)
{
#ifdef LINUX
    sleep(secs);
#endif
#ifdef WIN32
    Sleep(secs * 1000);
#endif
}

void on_msg(SOLHANDLE session, message_event* msg)
{
    std::cout << "MSG ON HANDLE "    << session << std::endl
              << "Msg("<<msg->buflen<<"): " << (const char*)msg->buffer <<std::endl
              << "\tRedelivered? " << msg->redelivered_flag <<std::endl
              << "\tDiscards? "    << msg->discard_flag <<std::endl
              << "\tDest: "        << (msg->destination != 0 ?  msg->destination : "") << std::endl
              << "\tDest-Type: "   << (msg->desttype == QUEUE ? "QUEUE" : "TOPIC") << std::endl
              << "\tFLOW: "        << msg->flow << std::endl
              << "\tID: "          << msg->id << std::endl ;
    if (msg->flow != 0)
        sol_ack_msg( session, msg->flow, msg->id );
}

void on_error(SOLHANDLE session, error_event* err)
{
    std::cerr << err->fn_name  
              << ": ReturnCode("  << err->return_code <<")=" << err->rc_str << std::endl
              << "\tSubCode("     << err->sub_code    <<")=" << err->sc_str << std::endl
              << "\tresponseCode="<< err->resp_code                         << std::endl
              << "\tInfo="        << err->err_str                           << std::endl;
}

const char* con2str(connectivity_event_type type) {
    switch(type) {
    case UP: return "UP";
    case DOWN: return "DOWN";
    case RECONNECTING: return "RECONNECTING";
    case RECONNECTED: return "RECONNECTED";
    }
    return "???";
}

void on_connectivity(SOLHANDLE session, connectivity_event* connevent)
{
    std::cerr << "Connectivity!" << con2str(connevent->type) << std::endl;
}

int main(int c, char** a) 
{
    if (c < 2) {
        std::cout << "USAGE: " << a[0] << " <PROPERTIES-FILE>" << std::endl;
        ::exit(1);
    }


    callback_functions cbs;
    cbs.err_cb = on_error;
    cbs.msg_cb = on_msg;
    cbs.pub_cb = NULL;
    cbs.con_cb = on_connectivity;
    cbs.user_data = NULL;

    SOLHANDLE session = sol_init( on_msg, on_error, NULL, on_connectivity, NULL );
    std::cout << "\tSESSION: " << session << std::endl;

    int rc = sol_connect( session, a[1] );

    const char* direct_topic = "t/1/2/3";
    void* buffer = (void*)"howdy doody";
    int buflen = strlen((const char*)buffer) + 1; // include the null-terminator

    for( int a = 0; a < 10; ++a) {

        sol_subscribe_topic( session, direct_topic );
        for(int i=0; i<10; ++i) {
            sol_send_direct( session, direct_topic, buffer, buflen );
        }
        sol_unsubscribe_topic( session, direct_topic );

    }
    mysleep(3);

    return 0;
}
