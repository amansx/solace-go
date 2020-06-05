#ifndef SOL_DATA_H
#define SOL_DATA_H
#if defined(_MSC_VER)
/* Windows-only includes */
#include <winsock2.h>
#else
/* Unix-only includes */
#include <unistd.h>
#endif

#ifdef __cplusplus
extern "C" {
#endif

    /**
    * Handle to a Solace session instance.
    **/
    // typedef void*   SOLHANDLE;
#if defined(_MSC_VER)
		typedef HANDLE SOLHANDLE;
#else
    typedef unsigned long long SOLHANDLE;
#endif
    /**
     * Handle to a Solace guaranteed-messaging flow instance;
     **/
    typedef unsigned long long FLOWHANDLE;
    typedef unsigned long long SOLMSGID;

    /**
     * Defines the range of possible forwarding modes for persistent messages.
     **/
    enum fwd_mode 
    {
    /**
     * STORE_FWD -- (DEFAULT) Store-and-forward mode of delivery for persistent messages 
     * persists the messages before delivery.
     **/
        STORE_FWD=1, 
    /**
     * CUT_THRU -- Cut-Through mode of delivery for persistent messages uses the direct 
     * data path to deliver messages in parallel while the messages are presisted.
     **/
        CUT_THRU=2 
    };

    /**
     * Defines the range of possible message-acknowledgment modes for persistent messages.
     **/
    enum ack_mode 
    { 
        /**
         * AUTO_ACK -- (DEFAULT) In auto-acknowledgment mode the underlying client-library 
         * before invoking the message-callback
         **/
        AUTO_ACK=1, 
        /**
         * MANUAL_ACK -- In manual-acknowledgment mode the client application is required 
         * to acknowleddge messages directly to the underlying library so that it can 
         * free the resources and acknowledge receipt of the message back to the appliance.
         **/
        MANUAL_ACK=2 
    };

    /**
    * Defines the destination type (topic or queue)
    **/
    enum dest_type
    {
        /**
        * TOPIC -- destination type is a direct topic
        **/
        TOPIC=1,
        /**
        * QUEUE -- destination type is a persistent queue
        **/
        QUEUE=2,
        /**
        * NONE -- destination is not set (this should not happen in Solace-callbacks)
        **/
        NONE =3
    };


    /**
    * Internal message event representation for messages that is non-Solace specific.
    *
    * Contains metadata from the message envelope, as well as the payload of binary messages 
    * (does not support structured messages).
    **/
    typedef struct message_event 
    {
        enum dest_type desttype;
        const char*    destination;
        FLOWHANDLE     flow;
        SOLMSGID       id;
        void*          buffer;
        unsigned int   buflen;
        int            req_id;
        int            redelivered_flag;
        int            discard_flag;
        void*          user_data;
    } message_event;
    /**
    * Message Callback function profile.
    **/
    typedef void (*message_cb)(SOLHANDLE handle, message_event*);




    /**
    * Internal error event representation for session errors. References Solace error codes values.
    **/
    typedef struct error_event 
    {
        const char* fn_name;
        int         return_code; 
        const char* rc_str;
        int         sub_code; 
        const char* sc_str; 
        int         resp_code;
        const char* err_str;
        void*       user_data;
    } error_event;
    /**
    * Error Callback function profile.
    **/
    typedef void (*error_cb)(SOLHANDLE handle, error_event* error);



    enum publisher_event_type { ACK=1, REJECT=2 };
    typedef struct publisher_event 
    {
        enum  publisher_event_type type;
        void* correlation_data;
        void* user_data;
    } publisher_event;
    typedef void (*pubevent_cb)(SOLHANDLE handle, publisher_event* pubevent);



    enum connectivity_event_type { UP=1, RECONNECTING=2, RECONNECTED=3, DOWN=4 };
    typedef struct connectivity_event 
    {
        enum  connectivity_event_type type;
        void* user_data;
    } connectivity_event;
    typedef void (*connectivity_cb)(SOLHANDLE handle, connectivity_event* reconnevent);



    typedef struct callback_functions 
    {
        message_cb      msg_cb;
        error_cb        err_cb;
        pubevent_cb     pub_cb;
        connectivity_cb con_cb;
        void*           user_data;
    } callback_functions;

#ifdef __cplusplus
}
#endif

#endif
