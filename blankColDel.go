package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

func main() {
	flag.Parse()

	//ログファイル準備
	/*
		logfile, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		failOnError(err)
		defer logfile.Close()

		log.SetOutput(logfile)
	*/

	log.Print("Start\r\n")

	// ファイルを読み込んで二次元配列に入れる
	records := readfile(flag.Arg(0))

	// 削除する列（カラム）を調査
	recFlags := surveyfile(records)

	// ファイルへ書き出す
	savefile(flag.Arg(0), records, recFlags)

	log.Print("Finesh !\r\n")

}

func readfile(filename string) [][]string {
	//入力ファイル準備
	infile, err := os.Open(filename)
	failOnError(err)
	defer infile.Close()

	reader := csv.NewReader(transform.NewReader(infile, japanese.ShiftJIS.NewDecoder()))
	reader.Comma = '\t'

	//CSVファイルを２次元配列に展開
	readrecords := make([][]string, 0)
	for {
		record, err := reader.Read() // 1行読み出す
		if err == io.EOF {
			break
		} else {
			failOnError(err)
		}

		readrecords = append(readrecords, record)
	}

	return readrecords
}

func savefile(filename string, saverecords [][]string, saveFlags []bool) {
	//出力ファイル準備
	outDir, outfileName := filepath.Split(filename)
	pos := strings.LastIndex(outfileName, ".")
	// outfile, err := os.Create(outDir + outfileName[:pos] + "d.txt")
	outfile, err := os.Create(outDir + outfileName[:pos] + ".txt")
	failOnError(err)
	defer outfile.Close()

	writer := csv.NewWriter(transform.NewWriter(outfile, japanese.ShiftJIS.NewEncoder()))
	writer.Comma = '\t'
	writer.UseCRLF = true

	for _, out_record := range saverecords {
		var out_rec []string
		for k, columnFlag := range saveFlags {
			if columnFlag {
				out_rec = append(out_rec, out_record[k])
			}
		}
		writer.Write(out_rec)
	}

	writer.Flush()

}

func surveyfile(surveyrecords [][]string) []bool {
	recLen := len(surveyrecords[0])
	surveyFlags := make([]bool, recLen)

	recordsLen := len(surveyrecords)
	for i := 1; i < recordsLen; i++ { //タイトル行は調査対象ではないので1から
		for j, flag := range surveyFlags {
			if !flag {
				if surveyrecords[i][j] != "" {
					surveyFlags[j] = true
				}
			}
		}
	}

	return surveyFlags
}
