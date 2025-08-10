package rickMorty

import (
	"1337b04rd/internal/domain"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

// mockHTTPTransport creates a mock HTTP transport that returns predefined responses
type mockHTTPTransport struct {
	responseBody string
	statusCode   int
	err          error
}

func (m *mockHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &http.Response{
		StatusCode: m.statusCode,
		Body:       io.NopCloser(bytes.NewBufferString(m.responseBody)),
		Header:     make(http.Header),
	}, nil
}

func TestRickMortyAPI_GenerateAvatarAndName(t *testing.T) {
	tests := []struct {
		name          string
		id            int
		mockResponse  string
		mockStatus    int
		mockError     error
		expected      *domain.UserOutlook
		expectedError string
	}{
		{
			name: "successful response with valid ID",
			id:   1,
			mockResponse: `{
				"id": 1,
				"name": "Rick Sanchez",
				"image": "https://rickandmortyapi.com/api/character/avatar/1.jpeg"
			}`,
			mockStatus: http.StatusOK,
			expected: &domain.UserOutlook{
				Name:      "Rick Sanchez",
				AvatarURL: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
			},
		},
		{
			name: "high ID wraps around",
			id:   827, // Should wrap to 1 (827 % 826 = 1)
			mockResponse: `{
				"id": 1,
				"name": "Rick Sanchez",
				"image": "https://rickandmortyapi.com/api/character/avatar/1.jpeg"
			}`,
			mockStatus: http.StatusOK,
			expected: &domain.UserOutlook{
				Name:      "Rick Sanchez",
				AvatarURL: "https://rickandmortyapi.com/api/character/avatar/1.jpeg",
			},
		},
		{
			name: "exactly 826 ID",
			id:   826,
			mockResponse: `{
				"id": 826,
				"name": "Butter Robot",
				"image": "https://rickandmortyapi.com/api/character/avatar/826.jpeg"
			}`,
			mockStatus: http.StatusOK,
			expected: &domain.UserOutlook{
				Name:      "Butter Robot",
				AvatarURL: "https://rickandmortyapi.com/api/character/avatar/826.jpeg",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock transport
			mockTransport := &mockHTTPTransport{
				responseBody: tt.mockResponse,
				statusCode:   tt.mockStatus,
				err:          tt.mockError,
			}

			// Create API instance with mock client
			api := NewRickMortyAPIWithClient(&http.Client{
				Transport: mockTransport,
			})

			result, err := api.GenerateAvatarAndName(tt.id)

			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.expectedError)
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("Expected error containing '%s', got '%v'", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result.Name != tt.expected.Name || result.AvatarURL != tt.expected.AvatarURL {
				t.Errorf("Expected %+v, got %+v", tt.expected, result)
			}

			// Verify the ID wrapping
			expectedWrappedID := tt.id % 826
			if expectedWrappedID == 0 {
				expectedWrappedID = 826
			}
			if !strings.Contains(mockTransport.responseBody, fmt.Sprintf(`"id": %d`, expectedWrappedID)) {
				t.Errorf("Expected API call with wrapped ID %d", expectedWrappedID)
			}
		})
	}
}

func TestNewRickMortyAPIWithClient(t *testing.T) {
	testClient := &http.Client{}
	api := NewRickMortyAPIWithClient(testClient)
	if api.client != testClient {
		t.Error("Expected provided HTTP client to be used")
	}
}
