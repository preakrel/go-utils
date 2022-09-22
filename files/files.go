package files

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// AbsolutePath 获取程序目录的绝对路径
func AbsolutePath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		log.Fatalln(err)
	}
	pwd, err := filepath.Abs(filepath.Dir(file))
	if err != nil {
		log.Fatalln(err)
	}
	pwd = strings.Replace(pwd, "\\", "/", -1)
	if !strings.HasSuffix(pwd, "/") {
		pwd = pwd + "/"
	}
	return pwd
}

// MakeDir 创建目录
func MakeDir(_dir string) bool {
	if IsExist(_dir) {
		return true
	}
	err := os.Mkdir(_dir, os.ModePerm)
	return err == nil
}

// IsExist 判断文件或目录是否存在
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

// IsDir 判断是否是目录
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// ReadFile 读取文件
func ReadFile(_file string) string {
	b, err := ioutil.ReadFile(_file)
	if err != nil {
		//Log(err)
		return ""
	}
	str := string(b)
	return str
}

// WriteFile 写文件
func WriteFile(path, data string) bool {
	if ok := MakeDir(filepath.Dir(path)); !ok {
		return false
	}
	return ioutil.WriteFile(path, []byte(data), 0644) == nil
}

// WriteFileAppend 追加的方式写文件
func WriteFileAppend(path, data string) bool {
	var err error
	fl, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return false
	}
	defer fl.Close()
	_, err = fl.Write([]byte(data))
	return err == nil
}

// ReadDir 读取目录下的文件
func ReadDir(path string) ([]string, error) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, item := range fs {
		if item.IsDir() {
			continue
		}
		files = append(files, item.Name())
	}
	return files, nil
}

// Copy 文件拷贝,将src拷贝到dst
func Copy(src, dst string) error {
	f1, err := os.OpenFile(src, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer f1.Close()
	//reader := bufio.NewReaderSize(f1, 1024*32)
	f2, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	//writer := bufio.NewWriterSize(f2, 1024*32)
	defer f2.Close()
	_, err = io.Copy(f2, f1)
	return err
}

// Zip 压缩srcFile为压缩文件
func Zip(srcFile string, destZip string) error {
	if ok := MakeDir(filepath.Dir(destZip)); !ok {
		return fmt.Errorf("create file fail")
	}
	zipFile, err := os.Create(destZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	return filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 如果是源路径，提前进行下一个遍历
		if path == srcFile {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+string(os.PathSeparator))
		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})
}

// UnZip 解压文件到指定文件夹
func UnZip(zipFile, dstDir string) (string, error) {
	file, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	var decodeName, zipDir string
	for _, f := range file.File {
		if err = func(f *zip.File) error {
			if f.Flags == 0 {
				//如果标致位是0  则是默认的本地编码   默认为gbk
				i := bytes.NewReader([]byte(f.Name))
				decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
				content, _ := ioutil.ReadAll(decoder)
				decodeName = string(content)
			} else {
				//如果标志为是 1 << 11也就是 2048  则是utf-8编码
				decodeName = f.Name
			}
			if zipDir == "" {
				zipDir = filepath.Base(decodeName)
			}
			decodeName = strings.ReplaceAll(decodeName, `\`, string(filepath.Separator))
			fpath := filepath.Join(dstDir, decodeName)
			if f.FileInfo().IsDir() {
				if err = os.MkdirAll(fpath, os.ModePerm); err != nil {
					return err
				}
			} else {
				if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
					return err
				}
				inFile, err := f.Open()
				if err != nil {
					return err
				}
				defer inFile.Close()

				outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
				if err != nil {
					return err
				}
				defer outFile.Close()

				_, err = io.Copy(outFile, inFile)
				if err != nil {
					return err
				}
			}
			return nil
		}(f); err != nil {
			return "", err
		}
	}
	return zipDir, nil
}

//ZipFiles  批量压缩文件
func ZipFiles(files []string, filename string) error {
	if ok := MakeDir(filepath.Dir(filename)); !ok {
		return fmt.Errorf("create file fail")
	}
	//创建输出文件目录
	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()
	//创建空的zip档案，可以理解为打开zip文件，准备写入
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	// Add files to zip
	for _, file := range files {
		//打开要压缩的文件
		fileToZip, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileToZip.Close()
		//获取文件的描述
		info, err := fileToZip.Stat()
		if err != nil {
			return err
		}
		//FileInfoHeader返回一个根据fi填写了部分字段的Header，可以理解成是将fileinfo转换成zip格式的文件信息
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = strings.TrimPrefix(file, filepath.Dir(file)+string(os.PathSeparator))
		/*
		   预定义压缩算法。
		   archive/zip包中预定义的有两种压缩方式。一个是仅把文件写入到zip中。不做压缩。一种是压缩文件然后写入到zip中。默认的Store模式。就是只保存不压缩的模式。
		   Store   unit16 = 0  //仅存储文件
		   Deflate unit16 = 8  //压缩文件
		*/
		header.Method = zip.Deflate
		//创建压缩包头部信息
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		//将源复制到目标，将fileToZip 写入writer   是按默认的缓冲区32k循环操作的，不会将内容一次性全写入内存中,这样就能解决大文件的问题
		_, err = io.Copy(writer, fileToZip)
		if err != nil {
			return err
		}
	}
	return nil
}

// FileSize filesize()
func FileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		return 0, err
	}
	return info.Size(), nil
}
