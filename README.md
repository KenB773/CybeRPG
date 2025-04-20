# 🛡️ CybeRPG: The Cybersecurity Training RPG

Welcome to **CybeRPG**, a terminal-based training RPG where you level up by answering cybersecurity questions across domains like OSINT, Threat Hunting, Incident Response, Cloud Security, Reverse Engineering, and more.

Designed for aspiring analysts and defenders, this interactive game tests your skills, tracks your XP by domain, and throws in boss fights every 10 questions to keep things spicy. Great for fun practice or as a conversation piece in your cybersecurity portfolio.

---

## Features

- 📚 100+ curated cybersecurity questions
- 🧠 Typo-tolerant answer checking with Levenshtein matching
- 🏆 XP system + level-up notifications per category
- 👹 Boss fights every 10 questions
- 📅 Auto-save progress to disk
- 🌟 Focused category play (or randomized mode)
- 📜 Achievement tracking
- 💻 Fully terminal-based — no external dependencies

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/KenB773/cybersec-rpg.git
cd cybersec-rpg
```

### 2. Run the game

```bash
go run ./cmd
```

> 💡 You can also run with the `--achievements` flag to see your unlocked trophies:
>
> ```bash
> go run ./cmd --achievements
> ```

---

## Folder Structure

```bash
cybersec-rpg/
├── assets/             # Optional: ASCII banner or future extensions
│   └── logo.txt        
├── cmd/                # Entrypoint: main.go
├── internal/           # Game engine, questions, player tracking
│   └── engine.go
├── examples/           # Save game data (JSON)
│   └── save.json       
├── go.mod / go.sum     # Module definitions
└── README.md           # You're here!
```

---

## ✍Built With

- [Go](https://golang.org/) — Fast and portable backend
- [Levenshtein](https://github.com/agnivade/levenshtein) — Forgiving answer matching
- 🧠 Your brain — For leveling up in real life

---

## Categories Covered

- OSINT  
- Threat Hunting  
- Digital Forensics  
- Incident Response  
- Networking  
- Cloud Security (AWS, Azure, GCP)  
- Reverse Engineering  
- Linux  
- Cryptography  
- Phishing & Social Engineering  
- Blue Team & SIEM  
- General Cybersecurity Knowledge  

---

## Why This Project?

This game was created as both a learning tool and a portfolio project — something fun, functional, and practical. Whether you’re a SOC analyst, red team hopeful, or just studying for your next cert, CybeRPG helps keep your reflexes sharp in a gamified way.

---

## 🛠TODOs & Future Features

- [ ] Multiplayer support (just kidding… maybe)
- [ ] Difficulty scaling (easy/medium/hard)
- [ ] Timed mode / streak-based scoring
- [ ] More achievements
- [ ] Web-based version?

---

## License

MIT — do whatever you'd like, just don't launch ransomware with it 🧃

---

## Author

Made with ✨love✨, caffeine, and ASCII hallucinations by [Ken Brigham](https://github.com/KenB773)


