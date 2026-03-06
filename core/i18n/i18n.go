package i18n

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

// Translator loads and resolves translated messages by locale.
type Translator struct {
	mu       sync.RWMutex
	messages map[string]map[string]string // locale -> key -> message
	fallback string
}

// NewTranslator returns an empty translator with the given fallback locale.
func NewTranslator(fallback string) *Translator {
	return &Translator{
		messages: make(map[string]map[string]string),
		fallback: fallback,
	}
}

// LoadFile reads a JSON translation file and stores it under the given locale.
func (t *Translator) LoadFile(locale, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var msgs map[string]string
	if err := json.Unmarshal(data, &msgs); err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.messages[locale] = msgs
	return nil
}

// LoadDir loads all *.json files from a directory. Each file's name
// (without extension) is used as the locale.
func (t *Translator) LoadDir(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		locale := strings.TrimSuffix(e.Name(), ".json")
		if err := t.LoadFile(locale, filepath.Join(dir, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

// Get resolves a translation key for the given locale. If the key is not
// found in the requested locale, the fallback locale is tried. If still
// not found, the raw key is returned. When args are provided, the first
// element is used as template data for text/template interpolation.
func (t *Translator) Get(locale, key string, args ...interface{}) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	msg := t.resolve(locale, key)
	if msg == "" {
		msg = t.resolve(t.fallback, key)
	}
	if msg == "" {
		return key
	}

	if len(args) > 0 {
		tmpl, err := template.New("").Parse(msg)
		if err != nil {
			return msg
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, args[0]); err != nil {
			return msg
		}
		return buf.String()
	}
	return msg
}

func (t *Translator) resolve(locale, key string) string {
	if msgs, ok := t.messages[locale]; ok {
		if val, ok := msgs[key]; ok {
			return val
		}
	}
	return ""
}
