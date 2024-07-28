package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"little-auth/crypto/hotp"
	"little-auth/models"
	"little-auth/utils"
	"little-auth/vault"
	"time"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show an OTP.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			chosenIndexer *models.Indexer
			mapStrIndexer map[string]*models.Indexer
			secretVault   vault.Vault
			otp           uint32
			err           error
		)
		if secretVault, err = vault.New(); err != nil {
			fmt.Printf("Failed to initialize vault: %v\n", err)
			panic(err)
		}

		indexes, err := secretVault.GetAllIndexes()
		if err != nil {
			fmt.Printf("Failed to get list of authenticate info: %v\n", err)
			panic(err)
		}

		if len(indexes) == 0 {
			fmt.Println("No entries found.")
			return
		}

		var indexStrs []string
		mapStrIndexer = make(map[string]*models.Indexer)
		for _, index := range indexes {
			name := ""
			if index.Issuer != "" && index.Path != "" {
				name = fmt.Sprintf("Path: %s - Issuer: %s", index.Path, index.Issuer)
			} else if index.Issuer == "" {
				name = fmt.Sprintf("Path: %s", index.Path)
			} else if index.Path == "" {
				name = fmt.Sprintf("Issuer: %s", index.Issuer)
			}
			if name != "" {
				indexStrs = append(indexStrs, name)
				mapStrIndexer[name] = index
			}
		}
		chosenIndexerStr := utils.PromptGetSelect(utils.PromptContent{
			ErrorMsg: "Please choose one to get OTP.",
			Label:    "Select an entry.",
		}, indexStrs)
		if chosenIndexer = mapStrIndexer[chosenIndexerStr]; chosenIndexer == nil {
			fmt.Println("Failed to get chosen indexer.")
			panic("failed to get chosen indexer")
		}

		authSecret, err := secretVault.Get(chosenIndexer)
		if err != nil {
			fmt.Printf("Failed to get authenticate info: %v\n", err)
			panic(err)
		}

		if otp, err = hotp.New(models.HashFnFromType(authSecret.HashType)).Generate(
			authSecret.Secret,
			utils.TotpCounterFromNow(authSecret.CountFactor),
			6,
		); err != nil {
			panic(err)
		}

		fmt.Printf("%s OTP: %d\n", time.Now().Format("2006/01/02 15:04:05"), otp)
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
