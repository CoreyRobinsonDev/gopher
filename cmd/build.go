package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	buildCmd = &cobra.Command{
		Use: "build [...args]",
		Short: "compile packages and dependencies",
		Long: fmt.Sprintf("compile packages and dependencies\n\n%s should be executed at the root of your module and will expect the entry point of your program to be main.go\n\n%s %s",
			lipgloss.NewStyle().Bold(true).Render("gopher build"),
			lipgloss.NewStyle().Foreground(YELLOW).Bold(true).Render("example:"),
			lipgloss.NewStyle().Italic(true).Render(
				fmt.Sprintf("gopher %s",
					lipgloss.NewStyle().Foreground(CYAN).Render("build"),
				),
			),
		),
		Run: func(cmd *cobra.Command, args []string) {
			var output []byte
			var e error
			if WebFlag {
				Expect(os.Setenv("GOOS", "js"))
				Expect(os.Setenv("GOARCH", "wasm"))
				modFile := Unwrap(os.Open("./go.mod"))
				defer modFile.Close()

				b := make([]byte, 128)
				Unwrap(modFile.Read(b))
				projectNameArr := strings.Split(strings.Split(string(b), "\n")[0], "/")
				projectName := strings.Split(projectNameArr[len(projectNameArr)-1], " ")[1]
				versionCmd := exec.Command("go", "version")
				goMajorVersion := Unwrap(strconv.Atoi(strings.Split(string(Unwrap(versionCmd.CombinedOutput())), " ")[2][4:6]))
				gorootCmd := exec.Command("go", "env", "GOROOT")
				goroot := strings.Trim(string(Unwrap(gorootCmd.CombinedOutput())), "\t\n ")

				if goMajorVersion > 23 {
					cpCmd := exec.Command("cp", goroot+"/lib/wasm/wasm_exec.js", ".")
					Unwrap(cpCmd.CombinedOutput())
				} else {
					cpCmd := exec.Command("cp", goroot+"/misc/wasm/wasm_exec.js", ".")
					Unwrap(cpCmd.CombinedOutput())
				}
				buildCmd := exec.Command("go", "build", "-o", projectName+".wasm", ".")
				output, e = buildCmd.CombinedOutput()
				buildHtml := Unwrap(os.Create(projectName+".html"))
				defer buildHtml.Close()
				Unwrap(buildHtml.WriteString(
					fmt.Sprintf(`<!DOCTYPE html>
<script src="wasm_exec.js"></script>
<script>
const go = new Go();
WebAssembly.instantiateStreaming(fetch("%s.wasm"), go.importObject).then(result => {
go.run(result.instance);
});
</script>`, projectName),
					))
				mainHtml := Unwrap(os.Create("main.html"))
				defer mainHtml.Close()
				Unwrap(mainHtml.WriteString(
					fmt.Sprintf(`<!DOCTYPE html>
<iframe src="%s.html" width="640" height="480" allow="autoplay"></iframe>`, projectName),
					))
			} else {
				buildCmd := exec.Command("go", "build")
				output, e = buildCmd.CombinedOutput()
			}
			if config.PrettyPrint && e != nil {
				outputLines := strings.Split(string(output), "\n")
				for _, line := range outputLines {
					if !strings.Contains(line, ":") {
						continue
					}
					arr := strings.Split(line, ":")
					file := arr[0]
					rownum := Unwrap(strconv.Atoi(arr[1]))
					colnum := Unwrap(strconv.Atoi(arr[2]))
					err := strings.Trim(strings.Join(arr[3:], ":"), " \t\n")
					f := Unwrap(os.Open("./" + file))
					previewLines := int(config.PrettyPrintPreviewLines)
					defer f.Close()
					reader := bufio.NewScanner(f)

					fmt.Printf("[%s]\n", lipgloss.NewStyle().Foreground(CYAN).Render(file))
					linenum := 0
					for reader.Scan() {
						linenum++
						if dif := linenum - rownum; dif <= previewLines && dif >= -previewLines {
							tabs := strings.Count(reader.Text(), "\t")
							if linenum == rownum {
								fmt.Printf("%s %s%s\n",
									lipgloss.NewStyle().Foreground(GRAY).Render(strconv.Itoa(linenum)),
									strings.Repeat("  ", tabs),
									lipgloss.NewStyle().Italic(true).Render(strings.Trim(reader.Text(), " \t")),
									)
							} else {
								fmt.Printf("%s %s%s\n",
									lipgloss.NewStyle().Foreground(GRAY).Render(strconv.Itoa(linenum)),
									strings.Repeat("  ", tabs),
									strings.Trim(reader.Text(), " \t"),
									)
							}
						}
						if linenum == rownum {
							tabs := strings.Count(reader.Text(), "\t")
							pad := strings.Repeat(" ", colnum+len(strconv.Itoa(linenum))-tabs)
							fmt.Printf("%s%s%s\n", 
								strings.Repeat("  ", tabs), 
								pad, 
								lipgloss.NewStyle().Foreground(RED).Render("^ "+err),
								)
						}
						if linenum > rownum+previewLines {
							fmt.Println()
							break
						}
					}
				}
			} else {
				fmt.Print(string(output))
			}		
		},
	}
)

func init() {
	buildCmd.
		PersistentFlags(). 
		BoolVar(&WebFlag, "web", false, "compile program to run in browser")	
}
