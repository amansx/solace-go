#ifndef SOL_MSG_H
#define SOL_MSG_H
#if defined(_MSC_VER)
/* Windows-only includes */
#include <winsock2.h>
#else
/* Unix-only includes */
#include <unistd.h>
#endif
#include "solclient/solClientMsg.h"

int sol_msg_alloc(solClient_opaqueMsg_pt *msg_p);

solClient_rxMsgCallback_returnCode_t 
on_msg_cb(solClient_opaqueSession_pt sess_p, solClient_opaqueMsg_pt msg_p, void *user_p);

solClient_rxMsgCallback_returnCode_t
on_flow_msg_cb(solClient_opaqueFlow_pt opaqueFlow_p, solClient_opaqueMsg_pt msg_p, void *user_p);

#endif
