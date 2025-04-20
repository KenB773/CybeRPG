// Game logic engine
package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/agnivade/levenshtein"
)

const levelThreshold = 30
const saveFilePath = "examples/save.json"

const bossFightIntro = "\n👹 BOSS FIGHT! Answer this to complete the chapter!"

type Player struct {
	Name string
	XP   map[string]int
}

type Question struct {
	Category string
	Prompt   string
	Answer   string
}

// normalize input for lenient matching
func normalize(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, "&", "and")
	s = strings.TrimSpace(s)
	s = strings.Join(strings.Fields(s), " ")
	return s
}

func answersMatch(input, correct string) bool {
	input = normalize(input)
	correct = normalize(correct)
	if input == correct {
		return true
	}
	return levenshtein.ComputeDistance(input, correct) <= 2
}

var bossQuestions = map[int]Question{
	3: {Category: "Boss", Prompt: "Final question of Chapter 3: What does XSS stand for?", Answer: "Cross Site Scripting"},
	5: {Category: "Boss", Prompt: "Chapter 5 Boss: Which framework classifies attacker techniques and tactics?", Answer: "MITRE ATT&CK"},
	7: {Category: "Boss", Prompt: "Final Showdown: What is the first step in the NIST incident response lifecycle?", Answer: "Preparation"},
}

// ... keep your long Questions array here (omitted for brevity) ...

func StartGame() {
	if len(os.Args) > 1 && os.Args[1] == "--achievements" {
		player := loadPlayerProgress()
		fmt.Printf("\n📜 Achievements for %s:\n", player.Name)
		printAchievements(player)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	player := loadPlayerProgress()
	if player.XP == nil {
		player.XP = make(map[string]int)
	}

	fmt.Println("Welcome to CybeRPG!")
	if player.Name == "Analyst" {
		fmt.Print("Enter your name: ")
		name, _ := reader.ReadString('\n')
		player.Name = strings.TrimSpace(name)
	}

	categories := getUniqueCategories(Questions)
	fmt.Println("\nChoose a category to focus on or press Enter for all:")
	for _, c := range categories {
		fmt.Printf("- %s\n", c)
	}
	fmt.Print("\n> ")
	catChoice, _ := reader.ReadString('\n')
	catChoice = strings.TrimSpace(catChoice)

	var selected []Question
	if catChoice == "" {
		selected = Questions
	} else {
		for _, q := range Questions {
			if strings.EqualFold(q.Category, catChoice) {
				selected = append(selected, q)
			}
		}
		if len(selected) == 0 {
			fmt.Println("No matching category found. Using all questions.")
			selected = Questions
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(selected), func(i, j int) { selected[i], selected[j] = selected[j], selected[i] })

	for i, q := range selected {
		fmt.Printf("[%d/%d] %s\n%s\n> ", i+1, len(selected), q.Category, q.Prompt)
		input, _ := reader.ReadString('\n')
		input = normalize(input)
		correct := normalize(q.Answer)

		if answersMatch(input, correct) {
			if input != correct {
				fmt.Println("✅ Close enough! We'll count that as correct.")
			} else {
				fmt.Println("✅ Correct! +10 XP")
			}
			player.XP[q.Category] += 10
			if player.XP[q.Category]%levelThreshold == 0 {
				fmt.Printf("🏅 You leveled up in %s!\n", q.Category)
			}
		} else {
			fmt.Printf("❌ Incorrect. Correct answer: %s\n", q.Answer)
		}
		fmt.Println("--------------------------------------")

		chapterCheckpoint := (i + 1) / 10
		if boss, exists := bossQuestions[chapterCheckpoint]; exists && (i+1)%10 == 0 {
			fmt.Println(bossFightIntro)
			fmt.Printf("%s\n> ", boss.Prompt)
			bossInput, _ := reader.ReadString('\n')
			bossInput = normalize(bossInput)
			bossAnswer := normalize(boss.Answer)

			if answersMatch(bossInput, bossAnswer) {
				if bossInput != bossAnswer {
					fmt.Println("🎉 Boss defeated with a near-match! +20 XP")
				} else {
					fmt.Println("🎉 Boss defeated! +20 XP")
				}
				player.XP["Boss"] += 20
			} else {
				fmt.Printf("💀 You failed the boss. Correct answer: %s\n", boss.Answer)
			}
			fmt.Println("--------------------------------------")
		}
	}

	savePlayerProgress(player)

	fmt.Println("\nGame Over. Here's your XP breakdown:")
	for category, xp := range player.XP {
		fmt.Printf("%s: %d XP\n", category, xp)
	}
	printAchievements(player)
}

func savePlayerProgress(player Player) {
	file, err := os.Create(saveFilePath)
	if err != nil {
		fmt.Println("⚠️  Could not save game progress.")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(map[string]interface{}{"player": player, "achievements": detectAchievements(player)})
}

func loadPlayerProgress() Player {
	file, err := os.Open(saveFilePath)
	if err != nil {
		return Player{Name: "Analyst"}
	}
	defer file.Close()

	var data map[string]Player
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return Player{Name: "Analyst"}
	}

	return data["player"]
}

func detectAchievements(player Player) []string {
	totalXP := 0
	perfect := true
	achievements := []string{}

	for _, xp := range player.XP {
		totalXP += xp
		if xp < 10 {
			perfect = false
		}
	}

	if xp, ok := player.XP["Boss"]; ok && xp >= 60 {
		achievements = append(achievements, "👑 Boss Slayer: Defeated all bosses")
	}
	if perfect {
		achievements = append(achievements, "🧠 Sharp Mind: No incorrect answers")
	}
	if len(player.XP) >= 8 {
		achievements = append(achievements, "🌐 Master of Domains: XP gained in all major categories")
	}
	if totalXP >= 500 {
		achievements = append(achievements, "💼 Veteran Analyst: 500+ XP earned")
	}

	return achievements
}

func getUniqueCategories(questions []Question) []string {
	seen := make(map[string]bool)
	var list []string
	for _, q := range questions {
		if !seen[q.Category] {
			seen[q.Category] = true
			list = append(list, q.Category)
		}
	}
	return list
}

func printAchievements(player Player) {
	fmt.Println("\n🏆 Achievements:")
	for _, a := range detectAchievements(player) {
		fmt.Printf("- %s\n", a)
	}
}
