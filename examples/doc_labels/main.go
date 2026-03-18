// 使用文档标签的例子
// 同时使用自定义评分规则

package main

import (
	"fmt"
	"flag"
	"log"
	"strconv"
	"strings"
	"os"
	"encoding/csv"

	"github.com/huichen/wukong/engine"
	"github.com/huichen/wukong/types"
)

var (
	// searcher是线程安全的
	searcher = engine.Engine{}

	options	= types.RankOptions{
		ScoringCriteria: RankByLPQ{},
	}
	LPQWeight float32
)

type LPQFields struct {
	LPQ float32
}

// Using Qury-Independent "Landing Page Quality" for scoring
// (similar in concept to Page Rank in web search)
type RankByLPQ struct {
}

func (criteria RankByLPQ) Score(doc types.IndexedDocument, fields interface{}) []float32 {
	lpqf := fields.(LPQFields)
	return []float32{lpqf.LPQ * LPQWeight + doc.BM25}
}

func main() {
	// Define flags: name, default value, and description
	keyword := flag.String("key", "百度中国", "Search Keyword")
	docCSV := flag.String("doc", "doc.csv", "CSV file for documents")
	dictFile := flag.String("dict", "dictionary.txt", "Dictionary file")
	float64LPQWeight := flag.Float64("lpqw", 1.0, "Weight of LPQ in scoring function")

	// Parse the flags from the command line
	flag.Parse()

	// Access the values using dereferencing (e.g., *name, *age, *keyword)
	LPQWeight = float32(*float64LPQWeight)

	// Dictionary files
	files := strings.Split(*dictFile, ",")
	for i, file := range files {
		files[i] = fmt.Sprintf("../../data/%s", file)
	}

	// 初始化
	searcher.Init(types.EngineInitOptions{
		SegmenterDictionaries: strings.Join(files, ","),
		DefaultRankOptions: &options,
	})
	defer searcher.Close()

	// Read documents from a CSV file
	csvFile, err := os.Open(*docCSV)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	rows, err := csvReader.ReadAll() // `rows` is of type [][]string
	if err != nil {
		panic(err)
	}

	// 将文档加入索引，docId 从1开始
	for i, row := range rows {
		//fmt.Println(row[0], row[1])
		float64LPQ, err := strconv.ParseFloat(row[2], 32)
		if err != nil {
			float64LPQ = 0
		}
		searcher.IndexDocument(
			uint64(i+1),	// docId
			types.DocumentIndexData{
				Content: row[0],
				Labels: strings.Split(row[1], ";"),
				Fields: LPQFields{LPQ: float32(float64LPQ),},
			},
			false)
	}

	// 等待索引刷新完毕
	searcher.FlushIndex()

	// 搜索输出格式见types.SearchResponse结构体
	log.Print(searcher.Search(types.SearchRequest{Text: *keyword}))
}
