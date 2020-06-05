#ifndef SOL_SESSEVENT_H
#define SOL_SESSEVENT_H
#if defined(_MSC_VER)
/* Windows-only includes */
#include <winsock2.h>
#else
/* Unix-only includes */
#include <unistd.h>
#endif

#include "solclient/solClient.h"

void 
on_event_cb(solClient_opaqueSession_pt sess_p, solClient_session_eventCallbackInfo_pt eventInfo_p, void *user_p);

void 
on_flow_event_cb(solClient_opaqueFlow_pt opaqueFlow_p, solClient_flow_eventCallbackInfo_pt eventInfo_p, void *user_p);

#endif
