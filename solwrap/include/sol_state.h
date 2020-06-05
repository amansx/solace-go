#ifndef SOL_STATE_H
#define SOL_STATE_H
#include "sol_data.h"
#include "solclient/solClient.h"
#include "solclient/solCache.h"

#include <map>
#include <string>

struct sol_flow_state 
{
    char        qname[128];
    const char* flp[15];
    solClient_opaqueFlow_pt flow;
};
typedef std::map<std::string,sol_flow_state*> queue2flow_map;
typedef std::map<std::string,solClient_opaqueCacheSession_pt> cache_map;

/** 
 * All the state associated with one Solace Session to handle send/receive/delivery 
 * of Solace messages.
 **/
struct sol_state
{
    /**
     * ctx_ -- the Solace context instance.
     **/
    solClient_opaqueContext_pt ctx_;

    /**
     * sess_ -- the Solace session instance.
     **/
    solClient_opaqueSession_pt sess_;

    /**
     * cachemap_ -- the maps SolCache name to cache-session instance.
     **/
    cache_map cachemap_;

    /**
     * Solace session properties (host, user, passwd, etc.)
     **/
    const char** props_;

    queue2flow_map flowmap_;

    /**
     * Session Message callback function-pointer.
     **/
    message_cb msg_cb_;

    /**
     * Cached message event delivered to subscription-callbacks when a message is received.
     **/
    message_event recvmsg_;

    /**
     * Cached Solace destination object used to extract topic/queue string for delivery to subscription-callbacks.
     **/
    solClient_destination_t recvdest_;

    /**
     * Cached Solace message object reused for sending messages to the appliance.
     **/
    solClient_opaqueMsg_pt sendmsg_;

    /**
     * Registered callback for session errors.
     **/
    error_cb err_cb_;
    /**
     * Cached error event delivered to session callbacks when a session error occurs.
     **/
    error_event error_;

    /**
     * Registered callback for Publication event callback function-pointer 
     * (ACK, REJECT).
     **/
    pubevent_cb pub_cb_;

    /**
     * Registered callback for Connectivity event callback function-pointer 
     * (UP, Reconnecting, Reconnected, DOWN).
     **/
    connectivity_cb conn_cb_;

    /**
     * User Data state passed back to the application in all callbacks.
     **/
    void* user_data_;
};

#endif
