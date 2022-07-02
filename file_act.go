package main

import (
	"fmt"
	"os-file/object"
	"strings"
	"time"
)

// Create 创建文件
func Create(name string) *object.Result {

	if name == "" {
		return object.FailMsg("[create file fail]: blank filename")
	}
	name = strings.TrimSpace(name)
	for _, v := range object.MMemory.CurDir.Children {
		if v.Fcb.Name == name && !v.Fcb.IsDir() {
			return object.FailMsg("[create file fail]: duplicate filename")
		}
	}
	ct := time.Now()
	fcb := object.NewFcb(false, name, 0, 0, object.All(), ct, ct)

	object.DCache.Disk.FcbList = append(object.DCache.Disk.FcbList, fcb)
	dir := object.NewDir(fcb, len(object.DCache.Disk.Dirs), object.MMemory.CurDir.IIndex())
	object.MMemory.CurDir.Children = append(object.MMemory.CurDir.Children, dir)
	return object.SuccessResult()
}

// Open 打开文件
func Open(path string) *object.Result {
	if path == "" {
		return object.FailMsg("[open file fail]: blank filename")
	}
	if object.MMemory.AcFile != nil {
		return object.FailMsg("[open file fail]: exist opened file")
	}
	res := PathResolve(path)
	if res.Success() {
		dir := res.Data().(*object.Dir)
		if !dir.Fcb.IsDir() {
			ress := GetRecord(dir.Fcb.StartBlock, dir.Fcb.BlockNum)
			fr := ress.Data().([]rune)
			if ress.Success() {
				acf := object.NewActiveFile(dir.Fcb, fr, 0, len(fr))
				// put in memory
				object.MMemory.AcFile = acf
				return object.SuccessData(acf)
			} else {
				return object.FailMsg("[open file fail]: " + ress.Msg())
			}
		}
	}
	return object.FailMsg("[open file fail]: not found")
}

func Read(recordNum int) *object.Result {
	if recordNum == 0 {
		return object.SuccessResult()
	}
	acf := object.MMemory.AcFile
	if acf == nil {
		return object.FailMsg("[read file fail]: plz open file at first")
	}
	sb := make([]rune, 0)
	if recordNum > 0 {
		if acf.ReadPtr+recordNum > len(acf.Record) {
			return object.FailMsg("[read file fail]: out of range")
		}
		// read record
		for i := 0; i < recordNum; i++ {
			sb = append(sb, []rune(string(acf.Record[acf.ReadPtr]))...)
			acf.ReadPtr += 1
		}
	} else {
		// 倒序
		if acf.ReadPtr+recordNum < 0 {
			return object.FailMsg("[read file fail]: smaller than len of file")
		}
		tmp := make([]rune, 0)
		for i := 0; i > recordNum; i-- {
			// todo
			copy(tmp[1:], tmp[0:])
			tmp[0] = acf.Record[acf.ReadPtr]
			acf.ReadPtr--
		}
		sb = append(tmp, sb...)
	}
	return object.SuccessData(string(sb))
}

func Write(record string) *object.Result {
	if record == "" {
		return object.SuccessResult()
	}
	acf := object.MMemory.AcFile
	if acf == nil {
		return object.FailMsg("[read file fail]: plz open file at first")
	}
	oriR := acf.Record
	insert := []rune(record)
	b := make([]rune, len(oriR[:acf.WritePtr])+len(insert))
	copy(b, oriR[:acf.WritePtr])
	for i := 0; i < len(insert); i++ {
		b[i] = insert[i]
	}
	b = append(b, oriR[acf.WritePtr:]...)
	oriR = b
	res := StoreRecord(acf.Fcb, oriR)
	if res.Success() {
		object.MMemory.AcFile.WritePtr += len(record)
		object.MMemory.AcFile.Record = oriR
		object.MMemory.AcFile.Fcb = res.Data().(*object.Fcb)
		return object.SuccessResult()
	} else {
		return object.FailMsg(fmt.Sprintf("[write file fail]: %s", res.Msg()))
	}
}

func Close() *object.Result {
	if object.MMemory.AcFile == nil {
		return object.FailMsg("[close file fail]: not found active file")
	}
	object.MMemory.AcFile.Fcb.UpdateTime = time.Now()
	object.MMemory.AcFile = nil
	return object.SuccessResult()
}

func Delete(path string) *object.Result {
	path = strings.TrimSpace(path)
	if path == "" {
		return object.FailMsg("[delete file fail]: blank filename")
	}
	res := PathResolve(path)
	dir := res.Data().(*object.Dir)
	if res.Success() && !dir.Fcb.IsDir() {
		if object.MMemory.AcFile != nil && object.MMemory.AcFile.Fcb == dir.Fcb {
			return object.FailMsg("[delete file fail]: the file is opened")
		}
		parDir := object.DCache.Disk.Dirs[dir.ParIndex]
		for i, v := range parDir.Children {
			if !v.Fcb.IsDir() && v.Fcb.Name == dir.Fcb.Name {
				copy(parDir.Children[i:], parDir.Children[i+1:])
				break
			}
		}
		object.DCache.Disk.Dirs[dir.Index] = nil
		FreeSpace(dir.Fcb.StartBlock, dir.Fcb.BlockNum)
		return object.SuccessResult()
	}
	return object.FailMsg("[delete file fail]: not found")
}

func Rename(path string, name string) *object.Result {
	path = strings.TrimSpace(path)
	name = strings.TrimSpace(name)
	if path == "" || name == "" {
		return object.FailMsg("[rename file fail]: blank name")
	}
	res := PathResolve(path)
	dir := res.Data().(*object.Dir)
	if res.Success() {
		parDir := object.DCache.Disk.Dirs[dir.ParIndex]
		for _, v := range parDir.Children {
			// todo: 文件夹与文件夹名分开
			if v.Fcb.Name == name {
				return object.FailMsg("[rename file fail]: existed name")
			}
		}
		dir.Fcb.Name = name
		dir.Fcb.UpdateTime = time.Now()
		return object.SuccessResult()
	}
	return object.FailMsg("[rename file fail]: not found")
}
