#include "sol_api.h"
#include "sol_state.h"
#include "sol_props.h"
#include "sol_msg.h"
#include "sol_sessevent.h"
#include "sol_error.h"

#include "solclient/solClient.h"

#include <cstring>

SOLHANDLE
sol_init(message_cb msg_cb, error_cb err_cb, pubevent_cb pub_cb, connectivity_cb con_cb, void* user_data)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;
    sol_state* state = new sol_state;

    state->err_cb_    = err_cb;
    state->msg_cb_    = msg_cb;
    state->pub_cb_    = pub_cb;
    state->conn_cb_   = con_cb;
    state->user_data_ = user_data;

    if( (rc = solClient_initialize(SOLCLIENT_LOG_DEFAULT_FILTER, NULL)) 
           != SOLCLIENT_OK ) {
        on_error((SOLHANDLE)state, rc, "solClient_initialize()" );
    }

    solClient_log_setFilterLevel(SOLCLIENT_LOG_CATEGORY_ALL, SOLCLIENT_LOG_ERROR);

    // Create a Context allowing solclient lib to create the context thread
    solClient_log(SOLCLIENT_LOG_INFO, "Creating solClient context");
    solClient_context_createFuncInfo_t ctx_fn_info = 
                        SOLCLIENT_CONTEXT_CREATEFUNC_INITIALIZER;
    if( (rc = solClient_context_create(SOLCLIENT_CONTEXT_PROPS_DEFAULT_WITH_CREATE_THREAD, 
                                      &(state->ctx_), &ctx_fn_info,
                                      sizeof(solClient_context_createFuncInfo_t)))
           != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_context_create()" );
    }

    // Pre-allocate cached instances used in per-message send/recv flows
    sol_msg_alloc( &(state->sendmsg_) );

    return (SOLHANDLE) state;
}

char* global_buffer = 0;
int
sol_test_cbs(message_cb msg_cb, error_cb err_cb, pubevent_cb pub_cb, connectivity_cb con_cb, void* user_data)
{
    if (global_buffer == 0) {
        global_buffer = new char[256];
        strcpy(global_buffer, "message buffer");
    }

    char udata[256];
    strcpy(udata, "user-data pointer");

    error_event err;
    err.fn_name      = "err fn-name";
    err.return_code  = 1;
    err.rc_str       = "err rc-str";
    err.sub_code     = 2;
    err.sc_str       = "err sc-str";
    err.resp_code    = 3;
    err.err_str      = "err err-str";
    err_cb((SOLHANDLE)global_buffer, &err);

    message_event msg;
    msg.buffer           = (void*)global_buffer;
    msg.buflen           = strlen((char*)msg.buffer);
    msg.destination      = "topic/string/1";
    msg.desttype         = TOPIC;
    msg.redelivered_flag = 0;
    msg.discard_flag     = 0;
    msg.flow             = 0;
    msg.id               = 98765;
    msg_cb((SOLHANDLE)global_buffer, &msg);

    publisher_event pub; 
    pub.type = REJECT;
    pub.correlation_data = global_buffer;
    pub.user_data = udata;
    pub_cb((SOLHANDLE)global_buffer, &pub);

    connectivity_event conn;
    conn.type = RECONNECTING;
    conn.user_data = udata;
    con_cb((SOLHANDLE)global_buffer, &conn);

    return 0;
}


int 
sol_connect(SOLHANDLE handle, const char* propsfile)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;
    sol_state* state = (sol_state*)handle;

    state->props_ = read_props( propsfile );

    solClient_session_createFuncInfo_t ss_fn_info = 
                                SOLCLIENT_SESSION_CREATEFUNC_INITIALIZER;
    ss_fn_info.eventInfo.callback_p = on_event_cb;
    ss_fn_info.eventInfo.user_p     = state;
    ss_fn_info.rxMsgInfo.callback_p = on_msg_cb;
    ss_fn_info.rxMsgInfo.user_p     = state;

    solClient_log(SOLCLIENT_LOG_INFO, "creating solClient session" );
    if( (rc = solClient_session_create(state->props_, state->ctx_, &(state->sess_), 
                                       &ss_fn_info,
                                       sizeof(solClient_session_createFuncInfo_t)))
           != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_session_create()" );
    }

    solClient_log(SOLCLIENT_LOG_INFO, "connecting solClient session" );
    if( (rc = solClient_session_connect(state->sess_)) != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_session_connect()" );
    }

    return rc;
}
int
sol_disconnect(SOLHANDLE handle)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;
    sol_state* state = (sol_state*)handle;

    if ( (rc = solClient_session_disconnect(state->sess_)) != SOLCLIENT_OK ) 
        on_error( (SOLHANDLE)state, rc, "solClient_session_disconnect()" );
    if ( (rc = solClient_session_destroy(&(state->sess_))) != SOLCLIENT_OK ) 
        on_error( (SOLHANDLE)state, rc, "solClient_session_destroy()" );
    state->sess_ = 0;

    return rc;
}

int 
sol_send_direct(SOLHANDLE handle, const char* topic, void* buffer, int buflen)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;
    sol_state* state = (sol_state*)handle;

    // set the payload
    solClient_msg_setBinaryAttachmentPtr( state->sendmsg_, buffer, buflen );

    // Set direct mode for the message
    if( (rc = solClient_msg_setDeliveryMode(state->sendmsg_, SOLCLIENT_SEND_FLAGS_DIRECT))
            != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_setDeliveryMode()" );
        return rc;
    }
    // Set the dest
    solClient_destination_t dest;
    dest.destType = SOLCLIENT_TOPIC_DESTINATION;
    dest.dest     = topic;
    if( (rc = solClient_msg_setDestination(state->sendmsg_, &dest, sizeof(solClient_destination_t)))
            != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_setDestination()" );
        return rc;
    }
    // Send the message
    if ((rc = solClient_session_sendMsg(state->sess_, state->sendmsg_)) != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_session_sendMsg()" );
        return rc;
    }
    return rc;
}

int sol_send_persistent(SOLHANDLE handle, const char* destination, enum dest_type desttype, void* buffer, int buflen, void* correlation_p, int corrlen)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;
    sol_state* state = (sol_state*)handle;

    // set the payload
    solClient_msg_setBinaryAttachmentPtr( state->sendmsg_, buffer, buflen );

    // Set persistent mode for the message
    if( (rc = solClient_msg_setDeliveryMode(state->sendmsg_, SOLCLIENT_SEND_FLAGS_PERSISTENT))
            != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_setDeliveryMode()" );
        return rc;
    }
    // Set the dest
    solClient_destination_t dest;
    if (desttype == QUEUE) 
    	dest.destType = SOLCLIENT_QUEUE_DESTINATION;
    else
    	dest.destType = SOLCLIENT_TOPIC_DESTINATION;
    dest.dest = destination;
    if( (rc = solClient_msg_setDestination(state->sendmsg_, &dest, sizeof(solClient_destination_t)))
            != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_setDestination()" );
        return rc;
    }
    // Set the correlation ptr
    if ( (rc = solClient_msg_setCorrelationTagPtr(state->sendmsg_, correlation_p, corrlen)) 
            != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_setCorrelationTagPtr()" );
        return rc;
    }

    // Send the message
    if ((rc = solClient_session_sendMsg(state->sess_, state->sendmsg_)) != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_session_sendMsg()" );
        return rc;
    }
    return rc;
}


int 
sol_subscribe_topic(SOLHANDLE handle, const char* topic)
{
    sol_state* state = (sol_state*)handle;
    return solClient_session_topicSubscribeExt(state->sess_, 
                            SOLCLIENT_SUBSCRIBE_FLAGS_WAITFORCONFIRM, topic);
}

int 
sol_unsubscribe_topic(SOLHANDLE handle, const char* topic)
{
    sol_state* state = (sol_state*)handle;
    return solClient_session_topicUnsubscribeExt(state->sess_, 
                            SOLCLIENT_SUBSCRIBE_FLAGS_WAITFORCONFIRM, topic);
}


void setup_flow_props(sol_flow_state* fstate, const char* queue, fwd_mode fm, ack_mode am)
{
    memcpy( (fstate->qname), queue, strlen(queue)+1 );
    // Flow properties object
    int p = 0;
    fstate->flp[p++] = SOLCLIENT_FLOW_PROP_BIND_NAME;          fstate->flp[p++] = queue;
    fstate->flp[p++] = SOLCLIENT_FLOW_PROP_BIND_ENTITY_ID;     fstate->flp[p++] = SOLCLIENT_FLOW_PROP_BIND_ENTITY_QUEUE;
    fstate->flp[p++] = SOLCLIENT_FLOW_PROP_BIND_BLOCKING;      fstate->flp[p++] = SOLCLIENT_PROP_ENABLE_VAL;
    if (am == MANUAL_ACK) {
        fstate->flp[p++] = SOLCLIENT_FLOW_PROP_ACKMODE;        fstate->flp[p++] = SOLCLIENT_FLOW_PROP_ACKMODE_CLIENT;
    }
    else {
        fstate->flp[p++] =SOLCLIENT_FLOW_PROP_ACKMODE;         fstate->flp[p++] = SOLCLIENT_FLOW_PROP_ACKMODE_AUTO;
    }
    // Disable to begin in the stopped state (e.g. if adding multiple flows)
    fstate->flp[p++] =SOLCLIENT_FLOW_PROP_START_STATE;         fstate->flp[p++] = SOLCLIENT_PROP_ENABLE_VAL;
    if (fm == CUT_THRU) {
        fstate->flp[p++] =SOLCLIENT_FLOW_PROP_FORWARDING_MODE; fstate->flp[p++] = SOLCLIENT_FLOW_PROP_FORWARDING_MODE_CUT_THROUGH;
    }
    else {
        fstate->flp[p++] =SOLCLIENT_FLOW_PROP_FORWARDING_MODE; fstate->flp[p++] = SOLCLIENT_FLOW_PROP_FORWARDING_MODE_STORE_AND_FORWARD;
    }
    fstate->flp[p++] =SOLCLIENT_FLOW_PROP_WINDOWSIZE;          fstate->flp[p++] = "255";
    fstate->flp[p] = 0;
}

int
sol_bind_queue(SOLHANDLE handle, const char* queue, fwd_mode fm, ack_mode am)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;

    sol_state* state = (sol_state*)handle;

    solClient_flow_createFuncInfo_t flw_fn_info = 
                                        SOLCLIENT_FLOW_CREATEFUNC_INITIALIZER;
    flw_fn_info.rxMsgInfo.callback_p = on_flow_msg_cb;
    flw_fn_info.rxMsgInfo.user_p     = state;
    flw_fn_info.eventInfo.callback_p = on_flow_event_cb;
    flw_fn_info.eventInfo.user_p     = state;

    sol_flow_state* fstate = new sol_flow_state;
    setup_flow_props( fstate, queue, fm, am );

    // Create the flow
    if ( (rc = solClient_session_createFlow(fstate->flp, state->sess_, &(fstate->flow), 
                                           &flw_fn_info, sizeof(flw_fn_info))) 
            != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_session_createFlow()" );
    }
    if ( (rc = solClient_flow_start(fstate->flow)) != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_flow_start()" );
    }

    state->flowmap_[fstate->qname] = fstate;

    return rc;
}

int
sol_unbind_queue(SOLHANDLE handle, const char* queue)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;

    sol_state* state = (sol_state*)handle;
    sol_flow_state* fstate = state->flowmap_[queue];
    
    if ( (rc = solClient_flow_destroy(&(fstate->flow))) != SOLCLIENT_OK )
        on_error( handle, rc, "solClient_flow_destroy()" );
    
    state->flowmap_.erase( queue );
    delete fstate;
    
    return rc;
}

int
sol_ack_msg(SOLHANDLE handle, FLOWHANDLE flow, SOLMSGID msg_id)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;
    if( (rc = solClient_flow_sendAck((solClient_opaqueFlow_pt)flow, msg_id)) != SOLCLIENT_OK )
        on_error(handle, rc, "solClient_flow_sendAck()" );
    return rc;
}

solClient_opaqueCacheSession_pt 
sol_cache_session(sol_state* state, const char* cache_name)
{
    cache_map::iterator it = state->cachemap_.find( cache_name );
    if (state->cachemap_.end() == it) {
        // Create the cache-session
        const char* props_p[3];
        int         pi = 0;
        props_p[pi++] = SOLCLIENT_CACHESESSION_PROP_CACHE_NAME;
        props_p[pi++] = cache_name;
        props_p[pi++] = 0;
        
        solClient_returnCode_t rc = SOLCLIENT_OK;
        solClient_opaqueCacheSession_pt csess_p;
        if ( (rc = solClient_session_createCacheSession ( (const char *const*)props_p,
                                                       state->sess_, &csess_p) ) != SOLCLIENT_OK ) {
            on_error ( (SOLHANDLE)state, rc, "solClient_session_createCacheSession()" );
            return 0;
        }
        state->cachemap_[cache_name] = csess_p;
        return csess_p;
    }
    return it->second;
}

int 
sol_cache_req(SOLHANDLE handle, const char* cache_name, const char* topic_sub, int request_id)
{
    sol_state* state = (sol_state*)handle;
    solClient_opaqueCacheSession_pt csess_p = sol_cache_session( state, cache_name );
    solClient_cacheRequestFlags_t flags = SOLCLIENT_CACHEREQUEST_FLAGS_LIVEDATA_FLOWTHRU;
    
    solClient_returnCode_t rc = SOLCLIENT_OK;
    if ( (rc = solClient_cacheSession_sendCacheRequest (csess_p, topic_sub, request_id, 
                                                    0 /* cache-event cbfn */, state, flags, 0)) != SOLCLIENT_OK ) {
        on_error ( handle, rc, "solClient_cacheSession_sendCacheRequest()" );
    }
    return rc;
}

