package use_case_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/VibeTeam/fitness-tracker-backend/user/auth"
	"github.com/VibeTeam/fitness-tracker-backend/user/models"
	"github.com/VibeTeam/fitness-tracker-backend/user/use_case"
)

/* -------------------------------------------------------------------------- */
/* In‑memory UserRepository stub – now includes Count so it satisfies the     */
/* full interface.                                                            */
/* -------------------------------------------------------------------------- */

type inMemUserRepo struct {
	mu    sync.Mutex
	store map[uint]*models.User
	next  uint
}

func newRepo() *inMemUserRepo {
	return &inMemUserRepo{store: make(map[uint]*models.User), next: 1}
}

func (r *inMemUserRepo) Create(_ context.Context, u *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u.ID = r.next
	r.next++
	cp := *u
	r.store[u.ID] = &cp
	return nil
}

func (r *inMemUserRepo) GetByEmail(_ context.Context, email string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range r.store {
		if u.Email == email {
			cp := *u
			return &cp, nil
		}
	}
	return nil, errors.New("not found")
}

func (r *inMemUserRepo) GetByID(_ context.Context, id uint) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if u, ok := r.store[id]; ok {
		cp := *u
		return &cp, nil
	}
	return nil, errors.New("not found")
}

func (r *inMemUserRepo) List(_ context.Context, _, _ int) ([]*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]*models.User, 0, len(r.store))
	for _, u := range r.store {
		cp := *u
		out = append(out, &cp)
	}
	return out, nil
}

func (r *inMemUserRepo) Update(_ context.Context, u *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	cp := *u
	r.store[u.ID] = &cp
	return nil
}

func (r *inMemUserRepo) Delete(_ context.Context, id uint) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, id)
	return nil
}

func (r *inMemUserRepo) Count(_ context.Context) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.store), nil
}

/* -------------------------------------------------------------------------- */

func makeService() *use_case.AuthService {
	repo := newRepo()
	tm := auth.NewManager("access‑key", "refresh‑key", time.Minute, time.Hour)
	return use_case.NewAuthService(repo, tm)
}

/* -------------------------------------------------------------------------- */
/* Tests                                                                      */
/* -------------------------------------------------------------------------- */

func TestRegisterAndLogin(t *testing.T) {
	svc := makeService()
	ctx := context.Background()

	// Register a user
	access, refresh, err := svc.Register(ctx, "x@y.com", "pwd")
	require.NoError(t, err)
	require.NotEmpty(t, access)
	require.NotEmpty(t, refresh)

	// Login with same credentials – should succeed and return tokens
	access2, refresh2, err := svc.Login(ctx, "x@y.com", "pwd")
	require.NoError(t, err)
	require.NotEmpty(t, access2)
	require.NotEmpty(t, refresh2)
}

func TestRefreshAndValidate(t *testing.T) {
	svc := makeService()
	ctx := context.Background()

	_, refresh, _ := svc.Register(ctx, "a@b.c", "pwd")

	// Refresh the pair
	newAcc, newRef, err := svc.Refresh(ctx, refresh)
	require.NoError(t, err)
	require.NotEmpty(t, newAcc)
	require.NotEmpty(t, newRef)

	// Validate new access token – should yield user ID 1
	uid, err := svc.Validate(ctx, newAcc)
	require.NoError(t, err)
	require.Equal(t, uint(1), uid)
}

func TestPasswordHashStored(t *testing.T) {
	repo := newRepo()
	tm := auth.NewManager("a", "b", time.Minute, time.Hour)
	svc := use_case.NewAuthService(repo, tm)

	_, _, err := svc.Register(context.Background(), "p@q.r", "secret")
	require.NoError(t, err)

	u, _ := repo.GetByEmail(context.Background(), "p@q.r")
	require.NotEmpty(t, u.PasswordHash)
	require.NotEqual(t, "secret", u.PasswordHash)
	require.NoError(t, bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte("secret")))
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc := makeService()
	_, _, _ = svc.Register(context.Background(), "dup@e.com", "pwd")

	_, _, err := svc.Register(context.Background(), "dup@e.com", "pwd")
	require.ErrorIs(t, err, use_case.ErrEmailAlreadyUsed)
}
