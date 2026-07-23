package main

/*
#include "bw_item.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"go.chrastecky.dev/bitwarden-client/bitwarden"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

//export BitwardenGetItem
func BitwardenGetItem(
	vault C.VaultHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	itemID C.UUID,
	outItem *C.BitwardenItem,
) C.BitwardenResult {
	if outItem == nil {
		setLastError(nullPointerError("outItem"))
		return BitwardenError
	}

	vaultGo, err := getHandleObj[bitwarden.Vault](handle(vault))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	ctxGo, err := getHandleObj[*contextHandle](handle(ctx))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	sessionGo, err := getHandleObj[*result.Session](handle(session))
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	itemIDGo := parseUUIDFromC(itemID)

	item, err := vaultGo.GetItem(ctxGo.ctx, sessionGo, itemIDGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	*outItem = bitwardenItemIntoC(item)

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenFreeItem
func BitwardenFreeItem(item *C.BitwardenItem) {
	freeBitwardenItem(item)
}

func bitwardenItemIntoC(item *result.Item) C.BitwardenItem {
	if item == nil {
		return C.BitwardenItem{}
	}

	return C.BitwardenItem{
		id:                  parseUUIDIntoC(item.ID),
		_type:               C.BitwardenItemType(item.Type),
		notes:               cStringFromPtr(item.Notes),
		organizationUseTotp: cBoolFromPtr(item.OrganizationUseTOTP),
		revisionDate:        cUnixMillis(item.RevisionDate),
		deletedDate:         cUnixMillisFromPtr(item.DeletedDate),
		favorite:            C.bool(item.Favorite),
		organizationId:      parseUUIDIntoC(item.OrganizationID),
		key:                 cStringFromPtr(item.Key),
		edit:                C.bool(item.Edit),
		permissions:         cItemPermissionsFromPtr(item.Permissions),
		collectionIds:       cUUIDSlice(item.CollectionIDs),
		archivedDate:        cUnixMillisFromPtr(item.ArchivedDate),
		folderId:            parseUUIDIntoC(item.FolderID),
		viewPassword:        C.bool(item.ViewPassword),
		name:                C.CString(item.Name),
		creationDate:        cUnixMillis(item.CreationDate),
		reprompt:            C.bool(item.Reprompt),
		fields:              cItemFieldSlice(item.Fields),
		login:               cItemLoginFromPtr(item.Login),
		card:                cItemCardFromPtr(item.Card),
		secureNote:          cItemSecureNoteFromPtr(item.SecureNote),
		identity:            cItemIdentityFromPtr(item.Identity),
		sshKey:              cItemSSHKeyFromPtr(item.SSHKey),
	}
}

func freeBitwardenItem(item *C.BitwardenItem) {
	if item == nil {
		return
	}

	C.free(unsafe.Pointer(item.notes))
	C.free(unsafe.Pointer(item.organizationUseTotp))
	C.free(unsafe.Pointer(item.deletedDate))
	C.free(unsafe.Pointer(item.key))
	C.free(unsafe.Pointer(item.permissions))
	freeUUIDSlice(item.collectionIds)
	C.free(unsafe.Pointer(item.archivedDate))
	C.free(unsafe.Pointer(item.name))
	freeItemFieldSlice(item.fields)
	freeItemLogin(item.login)
	freeItemCard(item.card)
	freeItemSecureNote(item.secureNote)
	freeItemIdentity(item.identity)
	freeItemSSHKey(item.sshKey)

	*item = C.BitwardenItem{}
}

func freeUUIDSlice(value C.UUIDSlice) {
	C.free(unsafe.Pointer(value.items))
}

func freeItemFieldSlice(value C.BitwardenItemFieldSlice) {
	fields := unsafe.Slice(value.items, int(value.len))
	for i := range fields {
		C.free(unsafe.Pointer(fields[i].name))
		C.free(unsafe.Pointer(fields[i].value))
		C.free(unsafe.Pointer(fields[i].linkedId))
	}
	C.free(unsafe.Pointer(value.items))
}

func freeItemLogin(login *C.BitwardenItemLogin) {
	if login == nil {
		return
	}

	C.free(unsafe.Pointer(login.uri))
	freeItemLoginURISlice(login.uris)
	C.free(unsafe.Pointer(login.username))
	C.free(unsafe.Pointer(login.password))
	C.free(unsafe.Pointer(login.passwordRevisionDate))
	C.free(unsafe.Pointer(login.totp))
	C.free(unsafe.Pointer(login))
}

func freeItemLoginURISlice(value C.BitwardenItemLoginUriSlice) {
	uris := unsafe.Slice(value.items, int(value.len))
	for i := range uris {
		C.free(unsafe.Pointer(uris[i].uri))
		C.free(unsafe.Pointer(uris[i].uriChecksum))
	}
	C.free(unsafe.Pointer(value.items))
}

func freeItemCard(card *C.BitwardenItemCard) {
	if card == nil {
		return
	}

	C.free(unsafe.Pointer(card.cardholderName))
	C.free(unsafe.Pointer(card.brand))
	C.free(unsafe.Pointer(card.number))
	C.free(unsafe.Pointer(card.expirationMonth))
	C.free(unsafe.Pointer(card.expirationYear))
	C.free(unsafe.Pointer(card.code))
	C.free(unsafe.Pointer(card))
}

func freeItemSecureNote(secureNote *C.BitwardenItemSecureNote) {
	C.free(unsafe.Pointer(secureNote))
}

func freeItemIdentity(identity *C.BitwardenItemIdentity) {
	if identity == nil {
		return
	}

	C.free(unsafe.Pointer(identity.firstName))
	C.free(unsafe.Pointer(identity.middleName))
	C.free(unsafe.Pointer(identity.lastName))
	C.free(unsafe.Pointer(identity.title))
	C.free(unsafe.Pointer(identity.passportNumber))
	C.free(unsafe.Pointer(identity.username))
	C.free(unsafe.Pointer(identity.email))
	C.free(unsafe.Pointer(identity.phone))
	C.free(unsafe.Pointer(identity.addressLine1))
	C.free(unsafe.Pointer(identity.addressLine2))
	C.free(unsafe.Pointer(identity.addressLine3))
	C.free(unsafe.Pointer(identity.city))
	C.free(unsafe.Pointer(identity.state))
	C.free(unsafe.Pointer(identity.postalCode))
	C.free(unsafe.Pointer(identity.country))
	C.free(unsafe.Pointer(identity.ssn))
	C.free(unsafe.Pointer(identity.company))
	C.free(unsafe.Pointer(identity))
}

func freeItemSSHKey(sshKey *C.BitwardenItemSshKey) {
	if sshKey == nil {
		return
	}

	C.free(unsafe.Pointer(sshKey.privateKey))
	C.free(unsafe.Pointer(sshKey.publicKey))
	C.free(unsafe.Pointer(sshKey.keyFingerprint))
	C.free(unsafe.Pointer(sshKey))
}
