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
		subject: "user", object: "/api/properties", action: "read",
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
		subject: "user", object: "/api/features", action: "read",
	},

	// api/property-logs
	// admin
	{
		subject: "admin", object: "/api/property-logs", action: "create",
	},
	{
		subject: "admin", object: "/api/property-logs", action: "read",
	},
	{
		subject: "admin", object: "/api/property-logs", action: "update",
	},
	{
		subject: "admin", object: "/api/property-logs", action: "delete",
	},

	// api/contacts
	// admin
	{
		subject: "admin", object: "/api/contacts", action: "create",
	},
	{
		subject: "admin", object: "/api/contacts", action: "read",
	},
	{
		subject: "admin", object: "/api/contacts", action: "update",
	},
	{
		subject: "admin", object: "/api/contacts", action: "delete",
	},

	// api/tasks
	// admin
	{
		subject: "admin", object: "/api/tasks", action: "create",
	},
	{
		subject: "admin", object: "/api/tasks", action: "read",
	},
	{
		subject: "admin", object: "/api/tasks", action: "update",
	},
	{
		subject: "admin", object: "/api/tasks", action: "delete",
	},

	// api/task-logs
	// admin
	{
		subject: "admin", object: "/api/task-logs", action: "create",
	},
	{
		subject: "admin", object: "/api/task-logs", action: "read",
	},
	{
		subject: "admin", object: "/api/task-logs", action: "update",
	},
	{
		subject: "admin", object: "/api/task-logs", action: "delete",
	},
}
