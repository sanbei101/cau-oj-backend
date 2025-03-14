package judge

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/panjf2000/ants/v2"
)

type CompiledJudge struct{}

var CompiledJudgeApp = new(CompiledJudge)

var (
	maxPoolSize       = 10
	CompiledJudgePool *ants.PoolWithFunc
)

type CompiledJudgeTask struct {
	executablePath string
	input          string
	output         string
	err            error
	done           chan bool
	timeout        time.Duration
}

func init() {
	CompiledJudgePool, _ = ants.NewPoolWithFunc(maxPoolSize, func(task interface{}) {
		t := task.(*CompiledJudgeTask)
		// 执行任务
		t.output, t.err = t.run()
		t.done <- true
	})
}

// 运行已编译的可执行文件并返回输出
func (t *CompiledJudgeTask) run() (string, error) {
	// 创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, t.executablePath)
	var stdout, stderr bytes.Buffer
	cmd.Stdin = strings.NewReader(t.input)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// 等待命令执行结果或者超时
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("超时错误: %v", err)
		}
		return "", fmt.Errorf("运行时错误: %v, %s", err, stderr.String())
	}

	return stdout.String(), nil
}
func (j *CompiledJudge) SubmitJudge(executablePath string, input string) (string, error) {
	task := &CompiledJudgeTask{
		executablePath: executablePath,
		input:          input,
		done:           make(chan bool),
		timeout:        time.Second,
	}

	err := CompiledJudgePool.Invoke(task)
	if err != nil {
		return "", fmt.Errorf("无法提交任务到协程池: %v", err)
	}

	<-task.done

	return task.output, task.err
}

// CompileCCode 编译C代码并返回可执行文件的路径
func (j *CompiledJudge) CompileCCode(codeContent []byte) (string, error) {
	// 生成唯一的临时文件名
	codeFile, err := os.CreateTemp("", "user_code_*.c")
	if err != nil {
		return "", fmt.Errorf("创建临时C文件失败: %v", err)
	}
	// 删除临时文件
	defer os.Remove(codeFile.Name())

	// 写入C代码
	if _, err := codeFile.Write(codeContent); err != nil {
		return "", fmt.Errorf("写入C代码失败: %v", err)
	}
	codeFile.Close()

	// 生成可执行文件路径
	executableName := fmt.Sprintf("c_out_%s", uuid.New().String())
	executablePath := filepath.Join(os.TempDir(), executableName)

	// 编译C代码
	cmd := exec.Command("gcc", codeFile.Name(), "-o", executablePath)
	var compileStderr bytes.Buffer
	cmd.Stderr = &compileStderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("编译失败: %v, %s", err, compileStderr.String())
	}

	return executablePath, nil
}

// CompileCppCode 编译C++代码并返回可执行文件的路径
func (j *CompiledJudge) CompileCppCode(codeContent []byte) (string, error) {
	// 生成唯一的临时文件名
	codeFile, err := os.CreateTemp("", "user_code_*.cpp")
	if err != nil {
		return "", fmt.Errorf("创建临时C++文件失败: %v", err)
	}
	// 删除临时文件
	defer os.Remove(codeFile.Name())

	// 写入C++代码
	if _, err := codeFile.Write(codeContent); err != nil {
		return "", fmt.Errorf("写入C++代码失败: %v", err)
	}
	codeFile.Close()

	// 生成可执行文件路径
	executableName := fmt.Sprintf("cpp_out_%s", uuid.New().String())
	executablePath := filepath.Join(os.TempDir(), executableName)

	// 编译C++代码
	cmd := exec.Command("g++", "-O3", codeFile.Name(), "-o", executablePath)
	var compileStderr bytes.Buffer
	cmd.Stderr = &compileStderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("编译失败: %v, %s", err, compileStderr.String())
	}

	return executablePath, nil
}
