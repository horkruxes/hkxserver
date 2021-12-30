package model

import (
	"testing"
	"time"
)

const (
	pubKey = "CErL14_ce3wEuivQrQSVUWNV4-O1vN1C89LF_bladi4="
	secKey = "Xhuky7Q7j87Bz6cE5udFJ_VOurWg_P2kjOQ5EhlGDO8ISsvXj9x7fAS6K9CtBJVRY1Xj47W83ULz0sX9uVp2Lg=="
)

func TestVerifyConstraints(t *testing.T) {
	correctMessage := Message{
		ID:              "",
		CreatedAt:       time.Now(),
		DisplayedName:   "Test Name - Message title",
		AuthorBase64:    pubKey,
		SignatureBase64: "GOE2De_UO2i9g9yl2wK2VBGenjCTy-fAscMqBjkBcVT2oqYhj-wNTcM67TsYqbD17nre7_fFzXUJDGjp5dTwDg==",
		Content:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	if err := correctMessage.VerifyConstraints(); err != nil {
		t.Errorf(err.Error())
	}
}

func TestVerifyOwnerShip(t *testing.T) {
	// Correct signature
	correctMessage := Message{
		ID:              "",
		CreatedAt:       time.Now(),
		DisplayedName:   "Test Name - Message title",
		AuthorBase64:    pubKey,
		SignatureBase64: "GOE2De_UO2i9g9yl2wK2VBGenjCTy-fAscMqBjkBcVT2oqYhj-wNTcM67TsYqbD17nre7_fFzXUJDGjp5dTwDg==",
		Content:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	if !correctMessage.VerifyOwnerShip() {
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

	if wrongMessage.VerifyOwnerShip() {
		t.Errorf("Wrong signature not detected")
	}
}
