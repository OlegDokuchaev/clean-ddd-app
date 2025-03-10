package repository

import (
	"context"
	"testing"
	domain "warehouse/internal/domain/common"
	outboxDomain "warehouse/internal/domain/outbox"
	productDomain "warehouse/internal/domain/product"
	"warehouse/internal/infrastructure/db"
	"warehouse/internal/infrastructure/db/migrations"
	outboxRepository "warehouse/internal/infrastructure/repository/outbox"
	"warehouse/internal/tests/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OutboxRepositoryTestSuite struct {
	suite.Suite
	ctx    context.Context
	testDB *testutils.TestDB
}

func (s *OutboxRepositoryTestSuite) SetupSuite() {
	config, err := db.NewConfig()
	require.NoError(s.T(), err)

	mConfig, err := migrations.NewConfig()
	require.NoError(s.T(), err)

	s.ctx = context.Background()

	testDB, err := testutils.NewTestDB(s.ctx, config, mConfig)
	require.NoError(s.T(), err)
	s.testDB = testDB

	err = s.testDB.Migrations.Up()
	require.NoError(s.T(), err)
}

func (s *OutboxRepositoryTestSuite) TearDownSuite() {
	if s.testDB != nil {
		err := s.testDB.Container.Terminate(s.ctx)
		require.NoError(s.T(), err)
	}
}

func (s *OutboxRepositoryTestSuite) getRepo() outboxDomain.Repository {
	return outboxRepository.New(s.testDB.DB)
}

func (s *OutboxRepositoryTestSuite) createTestOutboxMessage() *outboxDomain.Message {
	message, err := outboxDomain.Create(domain.NewEvent[productDomain.CreatedPayload, productDomain.CreateEvent](
		productDomain.CreatedPayload{
			ProductID: uuid.New(),
		},
	))
	require.NoError(s.T(), err)
	return message
}

func (s *OutboxRepositoryTestSuite) TestCreate() {
	tests := []struct {
		name          string
		setup         func(repo outboxDomain.Repository) *outboxDomain.Message
		expectedError error
	}{
		{
			name: "Success",
			setup: func(_ outboxDomain.Repository) *outboxDomain.Message {
				return s.createTestOutboxMessage()
			},
			expectedError: nil,
		},
		{
			name: "Failure: Message already exists",
			setup: func(repo outboxDomain.Repository) *outboxDomain.Message {
				message := s.createTestOutboxMessage()
				err := repo.Create(s.ctx, message)
				require.NoError(s.T(), err)
				return message
			},
			expectedError: outboxRepository.ErrOutboxMessageAlreadyExists,
		},
	}

	repo := s.getRepo()
	for _, test := range tests {
		s.Run(test.name, func() {
			message := test.setup(repo)

			err := repo.Create(s.ctx, message)

			if test.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), test.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				createdMessage, err := repo.GetByID(s.ctx, message.ID)
				require.NoError(s.T(), err)
				require.Equal(s.T(), message.ID, createdMessage.ID)
			}
		})
	}
}

func (s *OutboxRepositoryTestSuite) TestGetByID() {
	tests := []struct {
		name          string
		setup         func(repo outboxDomain.Repository) *outboxDomain.Message
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo outboxDomain.Repository) *outboxDomain.Message {
				message := s.createTestOutboxMessage()
				err := repo.Create(s.ctx, message)
				require.NoError(s.T(), err)
				return message
			},
			expectedError: nil,
		},
		{
			name: "Failure: Message not found",
			setup: func(_ outboxDomain.Repository) *outboxDomain.Message {
				return s.createTestOutboxMessage()
			},
			expectedError: outboxRepository.ErrOutboxMessageNotFound,
		},
	}

	repo := s.getRepo()
	for _, test := range tests {
		s.Run(test.name, func() {
			message := test.setup(repo)

			retrievedMessage, err := repo.GetByID(s.ctx, message.ID)

			if test.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), test.expectedError, err)
			} else {
				require.NoError(s.T(), err)
				require.Equal(s.T(), message.ID, retrievedMessage.ID)
			}
		})
	}
}

func (s *OutboxRepositoryTestSuite) TestGetAll() {
	tests := []struct {
		name          string
		setup         func(repo outboxDomain.Repository) []uuid.UUID
		expectedError error
	}{
		{
			name: "Success: One message",
			setup: func(repo outboxDomain.Repository) []uuid.UUID {
				message := s.createTestOutboxMessage()
				err := repo.Create(s.ctx, message)
				require.NoError(s.T(), err)
				return []uuid.UUID{message.ID}
			},
			expectedError: nil,
		},
		{
			name: "Success: Multiple messages",
			setup: func(repo outboxDomain.Repository) []uuid.UUID {
				var messageIDs []uuid.UUID
				for i := 0; i < 3; i++ {
					message := s.createTestOutboxMessage()
					err := repo.Create(s.ctx, message)
					require.NoError(s.T(), err)
					messageIDs = append(messageIDs, message.ID)
				}
				return messageIDs
			},
			expectedError: nil,
		},
		{
			name: "Success: No messages",
			setup: func(_ outboxDomain.Repository) []uuid.UUID {
				return []uuid.UUID{}
			},
			expectedError: nil,
		},
	}

	repo := s.getRepo()
	for _, test := range tests {
		s.Run(test.name, func() {
			messageIDs := test.setup(repo)

			retrievedMessages, err := repo.GetAll(s.ctx)

			if test.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), test.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				retrievedMessageIDs := make([]uuid.UUID, 0, len(retrievedMessages))
				for _, message := range retrievedMessages {
					retrievedMessageIDs = append(retrievedMessageIDs, message.ID)
				}
				for _, messageID := range messageIDs {
					require.Contains(s.T(), retrievedMessageIDs, messageID)
				}
			}
		})
	}
}

func (s *OutboxRepositoryTestSuite) TestDelete() {
	tests := []struct {
		name          string
		setup         func(repo outboxDomain.Repository) *outboxDomain.Message
		expectedError error
	}{
		{
			name: "Success",
			setup: func(repo outboxDomain.Repository) *outboxDomain.Message {
				message := s.createTestOutboxMessage()
				err := repo.Create(s.ctx, message)
				require.NoError(s.T(), err)
				return message
			},
			expectedError: nil,
		},
		{
			name: "Failure: Message not found",
			setup: func(_ outboxDomain.Repository) *outboxDomain.Message {
				return s.createTestOutboxMessage()
			},
			expectedError: outboxRepository.ErrOutboxMessageNotFound,
		},
	}

	repo := s.getRepo()
	for _, test := range tests {
		s.Run(test.name, func() {
			message := test.setup(repo)

			err := repo.Delete(s.ctx, message)

			if test.expectedError != nil {
				require.Error(s.T(), err)
				require.Equal(s.T(), test.expectedError, err)
			} else {
				require.NoError(s.T(), err)

				retrievedMessage, err := repo.GetByID(s.ctx, message.ID)
				require.Error(s.T(), err)
				require.Equal(s.T(), outboxRepository.ErrOutboxMessageNotFound, err)
				require.Nil(s.T(), retrievedMessage)
			}
		})
	}
}

func TestOutboxRepository(t *testing.T) {
	suite.Run(t, new(OutboxRepositoryTestSuite))
}
