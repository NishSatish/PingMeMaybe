package models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserCohortType string

const (
	CohortNonPremium        UserCohortType = "NON_PREMIUM"
	CohortActivePremium     UserCohortType = "ACTIVE_PREMIUM"
	CohortPremiumNearExpiry UserCohortType = "PREMIUM_NEAR_EXPIRY"
	CohortExpiredPremium    UserCohortType = "EXPIRED_PREMIUM"
)

// The schema is this bloated because we have a pre-seeded user table,
// and I'm too lazy to either make a separate user schema at the service level or a separate Cohorts table

// Will refactor in a future iteration to make this entire project an open-source module
type UserCohort struct {
	UserID            int            `json:"user_id"`
	Email             string         `json:"email"`
	Username          string         `json:"username"`
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_name"`
	SubscriptionTier  string         `json:"subscription_tier"`
	SubscriptionStart *time.Time     `json:"subscription_start_date"`
	SubscriptionEnd   *time.Time     `json:"subscription_end_date"`
	CohortType        UserCohortType `json:"cohort_type"`
	DaysUntilExpiry   *int           `json:"days_until_expiry"`
	LastLoginAt       *time.Time     `json:"last_login_at"`
	NotificationPrefs string         `json:"notification_preferences"`
	Timezone          string         `json:"timezone"`
}

type CohortFilters struct {
	SubscriptionTier *string          `json:"subscription_tier,omitempty"`
	CohortTypes      []UserCohortType `json:"cohort_types,omitempty"`
	Timezone         *string          `json:"timezone,omitempty"`
	IsActive         *bool            `json:"is_active,omitempty"`
	Limit            int              `json:"limit"`
	Offset           int              `json:"offset"`
}

type CohortStats struct {
	CohortType UserCohortType `json:"cohort_type"`
	Count      int            `json:"count"`
	Percentage float64        `json:"percentage"`
}

type UserCohortRepo struct {
	DB *pgxpool.Pool
}

type IUserCohortRepository interface {
	GetCohortUsers(ctx context.Context, cohortType UserCohortType, filters *CohortFilters) ([]UserCohort, error)
	GetCohortStats(ctx context.Context) ([]CohortStats, error)
	GetUsersNearExpiry(ctx context.Context, daysThreshold int) ([]UserCohort, error)
	GetUsersByCohorts(ctx context.Context, cohortTypes []UserCohortType, limit int) ([]UserCohort, error)
	GetCohortUserCount(ctx context.Context, cohortType UserCohortType) (int, error)
}

func NewUserCohortRepo(db *pgxpool.Pool) IUserCohortRepository {
	return &UserCohortRepo{
		DB: db,
	}
}

// GetCohortUsers returns users from a specific cohort with optional filters
func (r *UserCohortRepo) GetCohortUsers(ctx context.Context, cohortType UserCohortType, filters *CohortFilters) ([]UserCohort, error) {
	baseQuery := `
		SELECT 
			u.id,
			u.email,
			u.username,
			u.first_name,
			u.last_name,
			u.subscription_tier,
			u.subscription_start_date,
			u.subscription_end_date,
			u.last_login_at,
			u.notification_preferences,
			u.timezone,
			CASE 
				WHEN u.subscription_end_date IS NOT NULL 
				THEN EXTRACT(DAY FROM (u.subscription_end_date - CURRENT_DATE))::int
				ELSE NULL 
			END as days_until_expiry
		FROM users u
		WHERE u.is_active = true`

	var whereConditions []string
	var args []interface{}
	argIndex := 1

	// Add cohort-specific conditions
	switch cohortType {
	case CohortNonPremium:
		whereConditions = append(whereConditions, "u.is_premium_user = false")
		break
	case CohortActivePremium:
		whereConditions = append(whereConditions,
			"u.is_premium_user = true",
			"(u.subscription_end_date IS NULL OR u.subscription_end_date > CURRENT_DATE)")
		break
	case CohortPremiumNearExpiry:
		whereConditions = append(whereConditions,
			"u.is_premium_user = true",
			"u.subscription_end_date IS NOT NULL",
			"u.subscription_end_date > CURRENT_DATE",
			"u.subscription_end_date <= CURRENT_DATE + INTERVAL '30 days'")
		break
	case CohortExpiredPremium:
		whereConditions = append(whereConditions,
			"u.subscription_end_date IS NOT NULL",
			"u.subscription_end_date < CURRENT_DATE")
		break
	default:
		break // Just get all active users. These other cases are just templates, the whole query can be customised with the filters param.
	}

	// Add optional filters
	// Manually inject filters into the query, will refactor later to use a more structured approach
	if filters != nil {
		if filters.SubscriptionTier != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("u.subscription_tier = $%d", argIndex))
			args = append(args, *filters.SubscriptionTier)
			argIndex++
		}

		if filters.Timezone != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("u.timezone = $%d", argIndex))
			args = append(args, *filters.Timezone)
			argIndex++
		}

		if filters.IsActive != nil {
			whereConditions = append(whereConditions, fmt.Sprintf("u.is_active = $%d", argIndex))
			args = append(args, *filters.IsActive)
			argIndex++
		}
	}

	// Build final query
	query := baseQuery
	if len(whereConditions) > 0 {
		query += " AND " + fmt.Sprintf("(%s)", whereConditions[0])
		for _, condition := range whereConditions[1:] {
			query += " AND " + condition
		}
	}

	query += " ORDER BY u.created_at DESC"

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query cohort users: %w", err)
	}
	defer rows.Close()

	var users []UserCohort
	for rows.Next() {
		var user UserCohort
		err := rows.Scan(
			&user.UserID,
			&user.Email,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.SubscriptionTier,
			&user.SubscriptionStart,
			&user.SubscriptionEnd,
			&user.LastLoginAt,
			&user.NotificationPrefs,
			&user.Timezone,
			&user.DaysUntilExpiry,
		)
		if err != nil {
			fmt.Printf("Error scanning cohort user: %v\n", err)
			continue
		}
		user.CohortType = cohortType
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserCohortRepo) GetCohortStats(ctx context.Context) ([]CohortStats, error) {
	query := `
		WITH cohort_counts AS (
			SELECT 
				CASE 
					WHEN is_premium_user = false THEN 'NON_PREMIUM'
					WHEN is_premium_user = true AND (subscription_end_date IS NULL OR subscription_end_date > CURRENT_DATE) THEN 'ACTIVE_PREMIUM'
					WHEN is_premium_user = true AND subscription_end_date IS NOT NULL 
						 AND subscription_end_date > CURRENT_DATE 
						 AND subscription_end_date <= CURRENT_DATE + INTERVAL '30 days' THEN 'PREMIUM_NEAR_EXPIRY'
					WHEN subscription_end_date IS NOT NULL AND subscription_end_date < CURRENT_DATE THEN 'EXPIRED_PREMIUM'
					ELSE 'OTHER'
				END as cohort_type,
				COUNT(*) as count
			FROM users 
			WHERE is_active = true
			GROUP BY 1
		),
		total_users AS (
			SELECT COUNT(*) as total FROM users WHERE is_active = true
		)
		SELECT 
			cc.cohort_type,
			cc.count,
			ROUND((cc.count * 100.0 / tu.total), 2) as percentage
		FROM cohort_counts cc
		CROSS JOIN total_users tu
		ORDER BY cc.count DESC`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get cohort stats: %w", err)
	}
	defer rows.Close()

	var stats []CohortStats
	for rows.Next() {
		var stat CohortStats
		var cohortTypeStr string
		err := rows.Scan(&cohortTypeStr, &stat.Count, &stat.Percentage)
		if err != nil {
			fmt.Printf("Error scanning cohort stats: %v\n", err)
			continue
		}
		stat.CohortType = UserCohortType(cohortTypeStr)
		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

func (r *UserCohortRepo) GetUsersNearExpiry(ctx context.Context, daysThreshold int) ([]UserCohort, error) {
	query := `
		SELECT 
			u.id, u.email, u.username, u.first_name, u.last_name,
			u.subscription_tier, u.subscription_start_date, u.subscription_end_date,
			u.last_login_at, u.notification_preferences, u.timezone,
			EXTRACT(DAY FROM (u.subscription_end_date - CURRENT_DATE))::int as days_until_expiry
		FROM users u
		WHERE u.is_active = true 
		AND u.is_premium_user = true
		AND u.subscription_end_date IS NOT NULL
		AND u.subscription_end_date > CURRENT_DATE
		AND u.subscription_end_date <= CURRENT_DATE + INTERVAL '%d days'
		ORDER BY u.subscription_end_date ASC`

	formattedQuery := fmt.Sprintf(query, daysThreshold)
	rows, err := r.DB.Query(ctx, formattedQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get users near expiry: %w", err)
	}
	defer rows.Close()

	var users []UserCohort
	for rows.Next() {
		var user UserCohort
		err := rows.Scan(
			&user.UserID, &user.Email, &user.Username, &user.FirstName, &user.LastName,
			&user.SubscriptionTier, &user.SubscriptionStart, &user.SubscriptionEnd,
			&user.LastLoginAt, &user.NotificationPrefs, &user.Timezone,
			&user.DaysUntilExpiry,
		)
		if err != nil {
			fmt.Printf("Error scanning user near expiry: %v\n", err)
			continue
		}
		user.CohortType = CohortPremiumNearExpiry
		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserCohortRepo) GetUsersByCohorts(ctx context.Context, cohortTypes []UserCohortType, limit int) ([]UserCohort, error) {
	var allUsers []UserCohort

	for _, cohortType := range cohortTypes {
		filters := &CohortFilters{
			Limit: limit,
		}

		users, err := r.GetCohortUsers(ctx, cohortType, filters)
		if err != nil {
			fmt.Printf("Error getting users for cohort %s: %v\n", cohortType, err)
			continue
		}

		allUsers = append(allUsers, users...)
	}

	return allUsers, nil
}

func (r *UserCohortRepo) GetCohortUserCount(ctx context.Context, cohortType UserCohortType) (int, error) {
	var query string

	switch cohortType {
	case CohortNonPremium:
		query = "SELECT COUNT(*) FROM users WHERE is_active = true AND is_premium_user = false"
	case CohortActivePremium:
		query = `SELECT COUNT(*) FROM users WHERE is_active = true AND is_premium_user = true 
				 AND (subscription_end_date IS NULL OR subscription_end_date > CURRENT_DATE)`
	case CohortPremiumNearExpiry:
		query = `SELECT COUNT(*) FROM users WHERE is_active = true AND is_premium_user = true 
				 AND subscription_end_date IS NOT NULL AND subscription_end_date > CURRENT_DATE 
				 AND subscription_end_date <= CURRENT_DATE + INTERVAL '30 days'`
	case CohortExpiredPremium:
		query = `SELECT COUNT(*) FROM users WHERE is_active = true 
				 AND subscription_end_date IS NOT NULL AND subscription_end_date < CURRENT_DATE`
	default:
		return 0, fmt.Errorf("unknown cohort type: %s", cohortType)
	}

	var count int
	err := r.DB.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get cohort count: %w", err)
	}

	return count, nil
}
