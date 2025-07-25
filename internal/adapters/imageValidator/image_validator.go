package image_validator

import (
	"1337b04rd/internal/domain"
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

// MIMEValidator implements the ImageValidator port.
type MIMEValidator struct {
	allowedTypes map[string]struct{} // Set of allowed MIME types (e.g., "image/jpeg": {})
	maxSize      int64               // Optional max file size (in bytes)
}

// New creates a validator with allowed MIME types and optional max size.
func New(allowedTypes []string, maxSize int64) domain.ImageValidator {
	allowed := make(map[string]struct{})
	for _, t := range allowedTypes {
		allowed[t] = struct{}{}
	}
	return &MIMEValidator{
		allowedTypes: allowed,
		maxSize:      maxSize,
	}
}

// Validate checks:
// 1. File size (if maxSize > 0)
// 2. MIME type (against allowedTypes)
// 3. Magic numbers (to prevent spoofing)
func (v *MIMEValidator) Validate(image []byte) error {
	// 1. Check size
	if v.maxSize > 0 && int64(len(image)) > v.maxSize {
		return fmt.Errorf("image exceeds max size of %d bytes", v.maxSize)
	}

	// 2. Detect MIME type
	mimeType := http.DetectContentType(image)
	if _, allowed := v.allowedTypes[mimeType]; !allowed {
		return fmt.Errorf("invalid MIME type: %s (allowed: %v)", mimeType, v.AllowedTypes())
	}

	// 3. Verify magic numbers (optional but recommended)
	if err := v.validateMagicNumbers(image, mimeType); err != nil {
		return err
	}

	return nil
}

// AllowedTypes returns the allowed MIME types.
func (v *MIMEValidator) AllowedTypes() []string {
	types := make([]string, 0, len(v.allowedTypes))
	for t := range v.allowedTypes {
		types = append(types, t)
	}
	return types
}

// validateMagicNumbers checks file signatures to prevent spoofing.
func (v *MIMEValidator) validateMagicNumbers(image []byte, mimeType string) error {
	switch mimeType {
	case "image/jpeg":
		if !bytes.HasPrefix(image, []byte{0xFF, 0xD8, 0xFF}) {
			return errors.New("invalid JPEG signature")
		}
	case "image/png":
		if !bytes.HasPrefix(image, []byte{0x89, 0x50, 0x4E, 0x47}) {
			return errors.New("invalid PNG signature")
		}
		// Add more types as needed (e.g., GIF, WEBP)
	}
	return nil
}
