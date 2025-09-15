//go:build integration
// +build integration

package integration_test

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"ama/api/application"
	"ama/api/auth"
	"ama/api/constants"
)

func GetQuestionBaseUrl(secure bool) string {
	s := ""
	if secure {
		s = "s"
	}
	return fmt.Sprintf(
		"http%s://%s:%d%s",
		s,
		ResourceServerHost,
		ResourceServerPort,
		constants.QuestionBasePath,
	)
}

func QuestionSuite(t *testing.T) {
	questionsToCreate := len(testQuestions)
	client := &http.Client{}
	authToken, err := adminSignIn(client)
	if err != nil {
		t.Errorf("failed to sign in: %v", err)
		return
	}
	questionIds := make([]string, questionsToCreate)
	for i := range questionsToCreate {
		questionId, err := createQuestion(i, authToken, *client)
		if err != nil {
			t.Errorf("failed to create question %d: %v", i, err)
			return
		}
		t.Logf("Created question %d with ID: %s", i, questionId)
		questionIds[i] = questionId
	}
}

func QuestionTearDownSuite(t *testing.T) {
	// TODO: Delete all of the questions!
}

func adminSignIn(httpClient *http.Client) (string, error) {
	authClient, err := auth.NewAuthClient()
	if err != nil {
		return "", err
	}
	token, err := authClient.CustomTokenWithClaims(
		context.Background(),
		"integration-test-client",
		constants.GetAdminScopes(),
	)
	type adminToken struct {
		Token             string `json:"token"`
		ReturnSecureToken bool   `json:"returnSecureToken"`
	}
	aToken := adminToken{
		Token:             token,
		ReturnSecureToken: true,
	}
	secure := ""
	if IsSecure {
		secure = "s"
	}
	var validToken ReturnedToken
	validToken, err = HitApi(
		httpClient,
		fmt.Sprintf(
			"http%s://%s:%d/identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken?key=%s",
			secure,
			EmulatorHost,
			EmulatorPort,
			ApiKey,
		),
		http.MethodPost,
		"",
		aToken,
		validToken,
	)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v; err: %w", validToken, err)
	}
	return validToken.IdToken, nil
}

func createQuestion(i int, t string, client http.Client) (string, error) {
	newQuestion := application.NewQuestion{
		Prompt: testQuestions[i],
		Tags:   []string{"test"},
	}
	var question application.Question
	question, err := HitApi(
		&client,
		GetQuestionBaseUrl(IsSecure),
		http.MethodPost,
		t,
		newQuestion,
		question,
	)
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %s; err: %w", &question, err)
	}
	return question.ID, nil
}

var testQuestions = []string{
	"What do you hope for?",
	"What do you want?",
	"Who do you want to be?",
	"What are you doing about who you want to be?",
	"What kind of old person do you want to be?",
	"What are you afraid of?",
	"What do you hate about yourself?",
	"What do you want to learn?",
	"Bad habit?",
	"Favorite quote, book, movie, tv show, song right now?",
	"What do you do often that you know your parents wouldn't  like?",
	"Which of your parents do you more resemble, and in what ways are you like them both?",
	"How has your upbringing influenced your worldview?",
	"Least favorite color and why",
	"Favorite animal and why?",
	"Favorite element (earth wind fire water) and why?",
	"What is God doing right now when he looks at you?",
	"Dream job?",
	"Dream place to live?",
	"Weirdest/scariest/most memorable dream?",
	"What are you hoping will happen soon/what are you waiting for?",
	"Biggest lie you ever told? (besides the one you're about to tell)",
	"Worst job?",
	"Best job?",
	"Worst time you got in trouble?",
	"Hardest/weirdest/dumbest thing you've done?",
	"Where do you see yourself?",
	"How do you see yourself?",
	"How would others describe you and how does it differ from what you'd say?",
	"Why do you get out of bed in the morning?",
	"How easily do you fall asleep?",
	"What keeps you up at nights?",
	"What kinds of things do you forget intentionally?",
	"What do you wish never happened?",
	"What do you lie about most?",
	"Biggest regret?",
	"Who do you love?",
	"Who did you love?",
	"Who did you think you loved/loved?you?",
	"Any unrequited love?",
	"What is love?",
	"What's the meaning of life?",
	"What is happiness?",
	"What's reality ultimately?",
	`What do you do?
				Why
				Why
				Why
				Why
				Why`,
	"What do you struggle with?",
	"What have you yelled at God over?",
	"Ever feel like you're breathing but not really feeling alive? Tell me about it.",
	"How often do you trust your instincts or your natural thought progressions or feelings you get about things?",
	"Snake in your boot or spider in your toilet?",
	"Did you ever have a life threatening experience and what happened inside your head during it?",
	"Rather have one career or a ton of jobs?",
	"Are you more afraid of rejection or being lied to about being accepted?",
	"Favorite thing?",
	"Hate that sound?",
	"Best fart?",
	"Best cry and most recent?",
	"Last story that got you emotional?",
	"Something really nice you did for someone in secret that they or at least nobody else knows about?",
	"Unlimited cash for a week, after that it disappears and so does everything you spent it on, what do you do?",
	"Recent meaningful experience?",
	"What song did you sing in the shower",
	"What's something you do that if people saw they'd think you're crazy?",
	"Define success?",
	"What situations do you shut down emotionally for?",
	"What aspect of interacting with other people are you afraid of?",
	"What question do you not want to be asked? (not this one)",
	"Something nobody else knows about you?",
	"What environment do you best express yourself?",
	"How do you deal with stress?",
	"What do you hope heaven is like?",
	"What's the title of the biography of your life?",
	"Last thing people find out about you after getting to know you?",
	"Do people understand you and why or why not?",
	"Who understands you the best and why?",
	"What are you thinkin?",
	"You're on death row and get one food, song, and movie. What are they?",
	"How would you describe your outfit using only sound effects?",
	"What was your first impression of me?",
	"Do your best impression of everyone in the room right now.",
	"If there was one thing you could immediately know about everyone you ever met before you met them, what would that be?",
	"What are the 10 things you value the most/if you could only own 10 things for the rest of your life what are they?",
	"When is the last time you really took time to appreciate silence?",
	"Would you rather have 5 dollars or your local sports team win the national championship?",
	"Ideal banana?",
	"Google me the picture that represents you.",
	"If you were holding your soul right now, what is the texture and feel of it?",
	"Where would you rather be right now?",
	"Whole or skim and why?",
	"Pancakes or waffles? How bout French toast? ",
	"Without naming them, describe your 5 closest friends?",
	"If tomorrow you reduced the bible to 3 books, which are they?",
	"How would your friends describe you in 3 words?",
	"Best piece of practical advice you've heard?",
	"Something you hate that everyone else loves?",
	"Last piece of art that moved you?",
	"What makes you tick?",
	"10 words that describe what you're all about?",
	"How much money are you worth if someone was going to buy you?",
	"What's something that your parents saw as valuable that you still think is stupid?",
	"When did you feel the most powerful in your life?",
	"You've been siting at a red light and it still hasn't changed, you're the only one there and you're trying to turn left, how long do you wait before going?",
	"What word would you remove from existence never to have been spoken or be spoken/written ever again?",
	"What word that doesn't exist would you add and what would it mean?",
	"What was the last thing that defied your ability to explain it in words?",
	"Best kiss?",
	"Best date?",
	"Ideal date?",
	"Non-negotiable qualities you look for in a significant other?",
	"You have $1000 and have to spend it all now.",
	"Name the 3 non-survival related items you'd bring on a desert island.",
	"What is the weirdest place you masturbated?",
	"What is the hardest thing for you to give up to God?",
	"What do you value the most in a friend and what in life made you value those things?",
	"Best compliment you can give someone?",
	"What's the most messed up thing you can think of?",
	"What moments from your past do you remember that make you cringe?",
	"What's your favorite word?",
	"Besides text and talk, how do you spend most of the time on your phone?",
	"What do you do with your phone that most people don't?",
	"What song would you have your loved ones play at your funeral if you died tomorrow?",
	"If you were going to take out a big loan to start a business, of the people you know, whose name would you want to have signed next to yours?",
	"Given the choice of anyone in the world, whom would you want as a dinner guest?",
	"Would you like to be famous? In what way?",
	"Before making a telephone call, do you ever rehearse what you are going to say? Why?",
	"What would constitute a “perfect” day for you?",
	"When did you last sing to yourself? To someone else?",
	"If you were able to live to the age of 90 and retain either the mind or body of a 30-year-old for the last 60 years of your life, which would you want?",
	"Do you have a secret hunch about how you will die?",
	"Name three things you and your partner appear to have in common.",
	"For what in your life do you feel most grateful?",
	"If you could change anything about the way you were raised, what would it be?",
	"In four minutes and tell me your life story in as much detail as possible.",
	"If you could wake up tomorrow having gained any one quality or ability, what would it be?",
	"If a crystal ball could tell you the truth about yourself, your life, the future or anything else, what would you want to know?",
	"Is there something that you've dreamed of doing for a long time? Why haven't you done it?",
	"What is the greatest accomplishment of your life?",
	"What do you value most in a friendship?",
	"What is your most treasured memory?",
	"What is your most terrible memory?",
	"If you knew that in one year you would die suddenly, would you change anything about the way you are now living? Why?",
	"What does friendship mean to you?",
	"What roles do love and affection play in your life?",
	"We alternate sharing something we'd consider a positive characteristic of each other. Share a total of five items.",
	"How close and warm is your family? Do you feel your childhood was happier than most other people's?",
	"How do you feel about your relationship with your mother?",
	`Make three true "we" statements for each person. For instance, "We are both in this room feeling ... "`,
	"Complete this sentence: “I wish I had someone with whom I could share ... “",
	"If you were going to become a close friend with me, please share what would be important for me to know.",
	"Tell your me what you like about me; be very honest this time, saying things that you might not say to someone you've just met. ",
	"Share an embarrassing moment in your life.",
	"When did you last cry in front of another person? By yourself?",
	"Share something that you like about me already.",
	"What, if anything, is too serious to be joked about?",
	"If you were to die this evening with no opportunity to communicate with anyone, what would you most regret not having told someone? Why haven’t you told them yet?",
	"Your house, containing everything you own, catches fire. After saving your loved ones and pets, you have time to safely make a final dash to save any one item. What would it be? Why?",
	"Of all the people in your family, whose death would you find most disturbing? Why?",
	"Share a personal problem and ask the other person's advice on how he or she might handle it. Also, ask them to reflect back to you how you seem to be feeling about the problem you have chosen.",
	"What inner experiences do you find difficult to communicate accurately with words to other people?",
	"Toilet paper, folded or wadded up?",
	"Why are you the way that you are?",
	"What's the non-material thing you're most afraid of losing?",
	"What thing if I knew about you would prevent us from being friends?",
	"If both of your parents died right now, what would you regret the most not having said to them while they were still alive?",
	"If you were writing a book about your relationship to God and divided the chapters chronologically, when would those chapter breaks occur and what would the chapters be called?",
	"How many times have you urinated in public?",
	"Tell me a fun fact.",
	"What percentage of people like you, and what percentage of people do you like?",
	"Show me the last 3 google searches you made",
	"Show me the most recent 3 pictures on your phone",
	"When you instinctively open your phone for no reason, what apps do you check and in what order assuming there are no notifications?",
	"How many bumper stickers is okay for a car to have?",
}

func ReadQuestions(client *http.Client, userToken string, tags []string, limit int, finalId string, random bool) ([]application.Question, error) {
	var questions []application.Question
	params := []string{}
	if len(tags) > 0 {
		for _, tag := range tags {
			params = append(params, fmt.Sprintf("%s=%s", constants.TagParam, tag))
		}
	}
	if limit != 0 {
		params = append(params, fmt.Sprintf("%s=%d", constants.LimitParam, limit))
	}
	if finalId != "" {
		params = append(params, fmt.Sprintf("%s=%s", constants.FinalIdParam, finalId))
	}
	if random {
		params = append(params, fmt.Sprintf("%s=true", constants.RandomParam))
	}
	url := GetQuestionBaseUrl(IsSecure)
	if len(params) > 0 {
		queryParams := strings.Join(params, "&")
		url = fmt.Sprintf("%s?%s", url, queryParams)
	}
	return HitApi(
		client,
		url,
		http.MethodGet,
		userToken,
		nil,
		questions,
	)
}
