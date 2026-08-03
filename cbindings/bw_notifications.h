#ifndef BITWARDEN_NOTIFICATIONS
#define BITWARDEN_NOTIFICATIONS

#include "bw_common.h"

typedef enum {
    BitwardenNotificationStopped = 0,
    BitwardenNotificationConnecting = 1,
    BitwardenNotificationConnected = 2,
    BitwardenNotificationReconnecting = 3,
    BitwardenNotificationFailed = 4,
} NotificationState;

typedef enum {
    BitwardenNotificationSyncCipherUpdate = 0,
    BitwardenNotificationSyncCipherCreate = 1,
    BitwardenNotificationSyncLoginDelete = 2,
    BitwardenNotificationSyncFolderDelete = 3,
    BitwardenNotificationSyncCiphers = 4,
    BitwardenNotificationSyncVault = 5,
    BitwardenNotificationSyncOrgKeys = 6,
    BitwardenNotificationSyncFolderCreate = 7,
    BitwardenNotificationSyncFolderUpdate = 8,
    BitwardenNotificationSyncCipherDelete = 9,
    BitwardenNotificationSyncSettings = 10,
    BitwardenNotificationLogOut = 11,
    BitwardenNotificationSyncSendCreate = 12,
    BitwardenNotificationSyncSendUpdate = 13,
    BitwardenNotificationSyncSendDelete = 14,
    BitwardenNotificationAuthRequest = 15,
    BitwardenNotificationAuthRequestResponse = 16,
    BitwardenNotificationSyncOrganizations = 17,
    BitwardenNotificationSyncOrganizationStatusChanged = 18,
    BitwardenNotificationSyncOrganizationCollectionSettingChanged = 19,
    BitwardenNotificationNotification = 20,
    BitwardenNotificationNotificationStatus = 21,
    BitwardenNotificationRefreshSecurityTasks = 22,
    BitwardenNotificationOrganizationBankAccountVerified = 23,
    BitwardenNotificationProviderBankAccountVerified = 24,
    BitwardenNotificationSyncPolicy = 25,
    BitwardenNotificationAutoConfirmMember = 26,
    BitwardenNotificationPremiumStatusChanged = 27,
} NotificationType;

typedef struct {
    const char* contextId;
    NotificationType notificationType;
    const uint8_t* payload;
    size_t payloadLen;
} BitwardenNotification;

typedef BitwardenResult (*BitwardenNotificationCallback)(void* userData, const BitwardenNotification* notification);

static inline BitwardenResult bitwarden_call_notification_callback(
    BitwardenNotificationCallback callback,
    void* userData,
    const BitwardenNotification* notification
) {
    return callback(userData, notification);
}

#endif
