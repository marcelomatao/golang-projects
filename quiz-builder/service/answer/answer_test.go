package answer

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"quiz-builder/entity"
	"quiz-builder/logger"
	"quiz-builder/repository/quiz_builder"
)

// Define a test suite

type AnswerServiceTestSuite struct {
	suite.Suite
	service    AnswerService
	mockRepo   *quiz_builder.MockQuizBuilderRepository
	mockLogger *logger.MockLogger
}

func (suite *AnswerServiceTestSuite) SetupTest() {
	suite.mockRepo = new(quiz_builder.MockQuizBuilderRepository)
	suite.mockLogger = new(logger.MockLogger)
	suite.service = NewAnswerService(suite.mockLogger, suite.mockRepo)
}

func (suite *AnswerServiceTestSuite) TestGetOneAnswer() {
	id := uuid.New()
	answer := entity.Answer{UUID: id}

	suite.mockRepo.On("Find", mock.Anything, "UUID = ?", id.String()).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*entity.Answer)
		*arg = answer
	})

	result, err := suite.service.GetOneAnswer(id)

	suite.NoError(err)
	suite.Equal(answer, result)
}

func (suite *AnswerServiceTestSuite) TestGetAllAnswers() {
	answers := []entity.Answer{{UUID: uuid.New()}, {UUID: uuid.New()}}

	suite.mockRepo.On("Find", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*[]entity.Answer)
		*arg = answers
	})

	result, err := suite.service.GetAllAnswers()

	suite.NoError(err)
	suite.Equal(answers, result)
}

func (suite *AnswerServiceTestSuite) TestCreateAnswer() {
	answer := entity.Answer{UUID: uuid.New()}

	suite.mockRepo.On("Create", &answer).Return(nil)

	err := suite.service.CreateAnswer(answer)

	suite.NoError(err)
}

func (suite *AnswerServiceTestSuite) TestUpdateAnswer() {
	answer := entity.Answer{UUID: uuid.New()}

	suite.mockRepo.On("Save", &answer).Return(nil)

	err := suite.service.UpdateAnswer(answer)

	suite.NoError(err)
}

func (suite *AnswerServiceTestSuite) TestDeleteAnswer() {
	id := uuid.New()

	suite.mockRepo.On("Delete", &entity.Answer{}, "UUID = ?", id.String()).Return(nil)

	err := suite.service.DeleteAnswer(id)

	suite.NoError(err)
}

func TestAnswerServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AnswerServiceTestSuite))
}
