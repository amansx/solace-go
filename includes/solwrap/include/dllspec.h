#ifndef DLLSPEC_H
#define DLLSPEC_H

#if defined(_MSC_VER)
    #if defined(SOL_WRAP_API)
        #define SOL_API _declspec(dllexport)
    #else
        #define SOL_API _declspec(dllimport)
    #endif
#elif defined(_GCC)
    #if defined(SOL_WRAP_API)
        #define SOL_API __attribute__((visibility("default")))
    #else
        #define SOL_API 
    #endif
#else
    #define SOL_API 
    #pragma warning Unknown dynamic link import/export semantics.
#endif

#endif
