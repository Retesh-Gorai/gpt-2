package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiKey := viper.GetString("API_KEY")

	if apiKey == "" {
		log.Fatal("Missing API key")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	const inputFile = "./input_file_mainframe_code.txt"

	fileBytes, err := os.ReadFile(inputFile)

	if err != nil {
		log.Fatal("failed to read file %v", err)
	}

	msgPrefix := "Give me a short list of libraries that are used in the code \n```python\n"
	msgSuffix := "\n```"

	msg := msgPrefix + string(fileBytes) + msgSuffix

	outputBuilder := strings.Builder{}

	err = client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			msg,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(res *gpt3.CompletionResponse) {
		outputBuilder.WriteString(res.Choices[0].Text)
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}

	output := strings.TrimSpace(outputBuilder.String())
	const outputFile = "./output_file.txt"
	err = os.WriteFile(outputFile, []byte(output), os.ModePerm)

	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

}
