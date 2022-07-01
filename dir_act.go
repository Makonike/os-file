package main

import (
	"fmt"
	"os-file/object"
	"strings"
	"time"
)

// Mkdir 创建文件夹
func Mkdir(dirName string) *object.Result {
	if IsBlank(dirName) {
		return object.FailMsg("[create dir fail]: blank dirname")
	}
	dirName = strings.TrimSpace(dirName)
	for _, v := range object.MMemory.CurDir.Children {
		if v.Fcb.Name == dirName {
			return object.FailMsg("[create dir fail]: duplicate filename")
		}
	}
	curt := time.Now()
	// not assign fcb for dir
	fcb := object.NewFcb(true, dirName, 0, 0, nil, curt, curt)
	object.DCache.Disk.FcbList = append(object.DCache.Disk.FcbList, fcb)

	dir := object.NewDir(fcb, 0, object.MMemory.CurDir.IIndex())
	object.DCache.Disk.Dirs = append(object.DCache.Disk.Dirs, dir)
	dir.Index = len(object.DCache.Disk.Dirs) - 1
	object.MMemory.CurDir.Children = append(object.MMemory.CurDir.Children, dir)
	return object.SuccessResult()
}

func ChangeDir(path string) *object.Result {
	if IsBlank(path) {
		return object.SuccessResult()
	}
	res := PathResolve(path)
	if res.Success() {
		dir := res.Data().(*object.Dir)
		if dir.Fcb.IsDir() {
			object.MMemory.CurDir = res.Data().(*object.Dir)
			return object.SuccessResult()
		} else {
			return object.FailMsg("[change dir fail]: not found")
		}
	} else {
		return object.FailMsg("[change dir fail]: not found")
	}
}

func ShowDir(dir *object.Dir) *object.Result {
	if dir == nil {
		return object.FailMsg("[show dir fail]: internal server error")
	}
	if !dir.Fcb.IsDir() {
		return object.FailMsg(fmt.Sprintf("[show dir fail]: %s isn't dir", dir.Fcb.Name))
	}

	fcbs := make([]*object.Fcb, 0)
	for _, v := range dir.Children {
		child := v
		fcbs = append(fcbs, child.Fcb)
	}
	return object.SuccessData(fcbs)
}

// PathResolve solve the path and get current Dir
func PathResolve(path string) *object.Result {
	if IsBlank(path) {
		return object.FailMsg("[solve path fail]: blank path")
	}

	curDir := object.MMemory.CurDir
	pathL := strings.Split(path, object.PathSeparator)
	for i, v := range pathL {
		// back to pre
		if IsBackToPre(v) {
			if curDir.ParIndex != 0 {
				// up to pre
				curDir = object.DCache.Disk.Dirs[curDir.ParIndex]
			}
			continue
		}
		// 路径错误
		if IsBlank(pathL[i]) {
			if i != 0 && i != len(pathL)-1 {
				return object.FailMsg(fmt.Sprintf("[solve path fail]: %s not found", path))
			}
			continue
		}

		chdir := searchChildDir(curDir, pathL[i])
		if chdir == nil {
			return object.FailMsg(fmt.Sprintf("[solve path fail]: %s not found", path))
		} else {
			if i+1 < len(pathL) && !chdir.Fcb.IsDir() {
				return object.FailMsg(fmt.Sprintf("[solve path fail]: %s not found", path))
			}
			curDir = chdir
		}
	}
	return object.SuccessData(curDir)
}

func IsBackToPre(path string) bool {
	if path != "" && len(path) >= 2 {
		return strings.TrimSpace(path) == object.BackToPre
	}
	return false
}

func searchChildDir(dir *object.Dir, dirName string) *object.Dir {
	if dir == nil {
		return nil
	}
	for _, v := range dir.Children {
		if v.Fcb.Name == dirName {
			return v
		}
	}
	return nil
}
