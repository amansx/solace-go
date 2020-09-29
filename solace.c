#include "sol_api.h"

extern void gosol_on_msg(SOLHANDLE h, message_event* msg);
extern void gosol_on_err(SOLHANDLE h, error_event* err);
extern void gosol_on_pub(SOLHANDLE h, publisher_event* pub);
extern void gosol_on_con(SOLHANDLE h, connectivity_event* con);

static void on_msg(SOLHANDLE h, message_event* msg) {
	gosol_on_msg( h, msg );
}

static void on_err(SOLHANDLE h, error_event* err) {
	gosol_on_err( h, err );
}

static void on_pub(SOLHANDLE h, publisher_event* pub) {
	gosol_on_pub( h, pub );
}

static void on_con(SOLHANDLE h, connectivity_event* con) {
	gosol_on_con( h, con );
}

static SOLHANDLE gosol_init(void* gohandlers) {
	callback_functions cbs;
	cbs.err_cb = on_err;
	cbs.msg_cb = on_msg;
	cbs.pub_cb = on_pub;
	cbs.con_cb = on_con;
	cbs.user_data = gohandlers;
	return sol_init( on_msg, on_err, on_pub, on_con, gohandlers );
}