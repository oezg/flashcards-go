# Flashcards

## About
When learning a new language, it can be hard to remember all the new vocabulary, which is exactly where flashcards can help. Typically, flashcards show a hint (a task or a picture) on one side and the right answer on the other. Flashcards can be used to remember any sort of data, so this project is a useful tool to help your learning.

## Learning Outcomes
In this project,
- learn how to work with files and 
- call them from the command line.

## Stages
1. Display information about a single card on the screen.
2. Compare the lines and work with conditions: display the card and the user’s answer on the screen.
3. Practice arrays and loops: create a new card for the program to play with you.
4. Learn to use hash tables, display key values, and work with exceptions in order to fix the problem of repeating cards.
5. Work with files: create a menu that allows you to add, delete, save, and upload saved cards in your game.
6. Using statistics, set a correct answer for each card and teach the game to determine which card was the hardest to solve.
7. Enable the user to import files right upon starting the game, working with command-line arguments.

## Usage Examples
./flashcards --import_from=words13june.txt --export_to=words14june.txt

## Supported Actions
- `add`: add a card
- `remove`: remove a card
- `ask`: ask for definitions of some random cards
- `hardest card`: print the term or terms that the user makes most mistakes with
- `import`: load cards from file
- `export`: save cards to file
- `reset stats`: erase the mistake count for all cards
- `log`: save the application log to the given file
- `exit`: exit the program
  