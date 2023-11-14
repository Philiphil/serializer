# Serializer

This project is heavily inspired by the Symfony Serializer Component. 
Its purpose is to facilitate the serialization and deserialization of golang structures in various formats while selectively excluding or including fields using serialization groups.

## Overview

The Serializer library provides a flexible way to serialize and deserialize data structures in multiple formats. It draws inspiration from the Symfony Serializer Component, offering similar functionality in the Go programming language.

## Features

- **Format Support:** Serializer currently supports various serialization formats, including JSON, CSV & XML.
- **Serialization Groups:** You can use serialization groups to selectively include or exclude fields during the serialization process.
- **Inspired by Symfony:** Leveraging concepts from the Symfony Serializer Component ensures familiarity for users familiar with Symfony.
- **Deserialization and Merge:** Deserialize and merge functionality allows for the combination of serialized data with existing objects.

## Installation

To use the Serializer library in your Go project, simply import the package:

```shell
go get github.com/philiphil/serializer
```


## Example  
```
package main

import (
	"fmt"
	"github.com/philiphil/serializer" 
)

type MyStruct struct {
	Name  string `json:"name" group:"group1"`
	Age   int    `json:"age" group:"group2"`
	Email string `json:"email"`
}


func main() {
	// Create an instance of Serializer
	mySerializer := serializer.NewSerializer(serializer.JSON)
	// Data to serialize
	dataToSerialize := MyStruct{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	// Serialize data to JSON with the "group1" group
	serializedData, err := mySerializer.Serialize(dataToSerialize, "group1")
	if err != nil {
		fmt.Println("Serialization error:", err)
		return
	}
	//Serialized Data: {
	//  "name": "John Doe"
	//}
	fmt.Println("Serialized Data:", serializedData)

	// New structure for deserialization
	var deserializedData MyStruct

	// Deserialize the data
	err = mySerializer.Deserialize(serializedData, &deserializedData)
	if err != nil {
		fmt.Println("Deserialization error:", err)
		return
	}

	// Deserialized Data: {Name:John Doe Age:0 Email:}
	fmt.Printf("Deserialized Data: %+v\n", deserializedData)

	// Create another instance of MyStruct
	anotherData := MyStruct{
		Name:  "Jane Doe",
		Age:   25,
		Email: "jane@example.com",
	}

	// Merge the deserialized data with another instance of MyStruct
	if err := mySerializer.MergeObjects(&anotherData, &deserializedData); err != nil {
		fmt.Println("Merge error:", err)
		return
	}

	// Merged Data: {Name:John Doe Age:25 Email:jane@example.com}
	fmt.Printf("Merged Data: %+v\n", anotherData)
}
```


## State

- **JSON:** JSON format is fully supported.
- **XML and CSV:** XML and CSV formats are under testing and may have bugs.
- **Serialization:** Serialization functionality is working.
- **Serialization Groups:** Serialization groups are functioning as expected.
- **Deserialization:** Deserialization functionality is operational.
- **Deserialize and Merge:** The deserialize and merge feature is functional, providing an easy way to combine serialized data with existing objects.
