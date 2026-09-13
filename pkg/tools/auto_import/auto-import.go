package main

import (
	"bufio"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	packagePaths = make(map[string]bool)
	pkgMutex     sync.Mutex
	moduleName   string
	ignoreRules  []string
	ignoreREs    []*regexp.Regexp
)

// 自动读取 go.mod 中的 module 名称
func getModuleName() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		modPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(modPath); err == nil {
			file, err := os.Open(modPath)
			if err != nil {
				return ""
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			re := regexp.MustCompile(`^module\s+(\S+)`)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				match := re.FindStringSubmatch(line)
				if len(match) >= 2 {
					return match[1]
				}
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}

// 获取项目根目录（go.mod 所在目录）
func getProjectRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func main() {
	dirStr := flag.String("dir", ".", "scan directories, comma separated, e.g. -dir=./internal,./pkg,./app")
	output := flag.String("output", "auto_import.go", "output file")
	concurrency := flag.Int("concurrency", 8, "max parallel")
	ignore := flag.String("ignore", "vendor/**,.git/**,node_modules/**,testdata/**,*_test.go", "ignore patterns")
	manualModule := flag.String("module", "", "go module name (optional, auto read from go.mod)")
	flag.Parse()

	// 自动获取模块名：手动指定优先，否则自动读取
	moduleName = *manualModule
	if moduleName == "" {
		moduleName = getModuleName()
	}
	if moduleName == "" {
		fmt.Println("error: cannot find module name (not in go project or no go.mod)")
		os.Exit(1)
	}
	projectRoot := getProjectRoot()
	if projectRoot == "" {
		cwd, _ := os.Getwd()
		projectRoot = cwd
	}
	fmt.Printf("✅ using module: %s\n", moduleName)

	ignoreRules = strings.Split(*ignore, ",")
	compileIgnoreRules()
	sem := make(chan struct{}, *concurrency)
	var wg sync.WaitGroup

	dirs := strings.SplitSeq(*dirStr, ",")
	for d := range dirs {
		rootDir := strings.TrimSpace(d)
		if rootDir == "" {
			continue
		}

		absRootDir, err := filepath.Abs(rootDir)
		if err != nil {
			continue
		}

		fmt.Printf("🔍 scanning: %s\n", rootDir)

		_ = filepath.Walk(absRootDir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if f.IsDir() {
				if isIgnored(path, projectRoot) {
					return filepath.SkipDir
				}
				return nil
			}

			if !strings.HasSuffix(path, ".go") {
				return nil
			}
			if isIgnored(path, projectRoot) {
				return nil
			}

			sem <- struct{}{}
			wg.Add(1)
			go func(filePath string, rootAbs string, projRoot string) {
				defer func() { <-sem }()
				defer wg.Done()

				if hasInit(filePath) {
					fileDir := filepath.Dir(filePath)
					rel, err := filepath.Rel(rootAbs, fileDir)
					if err != nil {
						return
					}
					rel = filepath.ToSlash(rel)

					// 扫描根目录相对于项目根的相对路径（支持多级目录如 pkg/sdk）
					rootRel, err := filepath.Rel(projRoot, rootAbs)
					if err != nil {
						return
					}
					rootRel = filepath.ToSlash(rootRel)

					var pkgPath string
					if rel == "." {
						pkgPath = moduleName + "/" + rootRel
					} else {
						pkgPath = moduleName + "/" + rootRel + "/" + rel
					}

					pkgPath = filepath.Clean(pkgPath)
					pkgPath = filepath.ToSlash(pkgPath)

					pkgMutex.Lock()
					packagePaths[pkgPath] = true
					pkgMutex.Unlock()
				}
			}(path, absRootDir, projectRoot)
			return nil
		})
	}

	wg.Wait()
	writeOutput(*output)
	fmt.Printf("✅ success! generated %s, found %d init packages\n", *output, len(packagePaths))
}

func hasInit(file string) bool {
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, file, nil, 0)
	if err != nil {
		return false
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		if fn.Name.Name == "init" && fn.Type.Params.NumFields() == 0 {
			return true
		}
	}
	return false
}

// globToRegex 将 glob 模式转换为正则表达式，支持 ** 跨目录匹配。
//   - **/xxx  匹配任意层目录前缀下的 xxx
//   - xxx/**  匹配 xxx 目录本身及其下所有内容
//   - *       匹配单层内任意字符（不含路径分隔符）
func globToRegex(pattern string) string {
	pattern = regexp.QuoteMeta(pattern)
	// \*\*/ -> ([^/]*/)*  匹配任意层目录前缀
	pattern = strings.ReplaceAll(pattern, `\*\*/`, "([^/]*/)*")
	// /\*\* -> (/.*)?  匹配目录本身及其下所有内容（用于 SkipDir）
	pattern = strings.ReplaceAll(pattern, `/\*\*`, "(/.*)?")
	// 残留的 \*\* -> .*
	pattern = strings.ReplaceAll(pattern, `\*\*`, ".*")
	// \* -> [^/]*  匹配单层内任意字符
	pattern = strings.ReplaceAll(pattern, `\*`, "[^/]*")
	// \? -> [^/]
	pattern = strings.ReplaceAll(pattern, `\?`, "[^/]")
	return "^" + pattern + "$"
}

// compileIgnoreRules 预编译忽略规则为正则，在解析 ignoreRules 后调用。
func compileIgnoreRules() {
	ignoreREs = make([]*regexp.Regexp, 0, len(ignoreRules))
	for _, r := range ignoreRules {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		// 去掉相对路径前缀 ./，统一基于项目根进行匹配
		r = strings.TrimPrefix(r, "./")
		re, err := regexp.Compile(globToRegex(r))
		if err != nil {
			fmt.Printf("⚠️ invalid ignore pattern: %s\n", r)
			continue
		}
		ignoreREs = append(ignoreREs, re)
	}
}

// isIgnored 判断路径是否命中忽略规则。
// path 为绝对路径，projRoot 为项目根绝对路径；
// 内部转换为相对项目根的路径后再用正则匹配，
// 使 ./xxx/** 与 **/xxx/** 等规则均可对目录与文件生效。
func isIgnored(path, projRoot string) bool {
	rel, err := filepath.Rel(projRoot, path)
	if err != nil {
		rel = path
	}
	rel = filepath.ToSlash(rel)
	for _, re := range ignoreREs {
		if re.MatchString(rel) {
			return true
		}
	}
	return false
}

func writeOutput(out string) {
	f, _ := os.Create(out)
	defer f.Close()
	w := bufio.NewWriter(f)

	_, _ = w.WriteString("// Auto-generated by init scanner. DO NOT EDIT.\n\npackage main\n\nimport (\n")
	for pkg := range packagePaths {
		_, _ = fmt.Fprintf(w, "\t_ %q\n", pkg)
	}
	_, _ = w.WriteString(")\n")
	_ = w.Flush()
}
