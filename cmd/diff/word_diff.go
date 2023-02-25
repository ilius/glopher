package main

import "fmt"

func xmlWordSplit(s string) []int {
	idx := []int{}
	// last := -1
	tag := false
	// tagQuote := ""
	skip := false
	alphabet := false
	space := false
	n := len(s)
	for i := 0; i < n; i++ {
		c := s[i]
		if tag {
			//if tagQuote != "" {
			//}
			switch c {
			case '>':
				idx = append(idx, i+1)
				tag = false
				skip = true
				alphabet = false
				space = false
			}
			continue
		}
		switch c {
		case '<':
			idx = append(idx, i)
			tag = true
		case ' ', '\t', '\n':
			if alphabet {
				idx = append(idx, i)
				alphabet = false
			} else if !space {
				idx = append(idx, i)
			}
			space = true
		case '>', ',', '.', '/', '"', '\'', '\\', '!', '~', '#', '$', '%', '^',
			'&', '*', '(', ')', '_', '+', '=', ':':
			idx = append(idx, i)
			alphabet = false
			space = false
		default:
			if !alphabet && i > 0 {
				idx = append(idx, i)
			}
			alphabet = true
			space = false
		}
		if skip {
			if idx[len(idx)-1] == i {
				idx = idx[:len(idx)-1]
			}
			skip = false
		}
	}
	return idx
}

func extractWordMapFull(str string, idx []int) map[string][]int {
	wordMap := map[string][]int{}
	n := len(idx)
	for i := 0; i < n; i++ {
		switch str[idx[i]] {
		case ' ', '\t', '\n':
			continue
		}
		word := str[idx[i]:idx[i+1]]
		wordMap[word] = append(wordMap[word], i)
	}
	return wordMap
}

func xmlFormattedWordDiff(a_str string, b_str string) string {
	res := ""
	a_idx := xmlWordSplit(a_str)
	b_idx := xmlWordSplit(b_str)
	a_word_n := len(a_idx) - 1
	b_word_n := len(b_idx) - 1

	a_map := extractWordMapFull(a_str, a_idx)
	b_map := extractWordMapFull(b_str, b_idx)

	forward_a := func(a_word string, a_word_i int, b_word string, b_word_i int) bool {
		switch b_word[0] {
		case ' ', '\t', '\n':
			return false
		}
		switch a_word[0] {
		case ' ', '\t', '\n':
			return true
		}
		a_in_b := b_map[a_word]
		b_in_a := a_map[b_word]
		if len(a_in_b) == 0 {
			return true
		}
		if len(b_in_a) == 0 {
			return false
		}

		b_delta := a_in_b[0] - b_word_i
		if b_delta >= 0 && b_delta <= (b_word_n-a_word_n)*10/9 {
			return false
		}

		a_delta := b_in_a[0] - a_word_i
		if a_delta >= 0 && a_delta <= (a_word_n-b_word_n)*10/9 {
			return true
		}

		if len(a_in_b) < len(b_in_a)/2 {
			return true
		}
		//if len(b_in_a) < len(a_in_b)/2 {
		//	return false
		//}

		return false
	}

	a_word_i := 0
	b_word_i := 0
	for a_word_i < a_word_n && b_word_i < b_word_n {
		a_word := a_str[a_idx[a_word_i]:a_idx[a_word_i+1]]
		b_word := b_str[b_idx[b_word_i]:b_idx[b_word_i+1]]
		if a_word == b_word {
			res += a_word
			a_word_i++
			b_word_i++
			continue
		}
		if forward_a(a_word, a_word_i, b_word, b_word_i) {
			res += fmt.Sprintf("%s%s%s", red, a_word, reset)
			a_word_i++
			word_pos := a_map[a_word]
			if len(word_pos) > 0 {
				a_map[a_word] = word_pos[1:]
			}
			continue
		}
		res += fmt.Sprintf("%s%s%s", green, b_word, reset)
		b_word_i++
		word_pos := b_map[b_word]
		if len(word_pos) > 0 {
			b_map[b_word] = word_pos[1:]
		}
	}
	return res
}
