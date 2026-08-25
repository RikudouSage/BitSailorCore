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
typedef Handle NotificationSubscriptionHandle;
typedef struct {
    uint8_t bytes[16];
} UUID;

typedef struct {
   UUID *items;
   size_t len;
} UUIDSlice;

typedef struct {
    char** items;
    size_t len;
} BitwardenStringSlice;

typedef struct {
    uint8_t* items;
    size_t len;
} BitwardenByteSlice;

enum {
    BitwardenSuccess = 0,
    BitwardenError = 1,
};

#endif
