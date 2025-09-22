package openfga

import (
	"context"
	"fmt"
	"strings"

	"github.com/OtchereDev/ris-common-sdk/pkg/roles"
	openfga "github.com/openfga/go-sdk"
	"github.com/openfga/go-sdk/client"
	"gorm.io/gorm"
)

// OpenFGA Service using the current SDK
type OpenFGAService struct {
	client *client.OpenFgaClient
	db     *gorm.DB
}

// Initialize OpenFGA Service using current SDK patterns
func NewOpenFGAService(apiUrl, storeId, modelId string, db *gorm.DB) (*OpenFGAService, error) {
	fgaClient, err := client.NewSdkClient(&client.ClientConfiguration{
		ApiUrl:               apiUrl,
		StoreId:              storeId,
		AuthorizationModelId: modelId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create OpenFGA client: %w", err)
	}

	return &OpenFGAService{
		client: fgaClient,
		db:     db,
	}, nil
}

// Check if user has permission
func (s *OpenFGAService) CheckPermission(ctx context.Context, userID string, permission roles.Permission) (bool, error) {
	body := client.ClientCheckRequest{
		User:     fmt.Sprintf("user:%s", userID),
		Relation: "can_access",
		Object:   fmt.Sprintf("permission:%s", permission),
	}

	data, err := s.client.Check(ctx).Body(body).Execute()
	if err != nil {
		return false, fmt.Errorf("OpenFGA check failed: %w", err)
	}

	return data.GetAllowed(), nil
}

// Check permission on specific resource
func (s *OpenFGAService) CheckResourcePermission(ctx context.Context, userID, relation, resourceType, resourceID string) (bool, error) {
	body := client.ClientCheckRequest{
		User:     fmt.Sprintf("user:%s", userID),
		Relation: relation,
		Object:   fmt.Sprintf("%s:%s", resourceType, resourceID),
	}

	data, err := s.client.Check(ctx).Body(body).Execute()
	if err != nil {
		return false, fmt.Errorf("OpenFGA resource check failed: %w", err)
	}

	return data.GetAllowed(), nil
}

// Grant permission to user
func (s *OpenFGAService) GrantPermission(ctx context.Context, userID string, permission roles.Permission) error {
	body := client.ClientWriteRequest{
		Writes: []client.ClientTupleKey{
			{
				User:     fmt.Sprintf("user:%s", userID),
				Relation: "can_access",
				Object:   fmt.Sprintf("permission:%s", permission),
			},
		},
	}

	_, err := s.client.Write(ctx).Body(body).Execute()
	if err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}

	return nil
}

// Revoke permission from user
func (s *OpenFGAService) RevokePermission(ctx context.Context, userID string, permission roles.Permission) error {
	body := client.ClientWriteRequest{
		Deletes: []client.ClientTupleKeyWithoutCondition{
			{
				User:     fmt.Sprintf("user:%s", userID),
				Relation: "can_access",
				Object:   fmt.Sprintf("permission:%s", permission),
			},
		},
	}

	_, err := s.client.Write(ctx).Body(body).Execute()
	if err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	return nil
}

// Assign role to user with default permissions
func (s *OpenFGAService) AssignRole(ctx context.Context, userID, role string) error {
	// Get default permissions for role
	permissions, exists := roles.DefaultRolePermissions[role]
	if !exists {
		return fmt.Errorf("unknown role: %s", role)
	}

	// Create tuples for role membership and permissions
	var writes []client.ClientTupleKey

	// Assign user to role
	writes = append(writes, client.ClientTupleKey{
		User:     fmt.Sprintf("user:%s", userID),
		Relation: "member",
		Object:   fmt.Sprintf("role:%s", strings.ToLower(role)),
	})

	// Grant default permissions
	for _, permission := range permissions {
		writes = append(writes, client.ClientTupleKey{
			User:     fmt.Sprintf("user:%s", userID),
			Relation: "can_access",
			Object:   fmt.Sprintf("permission:%s", permission),
		})
	}

	body := client.ClientWriteRequest{
		Writes: writes,
	}

	_, err := s.client.Write(ctx).Body(body).Execute()
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// Remove role from user
func (s *OpenFGAService) RemoveRole(ctx context.Context, userID, role string) error {
	// Get permissions for role to remove
	permissions, exists := roles.DefaultRolePermissions[role]
	if !exists {
		return fmt.Errorf("unknown role: %s", role)
	}

	// Create deletes for role membership and default permissions
	var deletes []client.ClientTupleKeyWithoutCondition

	// Remove role membership
	deletes = append(deletes, client.ClientTupleKeyWithoutCondition{
		User:     fmt.Sprintf("user:%s", userID),
		Relation: "member",
		Object:   fmt.Sprintf("role:%s", strings.ToLower(role)),
	})

	// Remove default role permissions (but keep any custom grants)
	for _, permission := range permissions {
		deletes = append(deletes, client.ClientTupleKeyWithoutCondition{
			User:     fmt.Sprintf("user:%s", userID),
			Relation: "can_access",
			Object:   fmt.Sprintf("permission:%s", permission),
		})
	}

	body := client.ClientWriteRequest{
		Deletes: deletes,
	}

	_, err := s.client.Write(ctx).Body(body).Execute()
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	return nil
}

// Grant resource-specific permission
func (s *OpenFGAService) GrantResourcePermission(ctx context.Context, userID, relation, resourceType, resourceID string) error {
	body := client.ClientWriteRequest{
		Writes: []client.ClientTupleKey{
			{
				User:     fmt.Sprintf("user:%s", userID),
				Relation: relation,
				Object:   fmt.Sprintf("%s:%s", resourceType, resourceID),
			},
		},
	}

	_, err := s.client.Write(ctx).Body(body).Execute()
	if err != nil {
		return fmt.Errorf("failed to grant resource permission: %w", err)
	}

	return nil
}

// List user permissions
func (s *OpenFGAService) ListUserPermissions(ctx context.Context, userID string) ([]roles.Permission, error) {
	body := client.ClientReadRequest{
		User:     openfga.PtrString(fmt.Sprintf("user:%s", userID)),
		Relation: openfga.PtrString("can_access"),
		Object:   openfga.PtrString("permission:"),
	}

	response, err := s.client.Read(ctx).Body(body).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}

	var permissions []roles.Permission
	for _, tuple := range response.GetTuples() {
		if strings.HasPrefix(tuple.Key.GetObject(), "permission:") {
			permissionStr := strings.TrimPrefix(tuple.Key.GetObject(), "permission:")
			permissions = append(permissions, roles.Permission(permissionStr))
		}
	}

	return permissions, nil
}

// List objects user has access to
func (s *OpenFGAService) ListObjectsUserCanAccess(ctx context.Context, userID, relation, objectType string) ([]string, error) {
	body := client.ClientListObjectsRequest{
		User:     fmt.Sprintf("user:%s", userID),
		Relation: relation,
		Type:     objectType,
	}

	data, err := s.client.ListObjects(ctx).Body(body).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	return data.GetObjects(), nil
}

// Batch check multiple permissions
func (s *OpenFGAService) BatchCheckPermissions(ctx context.Context, userID string, permissions []roles.Permission) (map[roles.Permission]bool, error) {
	var checks []client.ClientBatchCheckItem
	results := make(map[roles.Permission]bool)

	for i, permission := range permissions {
		checks = append(checks, client.ClientBatchCheckItem{
			CorrelationId: fmt.Sprintf("check_%d", i),
			User:          fmt.Sprintf("user:%s", userID),
			Relation:      "can_access",
			Object:        fmt.Sprintf("permission:%s", permission),
		})
	}

	body := client.ClientBatchCheckRequest{
		Checks: checks,
	}

	data, err := s.client.BatchCheck(ctx).Body(body).Execute()
	if err != nil {
		return nil, fmt.Errorf("batch check failed: %w", err)
	}

	// Process results
	for i, permission := range permissions {
		correlationId := fmt.Sprintf("check_%d", i)
		if result, exists := data.GetResult()[correlationId]; exists {
			results[permission] = result.GetAllowed()
		} else {
			results[permission] = false
		}
	}

	return results, nil
}
