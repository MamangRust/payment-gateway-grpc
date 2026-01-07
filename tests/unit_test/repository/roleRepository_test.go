package test

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	mocks "MamangRust/paymentgatewaygrpc/internal/repository/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RoleRepositorySuite struct {
	suite.Suite
	mockCtrl *gomock.Controller
	mockRepo *mocks.MockRoleRepository
}

func (suite *RoleRepositorySuite) SetupTest() {
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockRepo = mocks.NewMockRoleRepository(suite.mockCtrl)
}

func (suite *RoleRepositorySuite) TearDownTest() {
	suite.mockCtrl.Finish()
}

func (suite *RoleRepositorySuite) TestFindAllRoles_Success() {
	req := &requests.FindAllRoles{Search: "", Page: 1, PageSize: 10}
	roles := []*record.RoleRecord{{ID: 1, Name: "Admin"}}
	total := 1

	suite.mockRepo.EXPECT().FindAllRoles(req).Return(roles, &total, nil)

	result, totalRes, err := suite.mockRepo.FindAllRoles(req)

	suite.NoError(err)
	suite.Equal(roles, result)
	suite.Equal(1, *totalRes)
}

func (suite *RoleRepositorySuite) TestFindByActiveRole_Success() {
	req := &requests.FindAllRoles{Search: "active", Page: 1, PageSize: 10}
	activeRoles := []*record.RoleRecord{{ID: 1, Name: "Active Role"}}
	total := 1

	suite.mockRepo.EXPECT().FindByActiveRole(req).Return(activeRoles, &total, nil)

	result, totalRes, err := suite.mockRepo.FindByActiveRole(req)

	suite.NoError(err)
	suite.Equal(activeRoles, result)
	suite.Equal(1, *totalRes)
}

func (suite *RoleRepositorySuite) TestFindByTrashedRole_Success() {
	req := &requests.FindAllRoles{Page: 1, PageSize: 10}
	trashedTime := time.Now().Format(time.RFC3339)
	trashedRoles := []*record.RoleRecord{{ID: 2, Name: "Deleted Role", DeletedAt: &trashedTime}}
	total := 1

	suite.mockRepo.EXPECT().FindByTrashedRole(req).Return(trashedRoles, &total, nil)

	result, totalRes, err := suite.mockRepo.FindByTrashedRole(req)

	suite.NoError(err)
	suite.Equal(trashedRoles, result)
	suite.Equal(1, *totalRes)
}

func (suite *RoleRepositorySuite) TestFindByIdRole_Success() {
	roleID := 1
	role := &record.RoleRecord{ID: roleID, Name: "Admin"}

	suite.mockRepo.EXPECT().FindById(roleID).Return(role, nil)

	result, err := suite.mockRepo.FindById(roleID)

	suite.NoError(err)
	suite.Equal(role, result)
}

func (suite *RoleRepositorySuite) TestFindByIdRole_NotFound() {
	roleID := 99
	expectedErr := errors.New("role not found")

	suite.mockRepo.EXPECT().FindById(roleID).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindById(roleID)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *RoleRepositorySuite) TestFindByName_Success() {
	name := "Admin"
	role := &record.RoleRecord{ID: 1, Name: name}

	suite.mockRepo.EXPECT().FindByName(name).Return(role, nil)

	result, err := suite.mockRepo.FindByName(name)

	suite.NoError(err)
	suite.Equal(role, result)
}

func (suite *RoleRepositorySuite) TestFindByName_NotFound() {
	name := "NonExistentRole"
	expectedErr := errors.New("role not found")

	suite.mockRepo.EXPECT().FindByName(name).Return(nil, expectedErr)

	result, err := suite.mockRepo.FindByName(name)

	suite.Error(err)
	suite.Nil(result)
	suite.Equal(expectedErr, err)
}

func (suite *RoleRepositorySuite) TestFindByUserId_Success() {
	userID := 1
	userRoles := []*record.RoleRecord{
		{ID: 1, Name: "Admin"},
		{ID: 2, Name: "Editor"},
	}

	suite.mockRepo.EXPECT().FindByUserId(userID).Return(userRoles, nil)

	result, err := suite.mockRepo.FindByUserId(userID)

	suite.NoError(err)
	suite.Equal(userRoles, result)
}

func (suite *RoleRepositorySuite) TestCreateRole_Success() {
	req := &requests.CreateRoleRequest{Name: "Viewer"}
	createdRole := &record.RoleRecord{ID: 10, Name: "Viewer"}

	suite.mockRepo.EXPECT().CreateRole(req).Return(createdRole, nil)

	result, err := suite.mockRepo.CreateRole(req)

	suite.NoError(err)
	suite.Equal(createdRole, result)
}

func (suite *RoleRepositorySuite) TestUpdateRole_Success() {
	roleID := 1
	req := &requests.UpdateRoleRequest{ID: &roleID, Name: "Super Admin"}
	updatedRole := &record.RoleRecord{ID: roleID, Name: "Super Admin"}

	suite.mockRepo.EXPECT().UpdateRole(req).Return(updatedRole, nil)

	result, err := suite.mockRepo.UpdateRole(req)

	suite.NoError(err)
	suite.Equal(updatedRole, result)
}

func (suite *RoleRepositorySuite) TestTrashedRole_Success() {
	roleID := 1
	trashedTime := time.Now().Format(time.RFC3339)
	trashedRole := &record.RoleRecord{ID: roleID, DeletedAt: &trashedTime}

	suite.mockRepo.EXPECT().TrashedRole(roleID).Return(trashedRole, nil)

	result, err := suite.mockRepo.TrashedRole(roleID)

	suite.NoError(err)
	suite.NotNil(result.DeletedAt) // Pastikan role telah di-trash
	suite.Equal(trashedRole, result)
}

func (suite *RoleRepositorySuite) TestRestoreRole_Success() {
	roleID := 1
	restoredRole := &record.RoleRecord{ID: roleID, DeletedAt: nil} // DeletedAt harus nil setelah restore

	suite.mockRepo.EXPECT().RestoreRole(roleID).Return(restoredRole, nil)

	result, err := suite.mockRepo.RestoreRole(roleID)

	suite.NoError(err)
	suite.Nil(result.DeletedAt) // Pastikan role telah di-restore
	suite.Equal(restoredRole, result)
}

func (suite *RoleRepositorySuite) TestDeleteRolePermanent_Success() {
	roleID := 1

	suite.mockRepo.EXPECT().DeleteRolePermanent(roleID).Return(true, nil)

	result, err := suite.mockRepo.DeleteRolePermanent(roleID)

	suite.NoError(err)
	suite.True(result)
}

func (suite *RoleRepositorySuite) TestRestoreAllRole_Success() {
	suite.mockRepo.EXPECT().RestoreAllRole().Return(true, nil)

	result, err := suite.mockRepo.RestoreAllRole()

	suite.NoError(err)
	suite.True(result)
}

func (suite *RoleRepositorySuite) TestDeleteAllRolePermanent_Success() {
	suite.mockRepo.EXPECT().DeleteAllRolePermanent().Return(true, nil)

	result, err := suite.mockRepo.DeleteAllRolePermanent()

	suite.NoError(err)
	suite.True(result)
}

func TestRoleRepositorySuite(t *testing.T) {
	suite.Run(t, new(RoleRepositorySuite))
}
