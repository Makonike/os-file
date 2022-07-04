package main

import (
	"fmt"
	"os-file/object"
	"strings"
)

func ShowBitMap(bitMap [][]int) {
	fmt.Println("0-free, 1-busy")
	fmt.Println("all block number:", object.BlockNum, "block size:", object.BlockSize, "total size:", object.DiskSize, "B")
	fmt.Printf("   ")
	for i := 0; i < object.BitMapLineLength; i++ {
		fmt.Printf("%-3d", i)
	}
	fmt.Println("  -")
	fmt.Printf("")
	for i := 0; i < object.BitMapLineLength; i++ {
		fmt.Printf("---")
	}
	fmt.Println()

	for i := 0; i < object.BitMapRowLength; i++ {
		fmt.Printf("%2d", i)
		fmt.Printf("|[")
		j := 0
		for j = 0; j < object.BitMapLineLength-1; j++ {
			fmt.Printf("%d, ", bitMap[i][j])
		}
		fmt.Printf("%d", bitMap[i][j])
		fmt.Printf("]")
	}

}

func ShowFileList() {
	fmt.Printf("%20s %10s %20s %20s\n", "filename", "size", "createTime", "updateTime")
}

func ParseCmd(raw string) (string, []string, error) {
	var command string
	args := make([]string, 0)

	resp := strings.Split(raw, " ")

	switch {
	case len(resp) == 2:
		command = strings.TrimSpace(resp[0])
		args = append(args, strings.TrimSpace(resp[1]))
	case len(resp) == 1:
		command = resp[0]
	case len(resp) > 2:
		command = strings.TrimSpace(resp[0])
		for i := 1; i < len(resp); i++ {
			args = append(args, strings.TrimSpace(resp[i]))
		}
	}
	return command, args, nil
}
