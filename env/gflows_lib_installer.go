package env

import (
	"fmt"
	"strings"

	"github.com/jbrunton/gflows/io"
	"github.com/jbrunton/gflows/io/content"
	"github.com/jbrunton/gflows/io/pkg"
	"github.com/spf13/afero"
)

type GFlowsLibInstaller struct {
	fs     *afero.Afero
	reader *content.Reader
	writer *content.Writer
	logger *io.Logger
}

func NewGFlowsLibInstaller(fs *afero.Afero, reader *content.Reader, writer *content.Writer, logger *io.Logger) *GFlowsLibInstaller {
	return &GFlowsLibInstaller{
		fs:     fs,
		reader: reader,
		writer: writer,
		logger: logger,
	}
}

func (installer *GFlowsLibInstaller) install(lib *GFlowsLib) ([]*pkg.PathInfo, error) {
	manifest, err := installer.loadManifest(lib.ManifestPath)
	if err != nil {
		return nil, err
	}

	rootPath, err := pkg.ParentPath(lib.ManifestPath)
	if err != nil {
		return nil, err
	}

	files := []*pkg.PathInfo{}
	for _, relPath := range manifest.Libs {
		localPath, err := installer.copyFile(lib, rootPath, relPath)
		if err != nil {
			return nil, err
		}

		pathInfo, err := lib.GetPathInfo(localPath)
		if err != nil {
			return nil, err
		}

		files = append(files, pathInfo)
	}
	return files, nil
}

func (installer *GFlowsLibInstaller) loadManifest(manifestPath string) (*GFlowsLibManifest, error) {
	manifestContent, err := installer.reader.ReadContent(manifestPath)
	if err != nil {
		return nil, err
	}
	return ParseManifest(manifestContent)
}

func (installer *GFlowsLibInstaller) copyFile(lib *GFlowsLib, rootPath string, relPath string) (string, error) {
	if !strings.HasPrefix(relPath, "libs/") && !strings.HasPrefix(relPath, "workflows/") {
		return "", fmt.Errorf("Unexpected directory %s, file must be in libs/ or workflows/", relPath)
	}
	sourcePath, err := pkg.JoinRelativePath(rootPath, relPath)
	if err != nil {
		return "", err
	}
	if pkg.IsRemotePath(rootPath) {
		installer.logger.Debugf("Downloading %s\n", sourcePath)
	} else {
		installer.logger.Debugf("Copying %s\n", sourcePath)
	}
	localPath, err := pkg.JoinRelativePath(lib.LocalDir, relPath)
	if err != nil {
		return "", err
	}
	sourceContent, err := installer.reader.ReadContent(sourcePath)
	if err != nil {
		return "", err
	}
	err = installer.writer.SafelyWriteFile(localPath, sourceContent)
	return localPath, err
}
