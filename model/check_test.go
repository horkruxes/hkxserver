package model

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/horkruxes/hkxserver/exceptions"
)

const (
	pubKey     = "CErL14_ce3wEuivQrQSVUWNV4-O1vN1C89LF_bladi4="
	secKey     = "Xhuky7Q7j87Bz6cE5udFJ_VOurWg_P2kjOQ5EhlGDO8ISsvXj9x7fAS6K9CtBJVRY1Xj47W83ULz0sX9uVp2Lg=="
	uuidTest   = "363de435-0268-4000-a333-fad3e1f6dd1f"
	loremIpsum = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi vulputate nunc sapien, eu ullamcorper lacus blandit nec. Mauris ac sem non dolor iaculis cursus. Duis volutpat dolor in diam accumsan, vel fermentum magna tincidunt. Maecenas vel ullamcorper ipsum. Vestibulum vel sem semper nunc imperdiet lacinia non ultricies metus. Quisque et erat condimentum, viverra felis nec, semper libero. Morbi congue auctor risus in volutpat. Vivamus sit amet neque vitae nisi luctus tempus. Suspendisse auctor, nulla sit amet sollicitudin consequat, mauris sapien lobortis felis, ut rhoncus diam ante sit amet augue." // length of 607
)

func TestVerify(t *testing.T) {
	tests := []struct {
		cause   string
		title   string
		content string
		want    error
	}{
		{
			cause:   "escaped html",
			title:   "escaped",
			content: `Hello &lt; World ! &gt; &amp; Happy &#34;New&#39; Year. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi vulputate nunc sapien, eu ullamcorper lacus `,
			want:    nil,
		},
		{
			cause:   "single quote must be escaped",
			title:   "",
			content: "Shouldn't store ' quotes of any type" + loremIpsum,
			want:    exceptions.ErrorContentWithHTML,
		},
		{
			cause:   "double quote must be escaped",
			title:   "",
			content: "Shouldn't store \" quotes of any type" + loremIpsum,
			want:    exceptions.ErrorContentWithHTML,
		},
		{
			cause:   "XSS attack",
			title:   "",
			content: " Lorem <script>alert('xss')</script>ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
			want:    exceptions.ErrorContentWithHTML,
		},
		{
			cause:   "content too short",
			title:   "",
			content: " Loro cia deserunt mollit anim id est laborum",
			want:    exceptions.ErrorContentTooShort,
		},
		{
			cause:   "content too long",
			title:   "",
			content: strings.Repeat(loremIpsum, 100),
			want:    exceptions.ErrorFieldsTooLong,
		},
		{
			cause:   "title too long",
			title:   loremIpsum[:51],
			content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi vulputate nunc sapien, eu ullamcorper lacus blandit nec. Mauris ac sem non dolor iaculis cursus. Duis volutpat dolor in diam accumsan, vel fermentum magna tincidunt. Maecenas vel ullamcorper ipsum. Vestibulum vel sem semper nunc imperdiet lacinia non ultricies metus. Quisque et erat condimentum, viverra felis nec, semper libero. Morbi congue auctor risus in volutpat. Vivamus sit amet neque vitae nisi luctus tempus. Suspendisse auctor, nulla sit amet sollicitudin consequat, mauris sapien lobortis felis, ut rhoncus diam ante sit amet augue. ",
			want:    exceptions.ErrorFieldsTooLong,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.cause, func(t *testing.T) {
			t.Parallel()
			m := Message{
				ID:              uuidTest,
				CreatedAt:       time.Now(),
				DisplayedName:   tc.title,
				AuthorBase64:    pubKey,
				Content:         tc.content,
				SignatureBase64: "  xO9KBOZaPxmT6IcVUazXTyj7mmCCmf8gnCXtmNd6GuGRW-naj6dubiIPYSEAyt6UE0rNCzV0G71w7xgfF5GcCA==",
				MessageID:       "",
			}

			// focusing on content, not the signature
			if err := m.Verify(); !errors.Is(err, exceptions.ErrorWrongSignature) && err != tc.want {
				t.Errorf("want '%v', got '%v'", tc.want, err)
			}
		})
	}
}

func BenchmarkVerify(b *testing.B) {
	m := Message{
		ID:              uuid.NewString(),
		CreatedAt:       time.Now(),
		DisplayedName:   "Test Name - Message title",
		AuthorBase64:    pubKey,
		SignatureBase64: "GOE2De_UO2i9g9yl2wK2VBGenjCTy-fAscMqBjkBcVT2oqYhj-wNTcM67TsYqbD17nre7_fFzXUJDGjp5dTwDg==",
		Content:         "Lorem  ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
		MessageID:       "",
	}

	for i := 0; i < b.N; i++ {
		m.Verify()
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

	if err := correctMessage.verifyConstraints(); err != nil {
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
