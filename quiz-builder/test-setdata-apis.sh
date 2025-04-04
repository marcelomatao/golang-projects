#!/bin/bash

: ${HOST=localhost}
: ${PORT=8080}


function registerUser() {
    local user=$1

    curl -X POST http://$HOST:$PORT/api/auth/register -H "Content-Type: application/json" --data "$user"
}

function createQuiz() {
  local quiz=$1
  local auth=$2

  curl -X POST http://$HOST:$PORT/api/quiz -H "Content-Type: application/json" $auth --data "$quiz"
}

function addQuestionToQuiz() {
  local quizUUID=$1
  local question=$2
  local auth=$3

  curl -X POST http://$HOST:$PORT/api/quiz/$quizUUID/question -H "Content-Type: application/json" $auth --data "$question"
}

function setupTestData() {
    random1=$((1 + RANDOM % 1000000))
    random2=$((1 + RANDOM % 1000000))
    random3=$((1 + RANDOM % 1000000))

    email1="test$random1@email.com"
    name1="Test 1"
    password1="pass1"

    email2="test$random2@email.com"
    name2="Test 2"
    password2="pass2"

    email3="test$random3@email.com"
    name3="Test 3"
    password3="pass3"

    bodyUser1="{\"email\":\"$email1\", \"name\":\"$name1\", \"password\":\"$password1\"}"
    bodyUser2="{\"email\":\"$email2\", \"name\":\"$name2\", \"password\":\"$password2\"}"
    bodyUser3="{\"email\":\"$email3\", \"name\":\"$name3\", \"password\":\"$password3\"}"
    registerUser "$bodyUser1"
    registerUser "$bodyUser2"
    registerUser "$bodyUser3"

    bodyLogin1="{\"email\":\"$email1\", \"password\":\"$password1\"}"
    bodyLogin2="{\"email\":\"$email2\", \"password\":\"$password2\"}"
    bodyLogin3="{\"email\":\"$email3\", \"password\":\"$password3\"}"
    tokenUser1=$(curl -X POST http://$HOST:$PORT/api/auth/login -H "Content-Type: application/json" --data "$bodyLogin1" | jq .token -r)
    tokenUser2=$(curl -X POST http://$HOST:$PORT/api/auth/login -H "Content-Type: application/json" --data "$bodyLogin2" | jq .token -r)
    tokenUser3=$(curl -X POST http://$HOST:$PORT/api/auth/login -H "Content-Type: application/json" --data "$bodyLogin3" | jq .token -r)
    authUser1="-H \"Authorization: Bearer $tokenUser1\""
    authUser2="-H \"Authorization: Bearer $tokenUser2\""
    authUser3="-H \"Authorization: Bearer $tokenUser3\""

    quiz1="Quiz 1"
    quiz2="Quiz 2"
    bodyQuiz1="{\"title\":\"$quiz1\"}"
    bodyQuiz2="{\"title\":\"$quiz2\"}"
    createQuiz "$bodyQuiz1" "$authUser1"
    createQuiz "$bodyQuiz2" "$authUser2"

    user1Quiz1UUID=$(curl -X GET http://$HOST:$PORT/api/quiz/all -H "Content-Type: application/json" $authUser1 | jq .data[0].uuid -r)
    user2Quiz1UUID=$(curl -X GET http://$HOST:$PORT/api/quiz/all -H "Content-Type: application/json" $authUser2 | jq .data[0].uuid -r)

    question1="Question 1"
    answerType1="Single"
    question2="Question 2"
    answerType2="Single"
    question3="Question 3"
    answerType3="Multiple"
    question4="Question 4"
    answerType4="Multiple"
    question5="Question 5"
    answerType5="Single"

    question1User1Quiz1="{\"text\":\"$question1\",\"answerType\":\"$answerType1\"}"
    addQuestionToQuiz "$user1Quiz1UUID" "$question1User1Quiz1" "$authUser1"
    question2User1Quiz1="{\"text\":\"$question2\",\"answerType\":\"$answerType2\"}"
    addQuestionToQuiz "$user1Quiz1UUID" "$question2User1Quiz1" "$authUser1"
    question3User1Quiz1="{\"text\":\"$question3\",\"answerType\":\"$answerType3\"}"
    addQuestionToQuiz "$user1Quiz1UUID" "$question3User1Quiz1" "$authUser1"

    question1User1Quiz1UUID=$(curl -X GET http://$HOST:$PORT/api/quiz/all -H "Content-Type: application/json" $authUser1 | jq .data[0].questions[0].uuid -r)
    question2User1Quiz1UUID=$(curl -X GET http://$HOST:$PORT/api/quiz/all -H "Content-Type: application/json" $authUser1 | jq .data[0].questions[1].uuid -r)
    question3User1Quiz1UUID=$(curl -X GET http://$HOST:$PORT/api/quiz/all -H "Content-Type: application/json" $authUser1 | jq .data[0].questions[2].uuid -r)

    question1User2Quiz1="{\"text\":\"$question1\",\"answerType\":\"$answerType1\"}"
    addQuestionToQuiz "$user2Quiz1UUID" "$question1User2Quiz1" "$authUser2"
    question2User2Quiz1="{\"text\":\"$question2\",\"answerType\":\"$answerType2\"}"
    addQuestionToQuiz "$user2Quiz1UUID" "$question2User2Quiz1" "$authUser2"
    question3User2Quiz1="{\"text\":\"$question3\",\"answerType\":\"$answerType3\"}"
    addQuestionToQuiz "$user2Quiz1UUID" "$question3User2Quiz1" "$authUser2"
    question4User2Quiz1="{\"text\":\"$question4\",\"answerType\":\"$answerType4\"}"
    addQuestionToQuiz "$user2Quiz1UUID" "$question4User2Quiz1" "$authUser2"
    question5User2Quiz1="{\"text\":\"$question5\",\"answerType\":\"$answerType5\"}"
    addQuestionToQuiz "$user1Quiz1UUID" "$question5User2Quiz1" "$authUser2"

    question1User2Quiz1UUID=$(curl -X GET http://$HOST:$PORT/api/quiz/all -H "Content-Type: application/json" $authUser2 | jq .data[0].questions[0].uuid -r)
    question2User2Quiz1UUID=$(curl -X GET http://$HOST:$PORT/api/quiz/all -H "Content-Type: application/json" $authUser2 | jq .data[0].questions[1].uuid -r)
    question3User2Quiz1UUID=$(curl -X GET http://$HOST:$PORT/api/quiz/all -H "Content-Type: application/json" $authUser2 | jq .data[0].questions[2].uuid -r)
    question4User2Quiz1UUID=$(curl -X GET http://$HOST:$PORT/api/quiz/all -H "Content-Type: application/json" $authUser2 | jq .data[0].questions[3].uuid -r)
    question5User2Quiz1UUID=$(curl -X GET http://$HOST:$PORT/api/quiz/all -H "Content-Type: application/json" $authUser2 | jq .data[0].questions[4].uuid -r)
  

}

echo "Start:" `date`

echo "HOST=${HOST}"
echo "PORT=${PORT}"

setupTestData