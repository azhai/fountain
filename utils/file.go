package utils

import (
	"io"
	"os"
	"os/exec"
)

func PathExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func PathNotExist(path string) bool {
	_, err := os.Stat(path)
	return os.IsNotExist(err)
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

// 通过Bash命令复制整个目录，只能运行于Linux或MacOS
// 当dst结尾带斜杠时，复制为dst下的子目录
func CopyDir(src, dst string) (err error) {
	if length := len(src); src[length-1] == '/' {
		src = src[:length-1] // 去掉结尾的斜杠
	}
	info, err := os.Stat(src)
	if err != nil || !info.IsDir() {
		return
	}
	cmd := exec.Command("cp", "-rf", src, dst)
	err = cmd.Run()
	return
}

// 删除目录及子目录下所有内容
func CleanDir(dst string) (err error) {
	cmd := exec.Command("rm", "-rf", dst)
	err = cmd.Run()
	return
}
