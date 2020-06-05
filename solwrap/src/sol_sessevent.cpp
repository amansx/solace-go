#include "sol_sessevent.h"
#include "sol_state.h"
#include "sol_error.h"
#include "solclient/solClient.h"

#ifdef PYTHON_SUPPORT
#include <Python.h>
#endif

connectivity_event conn;
publisher_event pub;

void 
on_event_cb(solClient_opaqueSession_pt sess_p,
            solClient_session_eventCallbackInfo_pt eventInfo_p, 
            void *user_p) 
{
    sol_state* state = (sol_state*) user_p;

    switch(eventInfo_p->sessionEvent) {
        // connectivity events
        case SOLCLIENT_SESSION_EVENT_UP_NOTICE:
            conn.type = UP;
            conn.user_data = state->user_data_;
            if (state->conn_cb_ != 0) {
#ifdef PYTHON_SUPPORT
                PyGILState_STATE gstate = PyGILState_Ensure();
#endif
                state->conn_cb_( (SOLHANDLE)state, &conn );
#ifdef PYTHON_SUPPORT
                PyGILState_Release(gstate);
#endif
            }
            return;
        case SOLCLIENT_SESSION_EVENT_RECONNECTING_NOTICE:
            conn.type = RECONNECTING;
            conn.user_data = state->user_data_;
            if (state->conn_cb_ != 0) {
#ifdef PYTHON_SUPPORT
                PyGILState_STATE gstate = PyGILState_Ensure();
#endif
                state->conn_cb_( (SOLHANDLE)state, &conn );
#ifdef PYTHON_SUPPORT
                PyGILState_Release(gstate);
#endif
            }
            return;
        case SOLCLIENT_SESSION_EVENT_RECONNECTED_NOTICE:
            conn.type = RECONNECTED;
            conn.user_data = state->user_data_;
            if (state->conn_cb_ != 0) {
#ifdef PYTHON_SUPPORT
                PyGILState_STATE gstate = PyGILState_Ensure();
#endif
                state->conn_cb_( (SOLHANDLE)state, &conn );
#ifdef PYTHON_SUPPORT
                PyGILState_Release(gstate);
#endif
            }
            return;
        // PUB events
        case SOLCLIENT_SESSION_EVENT_ACKNOWLEDGEMENT:
            pub.type = ACK;
            pub.correlation_data = eventInfo_p->correlation_p;
            pub.user_data = state->user_data_;
            if (state->pub_cb_ != 0) {
#ifdef PYTHON_SUPPORT
                PyGILState_STATE gstate = PyGILState_Ensure();
#endif
                state->pub_cb_( (SOLHANDLE)state, &pub );
#ifdef PYTHON_SUPPORT
                PyGILState_Release(gstate);
#endif
            }
            return;
        case SOLCLIENT_SESSION_EVENT_REJECTED_MSG_ERROR:
            pub.type = REJECT;
            pub.correlation_data = eventInfo_p->correlation_p;
            pub.user_data = state->user_data_;
            if (state->pub_cb_ != 0) {
#ifdef PYTHON_SUPPORT
                PyGILState_STATE gstate = PyGILState_Ensure();
#endif
                state->pub_cb_( (SOLHANDLE)state, &pub );
#ifdef PYTHON_SUPPORT
                PyGILState_Release(gstate);
#endif
            }
            return;
        // Happy days
        case SOLCLIENT_SESSION_EVENT_CAN_SEND:
        case SOLCLIENT_SESSION_EVENT_TE_UNSUBSCRIBE_OK:
        case SOLCLIENT_SESSION_EVENT_PROVISION_OK:
        case SOLCLIENT_SESSION_EVENT_SUBSCRIPTION_OK:
            /* Non error events are logged at the INFO level. */
            solClient_log(SOLCLIENT_LOG_INFO,
            "sol_session_event_triage(): %s\n",
                          solClient_session_eventToString(eventInfo_p->sessionEvent));
            return ;

        // Error events
        case SOLCLIENT_SESSION_EVENT_DOWN_ERROR:
        case SOLCLIENT_SESSION_EVENT_CONNECT_FAILED_ERROR:
            conn.type = DOWN;
            if (state->conn_cb_ != 0) {
#ifdef PYTHON_SUPPORT
                PyGILState_STATE gstate = PyGILState_Ensure();
#endif
                state->conn_cb_( (SOLHANDLE)state, &conn );
#ifdef PYTHON_SUPPORT
                PyGILState_Release(gstate);
#endif
            }
        case SOLCLIENT_SESSION_EVENT_SUBSCRIPTION_ERROR:
        case SOLCLIENT_SESSION_EVENT_RX_MSG_TOO_BIG_ERROR:
        case SOLCLIENT_SESSION_EVENT_TE_UNSUBSCRIBE_ERROR:
        case SOLCLIENT_SESSION_EVENT_PROVISION_ERROR:
            {
                on_error( (SOLHANDLE)state, -1, "session_event_cb()" );
                return ;
            }

        default:
            {
                // Unrecognized or deprecated events are output to STDOUT.
                on_error( (SOLHANDLE)state, -1, solClient_session_eventToString(eventInfo_p->sessionEvent) );
                return ;
            }
    }
}

void 
on_flow_event_cb(solClient_opaqueFlow_pt opaqueFlow_p, solClient_flow_eventCallbackInfo_pt eventInfo_p, void *user_p)
{

    switch(eventInfo_p->flowEvent) {
        case SOLCLIENT_FLOW_EVENT_UP_NOTICE:
        case SOLCLIENT_FLOW_EVENT_SESSION_DOWN:
            /* Non error events are logged at the INFO level. */
            solClient_log(SOLCLIENT_LOG_INFO,
            "sol_flow_event_triage(): %s\n",
                          solClient_flow_eventToString(eventInfo_p->flowEvent));
            return ;

        case SOLCLIENT_FLOW_EVENT_DOWN_ERROR:
        case SOLCLIENT_FLOW_EVENT_BIND_FAILED_ERROR:
        case SOLCLIENT_FLOW_EVENT_REJECTED_MSG_ERROR:
            {
            /* Error events are output to STDOUT. */
            solClient_errorInfo_pt errorInfo_p = solClient_getLastErrorInfo();
            printf(
            "sol_flow_event_triage(): %s; subCode %s, responseCode %d, reason %s\n",
                   solClient_flow_eventToString(eventInfo_p->flowEvent),
                   solClient_subCodeToString(errorInfo_p->subCode), 
                   errorInfo_p->responseCode, errorInfo_p->errorStr );
            return ;
            }

        default:
            /* Unrecognized or deprecated events are output to STDOUT. */
            printf("sol_flow_event_triage(): %s.  Unrecognized/deprecated event.\n",
                     solClient_flow_eventToString(eventInfo_p->flowEvent));
            return ;
    }
}

