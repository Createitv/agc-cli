// Package hvigor executes local HarmonyOS builds without shell interpolation.
package hvigor

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Request struct {
	Project    string `json:"project"`
	Executable string `json:"hvigorPath,omitempty"`
	Product    string `json:"product,omitempty"`
	BuildMode  string `json:"buildMode,omitempty"`
	Task       string `json:"task,omitempty"`
}

type Artifact struct {
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	SHA256 string `json:"sha256"`
}

type Result struct {
	Project   string     `json:"project"`
	Artifacts []Artifact `json:"artifacts"`
}

// Build 执行固定 assembleHap 任务；项目路径作为工作目录而非 shell 命令。
// 构建输出可能包含签名配置，不能原样写入 API 响应或错误信息。
func Build(ctx context.Context, request Request) (Result, error) {
	project, err := filepath.Abs(request.Project)
	if err != nil || request.Project == "" {
		return Result{}, fmt.Errorf("project directory required")
	}
	if info, err := os.Stat(filepath.Join(project, "build-profile.json5")); err != nil || !info.Mode().IsRegular() {
		return Result{}, fmt.Errorf("project must contain build-profile.json5")
	}
	executable := request.Executable
	if executable == "" {
		executable = "/Applications/DevEco-Studio.app/Contents/tools/hvigor/bin/hvigorw"
		if _, err := os.Stat(executable); err != nil {
			executable = "hvigorw"
		}
	}
	product := request.Product
	if product == "" {
		product = "default"
	}
	mode := request.BuildMode
	if mode == "" {
		mode = "debug"
	}
	if mode != "debug" && mode != "release" {
		return Result{}, fmt.Errorf("buildMode must be debug or release")
	}
	task, buildScope, extension := request.Task, "module", ".hap"
	if task == "" {
		task = "assembleHap"
	}
	if task == "assembleApp" {
		buildScope, extension = "project", ".app"
	} else if task != "assembleHap" {
		return Result{}, fmt.Errorf("task must be assembleHap or assembleApp")
	}
	command := exec.CommandContext(ctx, executable, "--mode", buildScope, "-p", "product="+product, "-p", "buildMode="+mode, task, "--no-daemon")
	command.Dir = project
	command.Env = os.Environ()
	if os.Getenv("JAVA_HOME") == "" {
		if _, err := os.Stat("/Applications/DevEco-Studio.app/Contents/jbr/Contents/Home"); err == nil {
			command.Env = append(command.Env, "JAVA_HOME=/Applications/DevEco-Studio.app/Contents/jbr/Contents/Home")
		}
	}
	if os.Getenv("DEVECO_SDK_HOME") == "" {
		if _, err := os.Stat("/Applications/DevEco-Studio.app/Contents/sdk"); err == nil {
			command.Env = append(command.Env, "DEVECO_SDK_HOME=/Applications/DevEco-Studio.app/Contents/sdk")
		}
	}
	if err := command.Run(); err != nil {
		if ctx.Err() != nil {
			return Result{}, ctx.Err()
		}
		return Result{}, fmt.Errorf("Hvigor build failed: %T (build output withheld)", err)
	}
	result := Result{Project: project, Artifacts: []Artifact{}}
	err = filepath.WalkDir(project, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && (entry.Name() == "oh_modules" || entry.Name() == ".git" || entry.Name() == ".hvigor") {
			return filepath.SkipDir
		}
		if entry.IsDir() || filepath.Ext(path) != extension {
			return nil
		}
		// 只返回所选产品的标准构建输出，排除其他产品缓存和项目中的样例 HAP。
		relative, err := filepath.Rel(project, path)
		if err != nil {
			return err
		}
		parts := strings.Split(relative, string(filepath.Separator))
		selected := false
		for index := 0; index+3 < len(parts); index++ {
			if (task == "assembleHap" && parts[index] == "build" && parts[index+1] == product && parts[index+2] == "outputs") || (task == "assembleApp" && parts[index] == "build" && parts[index+1] == "outputs" && parts[index+2] == product) {
				selected = true
				break
			}
		}
		if !selected {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		archive, err := zip.OpenReader(path)
		if err != nil {
			return fmt.Errorf("build artifact is not a valid ZIP archive")
		}
		archive.Close()
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		digest := sha256.New()
		if _, err := io.Copy(digest, file); err != nil {
			return err
		}
		result.Artifacts = append(result.Artifacts, Artifact{Path: path, Size: info.Size(), SHA256: hex.EncodeToString(digest.Sum(nil))})
		return nil
	})
	if err != nil {
		return Result{}, err
	}
	if len(result.Artifacts) == 0 {
		if task == "assembleHap" {
			return Result{}, fmt.Errorf("Hvigor exited successfully but produced no HAP artifacts")
		}
		return Result{}, fmt.Errorf("Hvigor exited successfully but produced no APP artifacts")
	}
	return result, nil
}
