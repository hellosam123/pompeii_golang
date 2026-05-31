// global variables
let currentVocab = null;
let nextVocab = null;
let isCorrect = false;
let correctAddTime = 0;
let correctScore = 0;
let maxGameTime = 0;
let gameTime = 0;
let gameScore = 0;
let gameOver = false;
let gameEnded = false;
let isProcessing = false;
let countdownFrame = undefined;
const fire = "🔥";
const form = document.getElementById("answer-form");
const userAnswer = document.getElementById("answer-input");
const vocabWordElement = document.getElementById("vocab-word");
const gameScoreElement = document.getElementById("game-score");
const timerElement = document.getElementById("timer");
const fireElement = document.getElementById("fire");
const gameImgElement = document.getElementById("game-img");
const difficultyModeContainer = document.getElementById("difficulty-mode-info");
const difficultyMode = difficultyModeContainer.dataset.difficulty_mode ?? "easy";
const BASE_URL = "pompeii";
async function fetchVocab() {
    // returns a new vocab {VocabID: int, VocabWord: string, EnglishTranslation: string, VocabGroup: string}
    if (gameOver || gameEnded) {
        return null;
    }
    const response = await fetch(`/${BASE_URL}/get_random_vocab`);
    const data = await response.json();
    if (data.gameOver) {
        gameOver = true;
        return null;
    }
    else {
        return data.vocab;
    }
}
async function fetchNextVocab() {
    // fetch the next vocab and update nextVocab variable
    if (gameOver || gameEnded) {
        return;
    }
    nextVocab = await fetchVocab();
}
async function fetchGameScore() {
    // returns game score
    const response = await fetch(`/${BASE_URL}/get_game_score`);
    const data = await response.json();
    return data.gameScore;
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
    if (gameOver || gameEnded) {
        return;
    }
    try {
        if (nextVocab) {
            // checks if the next vocab is ready
            currentVocab = nextVocab;
            nextVocab = null;
        }
        else {
            currentVocab = await fetchVocab();
        }
        if (currentVocab) {
            vocabWordElement.textContent = currentVocab.VocabWord;
        }
        await fetchNextVocab();
    }
    catch (error) {
        // try-catch immediately stops the rest of the code from running
        if (error instanceof Error && error.message == "GAME_OVER") {
            return;
        }
        console.error("unexpected error:", error);
    }
}
function setDifficultySettings() {
    // sets difficulty variables
    switch (difficultyMode) {
        case "easy":
            correctAddTime = 40;
            correctScore = 50;
            maxGameTime = 200;
            break;
        case "medium":
            correctAddTime = 30;
            correctScore = 100;
            maxGameTime = 150;
            break;
        case "hard":
            correctAddTime = 20;
            correctScore = 200;
            maxGameTime = 100;
            break;
    }
    gameTime = maxGameTime;
}
function loadFrame() {
    // loads the timer, score, and fire
    if (gameEnded) {
        return;
    }
    if (gameTime <= 0) {
        gameOver = true;
        checkGameEnd();
        return;
    }
    gameTime--;
    if (isCorrect) {
        timerElement.innerHTML = `Time Left:<br>${(gameTime / 10).toFixed(1)} (+${(correctAddTime / 10).toFixed(1)})`;
        gameScoreElement.textContent = `Score: ${gameScore} (+${correctScore})`;
    }
    else {
        timerElement.innerHTML = `Time Left:<br>${(gameTime / 10).toFixed(1)}`;
        gameScoreElement.textContent = `Score: ${gameScore}`;
    }
    let numFire = Math.floor(10 - (gameTime * 10) / maxGameTime);
    let fireString = fire.repeat(numFire);
    fireElement.textContent = fireString;
}
function startLoad() {
    // starts the countdown timer
    countdownFrame = setInterval(loadFrame, 100);
}
function endLoad() {
    // ends the countdown timer
    clearInterval(countdownFrame);
}
function updateFrame(isCorrect) {
    // updates the frame on the client, then creates a fetch request to check if the client matches the server
    if (gameEnded) {
        return;
    }
    if (isCorrect) {
        gameScore += correctScore;
        if (gameTime + correctAddTime <= maxGameTime) {
            gameTime += correctAddTime;
        }
        else {
            gameTime = maxGameTime;
        }
    }
    fetchGameScore().then((scoreData) => {
        // checks with server to ensure sync
        if (scoreData != gameScore) {
            gameScore = scoreData;
        }
    });
}
function checkGameEnd() {
    // checks gameOver condition, and throws an error if true
    if (gameOver && !gameEnded) {
        gameEnded = true;
        endLoad();
        currentVocab = nextVocab = null;
        vocabWordElement.textContent = "GAME OVER";
        gameImgElement.style.display = "none";
        userAnswer.disabled = true;
        setTimeout(redirectGameOver, 2000);
        throw new Error("GAME_OVER");
    }
}
function redirectGameOver() {
    // redirects to game over page
    window.location.href = `/${BASE_URL}/game_over`;
}
setDifficultySettings();
loadVocab(); // fetch on page load
document.addEventListener("DOMContentLoaded", () => {
    // 500ms delay before starting timer
    timerElement.innerHTML = `Time Left:<br>${(gameTime / 10).toFixed(1)}`;
    setTimeout(startLoad, 500);
});
form.addEventListener("submit", async function (event) {
    event.preventDefault(); // prevents refresh
    if (gameOver || gameEnded) {
        return;
    }
    if (isProcessing) {
        return; // prevents any double submissions which could desync the client
    }
    const userAnswerInput = userAnswer.value;
    if (userAnswerInput) {
        isProcessing = true;
        try {
            let answerId;
            if (currentVocab) {
                answerId = currentVocab.VocabID;
            }
            else {
                throw new Error("currentVocab not found");
            }
            userAnswer.value = "";
            await loadVocab();
            if (gameOver || gameEnded) {
                return;
            }
            isCorrect = await checkAnswer(answerId, userAnswerInput);
            updateFrame(isCorrect);
        }
        finally {
            // try-finally used to ensure isProcessing is never locked to true
            isProcessing = false;
        }
    }
});
export {};
