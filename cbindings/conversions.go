package main

/*
#include "bw_common.h"
#include "bw_item.h"
#include <stdlib.h>
#include <stdbool.h>
*/
import "C"
import (
	"time"
	"unsafe"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

func parseUUIDFromC(source C.UUID) uuid.UUID {
	var out uuid.UUID
	for i := range out {
		out[i] = byte(source.bytes[i])
	}

	return out
}

func parseUUIDIntoC(source uuid.UUID) C.UUID {
	var out C.UUID
	for i := range source {
		out.bytes[i] = C.uint8_t(source[i])
	}
	return out
}

func cItemPermissionsFromPtr(value *result.ItemPermissions) *C.BitwardenItemPermissions {
	if value == nil {
		return nil
	}

	out := (*C.BitwardenItemPermissions)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemPermissions{}))))
	*out = C.BitwardenItemPermissions{
		canDelete:  C.bool(value.Delete),
		canRestore: C.bool(value.Restore),
	}
	return out
}

func cItemFieldSliceParts(fields []*result.Field) (unsafe.Pointer, C.size_t) {
	if len(fields) == 0 {
		return nil, 0
	}

	items := (*C.BitwardenItemField)(C.malloc(C.size_t(len(fields)) * C.size_t(unsafe.Sizeof(C.BitwardenItemField{}))))
	out := unsafe.Slice(items, len(fields))
	for i, field := range fields {
		cItemFieldIntoC(&out[i], field)
	}

	return unsafe.Pointer(items), C.size_t(len(fields))
}

func cItemFieldIntoC(out *C.BitwardenItemField, field *result.Field) {
	if field == nil {
		clearC(out)
		return
	}

	putCValue(unsafe.Pointer(out), unsafe.Offsetof(out._type), C.BitwardenFieldType(field.Type))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.name), cStringPtr(field.Name))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.value), cStringPtrFromPtr(field.Value))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.linkedId), cIntPtr(field.LinkedID))
}

func cItemLoginFromPtr(login *result.ItemLogin) *C.BitwardenItemLogin {
	if login == nil {
		return nil
	}

	out := (*C.BitwardenItemLogin)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemLogin{}))))
	cItemLoginIntoC(out, login)
	return out
}

func cItemLoginIntoC(out *C.BitwardenItemLogin, login *result.ItemLogin) {
	if login == nil {
		clearC(out)
		return
	}

	var uris C.BitwardenItemLoginUriSlice
	urisItems, urisLen := cItemLoginURISliceParts(login.URIs)
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.uri), cStringPtr(login.URI))
	putCSlice(unsafe.Pointer(out), unsafe.Offsetof(out.uris), unsafe.Offsetof(uris.items), unsafe.Offsetof(uris.len), urisItems, urisLen)
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.username), cStringPtrFromPtr(login.Username))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.password), cStringPtrFromPtr(login.Password))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.passwordRevisionDate), cUnixMillisPtr(login.PasswordRevisionDate))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.totp), cStringPtrFromPtr(login.TOTP))
}

func cItemLoginURISliceParts(uris []*result.ItemLoginURI) (unsafe.Pointer, C.size_t) {
	if len(uris) == 0 {
		return nil, 0
	}

	items := (*C.BitwardenItemLoginUri)(C.malloc(C.size_t(len(uris)) * C.size_t(unsafe.Sizeof(C.BitwardenItemLoginUri{}))))
	out := unsafe.Slice(items, len(uris))
	for i, uri := range uris {
		cItemLoginURIIntoC(&out[i], uri)
	}

	return unsafe.Pointer(items), C.size_t(len(uris))
}

func cItemLoginURIIntoC(out *C.BitwardenItemLoginUri, uri *result.ItemLoginURI) {
	if uri == nil {
		clearC(out)
		return
	}

	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.uri), cStringPtr(uri.URI))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.uriChecksum), cStringPtr(uri.URIChecksum))
	putCValue(unsafe.Pointer(out), unsafe.Offsetof(out.match), C.BitwardenUriMatchType(uri.Match))
}

func cItemCardFromPtr(card *result.ItemCard) *C.BitwardenItemCard {
	if card == nil {
		return nil
	}

	out := (*C.BitwardenItemCard)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemCard{}))))
	cItemCardIntoC(out, card)
	return out
}

func cItemCardIntoC(out *C.BitwardenItemCard, card *result.ItemCard) {
	if card == nil {
		clearC(out)
		return
	}

	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.cardholderName), cStringPtr(card.CardholderName))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.brand), cStringPtr(card.Brand))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.number), cStringPtr(card.Number))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.expirationMonth), cStringPtr(card.ExpirationMonth))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.expirationYear), cStringPtr(card.ExpirationYear))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.code), cStringPtr(card.Code))
}

func cItemSecureNoteFromPtr(secureNote *result.ItemSecureNote) *C.BitwardenItemSecureNote {
	if secureNote == nil {
		return nil
	}

	out := (*C.BitwardenItemSecureNote)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemSecureNote{}))))
	cItemSecureNoteIntoC(out, secureNote)
	return out
}

func cItemSecureNoteIntoC(out *C.BitwardenItemSecureNote, secureNote *result.ItemSecureNote) {
	if secureNote == nil {
		clearC(out)
		return
	}

	putCValue(unsafe.Pointer(out), unsafe.Offsetof(out._type), C.int(secureNote.Type))
}

func cItemIdentityFromPtr(identity *result.ItemIdentity) *C.BitwardenItemIdentity {
	if identity == nil {
		return nil
	}

	out := (*C.BitwardenItemIdentity)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemIdentity{}))))
	cItemIdentityIntoC(out, identity)
	return out
}

func cItemIdentityIntoC(out *C.BitwardenItemIdentity, identity *result.ItemIdentity) {
	if identity == nil {
		clearC(out)
		return
	}

	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.firstName), cStringPtrFromPtr(identity.FirstName))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.middleName), cStringPtrFromPtr(identity.MiddleName))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.lastName), cStringPtrFromPtr(identity.LastName))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.title), cStringPtrFromPtr(identity.Title))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.passportNumber), cStringPtrFromPtr(identity.PassportNumber))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.username), cStringPtrFromPtr(identity.Username))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.email), cStringPtrFromPtr(identity.Email))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.phone), cStringPtrFromPtr(identity.Phone))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.addressLine1), cStringPtrFromPtr(identity.AddressLine1))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.addressLine2), cStringPtrFromPtr(identity.AddressLine2))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.addressLine3), cStringPtrFromPtr(identity.AddressLine3))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.city), cStringPtrFromPtr(identity.City))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.state), cStringPtrFromPtr(identity.State))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.postalCode), cStringPtrFromPtr(identity.PostalCode))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.country), cStringPtrFromPtr(identity.Country))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.ssn), cStringPtrFromPtr(identity.SSN))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.company), cStringPtrFromPtr(identity.Company))
}

func cItemSSHKeyFromPtr(sshKey *result.ItemSSHKey) *C.BitwardenItemSshKey {
	if sshKey == nil {
		return nil
	}

	out := (*C.BitwardenItemSshKey)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemSshKey{}))))
	cItemSSHKeyIntoC(out, sshKey)
	return out
}

func cItemSSHKeyIntoC(out *C.BitwardenItemSshKey, sshKey *result.ItemSSHKey) {
	if sshKey == nil {
		clearC(out)
		return
	}

	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.privateKey), cStringPtr(sshKey.PrivateKey))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.publicKey), cStringPtr(sshKey.PublicKey))
	putCPtr(unsafe.Pointer(out), unsafe.Offsetof(out.keyFingerprint), cStringPtr(sshKey.KeyFingerprint))
}

func cUUIDSliceParts(ids []uuid.UUID) (unsafe.Pointer, C.size_t) {
	if len(ids) == 0 {
		return nil, 0
	}

	items := (*C.UUID)(C.malloc(C.size_t(len(ids)) * C.size_t(unsafe.Sizeof(C.UUID{}))))
	out := unsafe.Slice(items, len(ids))
	for i, id := range ids {
		out[i] = parseUUIDIntoC(id)
	}

	return unsafe.Pointer(items), C.size_t(len(ids))
}

func cStringFromPtr(value *string) *C.char {
	if value == nil {
		return nil
	}
	return C.CString(*value)
}

func cBoolFromPtr(value *bool) *C.bool {
	if value == nil {
		return nil
	}

	out := (*C.bool)(C.malloc(C.size_t(unsafe.Sizeof(C.bool(false)))))
	*out = C.bool(*value)
	return out
}

func cIntFromPtr(value *int) *C.int {
	if value == nil {
		return nil
	}

	out := (*C.int)(C.malloc(C.size_t(unsafe.Sizeof(C.int(0)))))
	*out = C.int(*value)
	return out
}

func cUnixMillis(value time.Time) C.int64_t {
	if value.IsZero() {
		return 0
	}
	return C.int64_t(value.UnixMilli())
}

func cUnixMillisFromPtr(value *time.Time) *C.int64_t {
	if value == nil {
		return nil
	}

	out := (*C.int64_t)(C.malloc(C.size_t(unsafe.Sizeof(C.int64_t(0)))))
	*out = cUnixMillis(*value)
	return out
}

func cUnixMillisPtr(value *time.Time) unsafe.Pointer {
	return unsafe.Pointer(cUnixMillisFromPtr(value))
}
