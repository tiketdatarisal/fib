package ursa

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/tiket/DATA-RANGERS-URSA-BE/core_module/pkg/pb"
	"github.com/tiketdatarisal/stdres"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net/mail"
	"strings"
	"time"
)

type ursa struct {
	conn   *grpc.ClientConn
	auth   pb.AuthorizationClient
	inbox  pb.InboxClient
	users  pb.UsersClient
	groups pb.GroupsClient
}

// Auth authenticate user credential using URSA.
func (u *ursa) Auth(scopes ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		const bearer = "bearer"

		auth := c.Get(fiber.HeaderAuthorization)
		lb := len(bearer)
		if len(auth) <= lb {
			return c.Status(fiber.StatusUnauthorized).
				JSON(stdres.NewUnauthorizedResponse(ErrUserUnauthorized.Error(), nil))
		}

		if strings.ToLower(auth[:lb]) != bearer {
			return c.Status(fiber.StatusUnauthorized).
				JSON(stdres.NewUnauthorizedResponse(ErrUserUnauthorized.Error(), nil))
		}

		token := auth[lb+1:]
		res, err := u.auth.Authorize(context.Background(), &pb.AuthorizeReq{Token: token, AllowedScope: scopes})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).
				JSON(stdres.NewInternalServerErrorResponse(err.Error(), nil))
		}

		if res != nil {
			if res.Httpcode != 200 {
				return c.Status(int(res.Httpcode)).
					JSON(stdres.NewResponseBody(int(res.Httpcode), res.Message, nil))
			}

			if strings.ToLower(res.User.Status) != userStatusActive {
				return c.Status(fiber.StatusForbidden).
					JSON(stdres.NewForbiddenResponse(ErrUserDeactivated.Error(), nil))
			}
		}

		groups := Groups{}
		for _, group := range res.User.CombinedGroups {
			groups = append(groups, NewGroup(
				group.Id,
				group.Name,
			))
		}

		user := NewUser(
			res.User.Email,
			res.User.FullName,
			res.User.GoogleToken,
			res.User.OrganizationName,
			res.User.Status,
			groups,
		)

		c.Locals(authenticatedUser, &user)
		return c.Next()
	}
}

// New return a new URSA middleware.
func New(config ...Config) Authenticator {
	cfg := configDefault(config...)

	mutex.Lock()
	if u == nil || u.conn == nil {
		conn, err := grpc.Dial(cfg.Host, grpc.WithTransportCredentials(
			credentials.NewClientTLSFromCert(
				nil,
				"")))
		if err != nil {
			panic(err)
		}

		u = &ursa{
			conn:   conn,
			auth:   pb.NewAuthorizationClient(conn),
			inbox:  pb.NewInboxClient(conn),
			users:  pb.NewUsersClient(conn),
			groups: pb.NewGroupsClient(conn),
		}
	}
	mutex.Unlock()

	return u
}

// GetAuthenticator returns URSA auth middleware.
func GetAuthenticator() Authenticator {
	if u == nil {
		return &DefaultAuthenticator{}
	}

	return u
}

// Close cleanup all resources used by URSA.
func Close() {
	if u != nil {
		_ = u.conn.Close()
	}
}

// Conn return GRPC connection for URSA.
func Conn() *grpc.ClientConn {
	if u == nil || u.conn == nil {
		return nil
	}

	return u.conn
}

// Ping return error when failed.
func Ping(duration ...time.Duration) error {
	if u == nil || u.auth == nil {
		return ErrUrsaNotInitialized
	}

	d := 3 * time.Second
	if len(duration) > 0 && duration[0] > 0 {
		d = duration[0]
	}

	ctx, cancel := context.WithTimeout(bg, d)
	defer cancel()

	_, err := u.auth.Ping(ctx, &pb.PingMessage{})
	if err != nil {
		return ErrHostCannotBeReached
	}

	return nil
}

// GetUsers returns a list of users.
func GetUsers(param GetUsersParam) ([]User, error) {
	if u == nil || u.users == nil {
		return nil, ErrUrsaNotInitialized
	}

	res, err := u.users.GetAllUsers(bg, &pb.GetAllUsersReq{
		Keyword:          param.Keyword,
		OrganizationName: param.Organization,
		Status: func() string {
			if param.IsDeactivated {
				return userStatusInactive
			} else {
				return userStatusActive
			}
		}(),
	})
	if err != nil {
		return nil, err
	}

	var users []User
	for _, user := range res.Users {
		if _, err = mail.ParseAddress(user.Email); err != nil {
			continue
		}

		if strings.TrimSpace(user.FullName) == "" {
			continue
		}

		groups := Groups{}
		for _, group := range user.CombinedGroups {
			groups = append(groups, NewGroup(
				group.Id,
				group.Name,
			))
		}

		users = append(users, NewUser(
			user.Email,
			user.FullName,
			user.GoogleToken,
			user.OrganizationName,
			user.Status,
			groups,
		))
	}

	return users, nil
}

// GetUser return a user by its email.
func GetUser(email string) (*User, error) {
	if u == nil || u.users == nil {
		return nil, ErrUrsaNotInitialized
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, ErrEmailNotValid
	}

	res, err := u.users.GetUser(bg, &pb.GetUserReq{
		Email: email,
	})
	if err != nil {
		return nil, err
	}

	groups := Groups{}
	for _, group := range res.User.CombinedGroups {
		groups = append(groups, NewGroup(
			group.Id,
			group.Name,
		))
	}

	user := NewUser(
		res.User.Email,
		res.User.FullName,
		res.User.GoogleToken,
		res.User.OrganizationName,
		res.User.Status,
		groups,
	)

	return &user, nil
}

// GetAuthenticateUser return an authenticated user for active context.
func GetAuthenticateUser(c *fiber.Ctx) *User {
	val := c.Locals(authenticatedUser)
	if val == nil {
		return nil
	}

	user, ok := val.(*User)
	if !ok || user == nil {
		return nil
	}

	return user
}
