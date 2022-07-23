package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/phyrwork/benevolent-dictator/pkg/api/auth"
	"github.com/phyrwork/benevolent-dictator/pkg/api/database"
	"github.com/phyrwork/benevolent-dictator/pkg/api/graph/generated"
	"github.com/phyrwork/benevolent-dictator/pkg/api/graph/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserCreate is the resolver for the userCreate field.
func (r *mutationResolver) UserCreate(ctx context.Context, name string, email string, password string) (*model.User, error) {
	// Prepare password.
	key, salt, err := auth.Encode([]byte(password))
	if err != nil {
		return nil, fmt.Errorf("password encode error: %v", err)
	}
	// Create user.
	row := database.User{
		Name:  name,
		Email: email,
		Key:   key,
		Salt:  salt,
	}
	// TODO: check for unique email first?
	if err := r.DB.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	return &model.User{
		ID:   row.ID,
		Name: row.Name,
	}, nil
}

// UserUpdate is the resolver for the userUpdate field.
func (r *mutationResolver) UserUpdate(ctx context.Context, name *string) (*model.User, error) {
	userAuth := auth.ForContext(ctx)
	if userAuth == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	user := database.User{
		ID: userAuth.UserID,
	}
	do := false
	if name != nil {
		user.Name = *name
		do = true
	}
	if do {
		if err := r.DB.WithContext(ctx).Updates(&user).Error; err != nil {
			return nil, fmt.Errorf("database error: %w", err)
		}
	}
	return &model.User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

// UserLogin is the resolver for the userLogin field.
func (r *mutationResolver) UserLogin(ctx context.Context, email string, password string) (*model.UserToken, error) {
	user := database.User{Email: email}
	if err := r.DB.WithContext(ctx).Where(&user).First(&user).Error; err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			err = fmt.Errorf("user %s not found", email)
		default:
			err = fmt.Errorf("database error: %w", err)
		}
		return nil, err
	}
	key := auth.Key([]byte(password), user.Salt)
	if bytes.Compare(key, user.Key) != 0 {
		return nil, fmt.Errorf("password error")
	}
	now := time.Now()
	claims := model.UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(time.Hour * 24).Unix(),
			Issuer:    string(auth.Issuer),
			IssuedAt:  now.Unix(),
		},
		UserID: user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(auth.Issuer)
	return &model.UserToken{
		Token:     signedToken,
		ExpiresAt: int(claims.ExpiresAt),
	}, nil
}

// RuleCreate is the resolver for the ruleCreate field.
func (r *mutationResolver) RuleCreate(ctx context.Context, summary string, detail *string) (*model.Rule, error) {
	userAuth := auth.ForContext(ctx)
	if userAuth == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	row := database.Rule{
		UserID:  userAuth.UserID,
		Created: time.Now(),
		Summary: summary,
		Detail:  detail,
	}
	if err := r.DB.WithContext(ctx).Create(&row).Error; err != nil {
		return nil, fmt.Errorf("database error: %v", err)
	}
	return &model.Rule{
		ID:      row.ID,
		Created: row.Created.String(),
		Summary: row.Summary,
	}, nil
}

// RuleDelete is the resolver for the ruleDelete field.
func (r *mutationResolver) RuleDelete(ctx context.Context, id int) (*int, error) {
	userAuth := auth.ForContext(ctx)
	if userAuth == nil {
		return nil, fmt.Errorf("unauthorized")
	}
	var rows []database.Rule
	err := r.DB.WithContext(ctx).
		Clauses(clause.Returning{}).
		Where(&database.Rule{ID: id, UserID: userAuth.UserID}).
		Delete(&rows).Error
	if err != nil {
		return nil, err
	}
	switch len(rows) {
	case 0:
		return nil, nil
	case 1:
		return &rows[0].ID, nil
	default:
		log.Printf("deleted multiple rules with id=%d (impossible!): %v", id, rows)
		return &rows[0].ID, nil
	}
}

// LikesUpdate is the resolver for the likesUpdate field.
func (r *mutationResolver) LikesUpdate(ctx context.Context, add []int, remove []int) (*model.LikesUpdate, error) {
	userAuth := auth.ForContext(ctx)
	if userAuth == nil {
		return nil, fmt.Errorf("unauthorized")
	}

	addMap := make(map[int]struct{})
	for _, id := range add {
		addMap[id] = struct{}{}
	}

	var conflicts []int
	for _, id := range remove {
		if _, ok := addMap[id]; ok {
			conflicts = append(conflicts, id)
		}
	}
	if len(conflicts) != 0 {
		return nil, fmt.Errorf("request both add/remove ids=%v", conflicts)
	}

	var update model.LikesUpdate

	if err := r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		toRows := func(ids []int) []database.Like {
			likes := make([]database.Like, len(ids))
			for i, id := range ids {
				likes[i].UserID = userAuth.UserID
				likes[i].RuleID = id
			}
			return likes
		}
		fromRows := func(rows []database.Like) []int {
			ids := make([]int, len(rows))
			for i, row := range rows {
				ids[i] = row.RuleID
			}
			return ids
		}
		filterRows := func(rows []database.Like, filter func(row database.Like) bool) []database.Like {
			var out []database.Like
			for _, row := range rows {
				if filter(row) {
					out = append(out, row)
				}
			}
			return out
		}

		if add != nil {
			addRows := toRows(add)

			// TODO: Is there a way to make INSERT ... ON CONFLICT DO NOTHING
			//  return only newly inserted rows?
			var existRows []database.Like
			// TODO: Is there a way to do this with WHERE ... IN using structs?
			qry := tx
			for _, row := range addRows {
				qry = qry.Or(row)
			}
			if err := qry.Find(&existRows).Error; err != nil {
				return fmt.Errorf("find existing likes error: %w", err)
			}
			if len(existRows) != 0 {
				exist := make(map[int]struct{})
				for _, row := range existRows {
					exist[row.RuleID] = struct{}{}
				}
				addRows = filterRows(addRows, func(row database.Like) bool {
					_, ok := exist[row.RuleID]
					return !ok
				})
			}

			if len(addRows) > 0 {
				if err := tx.Create(&addRows).Error; err != nil {
					return fmt.Errorf("add likes error: %w", err)
				}
			}
			update.Added = fromRows(addRows)
		}
		if remove != nil {
			removeRows := toRows(remove)
			if err := tx.Clauses(clause.Returning{}).Delete(&removeRows).Error; err != nil {
				return fmt.Errorf("remove likes error: %w", err)
			}
			update.Removed = fromRows(removeRows)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &update, nil
}

// Users is the resolver for the users field.
func (r *queryResolver) Users(ctx context.Context, limit int, after int, name *string) (*model.UserPage, error) {
	page := PageReader[database.User]{
		Query: r.DB.WithContext(ctx),
		After: database.User{ID: after},
		Limit: limit,
	}
	if name != nil {
		page.Query = page.Query.Scopes(database.User{Name: *name}.NameLike())
	}
	if err := page.Read(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.UserPage{
		Users: MapPointersOf(page.Rows, func(row database.User) model.User {
			return model.User{
				ID:   row.ID,
				Name: row.Name,
			}
		}),
		PageInfo: page.Info(),
	}, nil
}

// Rules is the resolver for the rules field.
func (r *queryResolver) Rules(ctx context.Context, limit int, after int, userID *int) (*model.RulePage, error) {
	page := PageReader[database.Rule]{
		Query: r.DB.WithContext(ctx),
		After: database.Rule{ID: after},
		Limit: limit,
	}
	if userID != nil {
		page.Query = page.Query.Where(&database.Rule{UserID: *userID})
	}
	if err := page.Read(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.RulePage{
		Rules: MapPointersOf(page.Rows, func(row database.Rule) model.Rule {
			return model.Rule{
				ID:      row.ID,
				Summary: row.Summary,
				Detail:  row.Detail,
				Created: row.Created.String(),
			}
		}),
		PageInfo: page.Info(),
	}, nil
}

// User is the resolver for the user field.
func (r *ruleResolver) User(ctx context.Context, obj *model.Rule) (*model.User, error) {
	row := database.Rule{
		ID: obj.ID,
	}
	if err := r.DB.WithContext(ctx).Preload("User").Find(&row).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.User{
		ID:   row.User.ID,
		Name: row.User.Name,
	}, nil
}

// Likes is the resolver for the likes field.
func (r *ruleResolver) Likes(ctx context.Context, obj *model.Rule, limit int, after int) (*model.UserPage, error) {
	page := PageReader[database.UserLike]{
		Query: r.DB.WithContext(ctx).Preload("User").Where(database.UserLike{RuleID: obj.ID}),
		After: database.UserLike{UserID: after},
		Limit: limit,
	}
	if err := page.Read(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.UserPage{
		Users: MapPointersOf(page.Rows, func(row database.UserLike) model.User {
			return model.User{
				ID:   row.User.ID,
				Name: row.User.Name,
			}
		}),
		PageInfo: page.Info(),
	}, nil
}

// Rules is the resolver for the rules field.
func (r *userResolver) Rules(ctx context.Context, obj *model.User, limit int, after int) (*model.RulePage, error) {
	page := PageReader[database.Rule]{
		Query: r.DB.WithContext(ctx).Where(database.Rule{UserID: obj.ID}),
		After: database.Rule{ID: after},
		Limit: limit,
	}
	if err := page.Read(); err != nil {
		return nil, fmt.Errorf("database read error: %w", err)
	}
	return &model.RulePage{
		Rules: MapPointersOf(page.Rows, func(row database.Rule) model.Rule {
			return model.Rule{
				ID:      row.ID,
				Created: row.Created.String(),
				Summary: row.Summary,
				Detail:  row.Detail,
			}
		}),
		PageInfo: nil,
	}, nil
}

// Likes is the resolver for the likes field.
func (r *userResolver) Likes(ctx context.Context, obj *model.User, limit int, after int) (*model.RulePage, error) {
	page := PageReader[database.RuleLike]{
		Query: r.DB.WithContext(ctx).Preload("Rule").Where(database.RuleLike{UserID: obj.ID}),
		After: database.RuleLike{RuleID: after},
		Limit: limit,
	}
	if err := page.Read(); err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &model.RulePage{
		Rules: MapPointersOf(page.Rows, func(row database.RuleLike) model.Rule {
			return model.Rule{
				ID:      row.Rule.ID,
				Created: row.Rule.Created.String(),
				Summary: row.Rule.Summary,
				Detail:  row.Rule.Detail,
			}
		}),
		PageInfo: page.Info(),
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Rule returns generated.RuleResolver implementation.
func (r *Resolver) Rule() generated.RuleResolver { return &ruleResolver{r} }

// User returns generated.UserResolver implementation.
func (r *Resolver) User() generated.UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type ruleResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
