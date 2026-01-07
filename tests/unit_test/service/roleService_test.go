package test

import (
	"errors"
	"testing"

	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	mock_responseservice "MamangRust/paymentgatewaygrpc/internal/mapper/response/mocks"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"MamangRust/paymentgatewaygrpc/internal/service" // Sesuaikan path ini

	"MamangRust/paymentgatewaygrpc/pkg/errors/role_errors"
	mock_logger "MamangRust/paymentgatewaygrpc/pkg/logger/mocks" // Sesuaikan path ini

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type RoleServiceSuite struct {
	suite.Suite
	mockRepo   *mocks.MockRoleRepository
	mockLogger *mock_logger.MockLoggerInterface
	mockMapper *mock_responseservice.MockRoleResponseMapper
	mockCtrl   *gomock.Controller
}

func (suite *RoleServiceSuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRoleRepository(suite.mockCtrl)
	suite.mockLogger = mock_logger.NewMockLoggerInterface(suite.mockCtrl)
	suite.mockMapper = mock_responseservice.NewMockRoleResponseMapper(suite.mockCtrl)
}

func (suite *RoleServiceSuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *RoleServiceSuite) TestFindAll_Success() {
	req := &requests.FindAllRoles{Search: "admin", Page: 1, PageSize: 10}
	total := 1
	mockRoles := []*record.RoleRecord{{ID: 1, Name: "Admin"}}
	expectedResponse := []*response.RoleResponse{{ID: 1, Name: "Admin"}}

	suite.mockLogger.EXPECT().Debug("Fetching role", gomock.Any())
	suite.mockRepo.EXPECT().FindAllRoles(req).Return(mockRoles, &total, nil)
	suite.mockMapper.EXPECT().ToRolesResponse(mockRoles).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched role", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, errResp := service.FindAll(req)

	suite.Nil(errResp)
	suite.NotNil(result)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *RoleServiceSuite) TestFindAll_Failure() {
	req := &requests.FindAllRoles{Search: "admin", Page: 1, PageSize: 10}
	dbError := errors.New("database connection failed")

	suite.mockLogger.EXPECT().Debug("Fetching role", gomock.Any())
	suite.mockRepo.EXPECT().FindAllRoles(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve role list", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := service.FindAll(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(role_errors.ErrFailedFindAll, err)
}

func (suite *RoleServiceSuite) TestFindById_Success() {
	roleID := 1
	mockRole := &record.RoleRecord{ID: roleID, Name: "Admin"}
	expectedResponse := &response.RoleResponse{ID: roleID, Name: "Admin"}

	suite.mockLogger.EXPECT().Debug("Fetching role by ID", zap.Int("id", roleID))
	suite.mockRepo.EXPECT().FindById(roleID).Return(mockRole, nil)
	suite.mockMapper.EXPECT().ToRoleResponse(mockRole).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched role", zap.Int("id", roleID))

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, errResp := service.FindById(roleID)

	suite.Nil(errResp)
	suite.Equal(expectedResponse, result)
}

func (suite *RoleServiceSuite) TestFindById_NotFound() {
	roleID := 99
	repoError := errors.New("role not found")

	suite.mockLogger.EXPECT().Debug("Fetching role by ID", zap.Int("id", roleID))
	suite.mockRepo.EXPECT().FindById(roleID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve role details", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, err := service.FindById(roleID)

	suite.Nil(result)
	suite.Equal(role_errors.ErrRoleNotFoundRes, err)
}

func (suite *RoleServiceSuite) TestFindByUserId_Success() {
	userID := 123
	mockRoles := []*record.RoleRecord{{ID: 1, Name: "User"}}
	expectedResponse := []*response.RoleResponse{{ID: 1, Name: "User"}}

	suite.mockLogger.EXPECT().Debug("Fetching role by user ID", zap.Int("id", userID))
	suite.mockRepo.EXPECT().FindByUserId(userID).Return(mockRoles, nil)
	suite.mockMapper.EXPECT().ToRolesResponse(mockRoles).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched role by user ID", zap.Int("id", userID))

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, errResp := service.FindByUserId(userID)

	suite.Nil(errResp)
	suite.Equal(expectedResponse, result)
}

func (suite *RoleServiceSuite) TestFindByUserId_Failure() {
	userID := 123
	repoError := errors.New("user not found")

	suite.mockLogger.EXPECT().Debug("Fetching role by user ID", zap.Int("id", userID))
	suite.mockRepo.EXPECT().FindByUserId(userID).Return(nil, repoError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve role by user ID", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, err := service.FindByUserId(userID)

	suite.Nil(result)
	suite.Equal(role_errors.ErrRoleNotFoundRes, err)
}

func (suite *RoleServiceSuite) TestFindByActiveRole_Success() {
	req := &requests.FindAllRoles{Search: "active", Page: 1, PageSize: 10}
	total := 1
	mockRoles := []*record.RoleRecord{{ID: 1, Name: "Active Role", DeletedAt: nil}}
	expectedResponse := []*response.RoleResponseDeleteAt{{ID: 1, Name: "Active Role", DeletedAt: nil}}

	suite.mockLogger.EXPECT().Debug("Fetching active role", gomock.Any())
	suite.mockRepo.EXPECT().FindByActiveRole(req).Return(mockRoles, &total, nil)
	suite.mockMapper.EXPECT().ToRolesResponseDeleteAt(mockRoles).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched active role", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, errResp := service.FindByActiveRole(req)

	suite.Nil(errResp)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *RoleServiceSuite) TestFindByActiveRole_Failure() {
	req := &requests.FindAllRoles{Search: "active", Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching active role", gomock.Any())
	suite.mockRepo.EXPECT().FindByActiveRole(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to retrieve active roles from database", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := service.FindByActiveRole(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(role_errors.ErrFailedFindActive, err)
}

func (suite *RoleServiceSuite) TestFindByTrashedRole_Success() {
	req := &requests.FindAllRoles{Page: 1, PageSize: 10}
	total := 1
	trashedTime := "2024-01-01T00:00:00Z"
	mockRoles := []*record.RoleRecord{{ID: 2, Name: "Trashed Role", DeletedAt: &trashedTime}}
	expectedResponse := []*response.RoleResponseDeleteAt{{ID: 2, Name: "Trashed Role", DeletedAt: &trashedTime}}

	suite.mockLogger.EXPECT().Debug("Fetching trashed role", gomock.Any())
	suite.mockRepo.EXPECT().FindByTrashedRole(req).Return(mockRoles, &total, nil)
	suite.mockMapper.EXPECT().ToRolesResponseDeleteAt(mockRoles).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("Successfully fetched trashed role", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, errResp := service.FindByTrashedRole(req)

	suite.Nil(errResp)
	suite.Equal(expectedResponse, result)
	suite.Equal(1, *totalRes)
}

func (suite *RoleServiceSuite) TestFindByTrashedRole_Failure() {
	req := &requests.FindAllRoles{Page: 1, PageSize: 10}
	dbError := errors.New("database error")

	suite.mockLogger.EXPECT().Debug("Fetching trashed role", gomock.Any())
	suite.mockRepo.EXPECT().FindByTrashedRole(req).Return(nil, nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to fetch trashed role", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, totalRes, err := service.FindByTrashedRole(req)

	suite.Nil(result)
	suite.Nil(totalRes)
	suite.Equal(role_errors.ErrFailedFindTrashed, err)
}

func (suite *RoleServiceSuite) TestCreateRole_Success() {
	req := &requests.CreateRoleRequest{Name: "New Role"}
	createdRole := &record.RoleRecord{ID: 10, Name: "New Role"}
	expectedResponse := &response.RoleResponse{ID: 10, Name: "New Role"}

	suite.mockLogger.EXPECT().Debug("Starting CreateRole process", gomock.Any())
	suite.mockRepo.EXPECT().CreateRole(req).Return(createdRole, nil)
	suite.mockMapper.EXPECT().ToRoleResponse(createdRole).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("CreateRole process completed", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, errResp := service.CreateRole(req)

	suite.Nil(errResp)
	suite.Equal(expectedResponse, result)
}

func (suite *RoleServiceSuite) TestCreateRole_Failure() {
	req := &requests.CreateRoleRequest{Name: "New Role"}
	dbError := errors.New("failed to create role")

	suite.mockLogger.EXPECT().Debug("Starting CreateRole process", gomock.Any())
	suite.mockRepo.EXPECT().CreateRole(req).Return(nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to create new role record", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, err := service.CreateRole(req)

	suite.Nil(result)
	suite.Equal(role_errors.ErrFailedCreateRole, err)
}

func (suite *RoleServiceSuite) TestUpdateRole_Success() {
	roleID := 1
	req := &requests.UpdateRoleRequest{ID: &roleID, Name: "Updated Role"}
	updatedRole := &record.RoleRecord{ID: roleID, Name: "Updated Role"}
	expectedResponse := &response.RoleResponse{ID: roleID, Name: "Updated Role"}

	suite.mockLogger.EXPECT().Debug("Starting UpdateRole process", gomock.Any())
	suite.mockRepo.EXPECT().UpdateRole(req).Return(updatedRole, nil)
	suite.mockMapper.EXPECT().ToRoleResponse(updatedRole).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("UpdateRole process completed", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, errResp := service.UpdateRole(req)

	suite.Nil(errResp)
	suite.Equal(expectedResponse, result)
}

func (suite *RoleServiceSuite) TestUpdateRole_Failure() {
	roleID := 1
	req := &requests.UpdateRoleRequest{ID: &roleID, Name: "Updated Role"}
	dbError := errors.New("failed to update role")

	suite.mockLogger.EXPECT().Debug("Starting UpdateRole process", gomock.Any())
	suite.mockRepo.EXPECT().UpdateRole(req).Return(nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to update role record", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, err := service.UpdateRole(req)

	suite.Nil(result)
	suite.Equal(role_errors.ErrFailedUpdateRole, err)
}

func (suite *RoleServiceSuite) TestTrashedRole_Success() {
	roleID := 1
	trashedRole := &record.RoleRecord{ID: roleID, Name: "Trashed Role"}
	expectedResponse := &response.RoleResponseDeleteAt{ID: roleID, Name: "Trashed Role"}

	suite.mockLogger.EXPECT().Debug("Starting TrashedRole process", zap.Int("roleID", roleID))
	suite.mockRepo.EXPECT().TrashedRole(roleID).Return(trashedRole, nil)
	suite.mockMapper.EXPECT().ToRoleResponseDeleteAt(trashedRole).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("TrashedRole process completed", zap.Int("roleID", roleID))

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, errResp := service.TrashedRole(roleID)

	suite.Nil(errResp)
	suite.Equal(expectedResponse, result)
}

func (suite *RoleServiceSuite) TestTrashedRole_Failure() {
	roleID := 1
	dbError := errors.New("failed to trash role")

	suite.mockLogger.EXPECT().Debug("Starting TrashedRole process", zap.Int("roleID", roleID))
	suite.mockRepo.EXPECT().TrashedRole(roleID).Return(nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to move role to trash", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, err := service.TrashedRole(roleID)

	suite.Nil(result)
	suite.Equal(role_errors.ErrFailedTrashedRole, err)
}

func (suite *RoleServiceSuite) TestRestoreRole_Success() {
	roleID := 1
	restoredRole := &record.RoleRecord{ID: roleID, Name: "Restored Role", DeletedAt: nil}
	expectedResponse := &response.RoleResponseDeleteAt{ID: roleID, Name: "Restored Role", DeletedAt: nil}

	suite.mockLogger.EXPECT().Debug("Starting RestoreRole process", zap.Int("roleID", roleID))
	suite.mockRepo.EXPECT().RestoreRole(roleID).Return(restoredRole, nil)
	suite.mockMapper.EXPECT().ToRoleResponseDeleteAt(restoredRole).Return(expectedResponse)
	suite.mockLogger.EXPECT().Debug("RestoreRole process completed", zap.Int("roleID", roleID))

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, errResp := service.RestoreRole(roleID)

	suite.Nil(errResp)
	suite.Equal(expectedResponse, result)
}

func (suite *RoleServiceSuite) TestRestoreRole_Failure() {
	roleID := 1
	dbError := errors.New("failed to restore role")

	suite.mockLogger.EXPECT().Debug("Starting RestoreRole process", zap.Int("roleID", roleID))
	suite.mockRepo.EXPECT().RestoreRole(roleID).Return(nil, dbError)
	suite.mockLogger.EXPECT().Error("Failed to restore role from trash", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, err := service.RestoreRole(roleID)

	suite.Nil(result)
	suite.Equal(role_errors.ErrFailedRestoreRole, err)
}

func (suite *RoleServiceSuite) TestDeleteRolePermanent_Success() {
	roleID := 1

	suite.mockLogger.EXPECT().Debug("Starting DeleteRolePermanent process", zap.Int("roleID", roleID))
	suite.mockRepo.EXPECT().DeleteRolePermanent(roleID).Return(true, nil)
	suite.mockLogger.EXPECT().Debug("DeleteRolePermanent process completed", zap.Int("roleID", roleID))

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, errResp := service.DeleteRolePermanent(roleID)

	suite.Nil(errResp)

	suite.True(result)
}

func (suite *RoleServiceSuite) TestDeleteRolePermanent_Failure() {
	roleID := 1
	dbError := errors.New("failed to delete role permanently")

	suite.mockLogger.EXPECT().Debug("Starting DeleteRolePermanent process", zap.Int("roleID", roleID))
	suite.mockRepo.EXPECT().DeleteRolePermanent(roleID).Return(false, dbError)
	suite.mockLogger.EXPECT().Error("Failed to permanently delete role", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, err := service.DeleteRolePermanent(roleID)

	suite.False(result)
	suite.Equal(role_errors.ErrFailedDeletePermanent, err)
}

func (suite *RoleServiceSuite) TestRestoreAllRole_Success() {
	suite.mockLogger.EXPECT().Debug("Restoring all roles")
	suite.mockRepo.EXPECT().RestoreAllRole().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully restored all roles")

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, errResp := service.RestoreAllRole()

	suite.Nil(errResp)

	suite.True(result)
}

func (suite *RoleServiceSuite) TestRestoreAllRole_Failure() {
	dbError := errors.New("failed to restore all roles")

	suite.mockLogger.EXPECT().Debug("Restoring all roles")
	suite.mockRepo.EXPECT().RestoreAllRole().Return(false, dbError)
	suite.mockLogger.EXPECT().Error("Failed to restore all trashed roles", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, err := service.RestoreAllRole()

	suite.False(result)
	suite.Equal(role_errors.ErrFailedRestoreAll, err)
}

func (suite *RoleServiceSuite) TestDeleteAllRolePermanent_Success() {
	suite.mockLogger.EXPECT().Debug("Permanently deleting all roles")
	suite.mockRepo.EXPECT().DeleteAllRolePermanent().Return(true, nil)
	suite.mockLogger.EXPECT().Debug("Successfully deleted all roles permanently")

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, errResp := service.DeleteAllRolePermanent()

	suite.Nil(errResp)
	suite.True(result)
}

func (suite *RoleServiceSuite) TestDeleteAllRolePermanent_Failure() {
	dbError := errors.New("failed to delete all roles permanently")

	suite.mockLogger.EXPECT().Debug("Permanently deleting all roles")
	suite.mockRepo.EXPECT().DeleteAllRolePermanent().Return(false, dbError)
	suite.mockLogger.EXPECT().Error("Failed to permanently delete all trashed roles", gomock.Any())

	service := service.NewRoleService(suite.mockRepo, suite.mockLogger, suite.mockMapper)
	result, err := service.DeleteAllRolePermanent()

	suite.False(result)
	suite.Equal(role_errors.ErrFailedDeletePermanent, err)
}

func TestRoleServiceSuite(t *testing.T) {
	suite.Run(t, new(RoleServiceSuite))
}
