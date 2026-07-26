package main

/*
#include "bw_item.h"
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
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

	vaultGo, ctxGo, sessionGo, err := getCommonVaultHandles(vault, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	itemIDGo := parseUUIDFromC(itemID)

	item, err := vaultGo.GetItem(ctxGo, sessionGo, itemIDGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	bitwardenItemIntoC(outItem, item)

	clearLastError()
	return BitwardenSuccess
}

//export BitwardenFreeItem
func BitwardenFreeItem(item *C.BitwardenItem) {
	freeBitwardenItem(item)
}

func bitwardenItemIntoC(out *C.BitwardenItem, item *result.Item) {
	if item == nil {
		clearC(out)
		return
	}

	base := unsafe.Pointer(out)
	var collectionIds C.UUIDSlice
	var fields C.BitwardenItemFieldSlice
	collectionItems, collectionLen := cUUIDSliceParts(item.CollectionIDs)
	fieldItems, fieldLen := cItemFieldSliceParts(item.Fields)

	putCValue(base, unsafe.Offsetof(out.id), parseUUIDIntoC(item.ID))
	putCValue(base, unsafe.Offsetof(out._type), C.BitwardenItemType(item.Type))
	putCPtr(base, unsafe.Offsetof(out.notes), cStringPtrFromPtr(item.Notes))
	putCPtr(base, unsafe.Offsetof(out.organizationUseTotp), cBoolPtr(item.OrganizationUseTOTP))
	putCValue(base, unsafe.Offsetof(out.revisionDate), cUnixMillis(item.RevisionDate))
	putCPtr(base, unsafe.Offsetof(out.deletedDate), cUnixMillisPtr(item.DeletedDate))
	putCValue(base, unsafe.Offsetof(out.favorite), C.bool(item.Favorite))
	putCValue(base, unsafe.Offsetof(out.organizationId), parseUUIDIntoC(item.OrganizationID))
	putCPtr(base, unsafe.Offsetof(out.key), cStringPtrFromPtr(item.Key))
	putCValue(base, unsafe.Offsetof(out.edit), C.bool(item.Edit))
	putCPtr(base, unsafe.Offsetof(out.permissions), cPtr(cItemPermissionsFromPtr(item.Permissions)))
	putCSlice(base, unsafe.Offsetof(out.collectionIds), unsafe.Offsetof(collectionIds.items), unsafe.Offsetof(collectionIds.len), collectionItems, collectionLen)
	putCPtr(base, unsafe.Offsetof(out.archivedDate), cUnixMillisPtr(item.ArchivedDate))
	putCValue(base, unsafe.Offsetof(out.folderId), parseUUIDIntoC(item.FolderID))
	putCValue(base, unsafe.Offsetof(out.viewPassword), C.bool(item.ViewPassword))
	putCPtr(base, unsafe.Offsetof(out.name), cStringPtr(item.Name))
	putCValue(base, unsafe.Offsetof(out.creationDate), cUnixMillis(item.CreationDate))
	putCValue(base, unsafe.Offsetof(out.reprompt), C.bool(item.Reprompt))
	putCSlice(base, unsafe.Offsetof(out.fields), unsafe.Offsetof(fields.items), unsafe.Offsetof(fields.len), fieldItems, fieldLen)
	putCPtr(base, unsafe.Offsetof(out.login), cPtr(cItemLoginFromPtr(item.Login)))
	putCPtr(base, unsafe.Offsetof(out.card), cPtr(cItemCardFromPtr(item.Card)))
	putCPtr(base, unsafe.Offsetof(out.secureNote), cPtr(cItemSecureNoteFromPtr(item.SecureNote)))
	putCPtr(base, unsafe.Offsetof(out.identity), cPtr(cItemIdentityFromPtr(item.Identity)))
	putCPtr(base, unsafe.Offsetof(out.sshKey), cPtr(cItemSSHKeyFromPtr(item.SSHKey)))
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

	clearC(item)
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
