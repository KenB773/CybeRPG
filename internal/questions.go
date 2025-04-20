// Question bank storage
package internal

// Placeholder for question bank
var Questions = []Question{
	{Category: "OSINT", Prompt: "What tool can you use to find a domain's WHOIS information?", Answer: "whois"},
	{Category: "Incident Response", Prompt: "What port does SSH use by default?", Answer: "22"},
	// 48 more questions to be added
}

type Question struct {
	Category string
	Prompt   string
	Answer   string
}