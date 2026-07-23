package bitwarden

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/crypto"
	"go.chrastecky.dev/bitwarden-client/bitwarden/internal/dto"
	internalHttp "go.chrastecky.dev/bitwarden-client/bitwarden/internal/http"
	"go.chrastecky.dev/bitwarden-client/bitwarden/result"
)

func (receiver *vault) CreateSend(ctx context.Context, session *result.Session, item *result.Send) error {
	err, key := receiver.encryptSend(ctx, session, item)
	if err != nil {
		return fmt.Errorf("failed encrypting the send item: %w", err)
	}

	if item.Type == result.SendTypeText {
		return receiver.createTextSend(ctx, session, item)
	}

	return receiver.createFileSend(ctx, session, item, key)
}

func (receiver *vault) createTextSend(ctx context.Context, session *result.Session, item *result.Send) error {
	targetUri := new(*receiver.baseURL)
	targetUri.Path = "/sends"

	newItem, err := request[*result.Send](ctx, receiver.httpClient, http.MethodPost, targetUri, item, session)
	if err != nil {
		return fmt.Errorf("failed creating the send: %w", err)
	}

	*item = *newItem
	receiver.vaultData.Sends = append(receiver.vaultData.Sends, newItem)
	return nil

}

func (receiver *vault) createFileSend(ctx context.Context, session *result.Session, item *result.Send, key dto.Key) error {
	type uploadType int
	const (
		uploadTypeDirect uploadType = iota
		uploadTypeProvider
	)
	type metaResponse struct {
		URL            *string      `json:"url"`
		FileUploadType uploadType   `json:"fileUploadType"`
		SendResponse   *result.Send `json:"sendResponse"`
	}

	targetUri := new(*receiver.baseURL)
	targetUri.Path = "/sends/file/v2"

	if enc, ok := item.InputFile.(*crypto.File); ok {
		defer enc.Close()
	}

	err := receiver.encryptStruct(ctx, item.File, key, []string{"ID", "Size", "SizeName"})
	if err != nil {
		return fmt.Errorf("failed encrypting file metadata: %w", err)
	}

	meta, err := request[metaResponse](ctx, receiver.httpClient, http.MethodPost, targetUri, item, session)
	if err != nil {
		return fmt.Errorf("failed uploading file send metadata: %w", err)
	}

	if meta.FileUploadType == uploadTypeProvider {
		if meta.URL == nil {
			return errors.New("the target uri was returned as nil")
		}
		uri, err := url.Parse(*meta.URL)
		if err != nil {
			return fmt.Errorf("the returned uri was not a valid uri: %s", *meta.URL)
		}

		if err = receiver.uploadFileRemote(ctx, uri, item.InputFile, item.FileLength); err != nil {
			return fmt.Errorf("remote upload failed: %w", err)
		}
	} else if meta.FileUploadType == uploadTypeDirect {
		if err = receiver.uploadSendFileDirect(ctx, session, meta.SendResponse, item.InputFile); err != nil {
			return fmt.Errorf("direct upload failed: %w", err)
		}
	} else {
		return fmt.Errorf("unknown upload type: %d", meta.FileUploadType)
	}

	*item = *meta.SendResponse
	receiver.vaultData.Sends = append(receiver.vaultData.Sends, item)
	return nil
}

func (receiver *vault) uploadSendFileDirect(
	ctx context.Context,
	session *result.Session,
	metadata *result.Send,
	inputFile io.Reader,
) error {
	var body bytes.Buffer
	multipartWriter := multipart.NewWriter(&body)

	part, err := multipartWriter.CreateFormFile("data", metadata.File.FileName)
	if err != nil {
		return err
	}

	if _, err = io.Copy(part, inputFile); err != nil {
		return err
	}

	if err = multipartWriter.Close(); err != nil {
		return err
	}

	uri := new(*receiver.baseURL)
	uri.Path = fmt.Sprintf("/sends/%s/file/%s", metadata.ID, metadata.File.ID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri.String(), &body)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", session.Auth.TokenType+" "+session.Auth.AccessToken)
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	resp, err := receiver.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("direct upload failed: %d", resp.StatusCode)
	}

	return nil
}

func (receiver *vault) uploadFileRemote(ctx context.Context, uri *url.URL, file io.Reader, length int) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, uri.String(), file)
	if err != nil {
		return fmt.Errorf("failed creating request: %w", err)
	}
	req.Header.Set("X-MS-Date", time.Now().UTC().Format(http.TimeFormat))
	req.Header.Set("X-MS-Version", uri.Query().Get("sv"))
	req.Header.Set("x-MS-Blob-Type", "BlockBlob")
	req.Header.Set("Content-Length", strconv.Itoa(length))
	req.ContentLength = int64(length)

	resp, err := receiver.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed uploading to remote storage: %w", err)
	}
	defer internalHttp.DrainResponse(resp)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("remote upload failed: %d", resp.StatusCode)
	}

	return nil
}

func (receiver *vault) encryptSend(ctx context.Context, session *result.Session, item *result.Send) (error, dto.Key) {
	seed, err := crypto.GenerateRandomBytes(16)
	if err != nil {
		return fmt.Errorf("failed generating random encryption seed: %w", err), nil
	}
	key, err := crypto.DeriveSendKey(seed)
	if err != nil {
		return fmt.Errorf("failed deriving send encryption key: %w", err), nil
	}

	err = receiver.encryptStruct(ctx, item, key, []string{"Key", "Password", "Emails", "File"})
	if err != nil {
		return fmt.Errorf("failed encrypting fields: %w", err), nil
	}
	item.Key, err = crypto.EncryptBytes(seed, session.Encryption.UserKey)
	if err != nil {
		return fmt.Errorf("failed encrypting send seed: %w", err), nil
	}
	if item.Password != nil {
		item.Password = new(crypto.DeriveSendPassword(*item.Password, seed))
	}

	if item.InputFile != nil {
		encrypted, err := crypto.NewFile(key)
		if err != nil {
			return fmt.Errorf("failed initializing an encrypted file: %w", err), nil
		}
		_, err = io.Copy(encrypted, item.InputFile)
		if err != nil {
			return fmt.Errorf("failed encrypting the file: %w", err), nil
		}
		item.InputFile = encrypted
		item.FileLength = encrypted.Len()
	}

	return nil, key
}
