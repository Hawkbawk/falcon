package files

import (
	"os"

	"github.com/hawkbawk/falcon/lib/logger"
)

// OpenFile opens the file at the specified location. If it encounters any errors, it logs them
// and ends the program.
func OpenFile(path string) *os.File {
	file, err := os.Open(path)

	if err != nil {
		logger.LogError("Unable to open file %v.\n See the following error for more details: %v",
			path, err.Error())
	}

	return file
}

// CreateFile creates a file at the specified path. If it encounters any errors, it logs them
// and ends the program.
func CreateFile(path string) *os.File {
	file, err := os.Create(path)

	if err != nil {
		logger.LogError("Unable to create file %v. \n Error: ", path, err.Error())
	}

	return file
}

// DeleteFile deletes the file at the specified path. If it encounters any errors, it logs them
// and ends the program.
func DeleteFile(path string) {
	if err := os.Remove(path); err != nil {
		logger.LogError("Unable to delete file %v. Error: %v", path, err.Error())
	}
}

func Symlink(oldname string, newname string) {
	if err := os.Symlink(oldname, newname); err != nil {
		logger.LogError("Unable to create symlink to file %v at %v. Error: ",
			oldname, newname, err.Error())
	}
}

// ReadFile reads all of the contents of the specified file into a byte slice and returns that
// slice. If any errors are encountered along the way, it logs the error and ends the program.
func ReadFile(file *os.File) []byte {
	buffer := make([]byte, FileSize(file))

	_, err := file.Read(buffer)

	if err != nil {
		logger.LogError("Unable to read file: ", file.Name())
	}

	return buffer

}

// FileStats returns the info on the given file. If any errors are encountered, it logs the error
// and then ends the program.
func FileInfo(path string) os.FileInfo {
	stats, err := os.Stat(path)

	if err != nil {
		logger.LogError("Unable to get stats on file %v.\n See the following error for details: ",
			path, err.Error())
	}

	return stats
}

// FileSize returns the size of the specified file. If it encounters any errors, it logs them and
// ends the program.
func FileSize(file *os.File) int64 {
	fileinfo, err := file.Stat()

	if err != nil {
		logger.LogError("Unable to get stats on file: ", file.Name())
	}

	return fileinfo.Size()
}

// OverwriteFile overwrites the previous contents of the passed in file with the specified contents.
// If it encounters any errors, it logs them and ends the program.
func OverwriteFile(file *os.File, contents []byte) {
	err := file.Truncate(0)

	if err != nil {
		logger.LogError("Unable to clear out file %v for write. See the following error for more details: ",
			file.Name(), err.Error())
	}

	_, err = file.Seek(0, 0)

	if err != nil {
		logger.LogError("Unable to return to beginning of file %v for write. See the following error for more details: ",
			file.Name(), err.Error())
	}

	_, err = file.Write(contents)

	if err != nil {
		logger.LogError("Unable to write all data to the file %v. See the following error for details: ",
			file.Name(), err.Error())
	}

	err = file.Sync()
	if err != nil {
		logger.LogError("Unable to save data to disk for file %v. See the following error for details: ",
			file.Name(), err.Error())
	}
}
