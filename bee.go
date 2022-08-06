package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"unicode"
)

func main() {
	log.SetPrefix("bee: ")
	log.SetFlags(0)

	fmt.Println("Spelling Bee")
	resp, err := http.Get("https://raw.githubusercontent.com/sindresorhus/word-list/main/words.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	words, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	words = append(words, '\n')

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nWhat's the center letter? ")
	center, err := reader.ReadByte()
	if err != nil {
		log.Fatal(err)
	}
	center = byte(unicode.ToLower(rune(center)))

	reader.Reset(os.Stdin)
	fmt.Print("And the other letters? ")
	hive, err := reader.ReadBytes('\n')
	if err != nil {
		log.Fatal(err)
	}
	hive = hive[:len(hive)-1]

	if len(hive) != 6 {
		log.Fatal("there must be seven letters in the hive!")
	}

	hiveMap := map[byte]bool{
		center: true,
	}
	sort.Slice(hive, func(a, b int) bool {
		return hive[a] < hive[b]
	})
	for i, l := range hive {
		if l == center || (i > 0 && l == hive[i-1]) {
			log.Fatal("duplicate letters are not allowed!")
		}
		hiveMap[l] = false
	}

	answers := []string{}
	hasCenter := false
	n := 0
	b := 0
	for b < len(words) {
		if words[b] == '\n' {
			if hasCenter && n >= 4 {
				answers = append(answers, strings.TrimSpace(string(words[b-n:b])))
			}
			n = 0
			hasCenter = false
		} else {
			hc, ok := hiveMap[words[b]]
			if ok {
				if hc {
					hasCenter = true
				}
				n++
			} else {
				for words[b] != '\n' {
					b++
				}
				n = 0
				hasCenter = false
			}
		}
		b++
	}

	sort.Slice(answers, func(a, b int) bool {
		return len(answers[a]) < len(answers[b])
	})
	fmt.Println("\nResults:")
	for _, ans := range answers {
		fmt.Println(ans)
	}
}
