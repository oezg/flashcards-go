package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var history strings.Builder

type Flashcard struct {
	Term       string
	Definition string
	Mistakes   int
}

type Deck struct {
	cards []Flashcard
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	deck := Deck{}
	importFile := flag.String("import_from", "", "file name to import the initial card set from.")
	exportFile := flag.String("export_to", "", "file name to write into all cards after exit.")
	flag.Parse()
	if len(*importFile) > 0 {
		deck.ImportFile(*importFile)
	}
	for {
		action := readAction(reader)
		switch action {
		case "add":
			deck.Add(reader)
		case "remove":
			deck.Remove(reader)
		case "import":
			deck.Import(reader)
		case "export":
			deck.Export(reader)
		case "ask":
			deck.Ask(reader)
		case "log":
			Log(reader)
		case "hardest card":
			deck.Hardest()
		case "reset stats":
			deck.Reset()
		case "exit":
			if len(*exportFile) > 0 {
				deck.ExportFile(*exportFile)
			}
			fmt.Println("Bye bye!")
			return
		}
	}
}

func (deck *Deck) ImportFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		printLog("File not found.")
		return
	}
	var cardsFromFile []Flashcard
	err = json.Unmarshal(data, &cardsFromFile)
	if err != nil {
		log.Fatal(err)
	}
	printLog("%d cards have been loaded.", len(cardsFromFile))
	deck.Update(cardsFromFile)
}

func (deck *Deck) ExportFile(filename string) {
	data, err := json.Marshal(deck.cards)
	if err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile(filename, data, 0644); err != nil {
		log.Fatal(err)
	}
	printLog("%d cards have been saved.", len(deck.cards))
}

func readAction(reader *bufio.Reader) string {
	printLog("Input the action (add, remove, import, export, ask, exit, log, hardest card, reset stats):")
	return readLine(reader)
}

func printLog(message string, parameters ...interface{}) {
	message += "\n"
	_, err := fmt.Fprintf(&history, message, parameters...)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf(message, parameters...)
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func (deck *Deck) Add(reader *bufio.Reader) {
	card := deck.ReadCard(reader)
	deck.cards = append(deck.cards, card)
	printLog("The pair (\"%s\":\"%s\") has been added.", card.Term, card.Definition)
}

func (deck *Deck) ReadCard(reader *bufio.Reader) (card Flashcard) {
	card.Term = deck.ReadSide(reader, false)
	card.Definition = deck.ReadSide(reader, true)
	return
}

func (deck *Deck) ReadSide(reader *bufio.Reader, definition bool) string {
	if definition {
		printLog("The definition of the card:")
	} else {
		printLog("The card:")
	}
	for {
		text := readLine(reader)
		if alreadyExists, _ := deck.Contains(text, definition); !alreadyExists {
			return text
		}
		side := "card"
		if definition {
			side = "definition"
		}
		printLog("The %s \"%s\" already exists. Try again:", side, text)
	}
}

func (deck *Deck) Contains(text string, definition bool) (bool, int) {
	for i, card := range deck.cards {
		if definition && card.Definition == text {
			return true, i
		}
		if !definition && card.Term == text {
			return true, i
		}
	}
	return false, -1
}

func (deck *Deck) Remove(reader *bufio.Reader) {
	printLog("Which card?")
	text := readLine(reader)
	if exists, index := deck.Contains(text, false); !exists {
		printLog("Can't remove \"%s\": there is no such card.", text)
	} else {
		deck.cards[index] = deck.cards[len(deck.cards)-1]
		deck.cards = deck.cards[:len(deck.cards)-1]
		printLog("The card has been removed.")
	}
}

func (deck *Deck) Import(reader *bufio.Reader) {
	filename := readFilename(reader)
	deck.ImportFile(filename)
}

func readFilename(reader *bufio.Reader) string {
	printLog("File name:")
	return readLine(reader)
}

func (deck *Deck) Export(reader *bufio.Reader) {
	filename := readFilename(reader)
	deck.ExportFile(filename)
}

func (deck *Deck) Update(newCards []Flashcard) {
	for _, card := range newCards {
		if exists, index := deck.Contains(card.Term, false); !exists {
			deck.cards = append(deck.cards, card)
		} else {
			deck.cards[index].Definition = card.Definition
			deck.cards[index].Mistakes = card.Mistakes
		}
	}
}

func (deck *Deck) Ask(reader *bufio.Reader) {
	if len(deck.cards) == 0 {
		printLog("There are no cards available in memory")
	}
	printLog("How many times to ask?")
	number := readNumber(reader)
	for i := 0; i < number; i++ {
		deck.TestCard(reader, i%len(deck.cards))
	}
}

func readNumber(reader *bufio.Reader) int {
	line := readLine(reader)
	number, err := strconv.Atoi(line)
	if err != nil {
		log.Fatal(err)
	}
	return number
}

func (deck *Deck) TestCard(reader *bufio.Reader, index int) {
	card := deck.cards[index]
	printLog("Print the definition of \"%s\":", card.Term)
	text := readLine(reader)
	if text == card.Definition {
		printLog("Correct!")
	} else {
		deck.cards[index].Mistakes++
		if ok, i := deck.Contains(text, true); ok {
			template := "Wrong. The right answer is \"%s\", but your definition is correct for \"%s\"."
			printLog(template, card.Definition, deck.cards[i].Term)
		} else {
			printLog("Wrong. The right answer is \"%s\".\n", card.Definition)
		}
	}
}

func Log(reader *bufio.Reader) {
	filename := readFilename(reader)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	_, err = file.WriteString(history.String())
	if err != nil {
		log.Fatal(err)
	}
	printLog("The log has been saved.")
}

func (deck *Deck) Hardest() {
	maximum := 0
	for _, card := range deck.cards {
		if card.Mistakes > maximum {
			maximum = card.Mistakes
		}
	}
	if maximum == 0 {
		printLog("There are no cards with errors.")
		return
	}
	var hardestCards []Flashcard
	for _, card := range deck.cards {
		if card.Mistakes == maximum {
			hardestCards = append(hardestCards, card)
		}
	}
	template := "The hardest card%s %s. You have %d errors answering %s."
	plural := func() bool {
		return len(hardestCards) > 1
	}
	verbToBe := func() string {
		if plural() {
			return "s are"
		}
		return " is"
	}
	pronoun := func() string {
		if plural() {
			return "them"
		}
		return "it"
	}
	quotes := func(word string) string {
		return "\"" + word + "\""
	}
	names := func(cards []Flashcard) []string {
		out := make([]string, len(cards))
		for i, card := range cards {
			out[i] = card.Term
		}
		return out
	}
	quotedNames := func(names []string) []string {
		out := make([]string, len(names))
		for i, name := range names {
			out[i] = quotes(name)
		}
		return out
	}
	joinNames := func(quotedNames []string) string {
		return strings.Join(quotedNames, ", ")
	}
	hardestTerms := joinNames(quotedNames(names(hardestCards)))
	printLog(template, verbToBe(), hardestTerms, maximum, pronoun())
}

func (deck *Deck) Reset() {
	for i := 0; i < len(deck.cards); i++ {
		deck.cards[i].Mistakes = 0
	}
	printLog("Card statistics have been reset.")
}
