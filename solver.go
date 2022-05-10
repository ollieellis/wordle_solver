package main 

import (
	"fmt"
	"strings"
	"net/http"
	"log"
	"io/ioutil"
)
//code at home to get all wordle words

type word_intel struct {
	here string //default _____ and replace with "Green" letters (maybe this should be a map)
	not_here map[int]string//default empty nested array length 5, each inner empty and yellow letters will be appened
	not_at_all string //default empty string, append letters that dont exist here
	//somewhere can be got from not heres and therefore dangerous to store 2 copies 
}

func update_intel(intel word_intel, guess string, result string) word_intel {
	for i := 0; i< len(result); i++ { //I messed up the enumerate and was getting asci not the string
		result_char := string(result[i]) //wasnt the enumerate but kept for loop becasue this is learning exercise 
		guess_char := string(guess[i])
		switch result_char{
		case "_":
			intel = update_wrong_letter(intel, guess_char)
		case "Y":
			intel = update_yellow_letter(intel, guess_char, i)
		case "G":
			intel = update_green_letter(intel, guess_char, i)
		}
	}
	return intel
}

func get_possible_words(intel word_intel, current_possible_words []string) []string { //this really should return indexs of valid words
	var now_possible_words []string
	for _, i := range current_possible_words{
		if check_word_valid(intel, i){
			now_possible_words = append(now_possible_words, i)
		}
	}
	return now_possible_words
}

func check_word_valid(intel word_intel, word_2_check string) bool {
	if !check_letters_in_position(intel, word_2_check){return false} //double negation kinda gross
	if !check_letters_in_wrong_position(intel, word_2_check){return false} //naming wierd in context
	if !check_no_incorrect_letters(intel, word_2_check){return false}
	if !check_yellow_letters_in_word(intel, word_2_check){return false}
	return true
}

func check_letters_in_position(intel word_intel, word_2_check string) bool {
	for i, intel_asci := range intel.here {
		intel_char := string(intel_asci)
		word_2_check_char := string(word_2_check[i])
		if intel_char != "_"{
			if intel_char != word_2_check_char{
				return false
			}
		}
	}
	return true
}

func check_letters_in_wrong_position(intel word_intel, word_2_check string) bool {
	for i, intel_string := range intel.not_here {
		if strings.Contains(intel_string, string(word_2_check[i])){
			return false
		}
	}
	return true
}

func check_yellow_letters_in_word(intel word_intel, word_2_check string) bool{
	for _, intel_string := range intel.not_here {
		for _, intel_char_asci := range intel_string {
			if !strings.Contains(word_2_check, string(intel_char_asci)){
				return false
			}
		}
	}
	return true
}

func check_no_incorrect_letters(intel word_intel, word_2_check string) bool {
	for _, intel_char_asci := range intel.not_at_all {
		if strings.Contains(word_2_check, string(intel_char_asci)){
			return false
		}
	}
	return true
}


func update_green_letter(intel word_intel, letter string, position int) word_intel{
	intel.here = intel.here[:position] + letter + intel.here[position + 1:]
	return intel
}

func update_wrong_letter(intel word_intel, letter string) word_intel{
	intel.not_at_all = intel.not_at_all + letter
	return intel
}	

func update_yellow_letter(intel word_intel, letter string, position int) word_intel{
		shallow_map_copy := make(map[int]string)

		// Copy from the original map to the target map
		for key, value := range intel.not_here {
			shallow_map_copy[key] = value
		}
		
		shallow_map_copy[position] = shallow_map_copy[position] + letter
		intel_copy := word_intel{here: intel.here, not_here: shallow_map_copy, not_at_all: intel.not_at_all}
		return intel_copy
	}

func get_all_wordle_words() []string {
	res, err := http.Get("https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt")
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := ioutil.ReadAll(res.Body)
	words := strings.Split(string(body), "\r\n")

	wordle_words := []string{}
	for _, word := range words {
		if len(word) == 5 {
			wordle_words = append(wordle_words, strings.ToUpper(word))
		}
	}
	return wordle_words
}

func print_guess(guess string){
	fmt.Printf("Guess %v \nhow did I DO? reply with 5 char string G for green, Y for yellow, _ for none \n", guess)	
}
	
func main(){ //lunch time hence the copy and paste 
	possible_words := get_all_wordle_words()
	intel := word_intel{here: "_____", not_here: make(map[int]string), not_at_all: ""}
	
	fmt.Println("first: OCEAN how did I DO? reply with 5 char string G for green, Y for yellow, _ for none")
	var result string
	fmt.Scanln(&result) // 0 protection against incorrect input
	intel = update_intel(intel, "OCEAN", result)
	possible_words = get_possible_words(intel, possible_words)
	
	guess := possible_words[0] //I expected a pretty niave aproach to do well but even picking first available word does well
	print_guess(guess)
	fmt.Scanln(&result)
	intel = update_intel(intel, guess, result)
	possible_words = get_possible_words(intel, possible_words)
	
	guess = possible_words[0]
	print_guess(guess)
	fmt.Scanln(&result)
	intel = update_intel(intel, guess, result)
	possible_words = get_possible_words(intel, possible_words)

	guess := possible_words[0]
	print_guess(guess)
	fmt.Scanln(&result)
	intel = update_intel(intel, guess, result)
	possible_words = get_possible_words(intel, possible_words)
	
	guess = possible_words[0]
	print_guess(guess)
	fmt.Scanln(&result)
	intel = update_intel(intel, guess, result)
	possible_words = get_possible_words(intel, possible_words)

	guess := possible_words[0]
	print_guess(guess)
	fmt.Scanln(&result)
	intel = update_intel(intel, guess, result)
	possible_words = get_possible_words(intel, possible_words)
	
	guess = possible_words[0]
	fmt.Printf("Final Guess %v \n", guess)	
}
