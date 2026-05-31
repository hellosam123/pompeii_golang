// global variables
let currentVocab;
let nextVocab;
let result;
let numQuestions = 0;
let numCorrect = 0;
let userCheckAnswer = false;
let gameOver = false;
let isProcessing = false;
const form = document.getElementById("answer-form");
const userAnswer = document.getElementById("answer-input");
const vocabWordElement = document.getElementById("vocab-word");
const englishTranslationElement = document.getElementById("english-translation");
const feedbackElement = document.getElementById("feedback");
const scoreElement = document.getElementById("score");
const BASE_URL = "pompeii";

async function fetchVocab() {
	// returns a new vocab {VocabID: int, VocabWord: string, EnglishTranslation: string, VocabGroup: string}

	const response = await fetch(`/${BASE_URL}/get_vocab`);
	const data = await response.json();

	if (data.gameOver) {
		gameOver = true;
		return null;
	} else {
		return data.vocab;
	}
}

async function fetchNextVocab() {
	// fetch the next vocab and update nextVocab variable

	nextVocab = await fetchVocab();
}

async function fetchScore() {
	// returns num_questions and num_correct

	const response = await fetch(`/${BASE_URL}/get_score`);
	const data = await response.json();

	return data;
}

async function checkAnswer(vocabId, userAnswer) {
	// returns a boolean value based on the user input

	// encodeURIComponent required to prevent some errors
	const url = `/${BASE_URL}/check_answer?vocab_id=${vocabId}&answer=${encodeURIComponent(userAnswer)}`;

	const response = await fetch(url);
	const data = await response.json();
	return data.isCorrect;
}

async function loadVocab() {
	// loads the vocab by using the next vocab variables, or if they are unavailable run fetchVocabId
	// have to wrap in a function to use async await functionality
	try {
		checkGameEnd(); // gameOver should be true if fetchNextVocab was run, prevents unnecessary code from being run

		if (nextVocab) {
			// checks if the next vocab is ready
			currentVocab = nextVocab;

			nextVocab = null;
		} else {
			currentVocab = await fetchVocab();

			checkGameEnd();
		}

		vocabWordElement.textContent = currentVocab.VocabWord;
		englishTranslationElement.textContent = currentVocab.EnglishTranslation;

		await fetchNextVocab();
	} catch (error) {
		// try-catch immediately stops the rest of the code from running
		if (error.message == "GAME_OVER") {
			return;
		}
	}
}

function loadResult(result) {
	// loads the feedback
	// updates and loads the score on the client, then creates a fetch request to check if the client matches the server

	if (result) {
		feedbackElement.textContent = "✅Correct";
		numCorrect++;
	} else {
		feedbackElement.textContent = "❌Incorrect";
	}
	numQuestions++;

	scoreElement.textContent = `${numCorrect} / ${numQuestions}`;

	fetchScore().then((scoreData) => {
		// checks with server to ensure sync

		if (scoreData.numCorrect != numCorrect || scoreData.numQuestions != numQuestions) {
			numCorrect = scoreData.numCorrect;
			numQuestions = scoreData.numQuestions;
			scoreElement.textContent = `${numCorrect} / ${numQuestions}`;
		}
	});
}

function checkGameEnd() {
	// checks gameOver condition, and throws an error if true
	if (gameOver) {
		window.location.href = `/${BASE_URL}/game_over`;
		throw new Error("GAME_OVER");
	}
}

loadVocab(); // fetch on page load

form.addEventListener("submit", async function (event) {
	event.preventDefault(); // prevents refresh

	if (isProcessing) {
		return; // prevents any double submissions which could desync the client
	}

	const userAnswerInput = userAnswer.value;

	if (!userCheckAnswer && userAnswerInput) {
		// cleaner than using nested if statements
		isProcessing = true;
		try {
			result = await checkAnswer(currentVocab.VocabID, userAnswerInput);
			loadResult(result);

			userCheckAnswer = true;
			feedbackElement.classList.toggle("visible");
			englishTranslationElement.classList.toggle("visible"); // cool trick I learned from website project
		} finally {
			// try-finally used to ensure isProcessing is never locked to true
			isProcessing = false;
		}
	} else if (userCheckAnswer) {
		isProcessing = true;
		try {
			userCheckAnswer = false;
			feedbackElement.classList.toggle("visible");
			englishTranslationElement.classList.toggle("visible");
			userAnswer.value = "";
			feedbackElement.textContent = "";
			await loadVocab();
		} finally {
			isProcessing = false;
		}
	}
});
