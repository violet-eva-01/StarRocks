// Package ast @author: Violet-Eva @date  : 2025/5/21 @notes :
package ast

import (
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/violet-eva-01/StarRocks/sql/parser"
	"strings"
)

func GetObjectNameList(catalog, database string, ctx antlr.ParserRuleContext) []TableName {
	var tns []TableName
	switch ctx.(type) {
	case *parser.PrivObjectNameListContext:
		for _, objectName := range ctx.(*parser.PrivObjectNameListContext).AllPrivObjectName() {
			tn := NewTableNameWithQualifiedName(catalog, database, GetQualifiedName(objectName))
			tns = append(tns, *tn)
		}
	}
	return tns
}

type PrivilegeType int

const (
	All PrivilegeType = iota
	Alter
	Apply
	Blacklist
	Create
	Delete
	Drop
	Export
	File
	Impersonate
	Insert
	Grant
	Node
	Operate
	Plugin
	Repository
	Refresh
	Select
	Update
	Usage
)

var privilegeTypeNames = map[string]PrivilegeType{
	"ALL":         All,
	"ALTER":       Alter,
	"APPLY":       Apply,
	"BLACKLIST":   Blacklist,
	"CREATE":      Create,
	"DELETE":      Delete,
	"DROP":        Drop,
	"EXPORT":      Export,
	"FILE":        File,
	"IMPERSONATE": Impersonate,
	"INSERT":      Insert,
	"GRANT":       Grant,
	"NODE":        Node,
	"OPERATE":     Operate,
	"PLUGIN":      Plugin,
	"REPOSITORY":  Repository,
	"REFRESH":     Refresh,
	"SELECT":      Select,
	"UPDATE":      Update,
	"USAGE":       Usage,
}

func ParsePrivilegeType(str string) PrivilegeType {
	str = strings.ToUpper(str)
	if privilegeTypeName, ok := privilegeTypeNames[str]; ok {
		return privilegeTypeName
	} else {
		return -1
	}
}

func (a PrivilegeType) String() string {
	return findName(a, privilegeTypeNames)
}

type CreatePrivilegeType int

const (
	CreateDatabase CreatePrivilegeType = iota
	CreateTable
	CreateView
	CreateFunction
	CreateGlobalFunction
	CreateMaterializedView
	CreateResource
	CreateResourceGroup
	CreateExternalCatalog
	CreateStorageVolume
	CreateWarehouse
	CreateCNGroup
	CreatePipe
)

var createPrivilegeTypeNames = map[string]CreatePrivilegeType{
	"DATABASE":          CreateDatabase,
	"TABLE":             CreateTable,
	"VIEW":              CreateView,
	"FUNCTION":          CreateFunction,
	"GLOBAL FUNCTION":   CreateGlobalFunction,
	"MATERIALIZED VIEW": CreateMaterializedView,
	"RESOURCE":          CreateResource,
	"RESOURCE GROUP":    CreateResourceGroup,
	"EXTERNAL CATALOG":  CreateExternalCatalog,
	"STORAGE VOLUME":    CreateStorageVolume,
	"WAREHOUSE":         CreateWarehouse,
	"CN GROUP":          CreateCNGroup,
	"PIPE":              CreatePipe,
}

func ParseCreatePrivilegeType(str string) CreatePrivilegeType {
	str = strings.ToUpper(str)
	if privilegeType, ok := createPrivilegeTypeNames[str]; ok {
		return privilegeType
	} else {
		return -1
	}
}

func (a CreatePrivilegeType) String() string {
	return findName(a, createPrivilegeTypeNames)
}

func GetTypeList(ctx antlr.ParserRuleContext) map[PrivilegeType]CreatePrivilegeType {
	ptList := make(map[PrivilegeType]CreatePrivilegeType)
	switch ctx.(type) {
	case *parser.PrivilegeTypeListContext:
		for _, privType := range ctx.(*parser.PrivilegeTypeListContext).AllPrivilegeType() {
			var privNames []string
			for _, typeValue := range privType.GetChildren() {
				privNames = append(privNames, fmt.Sprintf("%s", typeValue))
			}
			pt := ParsePrivilegeType(privNames[0])
			if pt == Create {
				cptName := strings.Join(privNames[1:], " ")
				cpt := ParseCreatePrivilegeType(cptName)
				ptList[pt] = cpt
			} else {
				ptList[pt] = -1
			}
		}
	}
	return ptList
}

type privilegeOnType int

const (
	OnRole privilegeOnType = iota
	OnUser
	OnTableBrief
	OnFunc
	OnSystem
	OnPrimaryObj
	OnAll
)

var privilegeOnTypeNames = map[string]privilegeOnType{
	"ROLE":     OnRole,
	"USER":     OnUser,
	"TABLE":    OnTableBrief,
	"FUNCTION": OnFunc,
	"SYSTEM":   OnSystem,
	"PRIMARY":  OnPrimaryObj,
	"ALL":      OnAll,
}

func ParsePrivilegeOnType(str string) privilegeOnType {
	str = strings.ToUpper(str)
	if privilegeOnTypeName, ok := privilegeOnTypeNames[str]; ok {
		return privilegeOnTypeName
	} else {
		return -1
	}
}

func (a privilegeOnType) String() string {
	return findName(a, privilegeOnTypeNames)
}

type PrivilegeObjectType int

const (
	Catalog PrivilegeObjectType = iota
	Database
	MaterializedView
	Resource
	ResourceGroup
	StorageVolume
	System
	Table
	View
	Warehouse
	Pipe
)

var privilegeObjectTypeNames = map[string]PrivilegeObjectType{
	"CATALOG":           Catalog,
	"DATABASE":          Database,
	"MATERIALIZED VIEW": MaterializedView,
	"RESOURCE":          Resource,
	"RESOURCE GROUP":    ResourceGroup,
	"STORAGE VOLUME":    StorageVolume,
	"SYSTEM":            System,
	"TABLE":             Table,
	"VIEW":              View,
	"WAREHOUSE":         Warehouse,
	"PIPE":              Pipe,
}

func ParsePrivilegeObjectType(str string) PrivilegeObjectType {
	str = strings.ToUpper(str)
	if privilegeObjectTypeName, ok := privilegeObjectTypeNames[str]; ok {
		return privilegeObjectTypeName
	} else {
		return -1
	}
}

func (a PrivilegeObjectType) String() string {
	return findName(a, privilegeObjectTypeNames)
}

type PrivilegeObjectTypePlural int

const (
	Catalogs PrivilegeObjectTypePlural = iota
	Databases
	Functions
	GlobalFunctions
	MaterializedViews
	Policies
	Resources
	ResourceGroups
	StorageVolumes
	Tables
	Users
	Views
	Warehouses
	Pipes
)

var privilegeObjectTypePluralNames = map[string]PrivilegeObjectTypePlural{
	"CATALOGS":           Catalogs,
	"DATABASE":           Databases,
	"FUNCTIONS":          Functions,
	"GLOBALS":            GlobalFunctions,
	"MATERIALIZED VIEWS": MaterializedViews,
	"POLICIES":           Policies,
	"RESOURCES":          Resources,
	"RESOURCES GROUPS":   ResourceGroups,
	"STORAGE VOLUME":     StorageVolumes,
	"TABLES":             Tables,
	"USERS":              Users,
	"VIEWS":              Views,
	"WAREHOUSE":          Warehouses,
	"PIPES":              Pipes,
}

func ParsePrivilegeObjectTypePlural(str string) PrivilegeObjectTypePlural {
	str = strings.ToUpper(str)
	if privilegeObjectTypePlural, ok := privilegeObjectTypePluralNames[str]; ok {
		return privilegeObjectTypePlural
	} else {
		return -1
	}
}

func (a PrivilegeObjectTypePlural) String() string {
	return findName(a, privilegeObjectTypePluralNames)
}
