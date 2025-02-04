package cmd

import (
	"bufio"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed template/**/**/* template/* template/*/*/*/*
var templates embed.FS

var initCmd = &cobra.Command{
	Use:   "init [projectName]",
	Short: "Inicializa um novo Projeto Go.",
	Long: `O comando init cria a estrutura inicial do projeto,
tendo como padrão as pastas cmd, config, internal e arquivos
como padrões como Dockerfile, docker-compose e env.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]

		fmt.Print("Digite o nome do módulo Go (ex: github.com/seu-usuario/seu-projeto): ")
		reader := bufio.NewReader(os.Stdin)
		moduleName, _ := reader.ReadString('\n')
		moduleName = strings.TrimSpace(moduleName)

		fmt.Printf("Criando o projeto '%s' com módulo '%s'...\n", projectName, moduleName)

		createInitialProject(projectName, moduleName)

	},
}

func createInitialProject(projectName, moduleName string) {
	basePath := filepath.Join(".", projectName)

	folders := []string{
		"config/db",
		"internal/model",
		"internal/repository",
		"internal/service",
		"internal/api",
		"internal/api/router",
		"internal/api/handler",
		"cmd/server",
	}

	for _, folder := range folders {
		path := filepath.Join(basePath, folder)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			fmt.Printf("Erro ao criar pasta '%s': %v\n", path, err)
		} else {
			fmt.Printf("Criado: %s\n", path)
		}
	}

	data := struct {
		Module string
	}{
		Module: moduleName,
	}

	files := map[string]string{
		"go.mod":                              "template/go.mod.tmpl",
		".env":                                "template/.env.tmpl",
		"cmd/server/main.go":                  "template/cmd/server/main.go.tmpl",
		"config/config.go":                    "template/config/config.go.tmpl",
		"config/db/db.go":                     "template/config/db/db.go.tmpl",
		"internal/model/response.go":          "template/internal/model/response.go.tmpl",
		"internal/model/user.go":              "template/internal/model/user.go.tmpl",
		"internal/repository/userRepo.go":     "template/internal/repository/userRepo.go.tmpl",
		"internal/service/userService.go":     "template/internal/service/userService.go.tmpl",
		"internal/api/handler/userHandler.go": "template/internal/api/handler/userHandler.go.tmpl",
		"internal/api/router/router.go":       "template/internal/api/router/router.go.tmpl",
		"internal/api/api.go":                 "template/internal/api/api.go.tmpl",
		"docker-entrypoint.sh":                "template/docker-entrypoint.sh.tmpl",
		"Dockerfile":                          "template/Dockerfile.tmpl",
		"docker-compose.yaml":                 "template/docker-compose.yaml.tmpl",
		"makefile":                            "template/makefile.tmpl",
		".gitignore":                          "template/.gitignore.tmpl",
	}

	for file, templatePath := range files {
		outputPath := filepath.Join(basePath, file)
		generateFileFromTemplate(templatePath, outputPath, data)
	}

	runCommand(basePath, "go", "mod", "tidy")
}

func runCommand(projectPath, command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Erro ao executar '%s %v': %v\nSaída:\n%s", command, args, err, string(output))
	}

	fmt.Printf("Comando executado com sucesso: %s %v\n", command, args)
}

func generateFileFromTemplate(templatePath, outputPath string, data interface{}) {
	tmplContent, err := templates.ReadFile(templatePath)
	if err != nil {
		panic(fmt.Sprintf("Erro ao carregar template '%s': %v", templatePath, err))
	}

	tmpl, err := template.New("template").Parse(string(tmplContent))
	if err != nil {
		panic(fmt.Sprintf("Erro ao parsear template '%s': %v", templatePath, err))
	}

	file, err := os.Create(outputPath)
	if err != nil {
		panic(fmt.Sprintf("Erro ao criar arquivo '%s': %v", outputPath, err))
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		panic(fmt.Sprintf("Erro ao preencher template '%s': %v", outputPath, err))
	}

	fmt.Printf("Arquivo gerado: %s\n", outputPath)
}

func init() {
	rootCmd.AddCommand(initCmd)
}
