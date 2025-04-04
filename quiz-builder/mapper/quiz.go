package mapper

import (
	"quiz-builder/entity"
)

func GetAllPublishedQuizzesToSolve(quizzes []entity.Quiz) []entity.QuizRequestResponse {
	var quizzesToSolve []entity.QuizRequestResponse
	for _, quiz := range quizzes {
		quizToSolve := entity.QuizRequestResponse{
			UUID:      quiz.UUID,
			Title:     quiz.Title,
			Questions: getQuestionsToSolve(quiz.Questions),
		}

		quizzesToSolve = append(quizzesToSolve, quizToSolve)
	}

	return quizzesToSolve
}

func getQuestionsToSolve(questions []entity.Question) []entity.QuestionRequestResponse {
	var questionsToSolve []entity.QuestionRequestResponse
	for _, question := range questions {
		questionToSolve := entity.QuestionRequestResponse{
			UUID:       question.UUID,
			Text:       question.Text,
			AnswerType: question.AnswerType,
			Answers:    getAnswersToSolve(question.Answers),
		}

		questionsToSolve = append(questionsToSolve, questionToSolve)
	}

	return questionsToSolve
}

func getAnswersToSolve(answers []entity.Answer) []entity.AnswerRequestResponse {
	var answersToSolve []entity.AnswerRequestResponse
	for _, answer := range answers {
		answerToSolve := entity.AnswerRequestResponse{
			UUID: answer.UUID,
			Text: answer.Text,
		}

		answersToSolve = append(answersToSolve, answerToSolve)
	}

	return answersToSolve
}
