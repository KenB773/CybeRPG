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

var Questions = []Question{
	{Category: "OSINT", Prompt: "What platform allows searching for leaked credentials?", Answer: "Have I Been Pwned"},
	{Category: "OSINT", Prompt: "What website visualizes DNS data and passive DNS results?", Answer: "SecurityTrails"},
	{Category: "OSINT", Prompt: "Which site is often used for reverse image searching?", Answer: "TinEye"},
	{Category: "OSINT", Prompt: "What command line tool can be used to get DNS records?", Answer: "dig"},
	{Category: "Incident Response", Prompt: "What’s the first step in an incident response plan?", Answer: "Preparation"},
	{Category: "Incident Response", Prompt: "What command is used to find open ports on Linux?", Answer: "netstat"},
	{Category: "Incident Response", Prompt: "What file contains system logins in Linux?", Answer: "/var/log/auth.log"},
	{Category: "Incident Response", Prompt: "What port does SSH use by default?", Answer: "22"},
	{Category: "Threat Hunting", Prompt: "What tool is used to collect Windows event logs?", Answer: "Event Viewer"},
	{Category: "Threat Hunting", Prompt: "What is the Windows command for listing running tasks?", Answer: "tasklist"},
	{Category: "Threat Hunting", Prompt: "What tool detects persistence mechanisms in memory?", Answer: "Volatility"},
	{Category: "Threat Hunting", Prompt: "What is the process for establishing attacker objectives?", Answer: "MITRE ATT&CK"},
	{Category: "Threat Hunting", Prompt: "What does IOC stand for?", Answer: "Indicator of Compromise"},
	{Category: "Threat Hunting", Prompt: "Which log type is useful to detect brute force attempts?", Answer: "Authentication logs"},
	{Category: "Threat Hunting", Prompt: "What is Sysmon used for?", Answer: "Windows event logging"},
	{Category: "Threat Hunting", Prompt: "What tool can capture and filter live traffic?", Answer: "Wireshark"},
	{Category: "Networking", Prompt: "What does the OSI model layer 4 represent?", Answer: "Transport"},
	{Category: "Networking", Prompt: "Which command can test reachability via ICMP?", Answer: "ping"},
	{Category: "Networking", Prompt: "What port is used for HTTPS?", Answer: "443"},
	{Category: "Networking", Prompt: "Which tool performs network discovery scans?", Answer: "nmap"},
	{Category: "Phishing", Prompt: "What type of phishing targets a specific person?", Answer: "Spear phishing"},
	{Category: "Phishing", Prompt: "What protocol secures email delivery?", Answer: "SPF"},
	{Category: "Phishing", Prompt: "What tool is used to simulate phishing for training?", Answer: "GoPhish"},
	{Category: "Phishing", Prompt: "What does a link shortener often hide?", Answer: "Malicious destination"},
	{Category: "Phishing", Prompt: "What does DKIM help verify?", Answer: "Email authenticity"},
	{Category: "Phishing", Prompt: "Which social engineering tactic involves impersonating authority?", Answer: "Pretexting"},
	{Category: "Phishing", Prompt: "What protocol prevents spoofing by validating sender IPs?", Answer: "SPF"},
	{Category: "Phishing", Prompt: "What kind of phishing occurs over SMS?", Answer: "Smishing"},
	{Category: "Phishing", Prompt: "What security mechanism scans embedded links in emails?", Answer: "Safe Links"},
	{Category: "Cryptography", Prompt: "What does SSL stand for?", Answer: "Secure Sockets Layer"},
	{Category: "Cryptography", Prompt: "What is the modern replacement for SSL?", Answer: "TLS"},
	{Category: "Cryptography", Prompt: "What algorithm does bcrypt implement?", Answer: "Password hashing"},
	{Category: "Cryptography", Prompt: "What does asymmetric encryption require?", Answer: "Public and private key"},
	{Category: "Cryptography", Prompt: "What is a digital signature used for?", Answer: "Integrity and authenticity"},
	{Category: "Cryptography", Prompt: "What is the purpose of a hash function? Hint: Ensure ____ _________", Answer: "Ensure data integrity"},
	{Category: "Cryptography", Prompt: "Which hash algorithm is no longer considered secure?", Answer: "MD5"},
	{Category: "Cryptography", Prompt: "What is a rainbow table used for?", Answer: "Crack password hashes"},
	{Category: "Cryptography", Prompt: "What is the key length of AES-256?", Answer: "256 bits"},
	{Category: "Cryptography", Prompt: "What is elliptic curve cryptography known for?", Answer: "Efficiency"},
	{Category: "Linux", Prompt: "What command lists active network connections?", Answer: "netstat"},
	{Category: "Linux", Prompt: "How do you switch to root in terminal?", Answer: "sudo -i"},
	{Category: "Linux", Prompt: "What command shows current directory?", Answer: "pwd"},
	{Category: "Linux", Prompt: "What does chmod do?", Answer: "Changes file permissions"},
	{Category: "Linux", Prompt: "What log file shows kernel events (hint: full file path)?", Answer: "/var/log/kern.log"},
	{Category: "Digital Forensics", Prompt: "What file system is commonly used in Windows systems?", Answer: "NTFS"},
	{Category: "Digital Forensics", Prompt: "Which tool is widely used for memory forensics?", Answer: "Volatility"},
	{Category: "Digital Forensics", Prompt: "What does 'timeline analysis' help investigators determine?", Answer: "Event sequences"},
	{Category: "Digital Forensics", Prompt: "Which command in Linux shows all running processes?", Answer: "ps aux"},
	{Category: "Digital Forensics", Prompt: "What is the purpose of a disk image in forensics?", Answer: "Preserve evidence"},
	{Category: "Digital Forensics", Prompt: "What format is often used for forensic disk images?", Answer: "E01"},
	{Category: "Digital Forensics", Prompt: "Which Windows artifact tracks program execution?", Answer: "Prefetch"},
	{Category: "Digital Forensics", Prompt: "What tool extracts metadata from files?", Answer: "exiftool"},
	{Category: "Digital Forensics", Prompt: "What is a volatile artifact in the context of digital forensics?", Answer: "Short-lived evidence"},
	{Category: "Digital Forensics", Prompt: "What system file logs USB device connections in Windows?", Answer: "SYSTEM registry hive"},
	{Category: "Reverse Engineering", Prompt: "What tool is commonly used for analyzing malware statically?", Answer: "Ghidra"},
	{Category: "Reverse Engineering", Prompt: "Which tool can be used to unpack PE files?", Answer: "UPX"},
	{Category: "Reverse Engineering", Prompt: "What is IDA Pro used for?", Answer: "Static binary analysis"},
	{Category: "Reverse Engineering", Prompt: "What section of a PE file contains executable code?", Answer: ".text"},
	{Category: "Cloud Security", Prompt: "What AWS service handles IAM roles and users?", Answer: "IAM"},
	{Category: "Cloud Security", Prompt: "What AWS service is used to monitor logs?", Answer: "CloudWatch"},
	{Category: "Cloud Security", Prompt: "What service provides VPC flow logs?", Answer: "VPC"},
	{Category: "Cloud Security", Prompt: "What tool detects misconfigured cloud resources?", Answer: "ScoutSuite"},
	{Category: "Cloud Security", Prompt: "Which service scans AWS infrastructure for security risks?", Answer: "Inspector"},
	{Category: "Cloud Security", Prompt: "What GCP service provides threat detection?", Answer: "Security Command Center"},
	{Category: "Cloud Security", Prompt: "What Azure tool visualizes security posture?", Answer: "Defender for Cloud"},
	{Category: "Cloud Security", Prompt: "Which AWS service logs API calls?", Answer: "CloudTrail"},
	{Category: "Cloud Security", Prompt: "What AWS feature uses JSON-based policy documents?", Answer: "IAM policy"},
	{Category: "Blue Team", Prompt: "What is the goal of defense in depth?", Answer: "Layered security"},
	{Category: "Blue Team", Prompt: "What system centralizes logs?", Answer: "SIEM"},
	{Category: "Blue Team", Prompt: "What does EDR stand for?", Answer: "Endpoint Detection and Response"},
	{Category: "Blue Team", Prompt: "What tool correlates logs and alerts for analysis?", Answer: "SIEM"},
	{Category: "Blue Team", Prompt: "Which EDR product is made by CrowdStrike?", Answer: "Falcon"},
	{Category: "Blue Team", Prompt: "What is a common log format used in SIEMs?", Answer: "CEF"},
	{Category: "Blue Team", Prompt: "What detection type identifies previously unknown threats?", Answer: "Behavior-based"},
	{Category: "General", Prompt: "What does CVE stand for?", Answer: "Common Vulnerabilities and Exposures"},
	{Category: "General", Prompt: "What does SOC stand for?", Answer: "Security Operations Center"},
	{Category: "General", Prompt: "What does the CIA triad stand for?", Answer: "Confidentiality, Integrity, Availability"},
}

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
