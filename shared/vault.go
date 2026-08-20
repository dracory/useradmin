package shared

import (
	"context"

	"github.com/dracory/userstore"
)

// VaultTokenizer abstracts vault tokenization/untokenization of user
// fields (first name, last name, email, phone, business name).
//
// Host projects that use a vault store provide an implementation so
// useradmin can read and write tokenized user data without depending
// on any specific vault config or key management. If nil, useradmin
// treats user fields as plain text.
type VaultTokenizer interface {
	// Tokenize upserts tokens for the given user fields and returns
	// the resulting token strings to store on the user record.
	Tokenize(
		ctx context.Context,
		user userstore.UserInterface,
		firstName, lastName, email, phone, businessName string,
	) (firstNameToken, lastNameToken, emailToken, phoneToken, businessNameToken string, err error)

	// Untokenize resolves the tokenized fields on the given user back
	// to their plain-text values.
	Untokenize(
		ctx context.Context,
		user userstore.UserInterface,
	) (firstName, lastName, email, phone, businessName string, err error)
}
