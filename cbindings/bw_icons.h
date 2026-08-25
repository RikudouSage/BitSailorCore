#ifndef BITWARDEN_ICONS
#define BITWARDEN_ICONS

#include "bw_common.h"

typedef struct {
    char* hostname;
    BitwardenByteSlice icon;
} BitwardenIcon;

typedef struct {
    BitwardenIcon* items;
    size_t len;
} BitwardenIconSlice;

#endif
