#ifndef BITWARDEN_SEND
#define BITWARDEN_SEND

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include "bw_common.h"

typedef enum {
    BitwardenSendAuthTypeSpecificPeople = 0,
    BitwardenSendAuthTypePassword,
    BitwardenSendAuthTypeNoAuth,
} BitwardenSendAuthType;

typedef enum {
    BitwardenSendTypeText = 0,
    BitwardenSendTypeFile,
} BitwardenSendType;

typedef struct {
    char** items;
    size_t len;
} BitwardenStringSlice;

typedef struct {
    char* text;
    bool hidden;
} BitwardenSendText;

typedef struct {
    char* id;
    char* fileName;
    char* size;
    char* sizeName;
} BitwardenSendFile;

typedef struct {
    UUID id;
    char* accessId;
    BitwardenSendAuthType authType;
    char* name;
    bool disabled;
    int64_t revisionDate;
    int64_t deletionDate;
    bool hideEmail;
    char* notes;
    BitwardenSendFile* file;
    char* key;
    unsigned int accessCount;
    char* password;
    int64_t expirationDate;
    BitwardenSendType type;
    unsigned int* maxAccessCount;
    BitwardenStringSlice emails;
    BitwardenSendText* text;

    int fileLength;
    char* inputFilePath;
} BitwardenSend;

#endif
