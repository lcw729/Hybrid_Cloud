package util

type OidcIdentityProviderConfigRequest struct {
	_ struct{} `type:"structure"`

	// This is also known as audience. The ID for the client application that makes
	// authentication requests to the OpenID identity provider.
	//
	// ClientId is a required field
	ClientId *string `locationName:"clientId" type:"string" required:"true"`

	// The JWT claim that the provider uses to return your groups.
	GroupsClaim *string `locationName:"groupsClaim" type:"string"`

	// The prefix that is prepended to group claims to prevent clashes with existing
	// names (such as system: groups). For example, the valueoidc: will create group
	// names like oidc:engineering and oidc:infra.
	GroupsPrefix *string `locationName:"groupsPrefix" type:"string"`

	// The name of the OIDC provider configuration.
	//
	// IdentityProviderConfigName is a required field
	IdentityProviderConfigName *string `locationName:"identityProviderConfigName" type:"string" required:"true"`

	// The URL of the OpenID identity provider that allows the API server to discover
	// public signing keys for verifying tokens. The URL must begin with https://
	// and should correspond to the iss claim in the provider's OIDC ID tokens.
	// Per the OIDC standard, path components are allowed but query parameters are
	// not. Typically the URL consists of only a hostname, like https://server.example.org
	// or https://example.com. This URL should point to the level below .well-known/openid-configuration
	// and must be publicly accessible over the internet.
	//
	// IssuerUrl is a required field
	IssuerUrl *string `locationName:"issuerUrl" type:"string" required:"true"`

	// The key value pairs that describe required claims in the identity token.
	// If set, each claim is verified to be present in the token with a matching
	// value. For the maximum number of claims that you can require, see Amazon
	// EKS service quotas (https://docs.aws.amazon.com/eks/latest/userguide/service-quotas.html)
	// in the Amazon EKS User Guide.
	RequiredClaims map[string]*string `locationName:"requiredClaims" type:"map"`

	// The JSON Web Token (JWT) claim to use as the username. The default is sub,
	// which is expected to be a unique identifier of the end user. You can choose
	// other claims, such as email or name, depending on the OpenID identity provider.
	// Claims other than email are prefixed with the issuer URL to prevent naming
	// clashes with other plug-ins.
	UsernameClaim *string `locationName:"usernameClaim" type:"string"`

	// The prefix that is prepended to username claims to prevent clashes with existing
	// names. If you do not provide this field, and username is a value other than
	// email, the prefix defaults to issuerurl#. You can use the value - to disable
	// all prefixing.
	UsernamePrefix *string `locationName:"usernamePrefix" type:"string"`
}
