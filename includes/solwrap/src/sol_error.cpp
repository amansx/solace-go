#include "sol_error.h"
#include "sol_state.h"
#include "solclient/solClient.h"
#ifdef PYTHON_SUPPORT
#include <Python.h>
#endif

void 
on_error(SOLHANDLE handle, int rc, const char *fn_name)
{
    sol_state* state = (sol_state*) handle;

    solClient_errorInfo_pt info = solClient_getLastErrorInfo();

    error_event* err = &(state->error_);
    err->fn_name     = fn_name;
    err->err_str     = info->errorStr;
    err->return_code = rc;
    err->rc_str      = solClient_returnCodeToString( (solClient_returnCode_t)rc );
    err->sub_code    = info->subCode;
    err->resp_code   = info->responseCode;
    err->sc_str      = solClient_subCodeToString(info->subCode);
    err->user_data   = state->user_data_;
    
     solClient_log(SOLCLIENT_LOG_ERROR,
         "%s: ReturnCode=\"%s\", SubCode=\"%s\", ResponseCode=%d, Info=\"%s\"",
         fn_name, err->rc_str, err->sc_str, err->resp_code, err->err_str);
    solClient_resetLastErrorInfo();

    if (state->err_cb_ != 0) {
#ifdef PYTHON_SUPPORT
        PyGILState_STATE gstate = PyGILState_Ensure();
#endif
        state->err_cb_(handle, err);
#ifdef PYTHON_SUPPORT
        PyGILState_Release(gstate);
#endif
    }
}
