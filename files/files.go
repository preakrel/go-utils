package files

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"gopkg.in/iconv.v1"
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
	b, err := os.ReadFile(_file)
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
	return os.WriteFile(path, []byte(data), 0644) == nil
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
	fs, err := os.ReadDir(path)
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
func Zip(srcFile string, destZip any) error {
	var archive *zip.Writer
	switch dz := destZip.(type) {
	case *bytes.Buffer:
		archive = zip.NewWriter(dz)
		defer archive.Close()
	default:
		zipFile, err := os.Create(fmt.Sprint(dz))
		if err != nil {
			return err
		}
		defer zipFile.Close()
		archive = zip.NewWriter(zipFile)
		defer archive.Close()
	}
	return filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 忽略目录本身
		if path == srcFile {
			return nil
		}
		relPath, err := filepath.Rel(srcFile, path)
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(relPath)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		header.SetMode(0755) //修改权限
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

type charset string

const (
	Utf8    = charset("UTF-8")
	Gb18030 = charset("GB18030")
)

// 将byte转换为指定编码字符串
func ConvertByte(b []byte, charset charset) []byte {
	var byt = make([]byte, 0)
	switch charset {
	case Gb18030:
		byt, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(b)
	case Utf8:
		fallthrough
	default:
		byt = b
	}
	return byt
}

// 转换编码为
func Convert(src, dst string, b []byte) ([]byte, error) {
	if dst == "" {
		dst = "utf-8"
	}
	if src == "" {
		src = "utf-8"
	}
	if dst == src {
		return b, nil
	}
	cd, err := iconv.Open(dst, src) // convert src to dst
	if err != nil {
		return nil, err
	}
	defer cd.Close()
	r := iconv.NewReader(cd, ioutil.NopCloser(bytes.NewReader(b)), 0)
	var buffer = make([]byte, 0)
	for {
		buf := make([]byte, 1024)
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			io.Copy(ioutil.Discard, r) //释放掉未获取的数据
			return b, err
		}
		if n == 0 {
			break
		}
		buffer = append(buffer, buf[:n]...)
	}
	return buffer, nil
}

// 解压文件到指定文件夹
func UnZip(zipFile, dstDir string) error {
	file, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer file.Close()
	var decodeName string
	for _, f := range file.File {
		if err = func(f *zip.File) error {
			bname := []byte(f.Name)
			if !utf8.Valid(bname) {
				bd := ConvertByte(bname, "Gb18030")
				name, err := Convert("utf-8", "gbk", bd)
				if err != nil {
					return err
				}
				decodeName = string(name)
			} else {
				decodeName = f.Name
			}
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
			return err
		}
	}
	return nil
}

// UploadUnZip 上传压缩包解压文件到指定文件夹
func UploadUnZip(zipFile *multipart.FileHeader, dstDir string) (err error) {
	uf, err := zipFile.Open()
	if err != nil {
		return
	}
	defer uf.Close()

	file, err := zip.NewReader(uf, zipFile.Size)
	if err != nil {
		return
	}

	var decodeName string
	for _, f := range file.File {
		if err = func(f *zip.File) error {
			bName := []byte(f.Name)
			if f.Flags == 0 || !utf8.Valid(bName) {
				//如果标致位是0  则是默认的本地编码   默认为gbk
				bd := ConvertByte(bName, "Gb18030")
				name, err := Convert("utf-8", "gbk", bd)
				if err != nil {
					return err
				}
				decodeName = string(name)
			} else {
				//如果标志为是 1 << 11也就是 2048  则是utf-8编码
				decodeName = f.Name
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
			return err
		}
	}
	return nil
}

// 解压tar压缩文件到dst目录
func UnTar(tarFile, dstDir string) error {
	srcFile, err := os.Open(tarFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	tr := tar.NewReader(gr)
	var decodeName string
	for {
		hdr, err := tr.Next()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}
		if err := func(f *tar.Header) error {
			bname := []byte(f.Name)
			if !utf8.Valid(bname) {
				bd := ConvertByte(bname, "Gb18030")
				name, err := Convert("utf-8", "gbk", bd)
				if err != nil {
					return err
				}
				decodeName = string(name)
			} else {
				decodeName = f.Name
			}
			fpath := filepath.Join(dstDir, decodeName)
			if f.FileInfo().IsDir() {
				if err = os.MkdirAll(fpath, os.ModePerm); err != nil {
					return err
				}
			} else {
				if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
					return err
				}

				outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
				if err != nil {
					return err
				}
				defer outFile.Close()

				_, err = io.Copy(outFile, tr)
				if err != nil {
					return err
				}
			}
			return nil
		}(hdr); err != nil {
			return err
		}

	}
	return nil
}

// UploadUnTar 上传解压解压 tar.gz
func UploadUnTar(src *multipart.FileHeader, dst string) (err error) {
	// 打开准备解压的 tar 包
	fr, err := src.Open()
	if err != nil {
		return
	}
	defer fr.Close()

	// 将打开的文件先解压
	gr, err := gzip.NewReader(fr)
	if err != nil {
		return
	}
	defer gr.Close()

	// 通过 gr 创建 tar.Reader
	tr := tar.NewReader(gr)

	// 现在已经获得了 tar.Reader 结构了，只需要循环里面的数据写入文件就可以了
	for {
		hdr, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case hdr == nil:
			continue
		}
		if err := func(hdr *tar.Header) error {
			decodeName := hdr.Name
			bName := []byte(decodeName)
			if !utf8.Valid(bName) {
				bd := ConvertByte(bName, "Gb18030")
				name, err := Convert("utf-8", "gbk", bd)
				if err != nil {
					return err
				}
				decodeName = string(name)
			}
			// 处理下保存路径，将要保存的目录加上 header 中的 Name
			// 这个变量保存的有可能是目录，有可能是文件，所以就叫 FileDir 了……
			dstFileDir := filepath.Join(dst, decodeName)

			// 根据 header 的 Typeflag 字段，判断文件的类型
			switch hdr.Typeflag {
			case tar.TypeDir: // 如果是目录时候，创建目录
				// 判断下目录是否存在，不存在就创建
				if b := IsExist(dstFileDir); !b {
					// 使用 MkdirAll 不使用 Mkdir ，就类似 Linux 终端下的 mkdir -p，
					// 可以递归创建每一级目录
					if err := os.MkdirAll(dstFileDir, os.ModePerm); err != nil {
						return err
					}
				}
			case tar.TypeReg: // 如果是文件就写入到磁盘
				// 创建一个可以读写的文件，权限就使用 header 中记录的权限
				// 因为操作系统的 FileMode 是 int32 类型的，hdr 中的是 int64，所以转换下
				file, err := os.OpenFile(dstFileDir, os.O_CREATE|os.O_RDWR|os.O_TRUNC, os.FileMode(hdr.Mode))
				if err != nil {
					return err
				}
				// 不要忘记关闭打开的文件，因为它是在 for 循环中，不能使用 defer
				// 如果想使用 defer 就放在一个单独的函数中
				defer file.Close()
				// n, err := io.Copy(file, tr)
				if _, err := io.Copy(file, tr); err != nil {
					return err
				}
				// 将解压结果输出显示
				// fmt.Printf("成功解压： %s , 共处理了 %d 个字符\n", dstFileDir, n)
			}
			return nil
		}(hdr); err != nil {
			return err
		}
	}
}

// ZipFiles  批量压缩文件
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

func CopyFileIfExist(src, dst string) error {
	if !IsExist(src) {
		return nil
	}
	bt, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(bt)
	return err
}

// CopyFile 文件拷贝,将src拷贝到dst
func CopyFile(src, dst string) error {
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

// CopyDir copies a whole directory recursively, with force overwrite
func CopyDir(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	files, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, file := range files {
		srcFilePath := path.Join(src, file.Name())
		dstFilePath := path.Join(dst, file.Name())

		if file.IsDir() {
			if err = CopyDir(srcFilePath, dstFilePath); err != nil {
				return err
			}
		} else {
			if err = CopyFile(srcFilePath, dstFilePath); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadAllDir(pathname string) ([]string, error) {
	rd, err := os.ReadDir(pathname)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, fi := range rd {
		fullName := filepath.Join(pathname, fi.Name())
		if fi.IsDir() {
			downFiles, err := ReadAllDir(fullName)
			if err != nil {
				return files, err
			}
			files = append(files, downFiles...)
		} else {
			files = append(files, fullName)
		}
	}
	return files, nil
}

// ReadFirstDir 获取第一层文件夹路径
func ReadFirstDir(dirPath string) ([]string, error) {
	var subDirs []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != dirPath {
			subDirs = append(subDirs, path)
		}
		return nil
	})
	return subDirs, err
}

// 判断文件内是否为空
func FileIsNil(path string) (int64, error) {
	file, err := os.Stat(path)
	if err != nil {
		return 0, errors.New("所选择的文件不存在")
	}
	size := file.Size()
	if size == 0 {
		return size, errors.New("文件内容为空")
	}
	return size, nil
}

// ParseFileName windows下文件名不能包含\ / : * ? " < > |
func ParseFileName(dir string) string {
	sign := []string{`\`, `/`, `:`, `*`, `?`, `"`, `<`, `>`, `|`}
	for _, s := range sign {
		dir = strings.ReplaceAll(dir, s, "_")
	}
	for _, s := range []string{"\r", "\n", "\t", " "} {
		dir = strings.ReplaceAll(dir, s, "")
	}
	return dir
}
