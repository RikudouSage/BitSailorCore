package main

/*
#include "bw_common.h"
#include "bw_item.h"
*/
import "C"
import (
	"time"
	"unsafe"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/types"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

//export BitwardenCreateItem
func BitwardenCreateItem(
	vault C.VaultHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	item *C.BitwardenItem,
	outItem *C.BitwardenItem,
) C.BitwardenResult {
	if item == nil {
		setLastError(nullPointerError("item"))
		return BitwardenError
	}
	vaultGo, ctxGo, sessionGo, err := getCommonVaultHandles(vault, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	itemGo := bitwardenItemFromC(item)
	err = vaultGo.CreateItem(ctxGo, sessionGo, itemGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if outItem != nil {
		*outItem = bitwardenItemIntoC(itemGo)
	}

	clearLastError()
	return BitwardenSuccess
}

func bitwardenItemFromC(item *C.BitwardenItem) *result.Item {
	if item == nil {
		return nil
	}

	return &result.Item{
		ID:                  parseUUIDFromC(item.id),
		Type:                result.ItemType(item._type),
		Notes:               goStringFromCPtr(item.notes),
		OrganizationUseTOTP: goBoolFromCPtr(item.organizationUseTotp),
		RevisionDate:        goTimeFromCUnixMillis(item.revisionDate),
		DeletedDate:         goTimeFromCUnixMillisPtr(item.deletedDate),
		Favorite:            bool(item.favorite),
		OrganizationID:      parseUUIDFromC(item.organizationId),
		Key:                 goStringFromCPtr(item.key),
		Edit:                bool(item.edit),
		Permissions:         goItemPermissionsFromC(item.permissions),
		CollectionIDs:       goUUIDSliceFromC(item.collectionIds),
		ArchivedDate:        goTimeFromCUnixMillisPtr(item.archivedDate),
		FolderID:            parseUUIDFromC(item.folderId),
		ViewPassword:        bool(item.viewPassword),
		Name:                C.GoString(item.name),
		CreationDate:        goTimeFromCUnixMillis(item.creationDate),
		Reprompt:            types.NumBool(item.reprompt),
		Fields:              goItemFieldsFromC(item.fields),
		Login:               goItemLoginFromC(item.login),
		Card:                goItemCardFromC(item.card),
		SecureNote:          goItemSecureNoteFromC(item.secureNote),
		Identity:            goItemIdentityFromC(item.identity),
		SSHKey:              goItemSSHKeyFromC(item.sshKey),
	}
}

func goStringFromCPtr(value *C.char) *string {
	if value == nil {
		return nil
	}

	out := C.GoString(value)
	return &out
}

func goBoolFromCPtr(value *C.bool) *bool {
	if value == nil {
		return nil
	}

	return new(bool(*value))
}

func goIntFromCPtr(value *C.int) *int {
	if value == nil {
		return nil
	}

	out := int(*value)
	return &out
}

func goTimeFromCUnixMillis(value C.int64_t) time.Time {
	if value == 0 {
		return time.Time{}
	}

	return time.UnixMilli(int64(value))
}

func goTimeFromCUnixMillisPtr(value *C.int64_t) *time.Time {
	if value == nil {
		return nil
	}

	out := goTimeFromCUnixMillis(*value)
	return &out
}

func goItemPermissionsFromC(value *C.BitwardenItemPermissions) *result.ItemPermissions {
	if value == nil {
		return nil
	}

	return &result.ItemPermissions{
		Delete:  bool(value.canDelete),
		Restore: bool(value.canRestore),
	}
}

func goUUIDSliceFromC(value C.UUIDSlice) []uuid.UUID {
	if value.items == nil || value.len == 0 {
		return nil
	}

	items := unsafe.Slice(value.items, int(value.len))
	out := make([]uuid.UUID, len(items))
	for i, item := range items {
		out[i] = parseUUIDFromC(item)
	}

	return out
}

func goItemFieldsFromC(value C.BitwardenItemFieldSlice) []*result.Field {
	if value.items == nil || value.len == 0 {
		return nil
	}

	items := unsafe.Slice(value.items, int(value.len))
	out := make([]*result.Field, len(items))
	for i := range items {
		out[i] = &result.Field{
			Type:     result.FieldType(items[i]._type),
			Name:     C.GoString(items[i].name),
			Value:    goStringFromCPtr(items[i].value),
			LinkedID: goIntFromCPtr(items[i].linkedId),
		}
	}

	return out
}

func goItemLoginFromC(value *C.BitwardenItemLogin) *result.ItemLogin {
	if value == nil {
		return nil
	}

	return &result.ItemLogin{
		URI:                  C.GoString(value.uri),
		URIs:                 goItemLoginURIsFromC(value.uris),
		Username:             goStringFromCPtr(value.username),
		Password:             goStringFromCPtr(value.password),
		PasswordRevisionDate: goTimeFromCUnixMillisPtr(value.passwordRevisionDate),
		TOTP:                 goStringFromCPtr(value.totp),
	}
}

func goItemLoginURIsFromC(value C.BitwardenItemLoginUriSlice) []*result.ItemLoginURI {
	if value.items == nil || value.len == 0 {
		return nil
	}

	items := unsafe.Slice(value.items, int(value.len))
	out := make([]*result.ItemLoginURI, len(items))
	for i := range items {
		out[i] = &result.ItemLoginURI{
			URI:         C.GoString(items[i].uri),
			URIChecksum: C.GoString(items[i].uriChecksum),
			Match:       result.URIMatchType(items[i].match),
		}
	}

	return out
}

func goItemCardFromC(value *C.BitwardenItemCard) *result.ItemCard {
	if value == nil {
		return nil
	}

	return &result.ItemCard{
		CardholderName:  C.GoString(value.cardholderName),
		Brand:           C.GoString(value.brand),
		Number:          C.GoString(value.number),
		ExpirationMonth: C.GoString(value.expirationMonth),
		ExpirationYear:  C.GoString(value.expirationYear),
		Code:            C.GoString(value.code),
	}
}

func goItemSecureNoteFromC(value *C.BitwardenItemSecureNote) *result.ItemSecureNote {
	if value == nil {
		return nil
	}

	return &result.ItemSecureNote{Type: int(value._type)}
}

func goItemIdentityFromC(value *C.BitwardenItemIdentity) *result.ItemIdentity {
	if value == nil {
		return nil
	}

	return &result.ItemIdentity{
		FirstName:      goStringFromCPtr(value.firstName),
		MiddleName:     goStringFromCPtr(value.middleName),
		LastName:       goStringFromCPtr(value.lastName),
		Title:          goStringFromCPtr(value.title),
		PassportNumber: goStringFromCPtr(value.passportNumber),
		Username:       goStringFromCPtr(value.username),
		Email:          goStringFromCPtr(value.email),
		Phone:          goStringFromCPtr(value.phone),
		AddressLine1:   goStringFromCPtr(value.addressLine1),
		AddressLine2:   goStringFromCPtr(value.addressLine2),
		AddressLine3:   goStringFromCPtr(value.addressLine3),
		City:           goStringFromCPtr(value.city),
		State:          goStringFromCPtr(value.state),
		PostalCode:     goStringFromCPtr(value.postalCode),
		Country:        goStringFromCPtr(value.country),
		SSN:            goStringFromCPtr(value.ssn),
		Company:        goStringFromCPtr(value.company),
	}
}

func goItemSSHKeyFromC(value *C.BitwardenItemSshKey) *result.ItemSSHKey {
	if value == nil {
		return nil
	}

	return &result.ItemSSHKey{
		PrivateKey:     C.GoString(value.privateKey),
		PublicKey:      C.GoString(value.publicKey),
		KeyFingerprint: C.GoString(value.keyFingerprint),
	}
}
