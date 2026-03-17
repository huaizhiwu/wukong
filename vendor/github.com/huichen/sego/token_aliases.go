package sego

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

type TokenAliases struct {
	tokenAliases map[string]string
}

// 从aliasFile中读入同义词，一组同义词一行
// 文档索引建立时会把每个同义词映射到its canonical alias.
func (ta *TokenAliases) Init(aliasFile string) {
	ta.tokenAliases = make(map[string]string)
	if aliasFile == "" {
		return
	}

	csvFile, err := os.Open(aliasFile)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	rows, err := csvReader.ReadAll() // `rows` is of type [][]string
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		canonicalAlias := row[0]
		aliases := strings.Split(row[1], ";")
		for _, alias := range aliases {
			if value, found := ta.tokenAliases[alias]; !found {
				// New alias
				ta.tokenAliases[alias] = canonicalAlias
			} else {
				fmt.Printf("%s is already mapped to: %s, ignore the new mapping to: %s\n", alias, value, canonicalAlias)
			}
		}
		// Add the canonical alias to the map too
		// (so all the names involved can be found in the keys of the map)
		ta.tokenAliases[canonicalAlias] = canonicalAlias
	}
}

func (ta *TokenAliases) TokenAlias(token string) (string, bool) {
	alias, found := ta.tokenAliases[token]
	return alias, found
}

// 释放资源
func (ta *TokenAliases) Close() {
	ta.tokenAliases = nil
}
