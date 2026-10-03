#pragma once

typedef enum {
    BitwardenTfaKindAuthenticator = 0,
    BitwardenTfaKindEmail = 1,
} BitwardenTfaKind;

typedef struct {
    BitwardenTfaKind kind;
    const char *code;
} BitwardenTfaConfig;
