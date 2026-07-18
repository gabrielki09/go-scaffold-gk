package scaffold

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Input: Financial Account
// Output: financial_account
func normalizeWithUnderline(name string) string {
	name = globalTrimSpace(name)
	name = globalToLower(name)
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, " ", "_")

	return name
}

// Input: Financial Account
// Output: financialaccount
func normalizeNoWithUnderline(name string) string {
	name = globalTrimSpace(name)
	name = globalToLower(name)
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, " ", "")
	name = strings.ReplaceAll(name, "_", "")

	return name
}

func toPascalCase(value string) string {
	value = globalTrimSpace(value)
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, "-", "_")
	value = strings.ReplaceAll(value, " ", "_")

	parts := strings.Split(value, "_")

	var builder strings.Builder

	for _, part := range parts {
		part = globalTrimSpace(part)

		if part == "" {
			continue
		}

		builder.WriteString(strings.ToUpper(part[:1]))
		builder.WriteString(part[1:])
	}

	return builder.String()
}

func toKebabCase(value string) string {
	value = globalTrimSpace(value)
	value = globalToLower(value)
	value = strings.ReplaceAll(value, " ", "-")
	value = strings.ReplaceAll(value, "_", "-")

	return value
}

func buildRepoPatternNames(name string) RepoPatternNames {
	normalizedName := normalizeNoWithUnderline(name)
	snakeName := normalizeWithUnderline(name)
	pascalName := toPascalCase(name)

	return RepoPatternNames{
		NormalizedName: normalizedName,
		SnakeName:      snakeName,
		PascalName:     pascalName,

		RoutesPackage:     normalizedName + "routes",
		ControllerPackage: normalizedName + "controller",
		ServicePackage:    normalizedName + "service",
		RepositoryPackage: normalizedName + "repository",
		RequestPackage:    normalizedName + "request",
		ResponsePackage:   normalizedName + "response",

		RoutesFuncName:     fmt.Sprintf("Register%sRoutes", pascalName),
		ControllerFuncName: fmt.Sprintf("%sController", pascalName),
		ServiceFuncName:    fmt.Sprintf("%sService", pascalName),
		RepositoryFuncName: fmt.Sprintf("%sRepository", pascalName),
		RequestFuncName:    fmt.Sprintf("%sRequest", pascalName),
		ResponseFuncName:   fmt.Sprintf("%sResponse", pascalName),
	}
}

func buildRoutesContent(file File, names RepoPatternNames) string {
	baseImportPaths := fmt.Sprintf(
		"%s/%s",
		file.ModuleName,
		invertBarPath(file.RootDir),
	)

	return fmt.Sprintf(`package %s

import(
	"net/http"

	%s "%s/controller"
	%s "%s/repository"
	%s "%s/service"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

func %s(r *http.ServeMux, db *pgxpool.Pool) {
	repo := %s.New%s(db)
	service := %s.New%s(repo)
	controller := %s.New%s(service)

	r.HandleFunc("GET /%s", controller.GetAll)
	r.HandleFunc("GET /%s/{id}", controller.FindByID)
	r.HandleFunc("POST /%s", controller.Create)
	r.HandleFunc("PUT /%s/{id}", controller.Update)
	r.HandleFunc("DELETE /%s/delete/{id}", controller.Delete)
	r.HandleFunc("PATCH /%s/active/{id}", controller.Active)
}
	`,
		names.RoutesPackage,

		names.ControllerPackage,
		baseImportPaths,

		names.RepositoryPackage,
		baseImportPaths,

		names.ServicePackage,
		baseImportPaths,

		names.RoutesFuncName,
		names.RepositoryPackage,
		names.RepositoryFuncName,

		names.ServicePackage,
		names.ServiceFuncName,

		names.ControllerPackage,
		names.ControllerFuncName,

		toKebabCase(names.SnakeName),
		toKebabCase(names.SnakeName),
		toKebabCase(names.SnakeName),
		toKebabCase(names.SnakeName),
		toKebabCase(names.SnakeName),
		toKebabCase(names.SnakeName),
	)
}

func buildControllerContent(names RepoPatternNames) string {
	return fmt.Sprintf(`package %s

import (
	"context"

	"net/http"
)

type %s interface {
	GetAll(context.Context) (any, error) // add your_response in any
	Create(context.Context, any) (any, error) // add your_response and request in any
	Update(context.Context, any, int) (any, error) // add your_response and request in any
	FindByID(context.Context, int) (any, error) // add your_response in any
	Delete(context.Context, int) error
	Active(context.Context, int) error
}

type %s struct {
	service %s
}

func New%s(service %s) *%s {
	return &%s{
		service: service,
	}
}

func (c *%s) GetAll(w http.ResponseWriter, r *http.Request) {
}

func (c *%s) Create(w http.ResponseWriter, r *http.Request) {
}

func (c *%s) FindByID(w http.ResponseWriter, r *http.Request) {
}

func (c *%s) Update(w http.ResponseWriter, r *http.Request) {
}

func (c *%s) Delete(w http.ResponseWriter, r *http.Request) {
}

func (c *%s) Active(w http.ResponseWriter, r *http.Request) {
}
`,
		names.ControllerPackage,
		names.ServiceFuncName,
		names.ControllerFuncName,
		names.ServiceFuncName,
		names.ControllerFuncName,
		names.ServiceFuncName,
		names.ControllerFuncName,
		names.ControllerFuncName,
		names.ControllerFuncName,
		names.ControllerFuncName,
		names.ControllerFuncName,
		names.ControllerFuncName,
		names.ControllerFuncName,
		names.ControllerFuncName,
	)
}

func buildServiceContent(names RepoPatternNames) string {
	return fmt.Sprintf(`package %s
	
import (
	"context"

	"net/http"
)

type %s interface {
	GetAll(ctx context.Context) ([]any, error)
	Create(ctx context.Context, payload any) (any, error)
	Update(ctx context.Context, payload any, financialAccountId int) (any, error)
	FindByID(ctx context.Context, financialAccountId int) (any, error)
	Delete(ctx context.Context, financialAccountId int) error
	Active(ctx context.Context, financialAccountId int) error
}

type %s struct {
	repository %s
}

func New%s(repository %s) *%s {
	return &%s{
		repository: repository,
	}
}

func (c *%s) GetAll(w http.ResponseWriter, r *http.Request) {
}

func (c *%s) Create(w http.ResponseWriter, r *http.Request) {
}

func (c *%s) FindByID(w http.ResponseWriter, r *http.Request) {
}

func (c *%s) Update(w http.ResponseWriter, r *http.Request) {
}

func (c *%s) Delete(w http.ResponseWriter, r *http.Request) {
}

func (c *%s) Active(w http.ResponseWriter, r *http.Request) {
}
`,
		names.ServicePackage,
		names.RepositoryFuncName,
		names.ServiceFuncName,
		names.RepositoryFuncName,
		names.ServiceFuncName,
		names.RepositoryFuncName,
		names.ServiceFuncName,
		names.ServiceFuncName,
		names.ServiceFuncName,
		names.ServiceFuncName,
		names.ServiceFuncName,
		names.ServiceFuncName,
		names.ServiceFuncName,
		names.ServiceFuncName,
	)
}

func buildRepositoryContent(names RepoPatternNames) string {
	return fmt.Sprintf(`package %s

import (
	"context"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

type %s struct {
	db *pgxpool.Pool
}

func New%s(db *pgxpool.Pool) *%s {
	return &%s{
		db: db,
	}
}

func (f *%s) GetAll(ctx context.Context) ([]any, error) {
	return nil, nil
}

func (f *%s) Create(ctx context.Context, payload any) (any, error) {
	return nil, nil
}

func (f *%s) Update(ctx context.Context, payload any, financialAccountId int) (any, error) {
	return nil, nil
}

func (f *%s) FindByID(ctx context.Context, financialAccountId int) (any, error) {
	return nil, nil
}

func (f *%s) Delete(ctx context.Context, financialAccountId int) error {
	return nil
}

func (f *%s) Active(ctx context.Context, financialAccountId int) error {
	return nil
}


	`,
		names.RepositoryPackage,
		names.RepositoryFuncName,
		names.RepositoryFuncName,
		names.RepositoryFuncName,
		names.RepositoryFuncName,
		names.RepositoryFuncName,
		names.RepositoryFuncName,
		names.RepositoryFuncName,
		names.RepositoryFuncName,
		names.RepositoryFuncName,
		names.RepositoryFuncName,
	)
}

func createRepoPatternFiles(file File, option Options) error {
	repoPatternNames := buildRepoPatternNames(file.Name)

	routesFileName := fmt.Sprintf("%s_routes.go", normalizeWithUnderline(file.Name))
	routesFileContent := buildRoutesContent(file, repoPatternNames)
	routesFilePath := file.FilePaths["routes"]

	if err := createFileWithContent(routesFileName, routesFileContent, routesFilePath); err != nil {
		return err
	}

	controllerFileName := fmt.Sprintf("%s_controller.go", normalizeWithUnderline(file.Name))
	controllerFileContent := buildControllerContent(repoPatternNames)
	controllerFilePath := file.FilePaths["controller"]
	if err := createFileWithContent(controllerFileName, controllerFileContent, controllerFilePath); err != nil {
		return err
	}

	serviceFileName := fmt.Sprintf("%s_service.go", normalizeWithUnderline(file.Name))
	serviceFileContent := buildServiceContent(repoPatternNames)
	serviceFilePath := filepath.Join(file.FilePaths["service"])
	if err := createFileWithContent(serviceFileName, serviceFileContent, serviceFilePath); err != nil {
		return err
	}

	repositoryFileName := fmt.Sprintf("%s_repository.go", normalizeWithUnderline(file.Name))
	repositoryFileContent := buildRepositoryContent(repoPatternNames)
	repositoryFilePath := file.FilePaths["repo"]
	if err := createFileWithContent(repositoryFileName, repositoryFileContent, repositoryFilePath); err != nil {
		return err
	}

	return nil
}

func createControllerFile(file File, option Options) error {
	controllerFileName := fmt.Sprintf("%s_controller.go", normalizeWithUnderline(file.Name))
	controllerFileContent := buildControllerContent(buildRepoPatternNames(file.Name))

	if err := createFileWithContent(controllerFileName, controllerFileContent, file.FilePaths["controller"]); err != nil {
		return fmt.Errorf("erro ao criar o arquivo do controller: %w", err)
	}

	return nil
}

func buildModelContent(model, usageId string) string {
	modelFileName := normalizeWithUnderline(model)
	packageName := strings.ReplaceAll(modelFileName, "_", "") + "model"
	structName := toPascalCase(model) + "Model"

	content := fmt.Sprintf(`package %s

import "time"

type %s struct {
	%s
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
	`, packageName, structName, usageId)

	return content
}

func createModelFile(file File, option Options) error {
	var modelId string

	if option.Command["uuid_use"] {
		modelId = "UUID      string"
	} else {
		modelId = "ID        int"
	}

	modelFileName := fmt.Sprintf("%s_model.go", normalizeWithUnderline(file.Name))
	modelFileContet := buildModelContent(file.Name, modelId)
	modelFilePath := filepath.Join(file.FilePaths["m"], normalizeNoWithUnderline(file.Name))

	if err := createFileWithContent(modelFileName, modelFileContet, modelFilePath); err != nil {
		return fmt.Errorf("erro ao criar o arquivo do model: %w", err)
	}

	return nil
}

// 20060102150405_create_tableName.up.sql
// 20060102150405_create_tableName.down.sql
func buildMigrationContent(fileName, usageId string) (
	upFileName,
	downFileName,
	upContent,
	downContent string,
	err error,
) {
	version := time.Now().Format("20060102150405")

	upFileName = fmt.Sprintf("%s_create_%s.up.sql", version, fileName)
	downFileName = fmt.Sprintf("%s_create_%s.down.sql", version, fileName)

	upContent = fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
	%s,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	deleted_at TIMESTAMPTZ NULL
);`, fileName, usageId)

	downContent = fmt.Sprintf(`DROP TABLE IF EXISTS %s;`, fileName)

	return upFileName, downFileName, upContent, downContent, nil
}

func createMigrationFile(file File, option Options) error {
	var migrationId string

	if option.Command["uuid_use"] {
		migrationId = "id UUID PRIMARY KEY DEFAULT gen_random_uuid()"
	} else {
		migrationId = "id BIGSERIAL PRIMARY KEY"
	}

	migrationUpFileName, migrationDownFileName, migrationUpContent, migrationDownContent, err := buildMigrationContent(file.Name, migrationId)
	if err != nil {
		return err
	}

	if err := createFileWithContent(migrationUpFileName, migrationUpContent, file.FilePaths["migration"]); err != nil {
		return fmt.Errorf("erro ao criar o arquivo .up da migration: %w", err)
	}

	if err := createFileWithContent(migrationDownFileName, migrationDownContent, file.FilePaths["migration"]); err != nil {
		return fmt.Errorf("erro ao criar o arquivo .down da migration: %w", err)
	}

	return nil
}

func buildRequestContent(request string) string {
	packageName := normalizeNoWithUnderline(request) + "request"
	pascalRequest := toPascalCase(request)
	pascalRequestWithRequest := pascalRequest + "Request"

	content := fmt.Sprintf(`package %s

	type %s struct {
	}

	func (r %s) ValidatePayload() error {

		return nil
	}
	
`, packageName, pascalRequestWithRequest, pascalRequestWithRequest)

	return content
}

func createRequestFile(file File, option Options) error {
	requestFileName := fmt.Sprintf("%s_request.go", normalizeWithUnderline(file.Name))
	requestFileContent := buildRequestContent(file.Name)

	if err := createFileWithContent(requestFileName, requestFileContent, file.FilePaths["requests"]); err != nil {
		return fmt.Errorf("erro ao criar o arquivo do controller: %w", err)
	}

	return nil
}

func buildSeedContent(seed string) string {
	packageName := normalizeNoWithUnderline(seed) + "seed"
	pascalSeed := toPascalCase(seed)
	pascalSeedWithSeed := pascalSeed + "Seed"

	contentSQL := "`" + `
		INSERT INTO your_table_name ()
		VALUES ()
	` + "`"
	content := fmt.Sprintf(`package %s

	import (
		"context"

		"github.com/jackc/pgx/v5/pgxpool"
	)

	type %s struct {
		db *pgxpool.Pool
		ctx context.Context
	}

	func (s %s) %s() error {
		for i := 0; i < 50; i++ {
			if _, err := s.db.Exec(
				s.ctx,
				%s,
			); err != nil {
				return err
			}

		}

		return nil
	}
	
`, packageName, pascalSeedWithSeed, pascalSeedWithSeed, pascalSeedWithSeed, contentSQL)

	return content
}

func createSeedFile(file File, _ Options) error {
	seedFileName := fmt.Sprintf("%s_seed.go", normalizeWithUnderline(file.Name))
	seedFileContent := buildSeedContent(file.Name)

	if err := createFileWithContent(seedFileName, seedFileContent, file.FilePaths["seed"]); err != nil {
		return fmt.Errorf("erro ao criar o arquivo da seed: %w", err)
	}

	return nil
}

func buildResourceContent(seed string, usageId string) string {
	packageName := normalizeNoWithUnderline(seed) + "response"
	pascalResource := toPascalCase(seed)
	pascalWithResource := pascalResource + "Response"

	createdAtWithJson := fmt.Sprintf("CreatedAt time.Time `json:%s`", `"created_at"`)
	updatedAtWithJson := fmt.Sprintf("UpdatedAt time.Time `json:%s`", `"updated_at"`)
	deletedAtWithJson := fmt.Sprintf("DeletedAt *time.Time `json:%s`", `"deleted_at,omitempty"`)

	content := fmt.Sprintf(`package %s

import "time"

type %s struct {
	%s
	%s
	%s
	%s
}
	`, packageName, pascalWithResource, usageId, createdAtWithJson, updatedAtWithJson, deletedAtWithJson)

	return content
}

func createResourceFile(file File, option Options) error {
	var usageId string

	if option.Command["uuid_use"] {
		usageId = fmt.Sprintf("UUID      string    `json:%s`", `"uuid"`)
	} else {
		usageId = fmt.Sprintf("ID        int `json:%s`", `"id"`)
	}

	resourceFileName := fmt.Sprintf("%s_response.go", normalizeWithUnderline(file.Name))
	resourceFileContent := buildResourceContent(file.Name, usageId)

	if err := createFileWithContent(resourceFileName, resourceFileContent, file.FilePaths["resource"]); err != nil {
		return fmt.Errorf("erro ao criar o arquivo do resource: %w", err)
	}

	return nil
}

func createFileWithContent(fileName, fileContent, filePath string) error {
	if err := createInformedPath(filePath, 0755); err != nil {
		return err
	}

	fullPath := filepath.Join(filePath, fileName)

	osFile, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}

	defer osFile.Close()

	if _, err := osFile.WriteString(fileContent); err != nil {
		return err
	}

	return nil
}

func createFiles(file File, option Options) error {
	if err := createModelFile(file, option); err != nil {
		return err
	}

	if option.Command["repo"] {
		option.Command["controller"] = false
	}

	fileCreators := map[string]func(File, Options) error{
		"migration":  createMigrationFile,
		"requests":   createRequestFile,
		"seed":       createSeedFile,
		"resource":   createResourceFile,
		"repo":       createRepoPatternFiles,
		"controller": createControllerFile,
	}

	for key, enabled := range option.Command {
		if !enabled {
			continue
		}

		creator, exists := fileCreators[key]
		if !exists {
			continue
		}

		if err := creator(file, option); err != nil {
			return err
		}
	}

	return nil
}
