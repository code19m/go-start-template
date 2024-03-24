package mongostore

// Example of how to use the mongoStore struct and its methods.

// import (
// 	"context"
// 	"fmt"
// 	"go-start-template/internal/domain"

// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// // CreateUserLog inserts a new user log in the database.
// func (m *mongoStore) CreateUserLog(ctx context.Context, userLog domain.UserLog) error {
// 	const op = "mongoStore.CreateUserLog"

// 	_, err := m.userLogsColl.InsertOne(ctx, userLog)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }

// // GetPaginatedUserLogs returns a paginated list of user logs from the database.
// func (m *mongoStore) GetPaginatedUserLogs(ctx context.Context, params domain.GetUserLogsParams) (domain.GetUserLogsResponse, error) {
// 	const op = "mongoStore.GetPaginatedUserLogs"

// 	result := domain.GetUserLogsResponse{}

// 	content, err := m.getUserLogs(ctx, params)
// 	if err != nil {
// 		return result, fmt.Errorf("%s: %w", op, err)
// 	}

// 	count, err := m.getUserLogsCount(ctx, params)
// 	if err != nil {
// 		return result, fmt.Errorf("%s: %w", op, err)
// 	}

// 	result.Content = content
// 	result.PaginationResult = domain.NewPaginationResult(params.Page, params.Size, count)

// 	return result, nil
// }

// // getUserLogs returns a filtered list of user logs from the database.
// func (m *mongoStore) getUserLogs(ctx context.Context, params domain.GetUserLogsParams) ([]domain.UserLog, error) {
// 	const op = "mongoStore.GetUserLogs"

// 	filter := m.buildUserLogsFilter(params)

// 	sort := bson.M{"created_at": -1}
// 	skip := (params.Page - 1) * params.Size
// 	limit := params.Size

// 	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64(skip)).SetSort(sort)

// 	var userLogs = make([]domain.UserLog, 0)
// 	err := m.userLogsColl.SimpleFindWithCtx(ctx, &userLogs, filter, opts)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return userLogs, nil
// }

// // getUserLogsCount returns the total number of user logs that match the filter.
// func (m *mongoStore) getUserLogsCount(ctx context.Context, params domain.GetUserLogsParams) (int64, error) {
// 	const op = "mongoStore.GetUserLogsCount"

// 	filter := m.buildUserLogsFilter(params)

// 	count, err := m.userLogsColl.CountDocuments(ctx, filter)
// 	if err != nil {
// 		return 0, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return count, nil
// }

// // buildUserLogsFilter builds a filter for the user logs collection.
// func (m *mongoStore) buildUserLogsFilter(params domain.GetUserLogsParams) bson.M {
// 	filter := make(bson.M)

// 	if params.UserID != 0 {
// 		filter["user_id"] = params.UserID
// 	}

// 	if params.UserPinfl != "" {
// 		filter["user_pinfl"] = params.UserPinfl
// 	}

// 	if params.OrganizationId != 0 {
// 		filter["organization_id"] = params.OrganizationId
// 	}

// 	if params.OrganizationTin != "" {
// 		filter["organization_tin"] = params.OrganizationTin
// 	}

// 	if params.Action != "" {
// 		filter["action"] = params.Action
// 	}

// 	if params.DateFrom != "" || params.DateTo != "" {
// 		filter["created_at"] = m.buildDateFilter(params.DateFrom, params.DateTo)
// 	}

// 	return filter
// }

// // buildDateFilter builds a date filter for the user logs collection.
// func (m *mongoStore) buildDateFilter(dateFrom, dateTo string) bson.M {
// 	filter := make(bson.M)

// 	if dateFrom != "" {
// 		filter["$gte"] = dateFrom
// 	}

// 	if dateTo != "" {
// 		filter["$lte"] = dateTo
// 	}

// 	return filter
// }
