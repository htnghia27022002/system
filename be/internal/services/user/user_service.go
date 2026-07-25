package user

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"
	"time"

	apperrors "be/internal/common/errors"
	"be/pkg/hash"
	"be/pkg/query"
	userdto "be/internal/dto/user"
	usermodel "be/internal/models/user"
	searchpkg "be/internal/search"
	"be/internal/repository/interfaces"
	"be/internal/services/media"
)

type OutboxEnqueuer interface {
	EnqueueUpsert(ctx context.Context, entityType, entityID string) error
	EnqueueDelete(ctx context.Context, entityType, entityID string) error
}

type Service struct {
	repo   interfaces.UserRepository
	auth   interfaces.AuthRepository
	outbox OutboxEnqueuer
	media  *media.Service
}

func NewService(
	repo interfaces.UserRepository,
	auth interfaces.AuthRepository,
	outbox OutboxEnqueuer,
	mediaSvc *media.Service,
) *Service {
	return &Service{repo: repo, auth: auth, outbox: outbox, media: mediaSvc}
}

func (s *Service) Create(ctx context.Context, req userdto.CreateUserRequest) (*usermodel.User, error) {
	existing, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: email already exists", apperrors.ErrConflict)
	}

	hashed, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	normalized, err := userdto.ValidatePersonalFields(req.Phone, req.General, req.Birthday, req.Address, req.SocialLinks)
	if err != nil {
		return nil, err
	}

	status := usermodel.StatusActive
	if req.Status != "" {
		status = usermodel.Status(req.Status)
	}

	user := &usermodel.User{
		Email:        req.Email,
		PasswordHash: hashed,
		FullName:     strings.TrimSpace(req.Name),
		RoleID:       req.RoleID,
		Status:       status,
		SocialLinks:  []usermodel.SocialLink{},
	}
	normalized.Apply(user)
	if req.AvatarURL != nil {
		user.AvatarURL = strings.TrimSpace(*req.AvatarURL)
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	s.enqueueUpsert(ctx, searchpkg.EntityUser, user.ID)
	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*usermodel.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}
	return user, nil
}

func (s *Service) List(ctx context.Context, form userdto.ListUsersQuery) ([]usermodel.User, int64, int, int, error) {
	q := query.New(form.Page, form.PageSize).
		OrderBy("created_at DESC").
		WhereEqual("role_id", form.RoleID).
		WhereEqual("id", form.ID).
		WhereEqual("status", form.Status).
		WhereLikeAny([]string{"email", "full_name"}, form.Search)

	users, total, err := s.repo.List(ctx, q)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	return users, total, q.Page, q.PageSize, nil
}

func (s *Service) Update(ctx context.Context, id string, req userdto.UpdateUserRequest, sessionUserID string) (*usermodel.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}

	normalized, err := userdto.ValidatePersonalFields(req.Phone, req.General, req.Birthday, req.Address, req.SocialLinks)
	if err != nil {
		return nil, err
	}

	if req.Email != nil && *req.Email != user.Email {
		existing, err := s.repo.GetByEmail(ctx, *req.Email)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, fmt.Errorf("%w: email already exists", apperrors.ErrConflict)
		}
		user.Email = *req.Email
	}
	if req.Name != nil {
		user.FullName = strings.TrimSpace(*req.Name)
	}
	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}
	if req.Status != nil {
		if id == sessionUserID && *req.Status == string(usermodel.StatusInactive) {
			return nil, fmt.Errorf("%w: cannot deactivate your own account", apperrors.ErrForbidden)
		}
		user.Status = usermodel.Status(*req.Status)
	}
	if req.Password != nil && *req.Password != "" {
		hashed, err := hash.HashPassword(*req.Password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = hashed
	}
	normalized.Apply(user)
	if req.AvatarURL != nil {
		user.AvatarURL = strings.TrimSpace(*req.AvatarURL)
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	s.enqueueUpsert(ctx, searchpkg.EntityUser, user.ID)
	return user, nil
}

func (s *Service) UploadAvatar(ctx context.Context, id string, file multipart.File, header *multipart.FileHeader) (*usermodel.User, error) {
	if s.media == nil {
		return nil, fmt.Errorf("%w: media storage is not configured", apperrors.ErrBadRequest)
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrNotFound
	}

	publicPath, err := s.media.SaveAvatar(file, header)
	if err != nil {
		return nil, err
	}

	prev := user.AvatarURL
	user.AvatarURL = publicPath
	if err := s.repo.Update(ctx, user); err != nil {
		s.media.DeleteByPublicPath(publicPath)
		return nil, err
	}
	if prev != "" && prev != publicPath {
		s.media.DeleteByPublicPath(prev)
	}
	s.enqueueUpsert(ctx, searchpkg.EntityUser, user.ID)
	return user, nil
}

func (s *Service) Delete(ctx context.Context, id, sessionUserID string) error {
	if id == sessionUserID {
		return fmt.Errorf("%w: cannot delete your own account", apperrors.ErrForbidden)
	}
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.ErrNotFound
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.enqueueDelete(ctx, searchpkg.EntityUser, id)
	return nil
}

func (s *Service) enqueueUpsert(ctx context.Context, entityType, entityID string) {
	if s.outbox == nil {
		return
	}
	_ = s.outbox.EnqueueUpsert(ctx, entityType, entityID)
}

func (s *Service) enqueueDelete(ctx context.Context, entityType, entityID string) {
	if s.outbox == nil {
		return
	}
	_ = s.outbox.EnqueueDelete(ctx, entityType, entityID)
}

func ToResponse(user *usermodel.User, oauthProviders ...string) userdto.UserResponse {
	links := user.SocialLinks
	if links == nil {
		links = []usermodel.SocialLink{}
	}
	providers := oauthProviders
	if providers == nil {
		providers = []string{}
	}
	return userdto.UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		Name:           user.FullName,
		RoleID:         user.RoleID,
		Status:         string(user.Status),
		Phone:          user.Phone,
		AvatarURL:      user.AvatarURL,
		General:        user.General,
		Birthday:       userdto.FormatBirthday(user.Birthday),
		Address:        user.Address,
		SocialLinks:    userdto.SocialLinksToDTO(links),
		OAuthProviders: providers,
		CreatedAt:      user.CreatedAt.UTC().Format(time.RFC3339),
	}
}

// ResponsesForUsers builds UserResponse rows and batch-loads OAuth providers.
func (s *Service) ResponsesForUsers(ctx context.Context, users []usermodel.User) ([]userdto.UserResponse, error) {
	providersByUser := map[string][]string{}
	if s.auth != nil && len(users) > 0 {
		ids := make([]string, 0, len(users))
		for i := range users {
			ids = append(ids, users[i].ID)
		}
		var err error
		providersByUser, err = s.auth.ListProvidersByUserIDs(ctx, ids)
		if err != nil {
			return nil, err
		}
	}
	items := make([]userdto.UserResponse, 0, len(users))
	for i := range users {
		items = append(items, ToResponse(&users[i], providersByUser[users[i].ID]...))
	}
	return items, nil
}

// ResponseForUser builds one UserResponse including linked OAuth providers.
func (s *Service) ResponseForUser(ctx context.Context, user *usermodel.User) (userdto.UserResponse, error) {
	if user == nil {
		return userdto.UserResponse{}, apperrors.ErrNotFound
	}
	items, err := s.ResponsesForUsers(ctx, []usermodel.User{*user})
	if err != nil {
		return userdto.UserResponse{}, err
	}
	if len(items) == 0 {
		return ToResponse(user), nil
	}
	return items[0], nil
}
