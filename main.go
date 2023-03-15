package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

var words = []string{
	"satellite",
	"communication",
	"salvation",
	"eliminate",
	"mathematics",
	"ghostwriter",
	"dangerous",
	"girlfriend",
	"knowledge",
	"allocation",
	"conversation",
	"agreement",
	"convenience",
	"infection",
	"understand",
	"beginning",
	"photocopy",
	"battlefield",
	"performance",
	"houseplant",
	"conductor",
}

const clearScreenCharacter = "\033[H\033[2J"

var reader = bufio.NewReader(os.Stdin)

func main() {
	hangmanState := 0
	targetWord := getRandomWord(words)
	openCharacters := chooseDifficultyOfGame(targetWord)
	guessedLetters := getRandomMapOfCharacters(targetWord, openCharacters)

	for !isGameOver(targetWord, guessedLetters, hangmanState) {

		printGameState(targetWord, guessedLetters, hangmanState)

		userLetter := readUserInput()
		fmt.Print(clearScreenCharacter)
		if len([]rune(userLetter))-1 != 1 {
			fmt.Println("Неверный формат. Вы ввели больше одной буквы")
			continue
		}

		letter := rune(userLetter[0])
		if isCorrectLetter(targetWord, letter) {
			if !isLetterAlreadyUse(letter, guessedLetters) {
				guessedLetters[letter] = true
			} else {
				fmt.Println("Вы ввели букву, которая уже использована")
			}

		} else {
			hangmanState++
		}

	}

	fmt.Println("Игра закончена")
	if isWordGuessed(targetWord, guessedLetters) {
		fmt.Println("Поздравляю ты отгадал слово")
	} else if isHangmanComplete(hangmanState) {
		fmt.Println("Вы проиграли")
	} else {
		panic("Неизвестное состояние. Игра закончена но победителя нету")
	}
	fmt.Printf("Слово было %s\n\n", targetWord)

}

// Выбор случайного слова из массива слов
// Принимает массив слов
// Отдает слово
func getRandomWord(words []string) string {
	rand.Seed(time.Now().UnixNano())
	minIndex := 0
	maxIndex := len(words) - 1
	word := words[rand.Intn(maxIndex-minIndex+1)+minIndex]
	return word
}

// Показывает слово и угаданные буквы
// Принимает слово и словарь из угаданных букв
func printGameState(targetWord string, guessedLetters map[rune]bool, hangmanState int) {
	fmt.Println(getWordGuessingProgress(targetWord, guessedLetters))
	fmt.Println(getHangmanDrawing(hangmanState))
}

func getWordGuessingProgress(targetWord string, guessedLetters map[rune]bool) string {
	result := ""
	for _, ch := range targetWord {
		if guessedLetters[unicode.ToLower(ch)] == true {
			result += string(ch)
		} else if ch == ' ' {
			result += " "
		} else {
			result += "_"
		}
		result += " "
	}

	return result + "\n"
}

// Создает словарь с рандомно открытыми буквами
// Принимает загаданное слов и количество открытых букв
// Отдает словарь открытых букв
func getRandomMapOfCharacters(targetWord string, openCharacters int) map[rune]bool {
	rand.Seed(time.Now().UnixNano())
	minIndex := 0
	maxIndex := len(targetWord) - 1
	guessedLetters := map[rune]bool{}

	// Debug поменять условия цикла, ибо в случае когда все буквы открыты он бесконечен
	for i := 0; i != openCharacters; i++ {
		randomCharacter := rune(targetWord[rand.Intn(maxIndex-minIndex+1)+minIndex])
		if guessedLetters[unicode.ToLower(randomCharacter)] == false {
			guessedLetters[unicode.ToLower(randomCharacter)] = true
		} else {
			i--
			continue
		}
	}

	return guessedLetters
}

func getHangmanDrawing(hangmanState int) string {
	data, err := ioutil.ReadFile(fmt.Sprintf("states/hangman%v", hangmanState))
	if err != nil {
		panic(err)
	}

	return string(data)
}

func readUserInput() string {
	fmt.Print("Ваш выбор: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return input
}

func isCorrectLetter(targetWord string, letter rune) bool {
	return strings.ContainsRune(targetWord, letter)
}

func isWordGuessed(targetWord string, guessedLetters map[rune]bool) bool {
	for _, ch := range targetWord {
		if !guessedLetters[unicode.ToLower(ch)] {
			return false
		}
	}

	return true
}

func isHangmanComplete(hangmanState int) bool {
	return hangmanState >= 9
}

func isGameOver(targetWord string, guessedLetters map[rune]bool, hangmanState int) bool {
	return isWordGuessed(targetWord, guessedLetters) || isHangmanComplete(hangmanState)
}

func isLetterAlreadyUse(letter rune, guessedLetters map[rune]bool) bool {
	if guessedLetters[letter] {
		return true
	}

	return false
}

// Переписать под разные уровни сложности, от которых будет вычисляться
// количество открытых букв
func chooseDifficultyOfGame(targetWord string) int {
	for {
		fmt.Printf("Введите сложность игры\n1. Легкая \n2. Средняя \n3. Сложная \n4. Настоящая игра\n5. Разница в уровнях (help)\n")
		userInput := readUserInput()
		charactersInWord := utf8.RuneCountInString(targetWord)
		switch strings.ToLower(userInput) {
		case "1\n", "легкая\n":
			openCharacters := getPercentOfInt(charactersInWord, 50)
			return openCharacters
		case "2\n", "средняя\n":
			openCharacters := getPercentOfInt(charactersInWord, 30)
			return openCharacters
		case "3\n", "сложная\n":
			openCharacters := 1
			return openCharacters
		case "4\n", "настоящая игра\n":
			openCharacters := getPercentOfInt(charactersInWord, 0)
			return openCharacters
		case "5\n", "help\n":
			fmt.Print("Уровни сложности отличаются в количестве открытых букв.\n")
			fmt.Print("У легкого уровня это половина слова\n")
			fmt.Print("У среднего уровня это треть слова\n")
			fmt.Print("У сложного уровня одна буква\n")
			fmt.Print("У \"Настоящая игра\" уровня все слово скрыто\n")
		default:
			fmt.Println("Вы выбрали опцию, которой нету в списке.")
		}
	}
}

func getPercentOfInt(value int, percent int) int {
	return (value * percent) / 100
}
