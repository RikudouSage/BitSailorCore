package main

/*
#include "bw_common.h"
#include "bw_send.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"os"
	"unsafe"

	"go.chrastecky.dev/bitsailor-core/bitwarden/internal/types"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

//export BitwardenCreateSend
func BitwardenCreateSend(
	vault C.VaultHandle,
	ctx C.ContextHandle,
	session C.SessionHandle,
	send *C.BitwardenSend,
	outSend *C.BitwardenSend,
) C.BitwardenResult {
	if send == nil {
		setLastError(nullPointerError("send"))
		return BitwardenError
	}
	vaultGo, ctxGo, sessionGo, err := getCommonVaultHandles(vault, ctx, session)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	sendGo, closeInput, err := bitwardenSendFromC(send)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}
	if closeInput != nil {
		defer closeInput()
	}

	err = vaultGo.CreateSend(ctxGo, sessionGo, sendGo)
	if err != nil {
		setLastError(err)
		return BitwardenError
	}

	if outSend != nil {
		*outSend = bitwardenSendIntoC(sendGo)
	}

	clearLastError()
	return BitwardenSuccess
}

func bitwardenSendFromC(send *C.BitwardenSend) (*result.Send, func(), error) {
	if send == nil {
		return nil, nil, nil
	}

	out := &result.Send{
		ID:             parseUUIDFromC(send.id),
		AccessID:       goStringFromCPtr(send.accessId),
		AuthType:       result.SendAuthType(send.authType),
		Name:           goStringValueFromCPtr(send.name),
		Disabled:       bool(send.disabled),
		RevisionDate:   goTimeFromCUnixMillis(send.revisionDate),
		DeletionDate:   goTimeFromCUnixMillis(send.deletionDate),
		HideEmail:      bool(send.hideEmail),
		Notes:          goStringFromCPtr(send.notes),
		File:           goSendFileFromC(send.file),
		Key:            goStringValueFromCPtr(send.key),
		AccessCount:    uint(send.accessCount),
		Password:       goStringFromCPtr(send.password),
		ExpirationDate: goTimeFromCUnixMillis(send.expirationDate),
		Type:           result.SendType(send._type),
		MaxAccessCount: goUintFromCPtr(send.maxAccessCount),
		Emails:         goStringSliceFromC(send.emails),
		Text:           goSendTextFromC(send.text),
		FileLength:     int(send.fileLength),
	}

	closeInput, err := setSendInputFileFromC(out, send.inputFilePath)
	if err != nil {
		return nil, nil, err
	}

	return out, closeInput, nil
}

func goStringValueFromCPtr(value *C.char) string {
	if value == nil {
		return ""
	}

	return C.GoString(value)
}

func goUintFromCPtr(value *C.uint) *uint {
	if value == nil {
		return nil
	}

	out := uint(*value)
	return &out
}

func goStringSliceFromC(value C.BitwardenStringSlice) types.CSVSlice {
	if value.items == nil || value.len == 0 {
		return nil
	}

	items := unsafe.Slice(value.items, int(value.len))
	out := make(types.CSVSlice, len(items))
	for i, item := range items {
		out[i] = goStringValueFromCPtr(item)
	}

	return out
}

func goSendTextFromC(value *C.BitwardenSendText) *result.SendText {
	if value == nil {
		return nil
	}

	return &result.SendText{
		Text:   goStringValueFromCPtr(value.text),
		Hidden: bool(value.hidden),
	}
}

func goSendFileFromC(value *C.BitwardenSendFile) *result.SendFile {
	if value == nil {
		return nil
	}

	return &result.SendFile{
		ID:       goStringValueFromCPtr(value.id),
		FileName: goStringValueFromCPtr(value.fileName),
		Size:     goStringValueFromCPtr(value.size),
		SizeName: goStringValueFromCPtr(value.sizeName),
	}
}

func setSendInputFileFromC(send *result.Send, inputFilePath *C.char) (func(), error) {
	if inputFilePath == nil {
		if send.Type == result.SendTypeFile {
			return nil, nullPointerError("send.inputFilePath")
		}

		return nil, nil
	}

	filePath := C.GoString(inputFilePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed opening send input file %q: %w", filePath, err)
	}

	if send.FileLength == 0 {
		stat, err := file.Stat()
		if err != nil {
			_ = file.Close()
			return nil, fmt.Errorf("failed statting send input file %q: %w", filePath, err)
		}
		send.FileLength = int(stat.Size())
	}
	send.InputFile = file

	return func() { _ = file.Close() }, nil
}

func bitwardenSendIntoC(send *result.Send) C.BitwardenSend {
	if send == nil {
		return C.BitwardenSend{}
	}

	return C.BitwardenSend{
		id:             parseUUIDIntoC(send.ID),
		accessId:       cStringFromPtr(send.AccessID),
		authType:       C.BitwardenSendAuthType(send.AuthType),
		name:           C.CString(send.Name),
		disabled:       C.bool(send.Disabled),
		revisionDate:   cUnixMillis(send.RevisionDate),
		deletionDate:   cUnixMillis(send.DeletionDate),
		hideEmail:      C.bool(send.HideEmail),
		notes:          cStringFromPtr(send.Notes),
		file:           cSendFileFromPtr(send.File),
		key:            C.CString(send.Key),
		accessCount:    C.uint(send.AccessCount),
		password:       cStringFromPtr(send.Password),
		expirationDate: cUnixMillis(send.ExpirationDate),
		_type:          C.BitwardenSendType(send.Type),
		maxAccessCount: cUintFromPtr(send.MaxAccessCount),
		emails:         cStringSlice([]string(send.Emails)),
		text:           cSendTextFromPtr(send.Text),
		fileLength:     C.int(send.FileLength),
	}
}

func cUintFromPtr(value *uint) *C.uint {
	if value == nil {
		return nil
	}

	out := (*C.uint)(C.malloc(C.size_t(unsafe.Sizeof(C.uint(0)))))
	*out = C.uint(*value)
	return out
}

func cStringSlice(values []string) C.BitwardenStringSlice {
	if len(values) == 0 {
		return C.BitwardenStringSlice{}
	}

	items := (**C.char)(C.malloc(C.size_t(len(values)) * C.size_t(unsafe.Sizeof((*C.char)(nil)))))
	out := unsafe.Slice(items, len(values))
	for i, value := range values {
		out[i] = C.CString(value)
	}

	return C.BitwardenStringSlice{items: items, len: C.size_t(len(values))}
}

func cSendTextFromPtr(value *result.SendText) *C.BitwardenSendText {
	if value == nil {
		return nil
	}

	out := (*C.BitwardenSendText)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenSendText{}))))
	*out = C.BitwardenSendText{
		text:   C.CString(value.Text),
		hidden: C.bool(value.Hidden),
	}
	return out
}

func cSendFileFromPtr(value *result.SendFile) *C.BitwardenSendFile {
	if value == nil {
		return nil
	}

	out := (*C.BitwardenSendFile)(C.malloc(C.size_t(unsafe.Sizeof(C.BitwardenSendFile{}))))
	*out = C.BitwardenSendFile{
		id:       C.CString(value.ID),
		fileName: C.CString(value.FileName),
		size:     C.CString(value.Size),
		sizeName: C.CString(value.SizeName),
	}
	return out
}

//export BitwardenFreeSend
func BitwardenFreeSend(send *C.BitwardenSend) {
	freeBitwardenSend(send)
}

func freeBitwardenSend(send *C.BitwardenSend) {
	if send == nil {
		return
	}

	C.free(unsafe.Pointer(send.accessId))
	C.free(unsafe.Pointer(send.name))
	C.free(unsafe.Pointer(send.notes))
	freeSendFile(send.file)
	C.free(unsafe.Pointer(send.key))
	C.free(unsafe.Pointer(send.password))
	C.free(unsafe.Pointer(send.maxAccessCount))
	freeStringSlice(send.emails)
	freeSendText(send.text)
	C.free(unsafe.Pointer(send.inputFilePath))

	*send = C.BitwardenSend{}
}

func freeStringSlice(value C.BitwardenStringSlice) {
	items := unsafe.Slice(value.items, int(value.len))
	for i := range items {
		C.free(unsafe.Pointer(items[i]))
	}
	C.free(unsafe.Pointer(value.items))
}

func freeSendText(value *C.BitwardenSendText) {
	if value == nil {
		return
	}

	C.free(unsafe.Pointer(value.text))
	C.free(unsafe.Pointer(value))
}

func freeSendFile(value *C.BitwardenSendFile) {
	if value == nil {
		return
	}

	C.free(unsafe.Pointer(value.id))
	C.free(unsafe.Pointer(value.fileName))
	C.free(unsafe.Pointer(value.size))
	C.free(unsafe.Pointer(value.sizeName))
	C.free(unsafe.Pointer(value))
}
