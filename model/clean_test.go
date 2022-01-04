package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/horkruxes/hkxserver/exceptions"
)

func TestNormalize(t *testing.T) {
	correctMessage := Message{
		ID:              uuidTest,
		CreatedAt:       time.Now(),
		DisplayedName:   "  Test Name - Message title ",
		AuthorBase64:    pubKey,
		SignatureBase64: "  xO9KBOZaPxmT6IcVUazXTyj7mmCCmf8gnCXtmNd6GuGRW-naj6dubiIPYSEAyt6UE0rNCzV0G71w7xgfF5GcCA==",
		Content:         " Lorem <script>alert('xss')</script>ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	err := correctMessage.Normalize()
	if err != nil {
		t.Errorf(err.Error())
	}

	if correctMessage.DisplayedName != "Test Name - Message title" {
		t.Errorf("Expected: 'Test Name - Message title', got: %#v", correctMessage.DisplayedName)
	}

	if correctMessage.Verify() != exceptions.ErrorWrongSignature {
		t.Errorf("Expected: wrong signature, got: %#v", correctMessage.Verify())
	}
}

func BenchmarkNormalize(b *testing.B) {
	correctMessage := Message{
		ID:              uuid.NewString(),
		CreatedAt:       time.Now(),
		DisplayedName:   "Test Name - Message title",
		AuthorBase64:    pubKey,
		SignatureBase64: "GOE2De_UO2i9g9yl2wK2VBGenjCTy-fAscMqBjkBcVT2oqYhj-wNTcM67TsYqbD17nre7_fFzXUJDGjp5dTwDg==",
		Content:         "   Lorem <script>alert('XSS')</script> ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	for i := 0; i < b.N; i++ {
		correctMessage.Normalize()
	}
}
