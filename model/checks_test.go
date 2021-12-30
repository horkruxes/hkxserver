package model

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

const (
	pubKey   = "CErL14_ce3wEuivQrQSVUWNV4-O1vN1C89LF_bladi4="
	secKey   = "Xhuky7Q7j87Bz6cE5udFJ_VOurWg_P2kjOQ5EhlGDO8ISsvXj9x7fAS6K9CtBJVRY1Xj47W83ULz0sX9uVp2Lg=="
	uuidTest = "363de435-0268-4000-a333-fad3e1f6dd1f"
)

func TestSanitize(t *testing.T) {
	correctMessage := Message{
		ID:              uuidTest,
		CreatedAt:       time.Now(),
		DisplayedName:   "  Test Name - Message title ",
		AuthorBase64:    pubKey,
		SignatureBase64: "  xO9KBOZaPxmT6IcVUazXTyj7mmCCmf8gnCXtmNd6GuGRW-naj6dubiIPYSEAyt6UE0rNCzV0G71w7xgfF5GcCA==",
		Content:         " Lorem <script>alert('xss')</script>ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	err := correctMessage.Sanitize(false)
	if err != nil {
		t.Errorf(err.Error())
	}

	if correctMessage.DisplayedName != "Test Name - Message title" {
		t.Errorf("Expected: 'Test Name - Message title', got: %#v", correctMessage.DisplayedName)
	}

	if correctMessage.Content != "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum" {
		t.Errorf("Not protected against XSS")
	}
}

func BenchmarkSanitize(b *testing.B) {
	correctMessage := Message{
		ID:              uuid.NewString(),
		CreatedAt:       time.Now(),
		DisplayedName:   "Test Name - Message title",
		AuthorBase64:    pubKey,
		SignatureBase64: "GOE2De_UO2i9g9yl2wK2VBGenjCTy-fAscMqBjkBcVT2oqYhj-wNTcM67TsYqbD17nre7_fFzXUJDGjp5dTwDg==",
		Content:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	for i := 0; i < b.N; i++ {
		correctMessage.Sanitize(true)
	}
}

func TestVerifyConstraints(t *testing.T) {
	correctMessage := Message{
		ID:              uuid.NewString(),
		CreatedAt:       time.Now(),
		DisplayedName:   "Test Name - Message title",
		AuthorBase64:    pubKey,
		SignatureBase64: "GOE2De_UO2i9g9yl2wK2VBGenjCTy-fAscMqBjkBcVT2oqYhj-wNTcM67TsYqbD17nre7_fFzXUJDGjp5dTwDg==",
		Content:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	if err := correctMessage.verifyConstraints(true); err != nil {
		t.Errorf(err.Error())
	}
}

func TestVerifyOwnerShip(t *testing.T) {
	// Correct signature
	correctMessage := Message{
		ID:              uuid.NewString(),
		CreatedAt:       time.Now(),
		DisplayedName:   "Test Name - Message title",
		AuthorBase64:    pubKey,
		SignatureBase64: "GOE2De_UO2i9g9yl2wK2VBGenjCTy-fAscMqBjkBcVT2oqYhj-wNTcM67TsYqbD17nre7_fFzXUJDGjp5dTwDg==",
		Content:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	if !correctMessage.verifyOwnerShip() {
		t.Errorf("Wrong signature")
	}

	// Wrong signature
	wrongMessage := Message{
		ID:              "",
		CreatedAt:       time.Now(),
		DisplayedName:   "Test Name - Message title",
		AuthorBase64:    pubKey,
		SignatureBase64: "gOE2De_UO2i9g9yl2wK2VBGenjCTy-fAscMqBjkBcVT2oqYhj-wNTcM67TsYqbD17nre7_fFzXUJDGjp5dTwDg==",
		Content:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	if wrongMessage.verifyOwnerShip() {
		t.Errorf("Wrong signature not detected")
	}
}

func BenchmarkVerifyOwnership(b *testing.B) {
	correctMessage := Message{
		ID:              uuid.NewString(),
		CreatedAt:       time.Now(),
		DisplayedName:   "Test Name - Message title",
		AuthorBase64:    pubKey,
		SignatureBase64: "GOE2De_UO2i9g9yl2wK2VBGenjCTy-fAscMqBjkBcVT2oqYhj-wNTcM67TsYqbD17nre7_fFzXUJDGjp5dTwDg==",
		Content:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	for i := 0; i < b.N; i++ {
		correctMessage.verifyOwnerShip()
	}
}
