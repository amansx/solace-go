#include <iostream>
#include "json.hpp"
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
sol_msg_destination(solClient_opaqueMsg_pt msg_p, solClient_destination_t* dest, message_event* msg)
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
sol_msg_replyto(solClient_opaqueMsg_pt msg_p, solClient_destination_t* dest, message_event* msg)
{
    solClient_returnCode_t rc = solClient_msg_getReplyTo( msg_p, dest, sizeof(solClient_destination_t) );
    if ( rc == SOLCLIENT_OK ) {
        if (dest->destType == SOLCLIENT_TOPIC_DESTINATION) {
            msg->replyTo = dest->dest;
            msg->replyToDestType = TOPIC;
        }
        else if (dest->destType == SOLCLIENT_QUEUE_DESTINATION) {
            msg->replyTo = dest->dest;
            msg->replyToDestType = QUEUE;
        }
    }
    else {
        msg->replyTo = 0;
        msg->replyToDestType = NONE;
    }
    return rc;
}

// void *vp = static_cast<void*>(new std::string("it's easy to break stuff like this!"));
// int r = solClient_msg_getBinaryAttachmentPtr( msg_p, &(msg->buffer), &(msg->buflen) );
// int r = solClient_msg_getBinaryAttachmentString( msg_p, &(msg->buffer) );
// void *vp;
// unsigned int vpl;

// std::string *sp = static_cast<std::string*>(vp);
// std::string s = *sp;
// delete sp;

// msg->buffer = s.c_str();
// msg->buflen = vpl;
// return r;

// int 
// sol_msg_payload(solClient_opaqueMsg_pt msg_p, message_event* msg) 
// {
//     return solClient_msg_getBinaryAttachmentPtr(msg_p, &(msg->buffer), &(msg->buflen));
// }

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
    if ( (rc = solClient_msg_getCacheRequestId(msg_p, &reqid)) == SOLCLIENT_OK) {
        result = (int) reqid;
    }
    return result;
}

solClient_rxMsgCallback_returnCode_t 
on_msg_cb(solClient_opaqueSession_pt sess_p, solClient_opaqueMsg_pt msg_p, void *user_p)
{
    solClient_returnCode_t rc = SOLCLIENT_OK;

    sol_state* state = (sol_state*) user_p;
    message_event* recvmsg = &(state->recvmsg_);

    recvmsg->flow                 = 0;
    recvmsg->id                   = 0;
    recvmsg->req_id               = sol_msg_req_id( msg_p );
    recvmsg->redelivered_flag     = sol_msg_redelivered_flag( msg_p );
    recvmsg->discard_flag         = sol_msg_discard_flag( msg_p );
    recvmsg->user_data            = state->user_data_;

    solClient_msg_getApplicationMsgType(msg_p, &(recvmsg->application_message_type));
    solClient_msg_getCorrelationId(msg_p, &(recvmsg->correlationid) );
    sol_msg_replyto(msg_p, &(state->replytodest_), recvmsg);

    // ======================================
    // Populate Fields
    // ======================================

    solClient_opaqueContainer_pt ptr;
    if ((rc = solClient_msg_getUserPropertyMap(msg_p, &ptr)) != SOLCLIENT_OK) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_getUserPropertyMap()" );
    }

    nlohmann::json json;
    while (solClient_container_hasNextField(ptr)) {

        solClient_field_t field;
        const char *name_p = NULL;

        if ((rc = solClient_container_getNextField (ptr, &field, sizeof(solClient_field_t), &name_p)) != SOLCLIENT_OK) {
            solClient_log(SOLCLIENT_LOG_CRITICAL, "Bad Field");
            // on_error( (SOLHANDLE)state, rc, "solClient_msg_getUserPropertyMap()" );
            continue;
        }

        switch (field.type) {
            case SOLCLIENT_BOOL:
                // solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is bool", name_p);
                json["bool"][name_p] = field.value.boolean;
                break;
            case SOLCLIENT_UINT8:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is uint8", name_p);
                break;
            case SOLCLIENT_INT8:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_INT8", name_p);
                break;
            case SOLCLIENT_UINT16:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_UINT16", name_p);
                break;
            case SOLCLIENT_INT16:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_INT16", name_p);
                break;
            case SOLCLIENT_UINT32:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_UINT32", name_p);
                break;
            case SOLCLIENT_INT32:
                // solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_INT32", name_p);
                json["int32"][name_p] = field.value.int32;
                break;
            case SOLCLIENT_UINT64:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_UINT64", name_p);
                break;
            case SOLCLIENT_INT64:
                // solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_INT64", name_p);
                json["int64"][name_p] = field.value.int64;
                break;
            case SOLCLIENT_WCHAR:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_WCHAR", name_p);
                break;
            case SOLCLIENT_STRING:
                // solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_STRING", name_p);
                json["string"][name_p] = field.value.string;
                break;
            case SOLCLIENT_BYTEARRAY:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_BYTEARRAY", name_p);
                break;
            case SOLCLIENT_FLOAT:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_FLOAT", name_p);
                break;
            case SOLCLIENT_DOUBLE:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_DOUBLE", name_p);
                break;
            case SOLCLIENT_MAP:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_MAP", name_p);
                break;
            case SOLCLIENT_STREAM:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_STREAM", name_p);
                break;
            case SOLCLIENT_NULL:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_NULL", name_p);
                break;
            case SOLCLIENT_DESTINATION:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_DESTINATION", name_p);
                break;
            case SOLCLIENT_SMF:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_SMF", name_p);
                break;
            case SOLCLIENT_UNKNOWN:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s is SOLCLIENT_UNKNOWN", name_p);
                break;
            default:
                solClient_log(SOLCLIENT_LOG_CRITICAL, "Unmapped %s", name_p);
               break;
        }

    }

    const std::string user_properties_payload = json.dump();
    recvmsg->user_properties = user_properties_payload.c_str();

    // ======================================
    // END Populate Fields
    // ======================================

    recvmsg->buffer = 0;
    recvmsg->buflen = 0;

    if ( (rc = (solClient_returnCode_t) solClient_msg_getBinaryAttachmentPtr(msg_p, &(recvmsg->buffer), &(recvmsg->buflen)) ) != SOLCLIENT_OK ) {
        // on_error( (SOLHANDLE)state, rc, "solClient_msg_getBinaryAttachmentPtr()" );
    }

    if ( (rc = (solClient_returnCode_t) sol_msg_destination ( msg_p, &(state->recvdest_), recvmsg )) != SOLCLIENT_OK) {
        // on_error( (SOLHANDLE)state, rc, "sol_msg_destination()" );
    }



    if (state->msg_cb_ != 0) {
        (state->msg_cb_)( (SOLHANDLE)state, recvmsg );
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
    
    if ( (rc = (solClient_returnCode_t) sol_msg_id(msg_p, recvmsg)) != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_getMsgId()" );
    }
    

    // ======================================
    // Populate Fields
    // ======================================

    solClient_opaqueContainer_pt ptr;
    if ((rc = solClient_msg_getUserPropertyMap(msg_p, &ptr)) != SOLCLIENT_OK) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_getUserPropertyMap()" );
    }

    nlohmann::json json;
    while (solClient_container_hasNextField(ptr)) {
        solClient_field_t field;
        const char *name_p = NULL;

        if ((rc = solClient_container_getNextField (ptr, &field, sizeof(solClient_field_t), &name_p)) != SOLCLIENT_OK) {
            on_error( (SOLHANDLE)state, rc, "solClient_msg_getUserPropertyMap()" );
        }

        if (field.type == SOLCLIENT_BOOL) {
            json["bool"][name_p] = field.value.boolean;
        } else if (field.type == SOLCLIENT_INT64) {
            json["int64"][name_p] = field.value.int64;
        } else if (field.type == SOLCLIENT_STRING) {
            json["string"][name_p] = field.value.string;
        } else {
            on_error( (SOLHANDLE)state, rc, "unknown type" );
        }

    }

    const std::string user_properties_payload = json.dump();

    recvmsg->req_id               = sol_msg_req_id( msg_p );
    recvmsg->redelivered_flag     = sol_msg_redelivered_flag( msg_p );
    recvmsg->discard_flag         = sol_msg_discard_flag( msg_p );

    solClient_msg_getApplicationMsgType( msg_p, &(recvmsg->application_message_type) );
    solClient_msg_getCorrelationId( msg_p, &(recvmsg->correlationid) );
    
    sol_msg_replyto ( msg_p, &(state->replytodest_), recvmsg );

    if ( (rc = (solClient_returnCode_t) solClient_msg_getBinaryAttachmentPtr(msg_p, &(recvmsg->buffer), &(recvmsg->buflen)) ) != SOLCLIENT_OK ) {
        on_error( (SOLHANDLE)state, rc, "solClient_msg_getBinaryAttachmentPtr()" );
    }

    if ( (rc = (solClient_returnCode_t) sol_msg_destination( msg_p, &(state->recvdest_), recvmsg)) != SOLCLIENT_OK) {
        on_error( (SOLHANDLE)state, rc, "sol_msg_destination()" );
    }

    recvmsg->user_properties      = user_properties_payload.c_str();
    recvmsg->user_data            = state->user_data_;

    // ======================================
    // END Populate Fields
    // ======================================



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

