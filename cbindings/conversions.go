package main

/*
#include <bw_common.h>
#include <bw_item.h>
#include <stdlib.h>
#include <stdbool.h>
*/
import "C"
import (
	"time"
	"unsafe"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
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

func cItemFieldSlice(fields []*result.Field) C.BitwardenItemFieldSlice {
	if len(fields) == 0 {
		return C.BitwardenItemFieldSlice{}
	}

	items := (*C.BitwardenItemField)(C.malloc(C.size_t(len(fields)) * C.size_t(unsafe.Sizeof(C.BitwardenItemField{}))))
	out := unsafe.Slice(items, len(fields))
	for i, field := range fields {
		out[i] = cItemFieldFromPtr(field)
	}

	return C.BitwardenItemFieldSlice{items: items, len: C.size_t(len(fields))}
}

func cItemFieldFromPtr(field *result.Field) C.BitwardenItemField {
	if field == nil {
		return C.BitwardenItemField{}
	}

	return C.BitwardenItemField{
		_type:    C.BitwardenFieldType(field.Type),
		name:     C.CString(field.Name),
		value:    cStringFromPtr(field.Value),
		linkedId: cIntFromPtr(field.LinkedID),
	}
}

func cItemLoginFromPtr(login *result.ItemLogin) *C.BitwardenItemLogin {
	if login == nil {
		return nil
	}

	out := (*C.BitwardenItemLogin)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemLogin{}))))
	*out = C.BitwardenItemLogin{
		uri:                  C.CString(login.URI),
		uris:                 cItemLoginURISlice(login.URIs),
		username:             cStringFromPtr(login.Username),
		password:             cStringFromPtr(login.Password),
		passwordRevisionDate: cUnixMillisFromPtr(login.PasswordRevisionDate),
		totp:                 cStringFromPtr(login.TOTP),
	}
	return out
}

func cItemLoginURISlice(uris []*result.ItemLoginURI) C.BitwardenItemLoginUriSlice {
	if len(uris) == 0 {
		return C.BitwardenItemLoginUriSlice{}
	}

	items := (*C.BitwardenItemLoginUri)(C.malloc(C.size_t(len(uris)) * C.size_t(unsafe.Sizeof(C.BitwardenItemLoginUri{}))))
	out := unsafe.Slice(items, len(uris))
	for i, uri := range uris {
		out[i] = cItemLoginURIFromPtr(uri)
	}

	return C.BitwardenItemLoginUriSlice{items: items, len: C.size_t(len(uris))}
}

func cItemLoginURIFromPtr(uri *result.ItemLoginURI) C.BitwardenItemLoginUri {
	if uri == nil {
		return C.BitwardenItemLoginUri{}
	}

	return C.BitwardenItemLoginUri{
		uri:         C.CString(uri.URI),
		uriChecksum: C.CString(uri.URIChecksum),
		match:       C.BitwardenUriMatchType(uri.Match),
	}
}

func cItemCardFromPtr(card *result.ItemCard) *C.BitwardenItemCard {
	if card == nil {
		return nil
	}

	out := (*C.BitwardenItemCard)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemCard{}))))
	*out = C.BitwardenItemCard{
		cardholderName:  C.CString(card.CardholderName),
		brand:           C.CString(card.Brand),
		number:          C.CString(card.Number),
		expirationMonth: C.CString(card.ExpirationMonth),
		expirationYear:  C.CString(card.ExpirationYear),
		code:            C.CString(card.Code),
	}
	return out
}

func cItemSecureNoteFromPtr(secureNote *result.ItemSecureNote) *C.BitwardenItemSecureNote {
	if secureNote == nil {
		return nil
	}

	out := (*C.BitwardenItemSecureNote)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemSecureNote{}))))
	*out = C.BitwardenItemSecureNote{_type: C.int(secureNote.Type)}
	return out
}

func cItemIdentityFromPtr(identity *result.ItemIdentity) *C.BitwardenItemIdentity {
	if identity == nil {
		return nil
	}

	out := (*C.BitwardenItemIdentity)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemIdentity{}))))
	*out = C.BitwardenItemIdentity{
		firstName:      cStringFromPtr(identity.FirstName),
		middleName:     cStringFromPtr(identity.MiddleName),
		lastName:       cStringFromPtr(identity.LastName),
		title:          cStringFromPtr(identity.Title),
		passportNumber: cStringFromPtr(identity.PassportNumber),
		username:       cStringFromPtr(identity.Username),
		email:          cStringFromPtr(identity.Email),
		phone:          cStringFromPtr(identity.Phone),
		addressLine1:   cStringFromPtr(identity.AddressLine1),
		addressLine2:   cStringFromPtr(identity.AddressLine2),
		addressLine3:   cStringFromPtr(identity.AddressLine3),
		city:           cStringFromPtr(identity.City),
		state:          cStringFromPtr(identity.State),
		postalCode:     cStringFromPtr(identity.PostalCode),
		country:        cStringFromPtr(identity.Country),
		ssn:            cStringFromPtr(identity.SSN),
		company:        cStringFromPtr(identity.Company),
	}
	return out
}

func cItemSSHKeyFromPtr(sshKey *result.ItemSSHKey) *C.BitwardenItemSshKey {
	if sshKey == nil {
		return nil
	}

	out := (*C.BitwardenItemSshKey)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenItemSshKey{}))))
	*out = C.BitwardenItemSshKey{
		privateKey:     C.CString(sshKey.PrivateKey),
		publicKey:      C.CString(sshKey.PublicKey),
		keyFingerprint: C.CString(sshKey.KeyFingerprint),
	}
	return out
}

func cUUIDSlice(ids []uuid.UUID) C.UUIDSlice {
	if len(ids) == 0 {
		return C.UUIDSlice{}
	}

	items := (*C.UUID)(C.malloc(C.size_t(len(ids)) * C.size_t(unsafe.Sizeof(C.UUID{}))))
	out := unsafe.Slice(items, len(ids))
	for i, id := range ids {
		out[i] = parseUUIDIntoC(id)
	}

	return C.UUIDSlice{items: items, len: C.size_t(len(ids))}
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
