#ifndef SOL_PROPS_H
#define SOL_PROPS_H

const char** read_props(const char* propsfile);
const char** read_prop_params(const char* host, const char* vpn, const char* user, const char* pass, const char* windowsize);

#endif
