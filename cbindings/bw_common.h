#ifndef BITWARDEN_COMMON
#define BITWARDEN_COMMON

#include <stddef.h>
#include <stdint.h>

typedef int BitwardenResult;
typedef uint64_t Handle;
typedef Handle ContextHandle;
typedef Handle ClientHandle;
typedef Handle SessionHandle;
typedef Handle VaultHandle;
typedef struct {
    uint8_t bytes[16];
} UUID;

typedef struct {
   UUID *items;
   size_t len;
} UUIDSlice;

enum {
    BitwardenSuccess = 0,
    BitwardenError = 1,
};

#endif
