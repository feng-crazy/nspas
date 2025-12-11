#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Test script for neuroscience knowledge retriever
Allows interactive querying of the knowledge base from terminal
"""
import os

# 设置TOKENIZERS_PARALLELISM环境变量以避免huggingface/tokenizers的警告
os.environ["TOKENIZERS_PARALLELISM"] = "false"

from retrievers.neuroscience_retriever import NeuroscienceRetriever

def main():
    print("Neuroscience Knowledge Retriever Test")
    print("=====================================")
    print("Initializing retriever...")
    
    # Initialize the retriever
    retriever = NeuroscienceRetriever()
    
    print("Retriever initialized. Enter queries to test retrieval.")
    print("Type 'quit' or 'exit' to stop.\n")
    
    while True:
        try:
            # Get query from user
            query = input("Enter your query: ").strip()
            
            # Check for exit commands
            if query.lower() in ['quit', 'exit', 'q']:
                print("Goodbye!")
                break
            
            # Skip empty queries
            if not query:
                print("Please enter a valid query.\n")
                continue
            
            print(f"\nSearching for: '{query}'")
            print("-" * 40)
            
            # Retrieve relevant documents
            results = retriever.invoke(query)
            
            # Display results
            if results:
                for i, doc in enumerate(results, 1):
                    score = doc.metadata.get('score', doc.metadata.get('relevance_score', 'N/A'))
                    print(f"{i}. {doc.page_content}")
                    print(f"   Score: {score}")
                    print()
            else:
                print("No relevant documents found.\n")
            
            print("=" * 50 + "\n")
            
        except KeyboardInterrupt:
            print("\n\nInterrupted by user. Goodbye!")
            break
        except Exception as e:
            print(f"An error occurred: {e}\n")
            print("=" * 50 + "\n")

if __name__ == "__main__":
    main()