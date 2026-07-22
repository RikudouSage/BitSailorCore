#ifndef BITWARDEN_ITEM
#define BITWARDEN_ITEM

#include <stdbool.h>
#include <stddef.h>
#include <bw_common.h>

typedef enum {
    BitwardenItemTypeLogin = 1,
    BitwardenItemTypeSecureNote,
    BitwardenItemTypeCard,
    BitwardenItemTypeIdentity,
    BitwardenItemTypeSshKey,
    BitwardenItemTypeBankAccount,
    BitwardenItemTypeDriversLicense,
    BitwardenItemTypePassport,
} BitwardenItemType;

typedef enum {
    BitwardenFieldTypeText = 0,
    BitwardenFieldTypeHidden,
    BitwardenFieldTypeCheckbox,
    BitwardenFieldTypeLinkedId,
} BitwardenFieldType;

typedef enum {
    BitwardenUriMatchTypeDomain = 0,
    BitwardenUriMatchTypeHost,
    BitwardenUriMatchTypeStartsWith,
    BitwardenUriMatchTypeExact,
    BitwardenUriMatchTypeRegularExpression,
    BitwardenUriMatchTypeNever,
} BitwardenUriMatchType;

typedef struct {
    bool canDelete;
    bool canRestore;
} BitwardenItemPermissions;

typedef struct {
    BitwardenFieldType type;
    char* name;
    char* value;
    int* linkedId;
} BitwardenItemField;

typedef struct {
    BitwardenItemField* items;
    size_t len;
} BitwardenItemFieldSlice;

typedef struct {
    char* uri;
    char* uriChecksum;
    BitwardenUriMatchType match;
} BitwardenItemLoginUri;

typedef struct {
    BitwardenItemLoginUri* items;
    size_t len;
} BitwardenItemLoginUriSlice;

typedef struct {
    char* uri;
    BitwardenItemLoginUriSlice uris;
    char* username;
    char* password;
    int64_t* passwordRevisionDate;
    char* totp;
} BitwardenItemLogin;

typedef struct {
    char* cardholderName;
    char* brand;
    char* number;
    char* expirationMonth;
    char* expirationYear;
    char* code;
} BitwardenItemCard;

typedef struct {
    int type;
} BitwardenItemSecureNote;

typedef struct {
    char* firstName;
    char* middleName;
    char* lastName;
    char* title;
    char* passportNumber;
    char* username;
    char* email;
    char* phone;
    char* addressLine1;
    char* addressLine2;
    char* addressLine3;
    char* city;
    char* state;
    char* postalCode;
    char* country;
    char* ssn;
    char* company;
} BitwardenItemIdentity;

typedef struct {
    char* privateKey;
    char* publicKey;
    char* keyFingerprint;
} BitwardenItemSshKey;

typedef struct {
    UUID id;
    BitwardenItemType type;
    char* notes;
    bool* organizationUseTotp;
    int64_t revisionDate;
    int64_t* deletedDate;
    bool favorite;
    UUID organizationId;
    char* key;
    bool edit;
    BitwardenItemPermissions* permissions;
    UUIDSlice collectionIds;
    int64_t* archivedDate;
    UUID folderId;
    bool viewPassword;
    char* name;
    int64_t creationDate;
    bool reprompt;
    BitwardenItemFieldSlice fields;

    BitwardenItemLogin* login;
    BitwardenItemCard* card;
    BitwardenItemSecureNote* secureNote;
    BitwardenItemIdentity* identity;
    BitwardenItemSshKey* sshKey;
} BitwardenItem;

typedef struct {
    BitwardenItem* items;
    size_t len;
} BitwardenItemSlice;

#endif
