package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"little-auth/config"
	"little-auth/models"
	"little-auth/utils"
	"little-auth/vault"
	"strconv"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Add a new authentication entry.",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			path          string
			issuer        string
			secret        []byte
			algoType      models.AlgoType
			hashType      models.HashType
			counterFactor uint64
			secretVault   vault.Vault
			err           error
		)
		path = utils.PromptGetInput(utils.PromptContent{
			ErrorMsg: "Please provide a path:",
			Label:    "Enter the path:",
		})
		issuer = utils.PromptGetInput(utils.PromptContent{
			ErrorMsg: "Please provide an issuer:",
			Label:    "Enter the issuer:",
			ValidateFunc: func(input string) error {
				return nil
			},
		})

		secretStr := utils.PromptGetInput(utils.PromptContent{
			ErrorMsg: "Please provide a secret:",
			Label:    "Enter the secret:",
		})
		if secret, err = utils.StringBase32ToBytes(utils.NormalizeSecret(secretStr)); err != nil {
			fmt.Printf("Failed to convert secret to bytes: %v\n", err)
			panic(err)
		}

		var supportedTypes []string
		for _, t := range config.SUPPORTED_ALGO_TYPES {
			supportedTypes = append(supportedTypes, string(t))
		}
		algoType = models.AlgoType(utils.PromptGetSelect(utils.PromptContent{
			ErrorMsg: "Please provide a type:",
			Label:    "Select the type:",
		}, supportedTypes))

		hashType = models.HashType(utils.PromptGetSelect(utils.PromptContent{
			ErrorMsg: "Please provide a hash type:",
			Label:    "Select the hash type:",
		}, models.GetHashTypesStr()))

		counterFactorStr := utils.PromptGetInput(utils.PromptContent{
			ErrorMsg: "Please provide a counter factor. If default, enter 30 for TOTP, or 0 for HOTP:",
			Label:    "Enter the counter factor. For TOTP: period in seconds (default usually 30). For HOTP: start counter (default usually 0):",
			ValidateFunc: func(input string) error {
				if _, err := strconv.ParseUint(input, 10, 64); err != nil {
					return errors.New("not a valid number")
				}
				return nil
			},
		})
		counterFactor, _ = strconv.ParseUint(counterFactorStr, 10, 64)

		if secretVault, err = vault.New(); err != nil {
			fmt.Printf("Failed to initialize vault: %v\n", err)
			panic(err)
		}
		if err = secretVault.Set(
			&models.Indexer{
				Issuer: issuer,
				Path:   path,
			},
			&models.Secret{
				Secret:      secret,
				AlgoType:    algoType,
				CountFactor: counterFactor,
				HashType:    hashType,
			},
		); err != nil {
			fmt.Printf("Failed to save authenticate info. Please try again. Err: %v\n", err)
			panic(err)
		}

		fmt.Printf("Successfully saved authenticate info:\n"+
			"\tPath: %s\n"+
			"\tIssuer: %s\n"+
			"\tType: %s\n"+
			"\tHash Type: %s\n"+
			"\tCounter Factor: %d\n",
			path, issuer, algoType, hashType, counterFactor)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}
