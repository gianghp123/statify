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

	if status != "" {
		db = db.Where("projects.effective_status = ?", status)
	}

	// 4. Count
	// Using Scan is safer than Count() to avoid GORM mapping issues with complex WHERE clauses
	var totalCount int64
	if err := db.Select("count(distinct projects.id)").Scan(&totalCount).Error; err != nil {
		return result, err
	}
	result.Pagination.TotalCount = totalCount

	// 5. Retrieve Data
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
	return r.db.WithContext(ctx).Delete(project).Error
}

func (r *ProjectRepository) MarkReady(ctx context.Context, projectID uint, deploymentID uint) error {
	// Logic:
	// Only update EffectiveStatus if project is not currently pinned/live AND this deployment is the Latest.
	// We use a raw SQL UPDATE with WHERE clause to do this atomically and efficiently.

	return r.db.WithContext(ctx).Model(&models.Project{}).
		Where("id = ? AND latest_deployment_id = ? AND current_deployment_id = 0", projectID, deploymentID).
		Update("effective_status", enums.DeploymentStatusReady).Error
}

func (r *ProjectRepository) MarkFailed(ctx context.Context, projectID uint, deploymentID uint) error {
	// Logic: Same as MarkReady
	return r.db.WithContext(ctx).Model(&models.Project{}).
		Where("id = ? AND latest_deployment_id = ? AND current_deployment_id = 0", projectID, deploymentID).
		Update("effective_status", enums.DeploymentStatusFailed).Error
}

func (r *ProjectRepository) WithTx(tx *gorm.DB) IProjectRepository {
	if tx == nil {
		return r
	}
	return &ProjectRepository{db: tx}
}
