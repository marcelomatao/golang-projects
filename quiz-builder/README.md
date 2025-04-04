## Apis

### Authentication

- Register new user
```bash
curl -X POST -H "Content-type: application/json" http://localhost:8080/api/auth/register -d '{
    "email": "test10@gmail.com",
	"name": "Marcelo",
	"password": "user123"
}'
```

- Login with a registered user
```bash
curl -X POST -H "Content-type: application/json" http://localhost:8080/api/auth/login -d '{
    "email": "test10@gmail.com",
	"password": "user123"
}'
```

### Quiz creation and publishing

- Create Quiz
```bash
curl -X POST -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz -d '{
	"title": "Quiz1"
}'
```

- Get all quizzes created by the authenticated owner.
```bash
curl -X GET -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/all
```

- Publishing a quiz. This operation can be performed only by the quiz owner. The quiz must have 1 to 10 question and each question must have 1 to 5 answers where at least one answer must be correct.
```bash
curl -X POST -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/88f0cc76-9fa1-45ca-920c-6a145f3ef7b6/publish
```

- Add a question to a quiz. This operation can be performed only by the quiz owner. The current number of questions created already should be at maximum 9 in order to have one more addition, once 10 questions is the maximum. There are 2 types of questions, Single and Multiple, the first one expects only one correct answer. The second one can have one or more correct answers and multiple answers can be selected. If the quiz is published this operation can no longer be done.
```bash
curl -X POST -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/88f0cc76-9fa1-45ca-920c-6a145f3ef7b6/question -d '{
	"answerType": "Multiple",
	"text": "Question 2"
}'
```

- Add a answer to a question. This operation can be performed only by the quiz owner. The current number of answers created already to this question should be at maximum 4 in order to have one more addition, once 5 answers is the maximum. If the question expects Single selected answer, then this question can have only one isCorrect answer. If the quiz is published this operation can no longer be done.
```bash
curl -X POST -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/question/8ceb85f2-bbcf-4b64-9ac0-2abca0fdcb0f/answer -d '{
	"isCorrect": true,
	"text": "Answer 2"
}'
```

- Get a quiz created by the authenticated owner. If this quiz is not created by the authenticated user this operation will return error. This will return the details of a quiz created by the authenticated user.
```bash
curl -X GET -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/88f0cc76-9fa1-45ca-920c-6a145f3ef7b6
```

- Delete a quiz created by the authenticated user.
```bash
curl -X DELETE -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/111dacd2-b10f-48b5-ba4b-a88c59a0fef6
```

- Delete a question. This operation can be performed only by the quiz owner. If the quiz is published this operation can no longer be done.
```bash
curl -X DELETE -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/question/cd54b565-de62-40d2-b988-faa7445992f8
```

- Delete an answer. This operation can be performed only by the quiz owner. If the quiz is published this operation can no longer be done.
```bash
curl -X DELETE -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/answer/cecf0438-225c-4413-ad7c-5310a131b93e
```


### Quiz solutions submission and querying

- Get all published quizzes. This operation can be performed by any authenticated user to get published quizzes that can be answered.
```bash
curl -X GET -H "Content-type: application/json" -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/published/all
```

- Submit a quiz solution. This operation can be performed by any authenticated user to submit an solution to a quiz. If the quiz was already solved by the user this operation will return error.
```bash
curl -X POST -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/88f0cc76-9fa1-45ca-920c-6a145f3ef7b6/solution/submit -d '{
	"questions_solution": [{"question_uuid": "17a53db6-cef7-459b-a36f-55ecf8e98b81","selected_answers": [{"answer_uuid":"7264fff0-5aa2-4391-9621-e29e10c8020c"}]}, {"question_uuid": "8ceb85f2-bbcf-4b64-9ac0-2abca0fdcb0f","selected_answers": [{"answer_uuid":"d4196799-8f82-414f-b67a-f22a4c5b4369"}]}]
}'
```

- Get all quiz solutions for the published owner. Only the authenticated user that created this quiz can use performe this operation.
```bash
curl -X GET -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/88f0cc76-9fa1-45ca-920c-6a145f3ef7b6/solutions/all
```

- Get all quizzes solutions submited by an user to quizzes created by other users. Any authenticated user can perform this operation in order to get the solutions to all quizzes this user submit a solution.
```bash
curl -X GET -H "Content-type: application/json" -H "Authorization: Bearer <token get when signin>" http://localhost:8080/api/quiz/solutions
```

## Dependencies 

```bash
go mod init quiz-builder

go get go.uber.org/fx

go get github.com/spf13/viper

go get gorm.io/gorm/logger

go get github.com/gin-gonic/gin

go get github.com/joho/godotenv

go get go.uber.org/fx@v1.20.0

go get go.uber.org/zap@v1.23.0

go get gorm.io/driver/mysql

go get github.com/dgrijalva/jwt-go

go get github.com/google/uuid

go get github.com/rs/cors/wrapper/gin
```