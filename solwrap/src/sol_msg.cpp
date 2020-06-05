#include "sol_msg.h"
#include "sol_state.h"
#include "sol_error.h"
#ifdef PYTHON_SUPPORT
#include <Python.h>
#endif

int
sol_msg_alloc(solClient_opaqueMsg_pt *msg_p)
{
    return solClient_msg_alloc(msg_p);
}

int 
sol_msg_redelivered_flag(solClient_opaqueMsg_pt msg_p)
{
    return solClient_msg_isRedelivered( msg_p );
}

int 
sol_msg_discard_flag(solClient_opaqueMsg_pt msg_p)
{
    return solClient_msg_isDiscardIndication( msg_p );
}

int 
sol_msg_destination(solClient_opaqueMsg_pt msg_p, 
        			solClient_destination_t* dest, 
        			message_event* msg)
{
    solClient_returnCode_t rc = solClient_msg_getDestination( msg_p, dest, sizeof(solClient_destination_t) );
    if ( rc == SOLCLIENT_OK ) {
        if (dest->destType == SOLCLIENT_TOPIC_DESTINATION) {
        	msg->destination = dest->dest;
        	msg->desttype = TOPIC;
        }
        else if (dest->destType == SOLCLIENT_QUEUE_DESTINATION) {
        	msg->destination = dest->dest;
        	msg->desttype = QUEUE;
        }
    }
    else {
        msg->destination = 0;
        msg->desttype = NONE;
    }
    return rc;
}

int 
sol_msg_payload(solClient_opaqueMsg_pt msg_p, message_event* msg) 
{
    return solClient_msg_getBinaryAttachmentPtr( msg_p, &(msg->buffer), &(msg->buflen) );
}

int 
sol_msg_id(solClient_opaqueMsg_pt msg_p, message_event* msg)
{
    return solClient_msg_getMsgId( msg_p, &(msg->id) );
}

int 
sol_msg_req_id(solClient_opaqueMsg_pt msg_p)
{
    solClient_uint64_t reqid = 0;
    int result = -1;
    solClient_returnCode_t rc = SOLCLIENT_OK;
    if ( (rc = solClient_msg_getCacheRequestId(msg_p, &reqid))
        	== SOLCLIENT_OK) {
        result = (int) reqid;
    }
    return result;
}

void
populate_common_fields(sol_state* state, message_event* recvmsg, solClient_opaqueMsg_pt msg_p)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;
    recvmsg->discard_flag     = sol_msg_discard_flag( msg_p );
    recvmsg->redelivered_flag = sol_msg_redelivered_flag( msg_p );
    recvmsg->req_id           = sol_msg_req_id( msg_p );
    recvmsg->user_data        = state->user_data_;

    if ( (rc = (solClient_returnCode_t) sol_msg_payload(msg_p, recvmsg))
        	!= SOLCLIENT_OK) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_getBinaryAttachmentPtr()" );
    }
    
    if ( (rc = (solClient_returnCode_t) sol_msg_destination( msg_p, 
        		&(state->recvdest_), recvmsg))
        	!= SOLCLIENT_OK) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_getDestination()" );
    }
}

solClient_rxMsgCallback_returnCode_t 
on_msg_cb(solClient_opaqueSession_pt sess_p, 
          solClient_opaqueMsg_pt msg_p, 
          void *user_p)
{
    sol_state* state = (sol_state*) user_p;
    message_event* recvmsg    = &(state->recvmsg_);

    recvmsg->flow = 0;
    recvmsg->id = 0;
    populate_common_fields( state, recvmsg, msg_p );

    if (state->msg_cb_ != 0) {
#ifdef PYTHON_SUPPORT
        PyGILState_STATE gstate = PyGILState_Ensure();
#endif
        (state->msg_cb_)( (SOLHANDLE)state, recvmsg );
#ifdef PYTHON_SUPPORT
        PyGILState_Release(gstate);
#endif
    }

    return SOLCLIENT_CALLBACK_OK;
}

solClient_rxMsgCallback_returnCode_t
on_flow_msg_cb(solClient_opaqueFlow_pt opaqueFlow_p, solClient_opaqueMsg_pt msg_p, void *user_p)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;

    sol_state* state = (sol_state*) user_p;
    message_event* recvmsg = &(state->recvmsg_);
    
    recvmsg->flow = (FLOWHANDLE)opaqueFlow_p;
    if ( (rc = (solClient_returnCode_t) sol_msg_id(msg_p, recvmsg))
            != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_getMsgId()" );
    }
    populate_common_fields( state, recvmsg, msg_p );

    if (state->msg_cb_ != 0) {
#ifdef PYTHON_SUPPORT
        PyGILState_STATE gstate = PyGILState_Ensure();
#endif
        (state->msg_cb_)( (SOLHANDLE)state, recvmsg );
#ifdef PYTHON_SUPPORT
        PyGILState_Release(gstate);
#endif
    }

    return SOLCLIENT_CALLBACK_OK;
}

