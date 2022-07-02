package object

import "time"

type User struct {
	username string
	password string
}

func NewUser(username, password string) *User {
	return &User{
		username: username,
		password: password,
	}
}

func Login(username, password string) *Result {
	if MMemory.CurUser != nil {
		return FailMsg("[user login fail]: plz logout now")
	}

	u, ok := DCache.Disk.UMap[username]
	if !ok || u == nil {
		return FailMsg("[user login fail]: username not found")
	}
	if u.password != password {
		return FailMsg("[user login fail]: wrong password")
	}

	for _, v := range DCache.Disk.Dirs[0].Children {
		if v.Fcb.IsDir() && v.Fcb.Name == username {
			MMemory.CurUser = u
			return SuccessResult()
		}
	}
	// todo: lazy load to fix error
	return FailMsg("[user login fail]: internal server error, not found the user dir")
}

func Register(username, password string) *Result {
	if MMemory.CurUser == nil {
		return FailMsg("[register user fail]: plz logout now")
	}
	if username == "" {
		return FailMsg("[register user fail]: blank username")
	}
	if password == "" {
		return FailMsg("[register user fail]: blank password")
	}
	u, ok := DCache.Disk.UMap[username]
	if ok && u != nil {
		return FailMsg("[register user fail]: existed user")
	}
	DCache.Disk.UMap[username] = NewUser(username, password)
	ct := time.Now()
	fcb := NewFcb(true, username, 0, 0, nil, ct, ct)
	DCache.Disk.FcbList = append(DCache.Disk.FcbList, fcb)

	dir := NewDir(fcb, len(DCache.Disk.Dirs), 0)
	DCache.Disk.Dirs = append(DCache.Disk.Dirs, dir)
	DCache.Disk.Dirs[0].Children = append(DCache.Disk.Dirs[0].Children, dir)
	return SuccessMsg("[register user success]: register success")
}

func Logout() *Result {
	if MMemory.CurUser == nil {
		return FailMsg("[user logout fail]: no user logged now")
	}
	MMemory.CurUser = nil
	// change curDir
	MMemory.CurDir = DCache.Disk.Dirs[0]
	return SuccessResult()
}
