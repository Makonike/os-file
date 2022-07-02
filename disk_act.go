package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os-file/object"
	path2 "path"
	"strings"
)

// todo: deprecated
func SaveDisk(path string) *object.Result {
	var err error
	_, filename := path2.Split(path)
	fmt.Println(filename)
	err = os.MkdirAll(path, 0777)
	if err != nil {
		return object.FailMsg("create dir error")
	}
	return nil
}

// todo: deprecated
func LoadDisk(path string) *object.Result {
	f, err := os.Open(path)
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("close file error %s\n", err)
		}
	}(f)
	if err != nil {
		return object.FailMsg("[load disk fail]: not found disk file")
	}

	reader := bufio.NewReader(f)
	res := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Printf("read error %s", err)
			break
		}
		// not exist or not data in buf
		if n == 0 {
			break
		}
		// write into file
		res = append(res, buf[:n]...)
	}
	// json转换为二进制存储
	ress := string(res)
	var disk object.Disk
	dec := json.NewDecoder(strings.NewReader(ress))
	err = dec.Decode(&disk)
	if err != io.EOF {
		return object.FailMsg("[load disk fail]: open disk file fail")
	}
	// todo: test
	fmt.Println(disk)
	object.DCache.Disk = &disk
	object.MMemory.CurDir = object.DCache.Disk.Dirs[0]
	object.MMemory.BitMap = object.DCache.Disk.BitMap
	return object.SuccessResult()
}

// GetRecord 读取盘块中的信息至内存中
func GetRecord(sBlockId, blockNum int) *object.Result {
	record := make([]rune, 1024)
	if sBlockId == 0 {
		return object.SuccessData(record)
	}

	if sBlockId < object.RecordStartBlock || sBlockId >= object.BlockNum || blockNum <= 0 || blockNum > object.BlockNum {
		return object.FailMsg("[get record fail]: internal server error")
	}

	disk := object.DCache.Disk
	bP := 0
	for i := 0; i < blockNum; {
		if i == blockNum-1 {
			if bP >= len(disk.Disk[sBlockId+i]) || disk.Disk[sBlockId+i][bP] == 0 {
				break
			}
		}

		record = append(record, disk.Disk[sBlockId+i][bP])
		if bP >= object.BlockSize {
			bP = 0
			i++
		} else {
			bP++
		}
	}
	return object.SuccessData(record)
}

func StoreRecord(fcb *object.Fcb, record []rune) *object.Result {
	if len(record) == 0 {
		return object.SuccessResult()
	}

	setBitmapStatus(fcb.StartBlock, fcb.BlockNum, true)
	reqNum := DivWithUp(len(record), object.BlockSize)
	cnt, sBlockId := 0, object.RecordStartBlock

	for i := 0; i < object.BitMapRowLength; i++ {
		for j := 0; j < object.BitMapLineLength; j++ {
			if object.BitMapFree == object.MMemory.BitMap[i][j] {
				if cnt == 0 {
					sBlockId = i*object.BitMapLineLength + j
				}
				cnt++
				// 有足够的块可以分配
				if cnt == reqNum {
					storeDisk(sBlockId, record)
					setBitmapStatus(sBlockId, reqNum, false)
					fcb.StartBlock = sBlockId
					fcb.BlockNum = reqNum
					return object.SuccessData(fcb)
				}
			} else {
				// 重新计数
				cnt = 0
			}
		}
	}
	return object.FailMsg("[assign block fail]: out of disk space")
}

// 持久化至磁盘中
func storeDisk(sBlockId int, record []rune) {
	index, blockId := 0, sBlockId
	disk := object.DCache.Disk
	// 覆盖原先的信息
	for _, v := range record {
		k := v
		// may be using many blocks
		if index >= object.BlockSize {
			blockId++
			index = 0
		}
		copy(disk.Disk[blockId][index+1:], disk.Disk[blockId][index:])
		disk.Disk[blockId][index] = k
		index++
	}

	for index < object.BlockSize && len(disk.Disk[blockId]) > index {
		disk.Disk[blockId][index] = 0
		index++
	}
}

func FreeSpace(sBlockId, blockNum int) *object.Result {
	setBitmapStatus(sBlockId, blockNum, true)
	return object.SuccessResult()
}

// SetBitmapStatus 更改特定位置位示图的状态，用于连续分配
// status true-free false-busy
func setBitmapStatus(sBlockId, blockNum int, status bool) {
	// out of range
	if sBlockId == 0 || sBlockId < object.RecordStartBlock || sBlockId >= blockNum || blockNum <= 0 {
		return
	}
	// i = b / n
	row := sBlockId / object.BitMapLineLength
	// j = b % n
	line := sBlockId % object.BitMapLineLength
	for i := 0; i < blockNum; i++ {
		if status {
			object.MMemory.BitMap[row][line] = object.BitMapFree
		} else {
			object.MMemory.BitMap[row][line] = object.BitMapBusy
		}
		// 换行
		if line >= object.BitMapLineLength-1 {
			line = 0
			row++
		} else {
			line++
		}

	}

}
