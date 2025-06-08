package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	newCmd = &cobra.Command{
		Use: "new {module | github.com/user/module}",
		Short: "create new go module",
		Long: fmt.Sprintf("create new go module\n\n%s %s\n",
			lipgloss.NewStyle().Foreground(YELLOW).Bold(true).Render("example:"),
			lipgloss.NewStyle().Italic(true).Render(
			fmt.Sprintf("gopher %s github.com/user/mymodule",
				lipgloss.NewStyle().Foreground(CYAN).Render("new"),
				),
			),
		),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Fprintf(
					os.Stderr, 
					"%s %s\n",
					lipgloss.NewStyle().Foreground(GRAY).Render("gopher:"),
					"missing module name\nrun 'gopher new -h' for usage",
				)
				os.Exit(1)
			}
			name := ""
			path := args[0]
			if strings.Contains(path, "/") {
				pathArr := strings.Split(path, "/")
				name = pathArr[len(pathArr)-1]
			} else {
				name = path
			}

			Unwrap(fmt.Printf("%sCreated '%s' module\n", PAD, name))
			Expect(os.Mkdir(name, 0755))
			Expect(os.Chdir(name))
			goCmd := exec.Command("go", "mod", "init", path)
			gitCmd := exec.Command("git", "init")
			Expect(goCmd.Run())
			Expect(gitCmd.Run())

			f1 := Unwrap(os.Create("./main.go"))
			f2 := Unwrap(os.Create("./.gitignore"))
			f3 := Unwrap(os.Create("./README.md"))
			defer f1.Close()
			defer f2.Close()
			defer f3.Close()
			Unwrap(f1.WriteString("package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}"))
			Unwrap(f2.WriteString("bin/\n*.exe\n*.exe~\n*.dll\n*.so\n*.dylib\n\n*.test\n*.out\nvendor/\n\ngo.work\ngo.work.sum\n\n.env"))
			Unwrap(f3.WriteString(fmt.Sprintf("# %s", name)))
		},
	}
)


