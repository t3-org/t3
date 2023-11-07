package vcomplement

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"regexp"
)

//--------------------------------
// Social network identifier
//--------------------------------

// ErrSocialNetworkIdentifierInvalid is the default SocialNetworkIdentifier validation rules error.
var ErrSocialNetworkIdentifierInvalid = validation.NewError("validation_social_network_identifier_invalid", "Social network identifier is invalid")

// SocialNetworkIdentifier is the social social network identifier rule
// It's for telegram accounts.
var SocialNetworkIdentifier = validation.Match(regexp.MustCompile("^[A-Za-z]{2,}[_-]?[A-Za-z0-9]{2,}$")).ErrorObject(ErrSocialNetworkIdentifierInvalid)
