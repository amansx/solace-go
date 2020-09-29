#ifndef DLLSPEC_H
#define DLLSPEC_H

#if defined(_MSC_VER) || defined(__MINGW32__)
    #if defined(SOL_WRAP_API)
        #define SOL_API __attribute__ ((dllexport))
    #else
        #define SOL_API __attribute__ ((dllimport))
    #endif
#elif defined(_GCC) || defined(__GNUC__)
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


// #if defined _WIN32 || defined __CYGWIN__ || defined __MINGW32__
//     #ifdef BUILDING_DLL
//         #ifdef __GNUC__
//             #define DLL_PUBLIC __attribute__ ((dllexport))
//         #else
//             #define DLL_PUBLIC __declspec(dllexport) // Note: actually gcc seems to also supports this syntax.
//         #endif
//     #else
//         #ifdef __GNUC__
//             #define DLL_PUBLIC __attribute__ ((dllimport))
//         #else
//             #define DLL_PUBLIC __declspec(dllimport) // Note: actually gcc seems to also supports this syntax.
//         #endif
//     #endif
//     #define DLL_LOCAL
// #else
//     #if __GNUC__ >= 4
//         #define DLL_PUBLIC __attribute__ ((visibility ("default")))
//         #define DLL_LOCAL  __attribute__ ((visibility ("hidden")))
//     #else
//         #define DLL_PUBLIC
//         #define DLL_LOCAL
//     #endif
// #endif