package repository

import (
	"context"

	"github.com/gianghp/statify/internal/core/enums"
	"github.com/gianghp/statify/internal/core/repository"
	"github.com/gianghp/statify/internal/database/models"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) FindBySubdomain(ctx context.Context, subdomain string) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, "subdomain = ?", subdomain).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindByID(ctx context.Context, id uint) (*models.Project, error) {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, id).Error; err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindAllByUserID(ctx context.Context, userID uint, page int, limit int, status enums.DeploymentStatus) (repository.PaginatedEntities[*models.Project], error) {
	// 1. Initialize result
	result := repository.PaginatedEntities[*models.Project]{
		Entities: []*models.Project{},
		Pagination: repository.Pagination{
			TotalCount: 0,
			Page:       page,
			Limit:      limit,
		},
	}

	db := r.db.WithContext(ctx).Model(&models.Project{}).Where("projects.user_id = ?", userID)

	// 3. Conditional Status Filter
	if status != "" {
		// Logic:
		// 1. IF current_deployment_id > 0: Check if the deployment with that ID has the specific status.
		// 2. OR IF current_deployment_id == 0: Check if the LATEST deployment for this project has the specific status.
		db = db.Where(`
			? = (
				CASE 
					-- PRIORITY 1: If the project is currently running (Pinned ID is set AND its status is LIVE)
					-- Then the project status is forcefully 'LIVE', regardless of newer commits.
					WHEN projects.current_deployment_id > 0 AND (
						SELECT d.status 
						FROM deployments d 
						WHERE d.id = projects.current_deployment_id
					) = 'LIVE' THEN 'LIVE'

					-- PRIORITY 2: Otherwise, the status is determined by the strictly LATEST deployment.
					-- (e.g., READY, FAILED, BUILDING, or a LIVE deployment that isn't pinned yet)
					ELSE (
						SELECT d.status 
						FROM deployments d 
						WHERE d.project_id = projects.id 
						ORDER BY d.created_at DESC 
						LIMIT 1
					)
				END
			)
		`, status)
	}

	// 4. Count
	// Using Scan is safer than Count() to avoid GORM mapping issues with complex WHERE clauses
	var totalCount int64
	if err := db.Select("count(distinct projects.id)").Scan(&totalCount).Error; err != nil {
		return result, err
	}
	result.Pagination.TotalCount = totalCount

	// 5. Retrieve Data
	// Now we fetch the actual entities with Preload
	if err := db.Preload("Deployments").
		Distinct("projects.*").
		Order("projects.created_at DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&result.Entities).Error; err != nil {
		return result, err
	}

	return result, nil
}

func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *ProjectRepository) Update(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *ProjectRepository) Delete(ctx context.Context, project *models.Project) error {
	return r.db.WithContext(ctx).Unscoped().Delete(project).Error
}

func (r *ProjectRepository) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *ProjectRepository) WithTx(tx *gorm.DB) IProjectRepository {
	if tx == nil {
		return r
	}
	return &ProjectRepository{db: tx}
}
