package main

import (
	"os-file/object"
)

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
