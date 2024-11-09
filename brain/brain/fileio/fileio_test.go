package fileio_test

import (
	"os"
	"testing"

	"github.com/jacoblever/heating-controller/brain/brain/fileio"
	"github.com/stretchr/testify/assert"
)

func TestReadLastLine(t *testing.T) {
	t.Run("when the file does not exist", func(t *testing.T) {
		str, err := fileio.ReadLastLine("./tmp/no-file.txt")

		assert.Equal(t, "", str)
		assert.Equal(t, true, os.IsNotExist(err))
	})

	t.Run("when the file does exist and has several lines", func(t *testing.T) {
		err := fileio.AppendLineToFile("./test-file.txt", "first line")
		assert.NoError(t, err)
		err = fileio.AppendLineToFile("./test-file.txt", "second line")
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove("./test-file.txt")
		})

		str, err := fileio.ReadLastLine("./test-file.txt")

		assert.Equal(t, "second line", str)
		assert.NoError(t, err)
	})

	t.Run("when the file does exist but is empty", func(t *testing.T) {
		err := fileio.WriteToFile("./test-file.txt", "")
		assert.NoError(t, err)

		t.Cleanup(func() {
			os.Remove("./test-file.txt")
		})

		str, err := fileio.ReadLastLine("./test-file.txt")

		assert.Equal(t, "", str)
		assert.NoError(t, err)
	})
}
