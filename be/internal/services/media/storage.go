package media

import "mime/multipart"

// AvatarStorage is the persistence port used by auth/user services (DIP).
type AvatarStorage interface {
	SaveAvatar(file multipart.File, header *multipart.FileHeader) (string, error)
	DeleteByPublicPath(publicPath string)
}

var _ AvatarStorage = (*Service)(nil)
