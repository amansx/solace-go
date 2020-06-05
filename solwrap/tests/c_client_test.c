#include "sol_api.h"
#include <stdio.h>
#include <string.h>

void on_msg(SOLHANDLE handle, message_event* msg)
{
	printf("GOT IT, OKAY!?!?!?\n");
}

void on_pubevent(SOLHANDLE handle, publisher_event* pubevent) 
{
	printf("PUB EVENT!?\n");
}

void on_error(SOLHANDLE handle, error_event* err)
{
	printf("ERROR!?\n");
}

const char* con2str(enum connectivity_event_type type) {
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
	fprintf(stderr, "Connectivity! %s\n", con2str(connevent->type));
}

int main(int c, char** a) 
{
    SOLHANDLE pubhandle, subhandle;
    int rc, buflen, i, j;
    void* buffer;

    pubhandle = sol_init( on_msg, on_error, on_pubevent, on_connectivity, NULL );
    subhandle = sol_init( on_msg, on_error, on_pubevent, on_connectivity, NULL );

    rc = sol_connect( pubhandle, a[1] );
    rc = sol_connect( subhandle, a[1] );

    buffer = (void*)"howdy doody";
    buflen = strlen((const char*)buffer) + 1;

    for( j = 0; j < 10; ++j) {

        sol_subscribe_topic( subhandle, "t/1/2/3" );
        for(i=0; i<10; ++i) {
            sol_send_direct( pubhandle, "t/1/2/3", buffer, buflen );
        }
        sol_unsubscribe_topic( subhandle, "t/1/2/3" );
    }

    return 0;
}
