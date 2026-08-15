package scaffoldv2

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

type PackageName string
type StructName string
type KebabCase string

type ModelFile struct {
	PackageName PackageName
	StructName  StructName
	IDParamName string
}

type RouteFile struct {
	PackageName       PackageName
	StructName        StructName
	ControllerImport  string
	ServiceImport     string
	RepositoryImport  string
	RouteFuncName     string
	RepositoryPackage PackageName
	ControllerPackage PackageName
	ServicePackage    PackageName
	RouteName         KebabCase
}

type ControllerFile struct {
	PackageName          PackageName
	StructName           StructName
	ControllerName       string
	ServiceInterfaceName string
}

type ServiceFile struct {
	PackageName             PackageName
	StructName              StructName
	ServiceName             string
	RepositoryInterfaceName string
	IDParamName             string
}

type RepositoryFile struct {
	PackageName    PackageName
	StructName     StructName
	RepositoryName string
	IDParamName    string
}

func buildPackageName(fileName, fileModule string) PackageName {
	// m - Model
	// r - Routes
	// c - Controller
	// s - Service
	// R - Repository
	switch fileModule {
	case "m":
		return PackageName(strings.ReplaceAll(fileName, "_", "") + "model")

	case "r":
		return PackageName(globalNormalizeNoWithUnderline(fileName) + "routes")

	case "c":
		return PackageName(globalNormalizeNoWithUnderline(fileName) + "controller")

	case "s":
		return PackageName(globalNormalizeNoWithUnderline(fileName) + "service")

	case "R":
		return PackageName(globalNormalizeNoWithUnderline(fileName) + "repository")

	default:
		return ""
	}
}

func buildStructName(fileName string) StructName {
	return StructName(globalToPascalCase(fileName))
}

func processTemplate(templateForParse, fileModule string, structForExecute interface{}) (string, error) {
	tmpl, err := template.New(moduleEquivalence(fileModule)).Parse(templateForParse)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer

	if err := tmpl.Execute(&buffer, structForExecute); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func modelContent(fileName string, isId bool) (string, error) {
	modelFile := ModelFile{
		PackageName: buildPackageName(globalNormalizeWithUnderline(fileName), "m"),
		StructName:  buildStructName(fileName),
		IDParamName: func() string {
			if isId {
				return "ID        int"
			}
			return "UUID      string"
		}(),
	}

	const modelTemplate = `package {{.PackageName}}

import "time"

type {{.StructName}} struct {
	{{.IDParamName}}
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
`
	return processTemplate(modelTemplate, "m", modelFile)
}

func routeContent(fileName, baseImportPath string) (string, error) {
	normalizedName := globalNormalizeWithUnderline(fileName)

	routeFile := RouteFile{
		PackageName:       buildPackageName(normalizedName, "r"),
		StructName:        buildStructName(fileName),
		ControllerImport:  fmt.Sprintf("%s/controller", baseImportPath),
		ServiceImport:     fmt.Sprintf("%s/service", baseImportPath),
		RepositoryImport:  fmt.Sprintf("%s/repository", baseImportPath),
		RouteFuncName:     fmt.Sprintf("Register%sRoutes", buildStructName(fileName)),
		RepositoryPackage: buildPackageName(normalizedName, "R"),
		ControllerPackage: buildPackageName(normalizedName, "c"),
		ServicePackage:    buildPackageName(normalizedName, "s"),
		RouteName:         KebabCase(globalToKebabCase(fileName)),
	}

	const routeTemplat = `package {{.PackageName}}

import (
	"net/http"

	{{.ControllerPackage}} "{{.ControllerImport}}"
	{{.RepositoryPackage}} "{{.RepositoryImport}}"
	{{.ServicePackage}} "{{.ServiceImport}}"

	"github.com/jackc/pgx/v5/pgxpool"
)

func {{.RouteFuncName}}(r *http.ServeMux, db *pgxpool.Pool) {
	repo := {{.RepositoryPackage}}.New{{.StructName}}Repository(db)
	service := {{.ServicePackage}}.New{{.StructName}}Service(repo)
	controller := {{.ControllerPackage}}.New{{.StructName}}Controller(service)

	r.HandleFunc("GET /{{.RouteName}}", controller.GetAll)
	r.HandleFunc("GET /{{.RouteName}}/{id}", controller.FindByID)
	r.HandleFunc("POST /{{.RouteName}}", controller.Create)
	r.HandleFunc("PUT /{{.RouteName}}/{id}", controller.Update)
	r.HandleFunc("DELETE /{{.RouteName}}/delete/{id}", controller.Delete)
	r.HandleFunc("PATCH /{{.RouteName}}/active/{id}", controller.Active)

}
`

	return processTemplate(routeTemplat, "r", routeFile)
}

func controllerContent(fileName string) (string, error) {
	structName := buildStructName(fileName)

	controllerFile := ControllerFile{
		PackageName:          buildPackageName(fileName, "c"),
		StructName:           structName,
		ControllerName:       fmt.Sprintf("%sController", structName),
		ServiceInterfaceName: fmt.Sprintf("%sService", structName),
	}

	const controllerTemplate = `package {{.PackageName}}

import (
	"context"
	"net/http"
)

type {{.ServiceInterfaceName}} interface {
	GetAll(context.Context) (any, error)
	Create(context.Context, any) (any, error)
	Update(context.Context, any, int) (any, error)
	FindByID(context.Context, int) (any, error)
	Delete(context.Context, int) error
	Active(context.Context, int) error
}

type {{.ControllerName}} struct {
	service {{.ServiceInterfaceName}}
}

func New{{.ControllerName}}(service {{.ServiceInterfaceName}}) *{{.ControllerName}} {
	return &{{.ControllerName}}{
		service: service,
	}
}

func (c *{{.ControllerName}}) GetAll(w http.ResponseWriter, r *http.Request) {
}

func (c *{{.ControllerName}}) Create(w http.ResponseWriter, r *http.Request) {
}

func (c *{{.ControllerName}}) FindByID(w http.ResponseWriter, r *http.Request) {
}

func (c *{{.ControllerName}}) Update(w http.ResponseWriter, r *http.Request) {
}

func (c *{{.ControllerName}}) Delete(w http.ResponseWriter, r *http.Request) {
}

func (c *{{.ControllerName}}) Active(w http.ResponseWriter, r *http.Request) {
}
`

	return processTemplate(controllerTemplate, "c", controllerFile)
}

func serviceContent(fileName string) (string, error) {
	structName := buildStructName(fileName)

	serviceFile := ServiceFile{
		PackageName:             buildPackageName(globalNormalizeWithUnderline(fileName), "s"),
		StructName:              structName,
		ServiceName:             fmt.Sprintf("%sService", structName),
		RepositoryInterfaceName: fmt.Sprintf("%sRepository", structName),
		IDParamName:             buildIDParamName(fileName),
	}

	const serviceTemplate = `package {{.PackageName}}

import "context"

type {{.RepositoryInterfaceName}} interface {
	GetAll(ctx context.Context) ([]any, error)
	Create(ctx context.Context, payload any) (any, error)
	Update(ctx context.Context, payload any, {{.IDParamName}} int) (any, error)
	FindByID(ctx context.Context, {{.IDParamName}} int) (any, error)
	Delete(ctx context.Context, {{.IDParamName}} int) error
	Active(ctx context.Context, {{.IDParamName}} int) error
}

type {{.ServiceName}} struct {
	repository {{.RepositoryInterfaceName}}
}

func New{{.ServiceName}}(repository {{.RepositoryInterfaceName}}) *{{.ServiceName}} {
	return &{{.ServiceName}}{
		repository: repository,
	}
}

func (s *{{.ServiceName}}) GetAll(ctx context.Context) ([]any, error) {
	return s.repository.GetAll(ctx)
}

func (s *{{.ServiceName}}) Create(ctx context.Context, payload any) (any, error) {
	return s.repository.Create(ctx, payload)
}

func (s *{{.ServiceName}}) Update(ctx context.Context, payload any, {{.IDParamName}} int) (any, error) {
	return s.repository.Update(ctx, payload, {{.IDParamName}})
}

func (s *{{.ServiceName}}) FindByID(ctx context.Context, {{.IDParamName}} int) (any, error) {
	return s.repository.FindByID(ctx, {{.IDParamName}})
}

func (s *{{.ServiceName}}) Delete(ctx context.Context, {{.IDParamName}} int) error {
	return s.repository.Delete(ctx, {{.IDParamName}})
}

func (s *{{.ServiceName}}) Active(ctx context.Context, {{.IDParamName}} int) error {
	return s.repository.Active(ctx, {{.IDParamName}})
}
`

	return processTemplate(serviceTemplate, "s", serviceFile)
}

func repositoryContent(fileName string) (string, error) {
	structName := buildStructName(fileName)

	repositoryFile := RepositoryFile{
		PackageName:    buildPackageName(globalNormalizeWithUnderline(fileName), "R"),
		StructName:     structName,
		RepositoryName: fmt.Sprintf("%sRepository", structName),
		IDParamName:    buildIDParamName(fileName),
	}

	const repositoryTemplate = `package {{.PackageName}}

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type {{.RepositoryName}} struct {
	db *pgxpool.Pool
}

func New{{.RepositoryName}}(db *pgxpool.Pool) *{{.RepositoryName}} {
	return &{{.RepositoryName}}{
		db: db,
	}
}

func (r *{{.RepositoryName}}) GetAll(ctx context.Context) ([]any, error) {
	return nil, nil
}

func (r *{{.RepositoryName}}) Create(ctx context.Context, payload any) (any, error) {
	return nil, nil
}

func (r *{{.RepositoryName}}) Update(ctx context.Context, payload any, {{.IDParamName}} int) (any, error) {
	return nil, nil
}

func (r *{{.RepositoryName}}) FindByID(ctx context.Context, {{.IDParamName}} int) (any, error) {
	return nil, nil
}

func (r *{{.RepositoryName}}) Delete(ctx context.Context, {{.IDParamName}} int) error {
	return nil
}

func (r *{{.RepositoryName}}) Active(ctx context.Context, {{.IDParamName}} int) error {
	return nil
}
`

	return processTemplate(repositoryTemplate, "R", repositoryFile)
}

func createContent(
	fileName,
	fileModule,
	baseImportPath string,
	isId bool,
) (string, error) {
	switch fileModule {
	case "m":
		return modelContent(fileName, isId)

	case "r":
		return routeContent(fileName, baseImportPath)

	case "c":
		return controllerContent(fileName)

	case "s":
		return serviceContent(fileName)

	case "R":
		return repositoryContent(fileName)

	default:
		return "", fmt.Errorf("Módulo não especificado.")
	}
}
