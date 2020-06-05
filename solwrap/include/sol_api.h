#ifndef SOL_API_H
#define SOL_API_H
#include "sol_data.h"
#include "dllspec.h"

#ifdef __cplusplus
extern "C" {
#endif


    SOL_API SOLHANDLE sol_init(message_cb msg_cb, error_cb err_cb, pubevent_cb pub_cb, connectivity_cb con_cb, void* user_data);

    SOL_API int sol_test_cbs(message_cb msg_cb, error_cb err_cb, pubevent_cb pub_cb, connectivity_cb con_cb, void* user_data);

    SOL_API int sol_connect(SOLHANDLE handle, const char* propsfile);
    SOL_API int sol_disconnect(SOLHANDLE handle);

    SOL_API int sol_send_direct(SOLHANDLE handle, const char* topic, void* buffer, int buflen);
    SOL_API int sol_send_persistent(SOLHANDLE handle, const char* destination, enum dest_type type, void* buffer, int buflen, void* correlation_p, int corrlen);

    SOL_API int sol_subscribe_topic(SOLHANDLE handle, const char* topic);
    SOL_API int sol_unsubscribe_topic(SOLHANDLE handle, const char* topic);

    SOL_API int sol_bind_queue(SOLHANDLE handle, const char* queue, enum fwd_mode fm, enum ack_mode am);
    SOL_API int sol_unbind_queue(SOLHANDLE handle, const char* queue);

    SOL_API int sol_ack_msg(SOLHANDLE handle, FLOWHANDLE flow, SOLMSGID msg_id);

    SOL_API int sol_cache_req(SOLHANDLE handle, const char* cache_name, const char* topic_sub, int request_id);

#ifdef __cplusplus
}
#endif


#endif
