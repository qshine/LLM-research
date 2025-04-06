package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"os"
)

var authToken string

func init() {
	authToken = os.Getenv("OPENAI_API_KEY")
}

// test0 计算相似性
func test0() {
	client := openai.NewClient(authToken)
	// Create an EmbeddingRequest for the user query
	queryResponse, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{"How many chucks would a woodchuck chuck"},
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		log.Fatal("Error creating query embedding:", err)
	}

	// Create an EmbeddingRequest for the target text
	targetResponse, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{"How many chucks would a woodchuck chuck if the woodchuck could chuck wood"},
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		log.Fatal("Error creating target embedding:", err)
	}

	// Now that we have the embeddings for the user query and the target text, we
	// can calculate their similarity.
	queryEmbedding := queryResponse.Data[0]
	targetEmbedding := targetResponse.Data[0]

	similarity, err := queryEmbedding.DotProduct(&targetEmbedding)
	if err != nil {
		log.Fatal("Error calculating dot product:", err)
	}

	log.Printf("The similarity score between the query and the target is %f", similarity)
}

// test1 测试输出结果
func test1() {
	client := openai.NewClient(authToken)
	// Create an EmbeddingRequest for the user query
	queryResponse, err := client.CreateEmbeddings(context.Background(), openai.EmbeddingRequest{
		Input: []string{
			"hello",
			"hello world",
		},
		Model: openai.AdaEmbeddingV2,
	})
	if err != nil {
		log.Fatal("Error creating query embedding:", err)
	}
	fmt.Println(len(queryResponse.Data))

	for _, d := range queryResponse.Data {
		fmt.Println(d.Index)
		fmt.Println(d.Object)
	}
}

func main() {
	//test0()
	test1()
}
