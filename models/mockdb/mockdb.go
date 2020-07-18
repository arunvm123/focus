package mockdb

import (
	"github.com/arunvm/travail-backend/email"
	"github.com/arunvm/travail-backend/models"
	push "github.com/arunvm/travail-backend/push_notification"
	"github.com/stretchr/testify/mock"
)

func New() *MockDB {
	return &MockDB{}
}

type MockDB struct {
	mock.Mock
}

func (mock *MockDB) GetLists(args *models.GetListsArgs, user *models.User) (*[]models.ListInfo, error) {
	return nil, nil
}

func (mock *MockDB) UpdateList(args *models.UpdateListArgs, user *models.User) error {
	return nil
}

func (mock *MockDB) CreateList(args *models.CreateListArgs, user *models.User) (*models.List, error) {
	rets := mock.Called(args, user)
	return rets.Get(0).(*models.List), rets.Error(1)
}

func (mock *MockDB) CreateBoardColumn(args *models.CreateBoardColumnArgs) error {
	return nil
}

func (mock *MockDB) GetBoardColumns(boardID string) (*[]models.BoardColumn, error) {
	return nil, nil
}

func (mock *MockDB) UpdateBoardColumn(args *models.UpdateBoardColumnArgs) error {
	return nil
}

func (mock *MockDB) CheckIfColumnPartOfBoard(boardColumnID string, boardID string) bool {
	return false
}

func (mock *MockDB) CreateBoard(args *models.CreateBoardArgs, teamMember *models.User) error {
	return nil
}

func (mock *MockDB) UpdateBoard(args *models.UpdateBoardArgs) error {
	return nil
}

func (mock *MockDB) GetBoards(args *models.GetBoardsArgs) (*[]models.Board, error) {
	return nil, nil
}

func (mock *MockDB) CheckIfBoardPartOfTeam(boardID, teamID string) bool {
	return false
}

func (mock *MockDB) CreateBug(args *models.CreateBugArgs, user *models.User) error {
	return nil
}

func (mock *MockDB) GetBugs() (*[]models.BugInfo, error) {
	return nil, nil
}

func (mock *MockDB) UpdateBug(args *models.UpdateBugArgs) error {
	return nil
}

func (mock *MockDB) CreateColumnCard(args *models.CreateColumnCardArgs) error {
	return nil
}

func (mock *MockDB) GetColumnCards(columnID string) (*[]models.ColumnCard, error) {
	return nil, nil
}

func (mock *MockDB) UpdateColumnCard(args *models.UpdateColumnCardArgs, user *models.User) error {
	return nil
}

func (mock *MockDB) CreateEmailValidationToken(user *models.User, emailCLient email.Email) error {
	return nil
}

func (mock *MockDB) VerifyEmail(token string) error {
	return nil
}

func (mock *MockDB) InvalidateEmailTokens(userID int) error {
	return nil
}

func (mock *MockDB) AddNotificationToken(args *models.AddNotificationTokenArgs, user *models.User) error {
	return nil
}

func (mock *MockDB) GetNotificationTokens(user *models.User) ([]string, error) {
	return []string{}, nil
}

func (mock *MockDB) CreateForgotPasswordToken(user *models.User, emailCLient email.Email) error {
	return nil
}

func (mock *MockDB) ResetPassword(token, password string) error {
	return nil
}

func (mock *MockDB) CreateOrganisationInviteToken(args *models.InviteToOrganisationArgs, admin *models.User, emailClient email.Email) error {
	return nil
}

func (mock *MockDB) AcceptOrganisationInvite(args *models.AcceptOrganisationInviteArgs, user *models.User) error {
	return nil
}

func (mock *MockDB) GetOrganisationMembers(organisationID string) (*[]models.OrganisationMemberInfo, error) {
	return nil, nil
}

func (mock *MockDB) CheckIfOrganisationMember(organisationID string, user *models.User) bool {
	return false
}

func (mock *MockDB) CreateOrganisation(args *models.CreateOrganisationArgs, user *models.User) error {
	return nil
}

func (mock *MockDB) GetOrganisations(user *models.User) (*[]models.Organisation, error) {
	return nil, nil
}

func (mock *MockDB) UpdateOrganisation(args *models.UpdateOrganisationArgs, admin *models.User) error {
	return nil
}

func (mock *MockDB) CheckIfOrganisationAdmin(orgID string, user *models.User) bool {
	return false
}

func (mock *MockDB) GetOrganisationName(organisationID string) (string, error) {
	return "", nil
}

func (mock *MockDB) CreateTask(args *models.CreateTaskArgs, user *models.User) (*models.Task, error) {
	return nil, nil
}

func (mock *MockDB) GetTasks(args *models.GetTasksArgs, user *models.User) (*[]models.TaskInfo, error) {
	return nil, nil
}

func (mock *MockDB) UpdateTask(args models.UpdateTaskArgs, user *models.User) error {
	return nil
}

func (mock *MockDB) DeleteTasks(args *models.DeleteTasksArgs, user *models.User) error {
	return nil
}

func (mock *MockDB) AddTeamMember(args models.AddTeamMemberArgs) error {
	return nil
}

func (mock *MockDB) GetTeamMembers(teamID string) (*[]models.TeamMemberInfo, error) {
	return nil, nil
}

func (mock *MockDB) CheckIfTeamMember(teamID string, user *models.User) bool {
	return false
}

func (mock *MockDB) CreateTeam(args *models.CreateTeamArgs, user *models.User) error {
	return nil
}

func (mock *MockDB) UpdateTeam(args *models.UpdateTeamArgs, teamAdmin *models.User) error {
	return nil
}

func (mock *MockDB) GetPersonalTeamID(user *models.User) (string, error) {
	return "", nil
}

func (mock *MockDB) CheckIfTeamAdmin(teamID string, user *models.User) bool {
	return false
}

func (mock *MockDB) SendPushNotificationForTasksAboutToExpire(pushClient push.Notification) error {
	return nil
}

func (mock *MockDB) GetProfile(user *models.User) (*models.UserProfile, error) {
	return nil, nil
}

func (mock *MockDB) UpdateProfile(args models.UpdateProfileArgs, user *models.User) error {
	return nil
}

func (mock *MockDB) GetUserFromEmail(email string) (*models.User, error) {
	return nil, nil
}

func (mock *MockDB) GetUserFromID(userID int) (*models.User, error) {
	return nil, nil
}

func (mock *MockDB) CheckIfUserExists(email string) bool {
	return false
}

func (mock *MockDB) UserSignup(args *models.SignUpArgs, googleOauth bool, emailClient email.Email) (*models.User, error) {
	return nil, nil
}

func (mock *MockDB) UpdatePassword(args *models.UpdatePasswordArgs, user *models.User) error {
	return nil
}

// Transaction
func (mock *MockDB) Begin() models.DB {
	return nil
}

func (mock *MockDB) Commit() {
	return
}

func (mock *MockDB) Rollback() {
	return
}
