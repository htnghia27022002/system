package dependency

import (
	"be/internal/repository"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
)

func newAuthRepository(db *postgres.Postgres) interfaces.AuthRepository {
	return repository.NewAuthRepository(db)
}

func newUserRepository(db *postgres.Postgres) interfaces.UserRepository {
	return repository.NewUserRepository(db)
}

func newRoleRepository(db *postgres.Postgres) interfaces.RoleRepository {
	return repository.NewRoleRepository(db)
}

func newPermissionRepository(db *postgres.Postgres) interfaces.PermissionRepository {
	return repository.NewPermissionRepository(db)
}

func newWebhookInboxRepository(db *postgres.Postgres) interfaces.WebhookInboxRepository {
	return repository.NewWebhookInboxRepository(db)
}

func newWebhookRequestRepository(db *postgres.Postgres) interfaces.WebhookRequestRepository {
	return repository.NewWebhookRequestRepository(db)
}

func newCountryRepository(db *postgres.Postgres) interfaces.CountryRepository {
	return repository.NewCountryRepository(db)
}

func newDivisionRepository(db *postgres.Postgres) interfaces.DivisionRepository {
	return repository.NewDivisionRepository(db)
}

func newLocationRepository(db *postgres.Postgres) interfaces.LocationRepository {
	return repository.NewLocationRepository(db)
}

func newPlaceRepository(db *postgres.Postgres) interfaces.PlaceRepository {
	return repository.NewPlaceRepository(db)
}

func newCategoryRepository(db *postgres.Postgres) interfaces.CategoryRepository {
	return repository.NewCategoryRepository(db)
}

func newNewsRepository(db *postgres.Postgres) interfaces.NewsRepository {
	return repository.NewNewsRepository(db)
}

func newDataSourceRepository(db *postgres.Postgres) interfaces.DataSourceRepository {
	return repository.NewDataSourceRepository(db)
}

func newIngestRunRepository(db *postgres.Postgres) interfaces.IngestRunRepository {
	return repository.NewIngestRunRepository(db)
}
