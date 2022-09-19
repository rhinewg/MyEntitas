package generator

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const header string = `//////////////////////////////////////////////////////////////////////////
//
// Copyright (c) 2021 Vladislav Fedotov (Falldot)
// License: MIT License
// MIT License web page: https://opensource.org/licenses/MIT
//
//////////////////////////////////////////////////////////////////////////
//
// This file generated by Entitas-Go generator. PLEASE DO NOT EDIT IT.
//
// Entitas-Go: github.com/Falldot/Entitas-Go
//
//////////////////////////////////////////////////////////////////////////
package ecs
`

var files []string = []string{
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/entityPool.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/component.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/componentBitSet.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/componentPool.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/events.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/contexts.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/entity.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/entityBase.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/matcher.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/group.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/collector.go",
	"https://raw.githubusercontent.com/rhinewg/Entitas-Go/master/ecs/system.go",
}

func CreateEntitasLibFile() {
	os.Mkdir("./Entitas", 0777)

	file, _ := os.Create("./Entitas/Entitas-gen.go")
	defer file.Close()

	var EntitasFile []byte

	sliceByte := []byte(header)
	EntitasFile = append(EntitasFile, sliceByte...)

	for _, v := range files {
		resp, err := http.Get(v)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		fContent, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		EntitasFile = append(EntitasFile, fContent[147:]...)
	}

	file.WriteString(string(EntitasFile))
}

func CreateEntitasContextFile(context string, components []*Component, src []byte) {

	contextFile, _ := os.Create("./Entitas/" + context)
	defer contextFile.Close()

	newSRC := string(src)
	newSRC = strings.Replace(newSRC, "package", "//", -1)
	for _, v := range components {
		newSRC = strings.Replace(newSRC, v.Name, v.Name+"Component", -1)
	}

	contextData := header + componentConstansTemplate + newSRC

	for i, v := range components {
		if v.Ident {
			contextData += componentTemplate
			contextData += componentTemplateGetMethodSingleType
		} else {
			contextData += componentTemplate
			contextData += componentTemplateGetMethodStruct
		}

		// {name}
		contextData = strings.Replace(contextData, "{name}", v.Name, -1)

		// {const}
		contextData = strings.Replace(contextData, "{context}", context[:len(context)-3], -1)
		contextData = strings.Replace(contextData, "{componentCount}", fmt.Sprint(len(components)), -1)
		contextData = strings.Replace(contextData, "//go:generate go run github.com/rhinewg/Entitas-Go", "", -1)
		if i == 0 {
			contextData = strings.Replace(contextData, "{const}", v.Name, -1)
		} else {
			next := "\n" + v.Name + " //next"
			contextData = strings.Replace(contextData, " //next", next, -1)
		}

		var result, argsWithType, args []string
		for n, f := range v.Fields {
			var str, str2 string

			if v.Ident {
				contextData = strings.Replace(contextData, "{type}", f, -1)
				str = "c" + " = (*" + n + "Component" + ")(&" + strings.ToLower(n) + ")"
			} else {
				str = "c." + n + " = " + strings.ToLower(n) + "\n"
			}
			result = append(result, str)

			str2 = n + " " + f
			argsWithType = append(argsWithType, strings.ToLower(str2))

			args = append(args, strings.ToLower(n))
		}

		results := strings.Join(result, "")
		contextData = strings.Replace(contextData, "{result}", results, -1)

		argsWithTypes := strings.Join(argsWithType, ",")
		contextData = strings.Replace(contextData, "{argsWithType}", argsWithTypes, -1)

		arg := strings.Join(args, ",")
		contextData = strings.Replace(contextData, "{args}", arg, -1)
	}

	contextFile.WriteString(contextData)
}
