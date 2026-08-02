package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.chrastecky.dev/bitsailor-core/bitwarden"
	"go.chrastecky.dev/bitsailor-core/bitwarden/result"
)

type diagnostic struct {
	name    string
	err     error
	skipped bool
}

type runner struct {
	ctx         context.Context
	client      bitwarden.Client
	auth        bitwarden.Auth
	vault       bitwarden.Vault
	session     *result.Session
	email       string
	password    string
	keepCreated bool
	results     []diagnostic
}

func main() {
	var (
		server       = flag.String("server", "", "Bitwarden server base URL. Empty uses library defaults.")
		email        = flag.String("email", "", "Account email for password login, or for unlocking an API-key session.")
		password     = flag.String("password", "", "Master password for password login, or for unlocking an API-key session.")
		totp         = flag.String("totp", "", "Optional two-factor code for password login.")
		clientID     = flag.String("client-id", "", "API key client_id.")
		clientSecret = flag.String("client-secret", "", "API key client_secret.")
		deviceID     = flag.String("device-id", "", "Optional Bitwarden device UUID. Empty generates one for this run.")
		timeout      = flag.Duration("timeout", 2*time.Minute, "Overall diagnostic timeout.")
		keepCreated  = flag.Bool("keep-created", false, "Keep diagnostic vault objects instead of deleting them.")
	)
	flag.Parse()

	if err := validateFlags(*email, *password, *clientID, *clientSecret); err != nil {
		fmt.Fprintf(os.Stderr, "invalid arguments: %v\n\n", err)
		flag.Usage()
		os.Exit(2)
	}

	options := make([]bitwarden.Option, 0)
	if strings.TrimSpace(*server) != "" {
		options = append(options, bitwarden.WithBaseURL(strings.TrimSpace(*server)))
	}
	parsedDeviceID, err := diagnosticDeviceID(*deviceID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid arguments: %v\n\n", err)
		flag.Usage()
		os.Exit(2)
	}
	options = append(options, bitwarden.WithDeviceID(parsedDeviceID))

	client, err := bitwarden.NewClient(options...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed creating client: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	runnerInstance := &runner{
		ctx:         ctx,
		client:      client,
		auth:        client.Auth(),
		vault:       client.Vault(),
		email:       strings.TrimSpace(*email),
		password:    *password,
		keepCreated: *keepCreated,
	}

	fmt.Printf("bitsailor-core Bitwarden diagnostic\n")
	fmt.Printf("server: %s\n", displayServer(*server))
	fmt.Printf("login: %s\n", loginMode(*clientID, *clientSecret))
	fmt.Printf("device ID: %s\n", parsedDeviceID)
	fmt.Printf("started: %s\n\n", time.Now().Format(time.RFC3339))

	if *clientID != "" || *clientSecret != "" {
		runnerInstance.step("auth: login with API key", func() error {
			session, err := runnerInstance.auth.LoginApiKey(runnerInstance.ctx, *clientID, *clientSecret)
			runnerInstance.session = session
			return err
		})
		runnerInstance.step("auth: unlock API-key session with master password", func() error {
			if runnerInstance.session == nil {
				return fmt.Errorf("API-key login did not produce a session")
			}
			return runnerInstance.auth.UnlockSession(runnerInstance.ctx, runnerInstance.session, runnerInstance.email, runnerInstance.password)
		})
	} else {
		runnerInstance.step("auth: login with email and password", func() error {
			var twoFaCode *string
			if strings.TrimSpace(*totp) != "" {
				twoFaCode = totp
			}
			session, err := runnerInstance.auth.LoginPassword(runnerInstance.ctx, runnerInstance.email, runnerInstance.password, twoFaCode)
			runnerInstance.session = session
			return err
		})
	}

	runnerInstance.runClientDiagnostics()
	runnerInstance.runVaultDiagnostics()
	runnerInstance.printSummary()
}

func validateFlags(email, password, clientID, clientSecret string) error {
	hasAPIKeyPart := clientID != "" || clientSecret != ""
	if strings.TrimSpace(password) == "" {
		return fmt.Errorf("--password is required")
	}
	if strings.TrimSpace(email) == "" {
		return fmt.Errorf("--email is required")
	}
	if hasAPIKeyPart && (strings.TrimSpace(clientID) == "" || strings.TrimSpace(clientSecret) == "") {
		return fmt.Errorf("--client-id and --client-secret must be provided together")
	}
	return nil
}

func diagnosticDeviceID(value string) (uuid.UUID, error) {
	if strings.TrimSpace(value) == "" {
		return uuid.New(), nil
	}
	parsed, err := uuid.Parse(strings.TrimSpace(value))
	if err != nil {
		return uuid.Nil, fmt.Errorf("--device-id must be a UUID: %w", err)
	}
	if parsed == uuid.Nil {
		return uuid.Nil, fmt.Errorf("--device-id must not be nil")
	}
	return parsed, nil
}

func (receiver *runner) runClientDiagnostics() {
	receiver.step("generator: password with defaults", func() error {
		password, err := receiver.client.GeneratePassword(&bitwarden.PasswordGeneratorRequest{})
		if err != nil {
			return err
		}
		if password == "" {
			return fmt.Errorf("generated password is empty")
		}
		return nil
	})

	receiver.step("generator: passphrase with defaults", func() error {
		passphrase, err := receiver.client.GeneratePassphrase(&bitwarden.PassphraseGeneratorRequest{})
		if err != nil {
			return err
		}
		if passphrase == "" {
			return fmt.Errorf("generated passphrase is empty")
		}
		return nil
	})
}

func (receiver *runner) runVaultDiagnostics() {
	if receiver.session == nil {
		receiver.skip("auth: refresh token", "login did not produce a session")
		receiver.skip("vault: sync", "login did not produce a session")
		return
	}

	receiver.step("auth: refresh token", func() error {
		return receiver.auth.RefreshToken(receiver.ctx, receiver.session)
	})

	receiver.step("vault: sync", func() error {
		vault, err := receiver.vault.Sync(receiver.ctx, receiver.session)
		if err != nil {
			return err
		}
		receiver.vault = vault
		return nil
	})

	receiver.step("vault: inspect synced data", func() error {
		data := receiver.vault.GetVaultData()
		if data == nil {
			return fmt.Errorf("vault data is nil after sync")
		}
		fmt.Printf("      profile: %s, items: %d, folders: %d, collections: %d, sends: %d\n",
			nilSafeProfileEmail(data.Profile), len(data.Items), len(data.Folders), len(data.Collections), len(data.Sends))
		return nil
	})

	receiver.step("vault: clone with synced data", func() error {
		data := receiver.vault.GetVaultData()
		if data == nil {
			return fmt.Errorf("vault data is nil after sync")
		}
		clone := receiver.vault.WithVaultData(data)
		if clone.GetVaultData() == nil {
			return fmt.Errorf("cloned vault data is nil")
		}
		return nil
	})

	if receiver.session.Encryption == nil || len(receiver.session.Encryption.UserKey) == 0 {
		receiver.skip("vault: list items", "session is not unlocked")
		receiver.skip("vault: get first existing item", "session is not unlocked")
		receiver.skip("vault: create secure-note item", "session is not unlocked")
		receiver.skip("vault: get created item", "session is not unlocked")
		receiver.skip("vault: update created item", "session is not unlocked")
		receiver.skip("vault: list sends", "session is not unlocked")
		receiver.skip("vault: get first existing send", "session is not unlocked")
		receiver.skip("vault: create text send", "session is not unlocked")
		receiver.skip("vault: get created text send", "session is not unlocked")
		receiver.skip("vault: create file send", "session is not unlocked")
		receiver.skip("vault: get created file send", "session is not unlocked")
		receiver.skip("cleanup: delete created item", "session is not unlocked")
		receiver.skip("cleanup: delete created text send", "session is not unlocked")
		receiver.skip("cleanup: delete created file send", "session is not unlocked")
		return
	}

	receiver.step("vault: list items", func() error {
		items, err := receiver.vault.GetItems(receiver.ctx, receiver.session)
		if err != nil {
			return err
		}
		fmt.Printf("      decrypted items: %d\n", len(items))
		return nil
	})

	receiver.step("vault: get first existing item", func() error {
		items, err := receiver.vault.GetItems(receiver.ctx, receiver.session)
		if err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		_, err = receiver.vault.GetItem(receiver.ctx, receiver.session, items[0].ID)
		return err
	})

	itemName := fmt.Sprintf("bitsailor diagnostic item %s", time.Now().UTC().Format("20060102T150405Z"))
	testItem := &result.Item{
		Type:       result.ItemTypeSecureNote,
		Name:       itemName,
		SecureNote: &result.ItemSecureNote{Type: 0},
	}

	receiver.step("vault: create secure-note item", func() error {
		return receiver.vault.CreateItem(receiver.ctx, receiver.session, testItem)
	})
	receiver.step("vault: get created item", func() error {
		if testItem.ID == uuid.Nil {
			return fmt.Errorf("created item has nil ID")
		}
		_, err := receiver.vault.GetItem(receiver.ctx, receiver.session, testItem.ID)
		return err
	})
	receiver.step("vault: update created item", func() error {
		if testItem.ID == uuid.Nil {
			return fmt.Errorf("created item has nil ID")
		}
		notes := "updated by bitsailor diagnostic command"
		testItem.Name += " updated"
		testItem.Notes = &notes
		return receiver.vault.UpdateItem(receiver.ctx, receiver.session, testItem)
	})

	receiver.step("vault: list sends", func() error {
		sends, err := receiver.vault.GetSends(receiver.ctx, receiver.session)
		if err != nil {
			return err
		}
		fmt.Printf("      decrypted sends: %d\n", len(sends))
		return nil
	})

	receiver.step("vault: get first existing send", func() error {
		sends, err := receiver.vault.GetSends(receiver.ctx, receiver.session)
		if err != nil {
			return err
		}
		if len(sends) == 0 {
			return nil
		}
		_, err = receiver.vault.GetSend(receiver.ctx, receiver.session, sends[0].ID)
		return err
	})

	textSend := &result.Send{
		Type:         result.SendTypeText,
		AuthType:     result.SendAuthTypeNoAuth,
		Name:         fmt.Sprintf("bitsailor diagnostic text send %s", time.Now().UTC().Format("20060102T150405Z")),
		DeletionDate: time.Now().UTC().Add(24 * time.Hour),
		Text:         &result.SendText{Text: "bitsailor diagnostic send", Hidden: false},
	}
	receiver.step("vault: create text send", func() error {
		return receiver.vault.CreateSend(receiver.ctx, receiver.session, textSend)
	})
	receiver.step("vault: get created text send", func() error {
		if textSend.ID == uuid.Nil {
			return fmt.Errorf("created text send has nil ID")
		}
		_, err := receiver.vault.GetSend(receiver.ctx, receiver.session, textSend.ID)
		return err
	})

	fileContent := []byte("bitsailor diagnostic file send\n")
	fileSend := &result.Send{
		Type:         result.SendTypeFile,
		AuthType:     result.SendAuthTypeNoAuth,
		Name:         fmt.Sprintf("bitsailor diagnostic file send %s", time.Now().UTC().Format("20060102T150405Z")),
		DeletionDate: time.Now().UTC().Add(24 * time.Hour),
		File:         &result.SendFile{FileName: "bitsailor-diagnostic.txt"},
		InputFile:    bytes.NewReader(fileContent),
		FileLength:   len(fileContent),
	}
	receiver.step("vault: create file send", func() error {
		return receiver.vault.CreateSend(receiver.ctx, receiver.session, fileSend)
	})
	receiver.step("vault: get created file send", func() error {
		if fileSend.ID == uuid.Nil {
			return fmt.Errorf("created file send has nil ID")
		}
		_, err := receiver.vault.GetSend(receiver.ctx, receiver.session, fileSend.ID)
		return err
	})

	if receiver.keepCreated {
		receiver.skip("cleanup: delete created item", "--keep-created was set")
		receiver.skip("cleanup: delete created text send", "--keep-created was set")
		receiver.skip("cleanup: delete created file send", "--keep-created was set")
	} else {
		receiver.step("cleanup: delete created item", func() error {
			if testItem.ID == uuid.Nil {
				return nil
			}
			return receiver.vault.DeleteItem(receiver.ctx, receiver.session, testItem.ID)
		})
		receiver.step("cleanup: delete created text send", func() error {
			if textSend.ID == uuid.Nil {
				return nil
			}
			return receiver.vault.DeleteSend(receiver.ctx, receiver.session, textSend.ID)
		})
		receiver.step("cleanup: delete created file send", func() error {
			if fileSend.ID == uuid.Nil {
				return nil
			}
			return receiver.vault.DeleteSend(receiver.ctx, receiver.session, fileSend.ID)
		})
	}
}

func (receiver *runner) step(name string, fn func() error) {
	fmt.Printf("[RUN]  %s\n", name)
	err := fn()
	receiver.results = append(receiver.results, diagnostic{name: name, err: err})
	if err != nil {
		fmt.Printf("[FAIL] %s: %v\n\n", name, err)
		return
	}
	fmt.Printf("[OK]   %s\n\n", name)
}

func (receiver *runner) skip(name, reason string) {
	receiver.results = append(receiver.results, diagnostic{name: name, skipped: true})
	fmt.Printf("[SKIP] %s: %s\n\n", name, reason)
}

func (receiver *runner) printSummary() {
	var failed, skipped int
	for _, result := range receiver.results {
		if result.err != nil {
			failed++
		} else if result.skipped {
			skipped++
		}
	}

	fmt.Printf("summary: %d succeeded, %d failed, %d skipped\n", len(receiver.results)-failed-skipped, failed, skipped)
	if failed > 0 {
		os.Exit(1)
	}
}

func displayServer(server string) string {
	if strings.TrimSpace(server) == "" {
		return "library default"
	}
	return strings.TrimSpace(server)
}

func loginMode(clientID, clientSecret string) string {
	if clientID != "" || clientSecret != "" {
		return "api key + master password unlock"
	}
	return "email + master password"
}

func nilSafeProfileEmail(profile *result.Profile) string {
	if profile == nil {
		return "<nil>"
	}
	if profile.Email == "" {
		return "<empty>"
	}
	return profile.Email
}
