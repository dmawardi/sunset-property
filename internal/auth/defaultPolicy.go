package auth

type policySet struct {
	subject string
	object  string
	action  string
}

var DefaultPolicyList = []policySet{
	// User
	// api/me
	{
		subject: "user", object: "/api/me", action: "read",
	},
	{
		subject: "user", object: "/api/me", action: "update",
	},
	// Admin
	// api/me
	{
		subject: "admin", object: "/api/me", action: "read",
	},
	{
		subject: "admin", object: "/api/me", action: "create",
	},
	{
		subject: "admin", object: "/api/me", action: "update",
	},
	// api/users
	{
		subject: "admin", object: "/api/users", action: "create",
	},
	{
		subject: "admin", object: "/api/users", action: "read",
	},
	{
		subject: "admin", object: "/api/users", action: "update",
	},
	{
		subject: "admin", object: "/api/users", action: "delete",
	},
	// api/properties
	// admin
	{
		subject: "admin", object: "/api/properties", action: "create",
	},
	{
		subject: "admin", object: "/api/properties", action: "read",
	},
	{
		subject: "admin", object: "/api/properties", action: "update",
	},
	{
		subject: "admin", object: "/api/properties", action: "delete",
	},
	// user
	{
		subject: "user", object: "/api/properties", action: "create",
	},
	{
		subject: "user", object: "/api/properties", action: "read",
	},
	{
		subject: "user", object: "/api/properties", action: "update",
	},
	{
		subject: "user", object: "/api/properties", action: "delete",
	},
	// api/features
	// admin
	{
		subject: "admin", object: "/api/features", action: "create",
	},
	{
		subject: "admin", object: "/api/features", action: "read",
	},
	{
		subject: "admin", object: "/api/features", action: "update",
	},
	{
		subject: "admin", object: "/api/features", action: "delete",
	},
	// user
	{
		subject: "user", object: "/api/features", action: "create",
	},
	{
		subject: "user", object: "/api/features", action: "read",
	},
	{
		subject: "user", object: "/api/features", action: "update",
	},
	{
		subject: "user", object: "/api/features", action: "delete",
	},
}
