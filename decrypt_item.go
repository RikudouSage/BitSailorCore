package bitwarden

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/uuid"
	clone "github.com/huandu/go-clone/generic"
	"github.com/samber/lo"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/dto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func (receiver *vault) DecryptItem(ctx context.Context, session *result.Session, item *result.Item) (*result.Item, error) {
	key, err := receiver.getDecryptionKey(session, item)
	if err != nil {
		return nil, fmt.Errorf("failed fetching decryption key: %w", err)
	}

	resultItem := clone.Clone(item)
	resultItem.Notes, err = crypto.DecryptNullableString(resultItem.Notes, key)
	if err != nil {
		return nil, err
	}

	resultItem.Name, err = crypto.DecryptString(item.Name, key)
	if err != nil {
		return nil, err
	}

	if resultItem.Fields != nil {
		resultItem.Fields, err = lo.MapErr(resultItem.Fields, func(field *result.Field, _ int) (*result.Field, error) {
			field.Name, err = crypto.DecryptString(field.Name, key)
			if err != nil {
				return nil, err
			}
			field.Value, err = crypto.DecryptNullableString(field.Value, key)
			if err != nil {
				return nil, err
			}

			return field, nil
		})
		if err != nil {
			return nil, err
		}
	}

	switch item.Type {
	case result.ItemTypeLogin:
		err = receiver.decryptStruct(ctx, resultItem.Login, key)
	case result.ItemTypeSSHKey:
		err = receiver.decryptStruct(ctx, resultItem.SSHKey, key)
	case result.ItemTypeCard:
		err = receiver.decryptStruct(ctx, resultItem.Card, key)
	case result.ItemTypIdentity:
		err = receiver.decryptStruct(ctx, resultItem.Identity, key)
	case result.ItemTypeSecureNote:
		err = receiver.decryptStruct(ctx, resultItem.SecureNote, key)
	default:
		return nil, fmt.Errorf("unimplemented item type %d", item.Type)
	}
	if err != nil {
		return nil, err
	}

	return resultItem, nil
}

func (receiver *vault) decryptStruct(ctx context.Context, target any, key dto.Key) error {
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

			if !strings.HasPrefix(strVal, "2.") || strings.Count(strVal, "|") != 2 {
				continue
			}

			newVal, err := crypto.DecryptString(strVal, key)
			if err != nil {
				return err
			}

			value := reflect.ValueOf(target).Elem().FieldByName(field.Name)
			if value.Kind() == reflect.Pointer {
				value = value.Elem()
			}

			value.SetString(newVal)
		}
		if field.Type.Kind() == reflect.Slice {
			value := reflect.ValueOf(target).Elem().FieldByName(field.Name)
			for i := range value.Len() {
				elem := value.Index(i)
				if elem.Kind() == reflect.Pointer && elem.Elem().Kind() == reflect.Struct {
					err := receiver.decryptStruct(ctx, elem.Elem().Addr().Interface(), key)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (receiver *vault) getDecryptionKey(session *result.Session, item *result.Item) (dto.Key, error) {
	key := session.Encryption.UserKey
	if item.OrganizationID != uuid.Nil {
		orgKeys, err := receiver.getOrganizationKeys(session)
		if err != nil {
			return nil, fmt.Errorf("failed fetching organization keys: %w", err)
		}
		var ok bool
		key, ok = orgKeys[item.OrganizationID]
		if !ok {
			return nil, fmt.Errorf("failed fetching organization key for item %s (organization ID: %s)", item.ID, item.OrganizationID)
		}
	}

	if item.Key != nil {
		var err error
		key, err = crypto.DecryptBytes(*item.Key, key)
		if err != nil {
			return nil, fmt.Errorf("failed decrypting item key: %w", err)
		}
	}

	return key, nil
}
