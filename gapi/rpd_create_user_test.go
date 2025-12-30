package gapi

import (
	"context"
	mockdb "danielsxiong/simplebank/db/mock"
	db "danielsxiong/simplebank/db/sqlc"
	"danielsxiong/simplebank/pb"
	"danielsxiong/simplebank/util"
	"danielsxiong/simplebank/worker"
	mockworker "danielsxiong/simplebank/worker/mock"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArg, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(expected.password, actualArg.HashedPassword)
	if err != nil {
		return false
	}

	expected.arg.HashedPassword = actualArg.HashedPassword
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArg.CreateUserParams) {
		return false
	}

	err = actualArg.AfterCreate(expected.user)
	if err != nil {
		return false
	}

	return true
}

func (expected eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", expected.arg, expected.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{
		arg:      arg,
		password: password,
		user:     user,
	}
}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	return db.User{
		Username:       util.RandomOwner(),
		Role:           util.DepositorRole,
		HashedPassword: hashedPassword,
		Email:          util.RandomEmail(),
		FullName:       util.RandomOwner(),
	}, password
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor)
		checkResponse func(t *testing.T, response *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			body: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{User: user}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{Username: user.Username}

				taskDistributor.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, response *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, response)
				createdUser := response.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.FullName, createdUser.FullName)
				require.Equal(t, user.Email, createdUser.Email)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			taskDistributorCtrl := gomock.NewController(t)
			defer taskDistributorCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)
			taskDistributor := mockworker.NewMockTaskDistributor(taskDistributorCtrl)
			// build stubs
			tc.buildStubs(store, taskDistributor)

			server := newTestServer(t, store, taskDistributor)

			createUserResp, err := server.CreateUser(context.Background(), tc.body)
			if err != nil {
				return
			}
			// check response
			tc.checkResponse(t, createUserResp, err)
		})
	}
}
