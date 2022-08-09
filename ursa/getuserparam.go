package ursa

type GetUsersParam struct {
	// Keyword defines keyword that will be used for filtering URSA users.
	//
	// Optional.
	Keyword string

	// Organization defines user organization.
	//
	// Optional.
	Organization string

	// IsDeactivated defines whether a user is deactivated.
	//
	// Optional.
	IsDeactivated bool
}
