package main

import (
	"bufio"
	"fmt"
	"os"
	"os-file/object"
	"strconv"
	"strings"
	"time"
)

const (
	NotLoginError    = "[operate error]: plz login at first"
	MissingArgsError = "[operate error]: missing args"
)

func init() {
	fmt.Printf("init start\n")
	// init disk
	disk := object.NewDisk()
	disk.Disk = make([][]rune, 1024)
	disk.BitMap = make([][]int, object.BitMapRowLength)
	for i := 0; i < len(disk.BitMap); i++ {
		disk.BitMap[i] = make([]int, object.BitMapLineLength)
	}
	for i := 0; i < object.BlockNum; i++ {
		disk.Disk[i] = make([]rune, 1024)
	}
	// init user disk data
	for i := object.UserStartBlock; i < object.UserStartBlock+object.UserBlockNum; i++ {
		for j := 0; j < object.BlockSize; j++ {
			disk.Disk[i][j] = 'U'
		}
		disk.BitMap[0][i] = object.BitMapBusy
	}
	// init fcb disk data
	for i := object.FCBStartBlock; i < object.FCBStartBlock+object.FCBBlockNum; i++ {
		for j := 0; j < object.BlockSize; j++ {
			disk.Disk[i][j] = 'F'
		}
		disk.BitMap[0][i] = object.BitMapBusy
	}

	for i := object.DirStartBlock; i < object.DirStartBlock+object.DirBlockNum; i++ {
		for j := 0; j < object.BlockSize; j++ {
			disk.Disk[i][j] = 'D'
		}
		disk.BitMap[0][i] = object.BitMapBusy
	}

	for i := object.BitMapStartBlock; i < object.BitMapStartBlock+object.BitMapBlockNum; i++ {
		for j := 0; j < object.BlockSize; j++ {
			disk.Disk[i][j] = 'U'
		}
	}
	// init user
	disk.UMap["admin"] = object.NewUser("admin", "pass")
	// init fcb list
	ct := time.Now()
	root := object.NewFcb(true, "root", 0, 0, nil, ct, ct)
	disk.FcbList = append(disk.FcbList, root)
	admin := object.NewFcb(true, "admin", 0, 0, nil, ct, ct)
	disk.FcbList = append(disk.FcbList, admin)
	// init dir
	rootDir := object.NewDir(root, 0, -1)
	disk.Dirs = append(disk.Dirs, rootDir)
	adminDir := object.NewDir(admin, len(disk.Dirs), 0)
	disk.Dirs = append(disk.Dirs, adminDir)
	disk.Dirs[0].Children = append(disk.Dirs[0].Children, adminDir)
	// init DCache
	object.DCache.Disk = disk
	object.MMemory.CurDir = object.DCache.Disk.Dirs[0]
	object.MMemory.BitMap = object.DCache.Disk.BitMap
	fmt.Printf("init finish\n")
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var res *object.Result
	for {
		// todo: show()
		scanner.Scan()
		raw := scanner.Text()
		cmd, args, err := ParseCmd(raw)
		if err != nil {
			fmt.Printf("parse command error: %s\n", raw)
			continue
		}
		switch {
		case cmd == "register":
			res = object.Register(args[0], args[1])
			if !res.Success() {
				fmt.Println(res.Msg())
			}
		case cmd == "login":
			if len(args) < 2 {
				fmt.Println(MissingArgsError)
				continue
			}
			res = object.Login(args[0], args[1])
			if !res.Success() {
				fmt.Println(res.Msg())
			}
		case cmd == "mkdir":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			if len(args) < 1 {
				fmt.Println(MissingArgsError)
				continue
			}
			res = Mkdir(args[0])
			if !res.Success() {
				fmt.Println(res.Msg())
			}
		case cmd == "cd":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			if len(args) < 1 {
				fmt.Println(MissingArgsError)
				continue
			}
			res = ChangeDir(args[0])
			if !res.Success() {
				fmt.Println(res.Msg())
			}
		case cmd == "ls":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			res = ShowDir(object.MMemory.CurDir)
			if !res.Success() {
				fmt.Println(res.Msg())
			} else {
				// todo: show file list
				fcbl := res.Data().([]*object.Fcb)
				ShowFileList()
				for _, v := range fcbl {
					fmt.Println(v.String())
				}
			}
		case cmd == "touch":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			if len(args) < 1 {
				fmt.Println(MissingArgsError)
				continue
			}
			res = Create(args[0])
			if !res.Success() {
				fmt.Println(res.Msg())
			}
		case cmd == "open":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			if len(args) < 1 {
				fmt.Println(MissingArgsError)
				continue
			}
			res = Open(args[0])
			if !res.Success() {
				fmt.Println(res.Msg())
			}
		case cmd == "read":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			if len(args) < 1 {
				fmt.Println(MissingArgsError)
				continue
			}
			i, _ := strconv.ParseInt(args[0], 10, 32)
			res = Read(int(i))
			if !res.Success() {
				fmt.Println(res.Msg())
			} else {
				// will NPE
				fmt.Println(res.Data().(string))
			}
		case cmd == "write":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			sb := make([]rune, 0)
			for {
				scanner.Scan()
				line := scanner.Text()
				if strings.HasSuffix(line, "###") {
					sb = append(sb, []rune(line)...)
					break
				} else {
					sb = append([]rune(line), '\n')
				}
			}
			res = Write(string(sb))
			if res.Success() {
				fmt.Println()
			} else {
				fmt.Println(res.Msg())
			}
		case cmd == "delete":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			if len(args) < 1 {
				fmt.Println(MissingArgsError)
				continue
			}
			res = Delete(args[0])
			if !res.Success() {
				fmt.Println(res.Msg())
			}
		case cmd == "close":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			res = Close()
			if !res.Success() {
				fmt.Println(res.Msg())
			}
		case cmd == "rename":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			if len(args) < 2 {
				fmt.Println(MissingArgsError)
				continue
			}
			res = Rename(args[0], args[1])
			if !res.Success() {
				fmt.Println(res.Msg())
			}
		case cmd == "show":
			if object.MMemory.CurUser == nil {
				fmt.Println(NotLoginError)
				continue
			}
			ShowBitMap(object.MMemory.BitMap)
		case cmd == "exit":
			os.Exit(0)
		default:
			fmt.Printf("%s command not found\n", cmd)
		}
	}
}
