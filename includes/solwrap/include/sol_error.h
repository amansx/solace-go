#ifndef SOL_ERROR_H
#define SOL_ERROR_H
#if defined(_MSC_VER)
/* Windows-only includes */
#include <winsock2.h>
#else
/* Unix-only includes */
#include <unistd.h>
#endif

#include "sol_data.h"

void 
on_error(SOLHANDLE state, int rc, const char *fn_name);

#endif
