#include "sol_api.h"
#include <iostream>
#include <sstream>
#include <cstring>
#include <cstdlib>
#include <iomanip>

#include <time.h>
#include <sys/time.h>


void mysleep(int secs)
{
#ifdef LINUX
    sleep(secs);
#endif
#ifdef WIN32
    Sleep(secs * 1000);
#endif
}

#define MAX 1000000

void on_msg(SOLHANDLE handle, message_event* msg)
{
return ;
    std::cout << "MSG ON HANDLE "    << handle << std::endl
              << "Msg("<<msg->buflen<<"): " << (const char*)msg->buffer <<std::endl
              << "\tRedelivered? " << msg->redelivered_flag <<std::endl
              << "\tDiscards? "    << msg->discard_flag <<std::endl
              << "\tTopic: "       << (msg->desttype == TOPIC ? msg->destination : "") << std::endl
              << "\tQueue: "       << (msg->desttype == QUEUE ? msg->destination : "") << std::endl
              << "\tFLOW: "        << msg->flow << std::endl
              << "\tID: "          << msg->id << std::endl ;
    if (msg->flow != 0)
        sol_ack_msg( handle, msg->flow, msg->id );
}

void on_pubevent(SOLHANDLE handle, publisher_event* pubevent) 
{
return ;
    std::cout << "Pubevent! " 
              << (pubevent->type == ACK ? "ACK" : "REJECT")
              << " ID: " << *(int*)pubevent->correlation_data
              << std::endl;
}

void on_error(SOLHANDLE handle, error_event* err)
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

void on_connectivity(SOLHANDLE handle, connectivity_event* connevent)
{
    std::cerr << "Connectivity!" << con2str(connevent->type) << std::endl;
}

long time_delta(struct timeval& start, struct timeval& end)
{
 return ((end.tv_sec * 1000000 + end.tv_usec)
		  - (start.tv_sec * 1000000 + start.tv_usec));
}

int main(int c, char** a) 
{
    if (c < 3) {
        std::cout << "USAGE: " << a[0] << " <PROPERTIES-FILE> <#-subscriptions>" << std::endl;
        ::exit(1);
    }
    std::stringstream args( a[2] );
    int numsubs;
    args >> numsubs;

    SOLHANDLE subhandle = sol_init( on_msg, on_error, on_pubevent, on_connectivity, NULL );
    std::cout << "\tSUBHANDLE: " << subhandle << std::endl;

    int rc = sol_connect( subhandle, a[1] );

    sol_subscribe_topic( subhandle, "warm/up/topic" );
    for( int wu=0; wu < 1000; ++wu) {
        sol_send_direct( subhandle, "warm/up/topic", (void*)"0000", 5 );
    }
    sol_unsubscribe_topic( subhandle, "warm/up/topic" );

    mysleep(1);

    char* topics[ MAX ];
    for( int t = 0; t < numsubs; ++t) {
        std::stringstream ss;
        topics[t] = new char[ 30 ];
        ss << "topic/string/for/pico/" << std::setw(7) << std::setfill('0') << t;
        std::string tstr = ss.str();
        memcpy( topics[t], tstr.c_str(), tstr.length() + 1 );
    }

  struct timeval start, end;

    gettimeofday( &start, 0 );
    for(int a = 0; a < numsubs; ++a) {
        sol_subscribe_topic( subhandle, topics[a] );
    }
    gettimeofday( &end, 0 );

    mysleep(1);

    long usecs = time_delta( start, end );
    std::cout << numsubs << " subscribes = " << usecs << std::endl;

    gettimeofday( &start, 0 );
    for(int a = 0; a < numsubs; ++a) {
        sol_unsubscribe_topic( subhandle, topics[a] );
    }
    gettimeofday( &end, 0 );

    mysleep(1);

    usecs = time_delta( start, end );
    std::cout << numsubs << " unsubscribes = " << usecs << std::endl;

    rc = sol_disconnect( subhandle );

    return 0;
}
