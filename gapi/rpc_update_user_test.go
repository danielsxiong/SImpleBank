package gapi

import (
	"context"
	mockdb "danielsxiong/simplebank/db/mock"
	db "danielsxiong/simplebank/db/sqlc"
	"danielsxiong/simplebank/pb"
	"danielsxiong/simplebank/token"
	"danielsxiong/simplebank/util"
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"testing"
	"time"
)

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)

	newName := util.RandomOwner()
	newEmail := util.RandomEmail()

	testCases := []struct {
		name          string
		body          *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, response *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			body: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: newName,
				Email:    newEmail,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: sql.NullString{
						String: newName,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
				}

				updatedUser := db.User{
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					FullName:          newName,
					Email:             newEmail,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
					IsEmailVerified:   user.IsEmailVerified,
				}

				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				accessToken, _, err := tokenMaker.CreateToken(user.Username, time.Minute)
				require.NoError(t, err)
				bearerToken := fmt.Sprintf("%s %s", authorizationBearer, accessToken)
				md := metadata.MD{
					authorizationHeader: []string{
						bearerToken,
					},
				}
				return metadata.NewIncomingContext(context.Background(), md)
			},
			checkResponse: func(t *testing.T, response *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, response)
				updatedUser := response.GetUser()
				require.Equal(t, user.Username, updatedUser.Username)
				require.Equal(t, newName, updatedUser.FullName)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)
			// build stubs
			tc.buildStubs(store)

			server := newTestServer(t, store, nil)

			updatedUserResp, err := server.UpdateUser(tc.buildContext(t, server.tokenMaker), tc.body)
			if err != nil {
				return
			}
			// check response
			tc.checkResponse(t, updatedUserResp, err)
		})
	}
}
