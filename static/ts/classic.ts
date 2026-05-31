export {};

interface Vocab {
	VocabID: number;
	VocabWord: string;
	EnglishTranslation: string;
	VocabGroup: string;
}

interface Score {
	numQuestions: number;
	numCorrect: number;
}

// global variables
let currentVocab: Vocab | null = null;
let nextVocab: Vocab | null = null;
let isCorrect: boolean = false;
let numQuestions: number = 0;
let numCorrect: number = 0;
let userCheckAnswer: boolean = false;
let gameOver: boolean = false;
let isProcessing: boolean = false;
const form = document.getElementById("answer-form") as HTMLElement;
const userAnswer: HTMLInputElement = document.getElementById("answer-input") as HTMLInputElement;
const vocabWordElement = document.getElementById("vocab-word") as HTMLElement;
const englishTranslationElement = document.getElementById("english-translation") as HTMLElement;
const feedbackElement = document.getElementById("feedback") as HTMLElement;
const scoreElement = document.getElementById("score") as HTMLElement;
const BASE_URL: string = "pompeii";

async function fetchVocab(): Promise<Vocab | null> {
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

async function fetchNextVocab(): Promise<void> {
	// fetch the next vocab and update nextVocab variable

	nextVocab = await fetchVocab();
}

async function fetchScore(): Promise<Score> {
	// returns num_questions and num_correct

	const response = await fetch(`/${BASE_URL}/get_score`);
	const data = await response.json();

	return data;
}

async function checkAnswer(vocabId: number, userAnswer: string): Promise<boolean> {
	// returns a boolean value based on the user input

	// encodeURIComponent required to prevent some errors
	const url = `/${BASE_URL}/check_answer?vocab_id=${vocabId}&answer=${encodeURIComponent(userAnswer)}`;

	const response = await fetch(url);
	const data = await response.json();
	return data.isCorrect;
}

async function loadVocab(): Promise<void> {
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

		if (currentVocab) {
			vocabWordElement.textContent = currentVocab.VocabWord;
			englishTranslationElement.textContent = currentVocab.EnglishTranslation;
		} else {
			throw new Error("currentVocab not found");
		}

		await fetchNextVocab();
	} catch (error: unknown) {
		// try-catch immediately stops the rest of the code from running
		if (error instanceof Error && error.message == "GAME_OVER") {
			return;
		}
		console.error("unexpected error:", error);
	}
}

function loadResult(isCorrect: boolean): void {
	// loads the feedback
	// updates and loads the score on the client, then creates a fetch request to check if the client matches the server

	if (isCorrect) {
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

function checkGameEnd(): void {
	// checks gameOver condition, and throws an error if true
	if (gameOver) {
		window.location.href = `/${BASE_URL}/game_over`;
		throw new Error("GAME_OVER");
	}
}

loadVocab(); // fetch on page load

form.addEventListener("submit", async function (event): Promise<void> {
	event.preventDefault(); // prevents refresh

	if (isProcessing) {
		return; // prevents any double submissions which could desync the client
	}

	const userAnswerInput = userAnswer.value;

	if (!userCheckAnswer && userAnswerInput) {
		// cleaner than using nested if statements
		isProcessing = true;
		try {
			if (currentVocab) {
				isCorrect = await checkAnswer(currentVocab.VocabID, userAnswerInput);
			} else {
				throw new Error("currentVocab not found");
			}
			loadResult(isCorrect);

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
