package bitwarden

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/google/uuid"
	clone "github.com/huandu/go-clone/generic"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/dto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func (receiver *vault) CreateItem(ctx context.Context, session *result.Session, item *result.Item) error {
	if receiver.vaultData == nil {
		return ErrMissingVault
	}

	if item.OrganizationID != uuid.Nil {
		return errors.New("creating items inside organizations is not supported yet")
	}

	resultItem := clone.Clone(item)
	err := receiver.encryptStruct(ctx, resultItem, session.Encryption.UserKey)
	if err != nil {
		return fmt.Errorf("failed encrypting struct: %w", err)
	}

	targetUri := new(*receiver.baseURL)
	targetUri.Path = "/ciphers"
	newItemEnc, err := request[*result.Item](ctx, receiver.httpClient, http.MethodPost, targetUri, resultItem, session)
	if err != nil {
		return fmt.Errorf("failed creating the item: %w", err)
	}
	newItemDec, err := receiver.DecryptItem(ctx, session, newItemEnc)
	if err != nil {
		return fmt.Errorf("failed decrypting the item: %w", err)
	}
	*item = *newItemDec
	receiver.vaultData.Items = append(receiver.vaultData.Items, newItemEnc)

	return nil
}

func (receiver *vault) encryptStruct(ctx context.Context, target any, key dto.Key) error {
	typ := reflect.TypeOf(target)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct, %T given", target)
	}

	for field := range typ.Elem().Fields() {
		if field.Type.Kind() == reflect.String || (field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.String) {
			strVal, err := getStringValue(field, target)
			if err != nil {
				if errors.Is(err, errValIsNil) {
					continue
				}
				return err
			}

			newVal, err := crypto.EncryptString(strVal, key)
			if err != nil {
				return fmt.Errorf("failed encrypting string: %w", err)
			}

			value := reflect.ValueOf(target).Elem().FieldByName(field.Name)
			if value.Type().Kind() == reflect.Pointer {
				value = value.Elem()
			}

			value.SetString(newVal)
		} else if field.Type.Kind() == reflect.Pointer && field.Type.Elem().Kind() == reflect.Struct {
			value := reflect.ValueOf(target).Elem().FieldByName(field.Name)
			if value.IsNil() {
				continue
			}
			if err := receiver.encryptStruct(ctx, value.Interface(), key); err != nil {
				return err
			}
		} else if field.Type.Kind() == reflect.Slice {
			value := reflect.ValueOf(target).Elem().FieldByName(field.Name)
			for i := range value.Len() {
				elem := value.Index(i)
				if elem.Kind() == reflect.Pointer && elem.Elem().Kind() == reflect.Struct {
					err := receiver.encryptStruct(ctx, elem.Elem().Addr().Interface(), key)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}
