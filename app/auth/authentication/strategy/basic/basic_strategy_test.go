package basic_test

import (
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	kallax "gopkg.in/src-d/go-kallax.v1"

	"github.com/stretchr/testify/require"

	"github.com/ifreddyrondon/capture/app/auth/authentication/strategy/basic"
	"github.com/ifreddyrondon/capture/app/user"
	"github.com/ifreddyrondon/capture/database"
	"github.com/jinzhu/gorm"
)

var (
	once sync.Once
	db   *gorm.DB
)

const (
	testUserEmail    = "test@example.com"
	testUserPassword = "b4KeHAYy3u9v=ZQX"
)

func getDB() *gorm.DB {
	once.Do(func() {
		ds := database.Open("postgres://localhost/captures_app_test?sslmode=disable")
		db = ds.DB
	})
	return db
}

func setup(t *testing.T) (*basic.Basic, func()) {
	userStore := user.NewPGStore(getDB().Table("basic_auth-users"))
	userStore.Migrate()
	teardown := func() { userStore.Drop() }
	userService := user.NewService(userStore)

	// save a user to test
	u := user.User{Email: testUserEmail}
	err := u.SetPassword(testUserPassword)
	require.Nil(t, err)
	userService.Save(&u)

	return basic.New(userService), teardown
}

func TestValidateSuccess(t *testing.T) {
	strategy, teardown := setup(t)
	defer teardown()

	body := strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, testUserEmail, testUserPassword))
	req := httptest.NewRequest("GET", "/", body)

	u, err := strategy.Validate(req)
	assert.Nil(t, err)
	assert.Equal(t, testUserEmail, u.Email)
}

func TestValidateInvalidCredentials(t *testing.T) {
	strategy, teardown := setup(t)
	defer teardown()

	tt := []struct {
		name string
		body io.Reader
		errs []string
	}{
		{
			name: "invalid credentials",
			body: strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, testUserEmail, "123")),
			errs: []string{"invalid email or password"},
		},
		{
			name: "missing email",
			body: strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, "bla@example.com", "123")),
			errs: []string{"invalid email or password"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", tc.body)
			_, err := strategy.Validate(req)
			assert.Error(t, err)
			assert.True(t, strategy.IsErrCredentials(err))
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}

type MockUserGetterServiceFail struct{}

func (m *MockUserGetterServiceFail) GetByEmail(email string) (*user.User, error) {
	return nil, errors.New("test")
}

func (m *MockUserGetterServiceFail) GetByID(id kallax.ULID) (*user.User, error) {
	return nil, errors.New("test")
}

func TestValidateFailsDecoding(t *testing.T) {
	t.Parallel()

	userService := &MockUserGetterServiceFail{}
	strategy := basic.New(userService)

	tt := []struct {
		name string
		body io.Reader
		errs []string
	}{
		{
			name: "invalid json",
			body: strings.NewReader("{"),
			errs: []string{"unexpected EOF"},
		},
		{
			name: "missing data",
			body: strings.NewReader("{}"),
			errs: []string{"email must not be blank", "password must not be blank"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", tc.body)
			_, err := strategy.Validate(req)
			assert.Error(t, err)
			assert.True(t, strategy.IsErrDecoding(err))
			for _, v := range tc.errs {
				assert.Contains(t, err.Error(), v)
			}
		})
	}
}

func TestValidateFailsUnknowErr(t *testing.T) {
	t.Parallel()

	userService := &MockUserGetterServiceFail{}
	strategy := basic.New(userService)
	body := strings.NewReader(fmt.Sprintf(`{"email":"%v","password":"%v"}`, testUserEmail, testUserPassword))
	req := httptest.NewRequest("GET", "/", body)
	_, err := strategy.Validate(req)
	assert.EqualError(t, err, "test")
	assert.False(t, strategy.IsErrCredentials(err))
	assert.False(t, strategy.IsErrDecoding(err))
}
