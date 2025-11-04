package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type PolicySecurity struct {
	SignatureAlg   string `json:"signature_alg"`
	Signature      string `json:"signature"`
	EnforceHash    bool   `json:"enforce_hash"`
	OnHashMismatch string `json:"on_hash_mismatch"`
}

type PolicyStub struct {
	Security PolicySecurity `json:"security"`
}

func main() {
	// Default to data/policy.v1.json, or allow override via command line
	policyPath := "data/policy.v1.json"
	if len(os.Args) > 1 {
		policyPath = os.Args[1]
	}

	data, err := os.ReadFile(policyPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error reading policy: %v\n", err)
		os.Exit(1)
	}

	// Parse to get signature
	var stub PolicyStub
	if err := json.Unmarshal(data, &stub); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error parsing policy: %v\n", err)
		os.Exit(1)
	}

	// Parse full policy to remove signature field
	var rawPolicy map[string]interface{}
	if err := json.Unmarshal(data, &rawPolicy); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error parsing policy structure: %v\n", err)
		os.Exit(1)
	}

	// Remove signature from security section for hash calculation
	if sec, ok := rawPolicy["security"].(map[string]interface{}); ok {
		delete(sec, "signature")
	} else {
		fmt.Fprintf(os.Stderr, "❌ Policy missing security section\n")
		os.Exit(1)
	}

	// Calculate canonical JSON (sorted keys, consistent formatting)
	canonical, err := json.Marshal(rawPolicy)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error creating canonical JSON: %v\n", err)
		os.Exit(1)
	}

	// Calculate SHA256 hash
	hash := sha256.Sum256(canonical)
	calculatedHash := hex.EncodeToString(hash[:])

	// Check if signature is still placeholder
	if stub.Security.Signature == "REPLACE_WITH_SHA256_OR_SIGNATURE" || stub.Security.Signature == "" {
		fmt.Println("ℹ️  Policy signature is placeholder")
		fmt.Printf("📝 Calculated hash: %s\n", calculatedHash)
		fmt.Println("\n💡 To update policy with this signature:")
		fmt.Println("   1. Copy the hash above")
		fmt.Println("   2. Replace the 'signature' field in data/policy.v1.json")
		fmt.Println("   3. Run this script again to verify")
		os.Exit(0)
	}

	// Verify signature matches
	if calculatedHash != stub.Security.Signature {
		fmt.Fprintf(os.Stderr, "❌ Policy signature mismatch!\n")
		fmt.Fprintf(os.Stderr, "Expected: %s\n", stub.Security.Signature)
		fmt.Fprintf(os.Stderr, "Got:      %s\n", calculatedHash)
		fmt.Fprintf(os.Stderr, "\n⚠️  Policy file may have been modified or corrupted\n")
		os.Exit(1)
	}

	fmt.Println("✅ Policy signature valid")
	fmt.Printf("   Algorithm: %s\n", stub.Security.SignatureAlg)
	fmt.Printf("   Hash: %s\n", calculatedHash)
	fmt.Printf("   Enforcement: %v\n", stub.Security.EnforceHash)
}
