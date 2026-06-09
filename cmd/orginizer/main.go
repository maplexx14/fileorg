package main

import (
	"flag"
	//"fmt"
	//"os"

	"file-orginizer/internal/fileworker"
)

func main() {
	path := flag.String("path", " ", "Путь к папке") // считывание аргумента (по умолчанию рабочая директория)
	//dateSort := flag.Bool("d", false, "Рассортировать по датам")
	flag.Parse()

	fileworker.FileSorter(path)

}
